package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	vidbconfig "github.com/bcds/go-hpc-vidb/config"

	"github.com/bcds/go-hpc-vidb/common"

	"github.com/spf13/cobra"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"
)

// rangeQueryCmd represents the rangeQuery command
var rangeGetCmd = &cobra.Command{
	Use:   "range_get",
	Short: "range_get",
	Long:  `该命令用于执行范围查询操作。`,
	Run: func(cmd *cobra.Command, args []string) {
		executeRangeGet(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(rangeGetCmd)
	// 定义命令行参数
	rangeGetCmd.Flags().Uint64Var(&cacheCost, "cacheCost", 1<<30, "cache的内存空间大小")
	rangeGetCmd.Flags().IntVar(&VlogSize, "VlogSize", 1, "VlogSize文件的大小(GB)")
	rangeGetCmd.Flags().IntVar(&operationCount, "operationCount", 10*vidbsvc.M, "操作次数")
	rangeGetCmd.Flags().IntVar(&batchSize, "batchSize", 500, "数据写入的批次大小")
	rangeGetCmd.Flags().IntVar(&batchCount, "BatchCount", 20, "批次数量")
	rangeGetCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "值的大小")
	rangeGetCmd.Flags().IntVar(&keySize, "keySize", 32, "键的大小")
	rangeGetCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "存储路径")
	rangeGetCmd.Flags().IntSliceVar(&rs, "rs", []int{5, 50, 100, 200, 300, 400, 500, 1000, 2000}, "范围查询的范围列表，多个值用逗号分隔，例如 5,50,100")
}

func executeRangeGet(cmd *cobra.Command, args []string) {
	config := vidbconfig.GetDefaultConfig()
	config.DataPath = dataPath
	config.MaxCost = cacheCost
	config.VlogSize = uint64(VlogSize) * common.GiB

	if err := os.RemoveAll(config.DataPath); err != nil {
		panic(err)
	}

	instance, err := vidbsvc.GetVIDBInstance(config)
	if err != nil {
		panic(err)
	}

	duration, err := vidbsvc.MicroWrite(instance, operationCount, batchSize, keySize, valueSize)
	fmt.Println(fmt.Sprintf("Execute Write %d done. Lantency: %d ms", operationCount, duration.Milliseconds()))
	if err != nil {
		panic(err)
	}

	// warm up
	_, err = vidbsvc.RandomRead(instance, operationCount, 5000, keySize, valueSize)
	if err != nil {
		panic(err)
	}

	for _, r := range rs {
		for i := 0; i < batchCount; i++ {
			duration, err := vidbsvc.MicroRRead(instance, operationCount, keySize, r)
			if err != nil {
				panic(err)
			}
			throughput := float64(r) / duration.Seconds() // TPS 计算 总条目数 / 总时间
			fmt.Printf("Latency: %d ns (%.6f s), range size: %d, TPS: %.2f\n", duration.Nanoseconds(), duration.Seconds(), r, throughput)
		}
		time.Sleep(time.Second)
	}
	_ = instance.Close()
}
