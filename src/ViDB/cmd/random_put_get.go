/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	// "runtime"
	"time"

	"github.com/bcds/go-hpc-vidb/common"
	"github.com/spf13/cobra"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"
)

// writeCmd represents the write command
var randompgtCmd = &cobra.Command{
	Use:   "randompgt",
	Short: "randompgt",
	Long:  `letus 执行 vidb 单数据随机写入测试`,
	Run: func(cmd *cobra.Command, args []string) {
		executeRandomPGt(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(randompgtCmd)
	randompgtCmd.Flags().Uint64Var(&cacheCost, "cacheCost", 1<<30, "cache的内存空间大小")
	randompgtCmd.Flags().IntVar(&VlogSize, "VlogSize", 1, "VlogSize文件的大小(GB)")
	randompgtCmd.Flags().IntVar(&operationCount, "operationCount", 10*vidbsvc.M, "操作次数")
	randompgtCmd.Flags().IntVar(&batchCount, "BatchCount", 20, "批次数量")
	randompgtCmd.Flags().IntVar(&batchSize, "batchSize", 500, "批次大小")
	randompgtCmd.Flags().IntVar(&keySize, "keySize", 32, "键的大小")
	randompgtCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "值的大小")
	randompgtCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "存储路径")
}

func executeRandomPGt(cmd *cobra.Command, args []string) {
	// config := vidbconfig.GetDefaultConfig()
	loadBatchSize := 5000
	perTreeMetaNum := (operationCount) / loadBatchSize // 1 trees
	config := prepareConfig(uint64(perTreeMetaNum), uint64(loadBatchSize))
	config.DataPath = dataPath
	config.MaxCost = cacheCost
	config.VlogSize = uint64(VlogSize) * common.GiB
	config.VSize = uint32(valueSize)
	if err := os.RemoveAll(config.DataPath); err != nil {
		panic(err)
	}

	// var m runtime.MemStats
	// runtime.ReadMemStats(&m)
	// fmt.Println(fmt.Sprintf("Before load, memory usage: %d MB", m.Alloc/1024/1024))

	instance, err := vidbsvc.GetVIDBInstance(config)
	if err != nil {
		panic(err)
	}

	duration, err := vidbsvc.MicroWrite(instance, operationCount, loadBatchSize, keySize, valueSize)
	fmt.Println(fmt.Sprintf("Execute Write %d done. Lantency: %d ms", operationCount, duration.Milliseconds()))
	if err != nil {
		panic(err)
	}

	// warm up
	_, err = vidbsvc.RandomRead(instance, operationCount, loadBatchSize, keySize, valueSize)
	if err != nil {
		panic(err)
	}
	// runtime.ReadMemStats(&m)
	// fmt.Println(fmt.Sprintf("After load, memory usage: %d MB", m.Alloc/1024/1024))

	// time.Sleep(time.Second * 5)

	for i := 0; i < batchCount; i++ {
		duration, err := vidbsvc.RandomRead(instance, operationCount, batchSize, keySize, valueSize)
		fmt.Println(fmt.Sprintf("Execute Read %d done. Lantency: %d ns, qps: %.2f", batchSize, duration.Nanoseconds(), float64(batchSize)/float64(duration.Nanoseconds())*1e9))
		if err != nil {
			panic(err)
		}
	}

	// time.Sleep(time.Second * 5)

	for i := 0; i < batchCount; i++ {
		duration, err := vidbsvc.RandomWrite(instance, operationCount, batchSize, keySize, valueSize)
		fmt.Println(fmt.Sprintf("Execute Write %d done. Lantency: %d ns, tps: %.2f", batchSize, duration.Nanoseconds(), float64(batchSize)/float64(duration.Nanoseconds())*1e9))
		if err != nil {
			panic(err)
		}
	}

	_ = instance.Close()
	time.Sleep(time.Second * 5)

}
