/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	vidbconfig "github.com/bcds/go-hpc-vidb/config"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "read",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeRead()
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
	// 定义命令行参数
	readCmd.Flags().IntVar(&operationCount, "operationCount", 10*vidbsvc.M, "操作次数")
	readCmd.Flags().IntVar(&batchSize, "batchSize", 500, "批次大小")
	readCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "值的大小")
	readCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "存储路径")
	readCmd.Flags().IntSliceVar(&rs, "rs", []int{5, 50, 100, 200, 300, 400, 500, 1000, 2000}, "范围查询的范围列表，多个值用逗号分隔，例如 5,50,100")
}

func executeRead() {
	config := vidbconfig.GetDefaultConfig()
	config.DataPath = dataPath
	for _, r := range rs {
		instance, err := vidbsvc.GetVIDBInstance(config)
		if err != nil {
			panic(err)
		}
		duration, err := vidbsvc.MicroRead(instance, operationCount, batchSize)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Execute r %d done. Lantency: %d us", r, duration.Microseconds()))
		_ = instance.Close()
		time.Sleep(time.Second * 5)
	}
}
