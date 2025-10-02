/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"

	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"

	"github.com/spf13/cobra"
)

// LatencyBreakdownIteratorQueryCmd represents the LatencyBreakdownIteratorQuery command
var LatencyBreakdownIteratorQueryCmd = &cobra.Command{
	Use:   "LatencyBreakdownIteratorQuery",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeIteratorQuery(cmd)
	},
}

func init() {
	rootCmd.AddCommand(LatencyBreakdownIteratorQueryCmd)
	LatencyBreakdownIteratorQueryCmd.Flags().IntVar(&n1, "n_1", 20000, "单棵树的总量(树版本数=n_1 / batchSize)")
	LatencyBreakdownIteratorQueryCmd.Flags().IntVar(&n2, "n_2", 10*vidbsvc.M, "总数据写入量")
	LatencyBreakdownIteratorQueryCmd.Flags().IntVar(&n3, "n_3", 100, "数据查询的范围")
	LatencyBreakdownIteratorQueryCmd.Flags().IntVar(&batchSize, "batchSize", 5000, "批次大小,默认 500")
	LatencyBreakdownIteratorQueryCmd.Flags().IntVar(&keySize, "keySize", 32, "写入Key的大小")
	LatencyBreakdownIteratorQueryCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "写入Value的大小")
	LatencyBreakdownIteratorQueryCmd.Flags().StringVar(&dataPath, "dataPath", "testdata/paper", "存储路径")
}

func executeIteratorQuery(command *cobra.Command) {
	dataPath = filepath.Join(dataPath, "latency-iterator-query")
	// 这里需要去计算 数的版本
	currentMetaNum := n1 / batchSize
	config := prepareConfig(uint64(currentMetaNum), uint64(batchSize))
	config.DataPath = dataPath
	config.VSize = uint32(valueSize)
	ins, err := vidbsvc.GetVIDBInstance(config)
	if err != nil {
		panic(err)
	}

	totalW, err := vidbsvc.MicroWrite(ins, n2, batchSize, keySize, valueSize)
	if err != nil {
		panic(err)
	}
	tTps := float64(n2) / totalW.Seconds()

	rng := genPCGRand()
	idx := randInt(rng, 0, n2)
	startKey := []byte(fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", idx))
	endKey := []byte(fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", idx+n3))
	nextDuration, treeDuration, _ := vidbsvc.LatencyIterator(ins, startKey, endKey)

	printCmdSummary(command, []string{"write tps", "next(ns)", "tree_iterator(ns)"}, []float64{tTps, float64(nextDuration), float64(treeDuration)})
}
