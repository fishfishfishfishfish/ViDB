/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	vidbconfig "github.com/bcds/go-hpc-vidb/config"

	"github.com/bcds/go-hpc-vidb/common"

	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"

	"github.com/spf13/cobra"
)

// writeCmd represents the write command
var randomgetCmd = &cobra.Command{
	Use:   "randomget",
	Short: "randomget",
	Long:  `letus 执行 vidbconfig 单数据随机读取测试`,
	Run: func(cmd *cobra.Command, args []string) {
		executeRandomGet(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(randomgetCmd)
	randomgetCmd.Flags().BoolVar(&recover, "recover", true, "从现有数据恢复")
	randomgetCmd.Flags().Uint64Var(&cacheCost, "cacheCost", 1<<30, "cache的内存空间大小")
	randomgetCmd.Flags().IntVar(&VlogSize, "VlogSize", 1, "VlogSize文件的大小(GB)")
	randomgetCmd.Flags().IntVar(&operationCount, "operationCount", 10*vidbsvc.M, "操作次数")
	randomgetCmd.Flags().IntVar(&batchCount, "BatchCount", 20, "批次数量")
	randomgetCmd.Flags().IntVar(&batchSize, "batchSize", 500, "批次大小")
	randomgetCmd.Flags().IntVar(&keySize, "keySize", 32, "键的大小")
	randomgetCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "值的大小")
	randomgetCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "存储路径")
}

func executeRandomGet(cmd *cobra.Command, args []string) {
	config := vidbconfig.GetDefaultConfig()
	config.DataPath = dataPath
	config.MaxCost = cacheCost
	config.VlogSize = uint64(VlogSize) * common.GiB
	if !recover {
		if err := os.RemoveAll(config.DataPath); err != nil {
			panic(err)
		}
	}

	instance, err := vidbsvc.GetVIDBInstance(config)
	if err != nil {
		panic(err)
	}
	for i := 0; i < batchCount; i++ {
		duration, err := vidbsvc.RandomRead(instance, operationCount, batchSize, keySize, valueSize)
		fmt.Println(fmt.Sprintf("Execute Read %d done. Lantency: %d ns", batchSize, duration.Nanoseconds()))
		if err != nil {
			panic(err)
		}
	}
	_ = instance.Close()
	time.Sleep(time.Second * 5)

}
