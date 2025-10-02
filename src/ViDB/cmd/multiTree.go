/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	vidbconfig "github.com/bcds/go-hpc-vidb/config"
	"github.com/pingcap/go-ycsb/pkg/generator"
	"github.com/pingcap/go-ycsb/pkg/ycsb"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"

	"github.com/spf13/cobra"
)

// multiTreeCmd represents the multiTree command
var multiTreeCmd = &cobra.Command{
	Use:   "multiTree",
	Short: "mt",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		executeMultiTree()
	},
}

var multiTreeNum int
var n1 int
var n2 int
var n3 int
var n4 int
var treeHeight uint64
var bloomCap uint

func init() {
	rootCmd.AddCommand(multiTreeCmd)
	multiTreeCmd.Flags().Uint64Var(&cacheCost, "cacheCost", 1<<30, "全局Cache的内存空间大小")
	multiTreeCmd.Flags().IntVar(&multiTreeNum, "multiTree", 1, "Partition的数量")
	multiTreeCmd.Flags().IntVar(&n1, "n_1", 200000, "树的总量")
	multiTreeCmd.Flags().Uint64Var(&treeHeight, "treeHeight", 2, "树的高度")
	multiTreeCmd.Flags().UintVar(&bloomCap, "bloomCap", 5000, "布隆过滤器的大小")
	multiTreeCmd.Flags().IntVar(&VlogSize, "VlogSize", 1<<30, "VlogSize文件的大小(GB)")
	multiTreeCmd.Flags().IntVar(&n2, "n_2", 10*vidbsvc.M, "数据写入总数")
	multiTreeCmd.Flags().IntVar(&n3, "n_3", 500, "写入后的操作的 批次大小")
	multiTreeCmd.Flags().Float64Var(&p1, "zipf", 1.0, "zipf分布参数")
	multiTreeCmd.Flags().IntVar(&batchSize, "batchSize", 5000, "批次大小")
	multiTreeCmd.Flags().IntVar(&keySize, "keySize", 32, "写入Key的大小")
	multiTreeCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "写入Value的大小")
	multiTreeCmd.Flags().IntVar(&batchCount, "batchCount", 20, "测试的数量")
	multiTreeCmd.Flags().StringVar(&dataPath, "dataPath", "testdata/paper", "存储路径")
}

// var multiTreeHeader = "put[%d] latency(ms),put throughput,put[%d] latency(ms),put throughput,get[%d] latency(ms),get throughput\n"
var multiTreeHeader = "operations,latency(ms),throughput\n"
var multiTreeBodyT = "put[%d],%.2f,%.2f\n"
var multiTreeBodyW = "put[%d],%.2f,%.2f\n"
var multiTreeBodyR = "get[%d],%.2f,%.2f\n"

func executeMultiTree() {
	dataPath = filepath.Join(dataPath, "multiTree")
	// 这里需要去计算 数的版本
	currentMetaNum := n1 / batchSize
	config := prepareConfig(uint64(currentMetaNum), uint64(batchSize))
	// config.RotateStrategyTypeName = "adaptive"
	// config.RotateThRootFillPercent = 0.75
	// config.RotateThTreeHeight = treeHeight
	config.BloomCap = bloomCap
	config.BloomRate = 0.01
	config.DataPath = dataPath
	config.MaxCost = cacheCost
	config.VlogSize = uint64(VlogSize)
	config.VSize = uint32(valueSize)
	ins, err := vidbsvc.GetVIDBInstance(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total %d\n", n2)
	fmt.Printf("treeCap %d\n", n1)
	fmt.Printf("batchSize %d\n", batchSize)
	fmt.Printf("RotateThMetaNum %d\n", config.RotateThMetaNum)
	fmt.Printf("BloomCap %d\n", config.BloomCap)
	fmt.Printf("BloomRate %f\n", config.BloomRate)
	fmt.Printf("VlogSize %d\n", config.VlogSize)
	fmt.Printf("cacheCost %d\n", config.MaxCost)

	totalW, err := vidbsvc.MicroWrite(ins, n2, batchSize, keySize, valueSize)
	if err != nil {
		panic(err)
	}
	tTps := float64(n2) / float64(totalW.Nanoseconds()) * 1e9
	fmt.Printf(multiTreeHeader)
	fmt.Printf(multiTreeBodyT, n2, float64(totalW.Milliseconds()), tTps)

	var key_generator ycsb.Generator
	switch zipfType {
	case "zipfian":
		key_generator = generator.NewZipfianWithItems(int64(n1), p1)
	case "scrambled":
		key_generator = generator.NewScrambledZipfian(0, int64(n1), p1)
	default:
		key_generator = generator.NewZipfianWithItems(int64(n1), p1)
	}

	rTps_sum := 0.0
	wTps_sum := 0.0
	rLat_sum := 0.0
	wLat_sum := 0.0
	for i := 0; i < batchCount; i++ {
		// zipfKeys := genRandomKeywGenerator(key_generator, n3, n1, p1)
		zipfKeys := genRandomKeywGeneratorMoreOld(key_generator, n3)
		read, err := vidbsvc.ReadBatchKey(ins, zipfKeys)
		// read, err := vidbsvc.RandomRead(ins, n2, n3, keySize, valueSize)
		if err != nil {
			panic(err)
		}
		rTps := float64(n3) / float64(read.Nanoseconds()) * 1e9
		rTps_sum += rTps
		rLat_sum += float64(read.Milliseconds())
		fmt.Printf(multiTreeBodyR, n3, float64(read.Milliseconds()), rTps)
	}

	for i := 0; i < batchCount; i++ {
		write, err := vidbsvc.RandomWrite(ins, n2, n3, keySize, valueSize)
		if err != nil {
			panic(err)
		}
		wTps := float64(n3) / float64(write.Nanoseconds()) * 1e9
		wTps_sum += wTps
		wLat_sum += float64(write.Milliseconds())
		fmt.Printf(multiTreeBodyW, n3, float64(write.Milliseconds()), wTps)
	}
	fmt.Println("Average")
	fmt.Printf(multiTreeBodyR, n3, rLat_sum/float64(batchCount), rTps_sum/float64(batchCount))
	fmt.Printf(multiTreeBodyW, n3, wLat_sum/float64(batchCount), wTps_sum/float64(batchCount))
	// fmt.Printf(multiTreeHeaderBody, totalW.Milliseconds(), tTps, wLat_sum/float64(batchCount), wTps_sum/float64(batchCount), rLat_sum/float64(batchCount), rTps_sum/float64(batchCount))
}

func prepareConfig(metaNum, batchSize uint64) *vidbconfig.VidbConfig {
	start := time.Now()
	defer func() {
		log.Printf("Config preparation took %v", time.Since(start))
	}()
	defaultConfig := vidbconfig.GetDefaultConfig()
	// 根据当前的 版本数 来计算写入的总量
	defaultConfig.RotateThMetaNum = metaNum
	defaultConfig.BloomCap = uint(float64(metaNum*batchSize) * 1.1)
	defaultConfig.BloomRate = 0.01
	defaultConfig.MaxCost = 1 << 30
	return defaultConfig
}
