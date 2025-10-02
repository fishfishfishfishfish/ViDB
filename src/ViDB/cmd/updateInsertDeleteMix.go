/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"
)

// updateInsertDeleteMixCmd represents the updateInsertDeleteMix command
var updateInsertDeleteMixCmd = &cobra.Command{
	Use:   "updateInsertDeleteMix",
	Short: "uidm",
	Long:  `不同update-insert-delete比例`,
	Run: func(cmd *cobra.Command, args []string) {
		executeUpdateInsertDeleteMix()
	},
}

var p1, p2, p3 float64
var ibatchSize, ubatchSize, dbatchSize int

func init() {
	rootCmd.AddCommand(updateInsertDeleteMixCmd)
	updateInsertDeleteMixCmd.Flags().IntVar(&n1, "n_1", 10*vidbsvc.M, "写入N1条数据")
	updateInsertDeleteMixCmd.Flags().IntVar(&n2, "n_2", 500, "进行N2次操作")
	updateInsertDeleteMixCmd.Flags().IntVar(&n3, "n_3", 500, "均匀分布随机选取N3个key")
	updateInsertDeleteMixCmd.Flags().IntVar(&n4, "n_4", 500, "范围查询大小")
	updateInsertDeleteMixCmd.Flags().Float64Var(&p1, "p1", 0.5, "P1 操作 Insert [0,1]")
	updateInsertDeleteMixCmd.Flags().Float64Var(&p2, "p2", 0.5, "P2 操作 Update [0,1]")
	updateInsertDeleteMixCmd.Flags().Float64Var(&p3, "p3", 0.5, "P3 操作 Delete [0,1]")
	updateInsertDeleteMixCmd.Flags().IntVar(&metaNum, "meta_num", 1000, "默认的单棵树的版本大小")
	updateInsertDeleteMixCmd.Flags().IntVar(&batchSize, "batch_size", 5000, "默认的批次大小")
	updateInsertDeleteMixCmd.Flags().IntVar(&ibatchSize, "i_batch_size", 5000, "默认的批次大小")
	updateInsertDeleteMixCmd.Flags().IntVar(&ubatchSize, "u_batch_size", 5000, "默认的批次大小")
	updateInsertDeleteMixCmd.Flags().IntVar(&dbatchSize, "d_batch_size", 5000, "默认的批次大小")
	updateInsertDeleteMixCmd.Flags().IntVar(&keySize, "key_size", 32, "默认的Key大小")
	updateInsertDeleteMixCmd.Flags().IntVar(&valueSize, "value_size", 1024, "默认的Value大小")
	updateInsertDeleteMixCmd.Flags().StringVar(&dataPath, "dataPath", "testdata/paper", "存储路径")
}

var UpdateInsertDeleteMixHeader = "put[%d] latency(ms), put throughput, insert[%.2f] latency(ms), insert throughput, update[%.2f] latency(ms), update throughput, delete[%.2f] latency(ms), delete throughput, get[%d] latency, get[%d] throughput, iterator[%d] latency, iterator[%d] throughput\n"
var UpdateInsertDeleteMixHeaderBody = "%.2f, %.2f, %.2f, %.2f, %.2f, %.2f, %.2f, %.2f, %.2f, %.2f, %.2f, %.2f \n"

