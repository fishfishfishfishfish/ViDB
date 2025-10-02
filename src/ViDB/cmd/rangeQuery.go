package cmd

import (
	"fmt"
	"github.com/bcds/go-hpc-vidb/common"
	"os"
	"path/filepath"
	"time"

	vidbconfig "github.com/bcds/go-hpc-vidb/config"
	"github.com/spf13/cobra"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"
)

// rangeQueryCmd represents the rangeQuery command
var rangeQueryCmd = &cobra.Command{
	Use:   "range_query",
	Short: "range_query",
	Long:  `该命令用于执行范围查询操作 (使用循环get的方式)。`,
	Run: func(cmd *cobra.Command, args []string) {
		executeRangeQuery(cmd, args)
	},
}

var recover bool
var VlogSize int
var operationCount int
var batchSize int
var metaNum int
var batchCount int
var keySize int
var valueSize int
var dataPath string
var rs []int

func init() {
	rootCmd.AddCommand(rangeQueryCmd)
	// 定义命令行参数
	rangeQueryCmd.Flags().Uint64Var(&cacheCost, "cacheCost", 1<<30, "cache的内存空间大小")
	rangeQueryCmd.Flags().IntVar(&VlogSize, "VlogSize", 1, "VlogSize文件的大小(GB)")
	rangeQueryCmd.Flags().IntVar(&operationCount, "operationCount", 10*vidbsvc.M, "操作次数")
	rangeQueryCmd.Flags().IntVar(&batchSize, "batchSize", 500, "数据写入的批次大小")
	rangeQueryCmd.Flags().IntVar(&batchCount, "BatchCount", 20, "批次数量")
	rangeQueryCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "值的大小")
	rangeQueryCmd.Flags().IntVar(&keySize, "keySize", 32, "键的大小")
	rangeQueryCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "存储路径")
	rangeQueryCmd.Flags().IntSliceVar(&rs, "rs", []int{5, 50, 100, 200, 300, 400, 500, 1000, 2000}, "范围查询的范围列表，多个值用逗号分隔，例如 5,50,100")
}

func executeRangeQuery(cmd *cobra.Command, args []string) {
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

	// _ = instance.Close()
	// instance, err = vidbsvc.GetVIDBInstance(config)

	// warm up
	_, err = vidbsvc.RandomRead(instance, operationCount, 5000, keySize, valueSize)
	if err != nil {
		panic(err)
	}

	for _, r := range rs {
		for i := 0; i < batchCount; i++ {
			duration, err := vidbsvc.MicroRangeQuery(instance, operationCount, batchSize, keySize, r)
			if err != nil {
				panic(err)
			}
			throughput := float64(r) / duration.Seconds() // TPS 计算 总条目数 / 总时间
			fmt.Printf("Latency: %d ns (%.6f s), range size: %d, TPS: %.2f\n", duration.Nanoseconds(), duration.Seconds(), r, throughput)
			time.Sleep(time.Second)
		}
	}
	_ = instance.Close()
}
