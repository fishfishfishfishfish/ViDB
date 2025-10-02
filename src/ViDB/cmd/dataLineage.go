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
	"github.com/spf13/cobra"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"
)

// dataLineageCmd represents the dataLineage command
var dataLineageCmd = &cobra.Command{
	Use:   "dataLineage",
	Short: "dla",
	Long:  `历史版本追溯能力`,
	Run: func(cmd *cobra.Command, args []string) {
		executeDataLineAge(cmd, args)
	},
}

var metas []int

func init() {
	rootCmd.AddCommand(dataLineageCmd)
	dataLineageCmd.Flags().IntVar(&operationCount, "operationCount", 10*vidbsvc.M, "操作次数")
	dataLineageCmd.Flags().IntVar(&batchSize, "batchSize", 500, "批次大小")
	dataLineageCmd.Flags().IntVar(&keySize, "keySize", 32, "键的大小")
	dataLineageCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "值的大小")
	dataLineageCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "存储路径")
	dataLineageCmd.Flags().IntSliceVar(&metas, "metas", []int{2, 4, 10, 20, 40}, "读取的版本数")
}

func executeDataLineAge(cmd *cobra.Command, args []string) {
	config := vidbconfig.GetDefaultConfig()
	path := filepath.Join(dataPath, "dataLineage")
	_ = os.RemoveAll(path)
	config.DataPath = path
	ins, err := vidbsvc.GetVIDBInstance(config)
	if err != nil {
		panic(err)
	}
	_, err = vidbsvc.MicroWrite(ins, operationCount, batchSize, keySize, valueSize)
	if err != nil {
		panic(err)
	}
	// 只需要选择其中一个 Key
	no, err := ins.GetSeqNo() // 获取当前的版本
	if err != nil {
		panic(err)
	}
	key := fmt.Sprintf("-account%020d%012d", no-40, 1)

	nowTotal := time.Now()
	totalCount := uint64(0)
	for _, meta := range metas {
		now := time.Now()
		sumMeta := no - uint64(meta) // 800 2 -> 789 800
		totalCount += uint64(meta)
		for i := sumMeta + 1; i <= no; i++ {
			_, _ = ins.Proof([]byte(key), i)
		}
		duration := time.Since(now)
		Throughput := float64(sumMeta) / duration.Seconds()
		fmt.Println(fmt.Sprintf("n=%d Lantency: %d ns TPS: %.2f(n/s)", meta, duration.Nanoseconds(), Throughput))
	}
	totalDuration := time.Since(nowTotal)
	totalTPS := float64(totalCount) / totalDuration.Seconds()
	fmt.Println(fmt.Sprintf("All Lantency: %d ms TPS: %.2f(n/s)", totalDuration.Milliseconds(), totalTPS))
}
