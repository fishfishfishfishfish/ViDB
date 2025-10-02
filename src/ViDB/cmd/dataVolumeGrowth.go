/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bcds/go-hpc-vidb/common"
	"github.com/spf13/cobra"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"
)

// dataVolumeGrowthCmd represents the dataVolumeGrowth command
var dataVolumeGrowthCmd = &cobra.Command{
	Use:   "growth",
	Short: "growth",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		executeDataVolumeGrowth()
	},
}

var loadRound int = 0
func init() {
	rootCmd.AddCommand(dataVolumeGrowthCmd)
	dataVolumeGrowthCmd.Flags().Uint64Var(&cacheCost, "cacheCost", 1<<30, "全局Cache的内存空间大小")
	dataVolumeGrowthCmd.Flags().IntVar(&VlogSize, "VlogSize", 1, "VlogSize文件的大小(GB)")
	dataVolumeGrowthCmd.Flags().IntVar(&operationCount, "operationCount", 10*vidbsvc.M, "一轮数据写入总条目数")
	dataVolumeGrowthCmd.Flags().IntVar(&loadRound, "loadRound", 10, "数据写入的轮次")
	dataVolumeGrowthCmd.Flags().Float64Var(&p1, "zipf", 1.0, "zipf分布参数")
	dataVolumeGrowthCmd.Flags().IntVar(&batchSize, "batchSize", 500, "批次大小")
	dataVolumeGrowthCmd.Flags().IntVar(&batchCount, "batchCount", 20, "批次数量")
	dataVolumeGrowthCmd.Flags().IntVar(&keySize, "keySize", 32, "写入Key的大小")
	dataVolumeGrowthCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "写入Value的大小")
	dataVolumeGrowthCmd.Flags().StringVar(&dataPath, "dataPath", "testdata/paper", "存储路径")
}

// const header = "TotalCount\tTotalPutLan(ms)\tTotalTPS\tBatchSize\tPutLan(ms)\tPutTPS\tGetLan(ms)\tGetTPS\n"
// const res = "%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n"

const header = "Operation,Latency(s),TPS\n"
const tRes = "Load[%d],%.2f,%.2f\n"
const rRes = "Read[%d],%.2f,%.2f\n"
const wRes = "Write[%d],%.2f,%.2f\n"

// 执行 数据量增长测试
func executeDataVolumeGrowth() {
	dataPath = filepath.Join(dataPath, "dataGrowth")
	// defaultConfig := vidbconfig.GetDefaultConfig()
	loadBatchSize := 5000
	perTreeMetaNum := (operationCount*loadRound / 5) / loadBatchSize // 5 trees
	defaultConfig := prepareConfig(uint64(perTreeMetaNum), uint64(loadBatchSize))
	defaultConfig.DataPath = dataPath
	defaultConfig.MaxCost = cacheCost
	defaultConfig.VlogSize = uint64(VlogSize) * common.GiB
	defaultConfig.VSize = uint32(valueSize)

	ins, err := vidbsvc.GetVIDBInstance(defaultConfig)
	if err != nil {
		panic(err)
	}
	
	fmt.Print(header)
	for round := 0; round < loadRound; round++ {
		duration, err := vidbsvc.MicroWrite4Growth(ins, round*operationCount, operationCount, loadBatchSize, keySize, valueSize, dataPath)
		if err != nil {
			panic(err)
		}
		totalWriteTps := int(float64(operationCount) / duration.Seconds())
		fmt.Println(header)
		fmt.Printf(tRes, operationCount, duration.Seconds(), float64(totalWriteTps))
		// 计算数据目录大小
		dataDirSize := calculateDirSize(dataPath)
		fmt.Printf("DataDirSize: %d bytes\n", dataDirSize)

		// [TODO]: 获取LRU cache\Node cache\FreeMap\Pending list\allocating list\Vlog index内存大小

		if batchCount <= 0 {
			fmt.Printf("按任意键退出...\n")
			b := make([]byte, 1)
			os.Stdin.Read(b)
		}
		
		readTPS_sum := 0.0
		readLat_sum := int64(0)
		for i := 0; i < batchCount; i++ {
			zipfKeys := genRandomKeyByZipf(zipfType, batchSize, operationCount, p1)
			read_duration, err := vidbsvc.ReadBatchKey(ins, zipfKeys)
			// read_duration, err := vidbsvc.RandomRead(ins, operationCount, batchSize, keySize, valueSize)
			if err != nil {
				panic(err)
			}
			readTPS := float64(batchSize) / (float64(read_duration.Nanoseconds())) * 1e9
			readTPS_sum += readTPS
			readLat_sum += read_duration.Nanoseconds()
			fmt.Printf(rRes, batchSize, read_duration.Seconds(), readTPS)
		}

		
		writeTPS_sum := 0.0
		writeLat_sum := int64(0)
		for i := 0; i < batchCount; i++ {
			write_duration, err := vidbsvc.RandomWrite(ins, operationCount, batchSize, keySize, valueSize)
			if err != nil {
				panic(err)
			}
			writeTPS := float64(batchSize) / (float64(write_duration.Nanoseconds())) * 1e9
			writeTPS_sum += writeTPS
			writeLat_sum += write_duration.Nanoseconds()
			fmt.Printf(wRes, batchSize, write_duration.Seconds(), writeTPS)
		}
		fmt.Println("Average")
		fmt.Printf(rRes, batchSize, float64(readLat_sum)/float64(batchCount), readTPS_sum/float64(batchCount))
		fmt.Printf(wRes, batchSize, float64(writeLat_sum)/float64(batchCount), writeTPS_sum/float64(batchCount))
		// fmt.Printf(res, operationCount, duration.Milliseconds(), totalWriteTps, batchSize, write.Milliseconds(), putTPS, read.Milliseconds(), readTPS)
	}

}


func calculateDirSize(path string) int64 {
 var totalSize int64
   err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
      if !info.IsDir() {
         totalSize += info.Size()
      }
      return nil
   })
   if err != nil {
      fmt.Println(err)
      return 0
   }
   return totalSize
}