func executeUpdateInsertDeleteMix() {
	dataPath := filepath.Join(dataPath, "updateInsertDelete")
	perTreeMetaNum := ((n1 + n2) / 5) / batchSize // 5 trees
	config := prepareConfig(uint64(perTreeMetaNum), uint64(batchSize))
	// config := prepareConfig(uint64(200), uint64(batchSize))
	config.DataPath = dataPath
	config.VSize = uint32(valueSize)
	instance, err := vidbsvc.GetVIDBInstance(config)
	if err != nil {
		panic(err)
	}
	// 预写数据量
	writeDuration, err := vidbsvc.MicroWrite(instance, n1, batchSize, keySize, valueSize)
	if err != nil {
		panic(err)
	}
	wTPS := float64(n1) / writeDuration.Seconds()

	// 计算各个操作的量级
	updateCount := int(float64(n2) * p2) // 取整
	insertCount := int(float64(n2) * p1) // 取整
	deleteCount := int(float64(n2) * p3) // 取整
	fmt.Printf("updateCount: %d, insertCount: %d, deleteCount: %d\n", updateCount, insertCount, deleteCount)

	// 测试更新操作
	uBatchCount := updateCount / ubatchSize
	uLat := 0
	for i := 0; i < uBatchCount; i++ {
		// updateKey := genRandomKeyByNormal(ubatchSize)
		updateKey := genKey(n1-ubatchSize, ubatchSize)
		updateDuration, err := vidbsvc.UpdateBatchKey(instance, updateKey, valueSize, false)
		if err != nil {
			panic(err)
		}
		uLat += int(updateDuration.Nanoseconds())
		time.Sleep(time.Millisecond * 5)
	}
	// updateKey := genRandomKeyByNormal(updateCount % ubatchSize)
	updateKey := genKey(n1-(updateCount % ubatchSize), updateCount % ubatchSize)
	updateDuration, err := vidbsvc.UpdateBatchKey(instance, updateKey, valueSize, false)
	if err != nil {
		panic(err)
	}
	uLat += int(updateDuration.Nanoseconds())
	uTPS := float64(n2) / (float64(uLat) * 1e-9)
	fmt.Printf("update tps %.2f\n", uTPS)

	// 测试删除的操作
	dBatchCount := deleteCount / dbatchSize
	dLat := 0
	for i := 0; i < dBatchCount; i++ {
		// deleteKey := genRandomKeyByNormal(dbatchSize)
		deleteKey := genKey(n1-ubatchSize, ubatchSize)
		deleteDuration, err := vidbsvc.UpdateBatchKey(instance, deleteKey, valueSize, true)
		if err != nil {
			panic(err)
		}
		dLat += int(deleteDuration.Nanoseconds())
		time.Sleep(time.Millisecond * 5)
	}
	// deleteKey := genRandomKeyByNormal(deleteCount % dbatchSize)
	deleteKey := genKey(n1-(updateCount % ubatchSize), updateCount % ubatchSize)
	deleteDuration, err := vidbsvc.UpdateBatchKey(instance, deleteKey, valueSize, true)
	if err != nil {
		panic(err)
	}
	dLat += int(deleteDuration.Nanoseconds())
	dTPS := float64(n2) / (float64(dLat) * 1e-9)
	fmt.Printf("delete tps %.2f\n", dTPS)

	// 测试Insert操作
	iBatchCount := insertCount / ibatchSize
	iLat := 0
	for i := 0; i < iBatchCount; i++ {
		insertKey := genKey(n1, ibatchSize)
		insertDuration, err := vidbsvc.UpdateBatchKey(instance, insertKey, valueSize, false)
		if err != nil {
			panic(err)
		}
		iLat += int(insertDuration.Nanoseconds())
		time.Sleep(time.Millisecond * 5)
	}
	insertKey := genKey(n1, insertCount%ibatchSize)
	insertDuration, err := vidbsvc.UpdateBatchKey(instance, insertKey, valueSize, false)
	if err != nil {
		panic(err)
	}
	iLat += int(insertDuration.Nanoseconds())
	iTPS := float64(n2) / (float64(iLat) * 1e-9)
	fmt.Printf("insert tps %.2f\n", iTPS)

	// N3 个操作
	n3Keys := genRandomKeyByNormal(n3)
	duration, _ := vidbsvc.ReadBatchKey(instance, n3Keys)
	sLat := int(duration.Nanoseconds())
	sQPS := float64(n3) / (float64(sLat) * 1e-9)
	fmt.Printf("point query qps %.2f\n", sQPS)

	// N3 迭代器
	idx := randInt(genPCGRand(), 0, n1)
	startKey := fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", idx)
	endKey := fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", idx+n4)
	n3IteratorDuration, _ := vidbsvc.ReadIteratorKey(instance, []byte(startKey), []byte(endKey))
	rLat := int(n3IteratorDuration.Nanoseconds())
	rQPS := float64(n4) / (float64(rLat) * 1e-9)
	fmt.Printf("range query qps %.2f\n", rQPS)

	fmt.Printf(UpdateInsertDeleteMixHeader, n1, p1, p2, p3, n3, n3)
	fmt.Printf(UpdateInsertDeleteMixHeaderBody, writeDuration.Seconds(), wTPS,
		(float64(uLat) * 1e-9), uTPS,
		(float64(dLat) * 1e-9), dTPS,
		(float64(iLat) * 1e-9), iTPS,
		(float64(sLat) * 1e-9), sQPS,
		(float64(rLat) * 1e-9), rQPS,
	)
}
