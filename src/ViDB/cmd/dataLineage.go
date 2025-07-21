/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	vidbconfig "github.com/bcds/go-hpc-vidb"
	"github.com/duke-git/lancet/v2/random"
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

var queryVerStart []int
var queryVerCount []int
var updateCntBfLoad int
var updateCntAfLoad int

func init() {
	rootCmd.AddCommand(dataLineageCmd)
	dataLineageCmd.Flags().Uint64Var(&cacheCost, "cacheCost", 1<<30, "cache的内存空间大小")
	dataLineageCmd.Flags().IntVar(&VlogSize, "VlogSize", 1, "VlogSize文件的大小(GB)")
	dataLineageCmd.Flags().IntVar(&operationCount, "operationCount", 10*vidbsvc.M, "操作次数")
	dataLineageCmd.Flags().IntVar(&batchSize, "batchSize", 500, "批次大小")
	dataLineageCmd.Flags().IntVar(&keySize, "keySize", 32, "键的大小")
	dataLineageCmd.Flags().Uint32Var(&valueSize, "valueSize", 1024, "值的大小")
	dataLineageCmd.Flags().IntVar(&batchCount, "BatchCount", 20, "测试数量")
	dataLineageCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "存储路径")
	dataLineageCmd.Flags().IntSliceVar(&queryVerStart, "queryVerStart", []int{2, 4, 10, 20, 40}, "查询的起始版本")
	dataLineageCmd.Flags().IntSliceVar(&queryVerCount, "queryVerCount", []int{2, 4, 10, 20, 40}, "查询的版本数")
	dataLineageCmd.Flags().IntVar(&updateCntBfLoad, "bfLoad", 0, "全量加载前运行的更新")
	dataLineageCmd.Flags().IntVar(&updateCntAfLoad, "afLoad", 40, "全量加载后运行的更新")
}

func executeDataLineAge(cmd *cobra.Command, args []string) {
	config := vidbconfig.GetDefaultConfig()
	path := filepath.Join(dataPath, "dataLineage")
	_ = os.RemoveAll(path)
	config.DataPath = path
	config.MaxCost = cacheCost
	config.VSize = valueSize
	
	// generate test keys
	keys := make([][]byte, 0)
	for b := 0; b < batchCount; b++ {
		num := rand.Intn(operationCount)
		key := []byte(fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", num))
		keys = append(keys, key)
	}
	
	ins, err := vidbsvc.GetVIDBInstance(config)
	if err != nil {
		panic(err)
	}

	// updates
	var v uint64
	for v = 0; v < uint64(updateCntBfLoad); v++ {
		now := time.Now()
		tx, err := ins.NewBatchWithEngine()
		if err != nil {
			panic(err)
		}
		for _, key := range keys {
			value := random.RandBytes(int(valueSize))
			if err := tx.Put(key, value); err != nil {
				panic(err)
			}
		}
		if err := tx.Hash(uint64(v)); err != nil {
			panic(err)
		}

		if err := tx.Write(uint64(v)); err != nil {
			panic(err)
		}
		duration := time.Since(now)
		Throughput := float64(batchCount) / duration.Seconds()
		fmt.Println(fmt.Sprintf("version=%d, update_count=%d, Latency: %d ns TPS: %.2f(n/s)", v, batchCount, duration.Nanoseconds(), Throughput))
	}
	_ = ins.Commit(v)

	_, err = vidbsvc.MicroWrite(ins, operationCount, batchSize, keySize, valueSize)
	if err != nil {
		panic(err)
	}

	// updates
	seq, err := ins.GetSeqNo()
	if err != nil {
		panic(err)
	}
	for v = seq+1; v < seq + uint64(updateCntAfLoad); v++ {
		now := time.Now()
		tx, err := ins.NewBatchWithEngine()
		if err != nil {
			panic(err)
		}
		for _, key := range keys {
			value := random.RandBytes(int(valueSize))
			if err := tx.Put(key, value); err != nil {
				panic(err)
			}
		}
		if err := tx.Hash(uint64(v)); err != nil {
			panic(err)
		}

		if err := tx.Write(uint64(v)); err != nil {
			panic(err)
		}
		duration := time.Since(now)
		Throughput := float64(batchCount) / duration.Seconds()
		fmt.Println(fmt.Sprintf("version=%d, update_count=%d, Latency: %d ns TPS: %.2f(n/s)", v, batchCount, duration.Nanoseconds(), Throughput))
	}
	_ = ins.Commit(v)

	// warm up
	_, err = vidbsvc.RandomRead(ins, operationCount, 5000, keySize, valueSize)
	if err != nil {
		panic(err)
	}

	no, err := ins.GetSeqNo() // 获取当前的版本
	if err != nil {
		panic(err)
	}
	
	nowTotal := time.Now()
	totalCount := uint64(0)
	for _, nv := range queryVerCount {
		for _, ver_ago := range queryVerStart {
			ver_start := no - uint64(ver_ago) // 800 2 -> 789 800
			ver_end := ver_start + uint64(nv)
			if nv == -1 || ver_end > no + 1 { ver_end = no + 1 }
			for _, key := range keys {
				// num := rand.Intn(operationCount)
				// key := []byte(fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", num))
				// totalCount += uint64(meta)
				now := time.Now()
				for i := ver_start; i < ver_end; i++ {
					_, _ = ins.Proof(key, i)
				}
				duration := time.Since(now)
				Throughput := float64(nv) / duration.Seconds()
				fmt.Println(fmt.Sprintf("query %d versions starting from %d versions ago [%d , %d), Latency: %d ns (%f s) TPS: %.2f(n/s)", 
					nv, ver_ago, ver_start, ver_end, duration.Nanoseconds(), duration.Seconds(), Throughput))
			}
		}
	}
	totalDuration := time.Since(nowTotal)
	totalTPS := float64(totalCount) / totalDuration.Seconds()
	fmt.Println(fmt.Sprintf("All Latency: %d ms TPS: %.2f(n/s)", totalDuration.Milliseconds(), totalTPS))
}
