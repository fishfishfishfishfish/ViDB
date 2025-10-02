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

// LatencyBreakdownSingleQueryCmd represents the LatencyBreakdownSingleQuery command
var LatencyBreakdownSingleQueryCmd = &cobra.Command{
	Use:   "LatencyBreakdownSingleQuery",
	Short: "lbsq",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		executeLatencyBreakdownSingleQuery(cmd)
	},
}

func init() {
	rootCmd.AddCommand(LatencyBreakdownSingleQueryCmd)
	LatencyBreakdownSingleQueryCmd.Flags().Uint64Var(&cacheCost, "cacheCost", 1<<30, "全局Cache的内存空间大小")
	LatencyBreakdownSingleQueryCmd.Flags().Uint64Var(&cacheCounters, "cache_count", 1000, "Cache可写入的数量")
	LatencyBreakdownSingleQueryCmd.Flags().Uint64Var(&treeHeight, "treeHeight", 2, "树的高度")
	LatencyBreakdownSingleQueryCmd.Flags().UintVar(&bloomCap, "bloomCap", 5000, "布隆过滤器的大小")
	LatencyBreakdownSingleQueryCmd.Flags().IntVar(&n1, "n_1", 20000, "单棵树的总量(树版本数=n_1 / batchSize)")
	LatencyBreakdownSingleQueryCmd.Flags().IntVar(&n2, "n_2", 10*vidbsvc.M, "总数据写入量")
	LatencyBreakdownSingleQueryCmd.Flags().IntVar(&batchSize, "batchSize", 500, "批次大小,默认 500")
	LatencyBreakdownSingleQueryCmd.Flags().IntVar(&keySize, "keySize", 32, "写入Key的大小")
	LatencyBreakdownSingleQueryCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "写入Value的大小")
	LatencyBreakdownSingleQueryCmd.Flags().StringVar(&dataPath, "dataPath", "testdata/paper", "存储路径")
}

func executeLatencyBreakdownSingleQuery(cmd *cobra.Command) {
	dataPath = filepath.Join(dataPath, "latency-single-query")
	// 这里需要去计算 数的版本
	currentMetaNum := n1 / batchSize
	config := prepareConfig(uint64(currentMetaNum), uint64(batchSize))
	// config.RotateStrategyTypeName = "adaptive"
	// config.RotateThRootFillPercent = 0.75
	// config.RotateThTreeHeight = treeHeight
	// config.BloomCap = bloomCap
	// config.BloomRate = 0.01
	config.DataPath = dataPath
	config.MaxCost = cacheCost
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
	fmt.Printf("Write %d keys, tps:%.2f\n", n2, tTps)

	// rng := genPCGRand()
	// idx := randInt(rng, 0, n2)
	fmt.Printf("key,bloom(us),cache(us),tree(us)\n")
	format := "%d,%.2f,%.2f,%.2f\n"
	// for i := 0; i < 20; i++ {
	// 	idx := 0
	// 	bloomDuration, cacheDuration, treeDuration, err := vidbsvc.LatencyGet(ins, []byte(fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", idx)))
	// 	if err != nil {
	// 		fmt.Printf("error: %v\n", err)
	// 	}
	// 	fmt.Printf(format,
	// 		idx,
	// 		float64(bloomDuration.Microseconds()),
	// 		float64(cacheDuration.Microseconds()),
	// 		float64(treeDuration.Microseconds()),
	// 	)
	// }
	// for idx := 0; idx < n1; idx += n1 / 20 {
	// 	bloomDuration, cacheDuration, treeDuration, err := vidbsvc.LatencyGet(ins, []byte(fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", idx)))
	// 	if err != nil {
	// 		fmt.Printf("error: %v\n", err)
	// 		continue
	// 	}
	// 	fmt.Printf(format,
	// 		idx,
	// 		float64(bloomDuration.Microseconds()),
	// 		float64(cacheDuration.Microseconds()),
	// 		float64(treeDuration.Microseconds()),
	// 	)
	// }
	for idx := 0; idx < n2; idx += n2 / 20 {
		bloomDuration, cacheDuration, treeDuration, err := vidbsvc.LatencyGet(ins, []byte(fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", idx)))
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		fmt.Printf(format,
			idx,
			float64(bloomDuration.Microseconds()),
			float64(cacheDuration.Microseconds()),
			float64(treeDuration.Microseconds()),
		)
	}
}
