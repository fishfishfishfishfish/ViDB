/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bcds/go-hpc-vidb/common"
	vidbconfig "github.com/bcds/go-hpc-vidb/config"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"

	"github.com/spf13/cobra"
)

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "w",
	Long:  `letus 执行 vidb 单数据写入测试`,
	Run: func(cmd *cobra.Command, args []string) {
		executeWrite(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(writeCmd)
	writeCmd.Flags().Uint64Var(&cacheCost, "cacheCost", 1<<30, "cache的内存空间大小")
	writeCmd.Flags().IntVar(&VlogSize, "VlogSize", 1, "VlogSize文件的大小(GB)")
	writeCmd.Flags().IntVar(&operationCount, "operationCount", 10*vidbsvc.M, "操作次数")
	writeCmd.Flags().IntVar(&batchSize, "batchSize", 500, "批次大小")
	writeCmd.Flags().IntVar(&keySize, "keySize", 32, "键的大小")
	writeCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "值的大小")
	writeCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "存储路径")
}

func executeWrite(cmd *cobra.Command, args []string) {
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
	_ = instance.Close()
	time.Sleep(time.Second * 5)

}
