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
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"

	"github.com/spf13/cobra"
)

// writeCmd represents the write command
var randomputCmd = &cobra.Command{
	Use:   "randomput",
	Short: "randomput",
	Long:  `letus 执行 vidbconfig 单数据随机写入测试`,
	Run: func(cmd *cobra.Command, args []string) {
		executeRandomPut(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(randomputCmd)
	randomputCmd.Flags().IntVar(&operationCount, "operationCount", 10*vidbsvc.M, "操作次数")
	randomputCmd.Flags().IntVar(&batchCount, "BatchCount", 20, "批次数量")
	randomputCmd.Flags().IntVar(&batchSize, "batchSize", 500, "批次大小")
	randomputCmd.Flags().IntVar(&keySize, "keySize", 32, "键的大小")
	randomputCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "值的大小")
	randomputCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "存储路径")
}

func executeRandomPut(cmd *cobra.Command, args []string) {
	config := vidbconfig.GetDefaultConfig()
	config.DataPath = dataPath
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

	for i := 0; i < batchCount; i++ {
		duration, err := vidbsvc.RandomWrite(instance, operationCount, batchSize, keySize, valueSize)
		fmt.Println(fmt.Sprintf("Execute Write %d done. Lantency: %d ns", batchSize, duration.Nanoseconds()))
		if err != nil {
			panic(err)
		}
	}
	_ = instance.Close()
	time.Sleep(time.Second * 5)

}
