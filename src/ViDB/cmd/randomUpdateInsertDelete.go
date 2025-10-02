/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math/rand/v2"
	"path/filepath"
	"strconv"
	"time"

	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"

	"github.com/spf13/cobra"
)

// randomUpdateInsertDeleteCmd represents the randomUpdateInsertDelete command
var randomUpdateInsertDeleteCmd = &cobra.Command{
	Use:   "randomUpdateInsertDelete",
	Short: "ruid", // random_update_insert_delete
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		executeRandomUpdateInsertDelete()
	},
}

func init() {
	rootCmd.AddCommand(randomUpdateInsertDeleteCmd)
	randomUpdateInsertDeleteCmd.Flags().IntVar(&n1, "n_1", 10*vidbsvc.M, "写入N1条数据")
	randomUpdateInsertDeleteCmd.Flags().IntVar(&n2, "n_2", 500, "均匀分布随机选取N2个key")
	randomUpdateInsertDeleteCmd.Flags().IntVar(&metaNum, "meta_num", 1000, "默认的单棵树的版本大小")
	randomUpdateInsertDeleteCmd.Flags().IntVar(&batchSize, "batch_size", 5000, "默认的批次大小")
	randomUpdateInsertDeleteCmd.Flags().IntVar(&keySize, "key_size", 32, "默认的Key大小")
	randomUpdateInsertDeleteCmd.Flags().IntVar(&valueSize, "value_size", 1024, "默认的Value大小")
	randomUpdateInsertDeleteCmd.Flags().StringVar(&dataPath, "dataPath", "testdata/paper", "存储路径")
}

var randomUpdateInsertDeleteHeader = "put[%d] latency(ms)\t put throughput\t insert[%d] latency(ms)\t insert throughput\t update[%d] latency(ms)\t update throughput\t delete[%d] latency(ms)\t delete throughput\t\n"
var randomUpdateInsertDeleteHeaderBody = "%d\t %.2f\t %d\t %.2f\t %d\t %.2f\t %d\t %.2f\t \n"

func executeRandomUpdateInsertDelete() {
	dataPath := filepath.Join(dataPath, "updateInsertDelete")
	LoadBatchSize := 5000
	config := prepareConfig(uint64(metaNum), uint64(LoadBatchSize))
	config.DataPath = dataPath
	config.VSize = uint32(valueSize)
	instance, err := vidbsvc.GetVIDBInstance(config)
	if err != nil {
		panic(err)
	}
	// 预写数据量
	writeDuration, err := vidbsvc.MicroWrite(instance, n1, LoadBatchSize, keySize, valueSize)
	if err != nil {
		panic(err)
	}
	wTPS := float64(n1) / writeDuration.Seconds()

	// 测试更新操作
	// updateKey := genRandomKeyByNormal(n2)
	updateKey := genKey(n1-n2, n2)
	updateDuration, err := vidbsvc.UpdateBatchKey(instance, updateKey, valueSize, false)
	if err != nil {
		panic(err)
	}
	uTPS := float64(n2) / updateDuration.Seconds()

	// 测试Insert操作
	insertKey := genKey(n1, n2)
	insertDuration, err := vidbsvc.UpdateBatchKey(instance, insertKey, valueSize, false)
	if err != nil {
		panic(err)
	}
	iTPS := float64(n2) / insertDuration.Seconds()

	// 测试删除的操作
	deleteKey := genRandomKeyByNormal(n2)
	deleteDuration, err := vidbsvc.UpdateBatchKey(instance, deleteKey, valueSize, true)
	if err != nil {
		panic(err)
	}
	dTPS := float64(n2) / deleteDuration.Seconds()

	fmt.Printf(randomUpdateInsertDeleteHeader, n1, n2, n2, n2)
	fmt.Printf(randomUpdateInsertDeleteHeaderBody, writeDuration.Milliseconds(), wTPS,
		insertDuration.Milliseconds(), iTPS,
		updateDuration.Milliseconds(), uTPS,
		deleteDuration.Milliseconds(), dTPS,
	)
}

func genRandomKeyByNormal(num int) [][]byte {
	data := make([][]byte, 0)
	pcgRand := genPCGRand()
	for i := 0; i < num; i++ {
		idx := randInt(pcgRand, 0, n1)
		data = append(data, getKey(idx))
	}
	return data
}

func genKey(start int, num int) [][]byte {
	data := make([][]byte, 0)
	for i := start; i < start+num+1; i++ {
		data = append(data, getKey(i))
	}
	return data
}

func getKey(idx int) []byte {
	return []byte(fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", idx))
}

func genPCGRand() *rand.Rand {
	seed1 := uint64(time.Now().UnixNano())
	seed2 := uint64(time.Now().UnixNano() ^ 0x5EED) // 稍微混淆第二个种子
	rng := rand.New(rand.NewPCG(seed1, seed2))
	return rng
}

// 生成 [min, max] 范围内的均匀分布整数
func randInt(rng *rand.Rand, min, max int) int {
	if min >= max {
		return min
	}
	return min + rng.IntN(max-min+1)
}
