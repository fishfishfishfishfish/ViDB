/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	// "sort"
	// vidbconfig "github.com/bcds/go-hpc-vidb/config"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/pingcap/go-ycsb/pkg/generator"
	"github.com/pingcap/go-ycsb/pkg/ycsb"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"
)

// getZipfCmd represents the getZipf command
var getZipfCmd = &cobra.Command{
	Use:   "get_zipf",
	Short: "gz",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		executeGetZIPF(cmd)
	},
}
var zipfType string

func init() {
	rootCmd.AddCommand(getZipfCmd)
	getZipfCmd.Flags().IntVar(&n1, "n_1", 100000, "写入N1条数据")
	getZipfCmd.Flags().IntVar(&n3, "n_3", 500, "随机生成的Key的个数")
	getZipfCmd.Flags().Float64Var(&p1, "p1", 1.0, "zipf分布参数")
	getZipfCmd.Flags().StringVar(&zipfType, "zType", "zipfian", "zipf 的类型[zipfian/scrambled],默认 zipfian")
	getZipfCmd.Flags().IntVar(&metaNum, "meta_num", 1000, "默认的单棵树的版本大小")
	getZipfCmd.Flags().IntVar(&batchSize, "batch_size", 5000, "默认的批次大小")
	getZipfCmd.Flags().IntVar(&keySize, "key_size", 32, "默认的Key大小")
	getZipfCmd.Flags().IntVar(&valueSize, "value_size", 1024, "默认的Value大小")
	getZipfCmd.Flags().IntVar(&batchCount, "batch_count", 20, "测试的数量")
	getZipfCmd.Flags().StringVar(&dataPath, "dataPath", "testdata/paper", "存储路径")
}

var zipfHeaderStringArr = []string{"总基数，总基数写入延迟"}
var zipfHeader = "operations,latency(ms),throughput\n"

// var zipfBodyT = "put[%d],%.2f,%.2f\n"
var zipfBodyW = "put[%d],%.2f,%.2f\n"
var zipfBodyR = "get[%d],%.2f,%.2f\n"

func executeGetZIPF(cmd *cobra.Command) {
	dataPath := filepath.Join(dataPath, "zipf")
	// config := vidbconfig.GetDefaultConfig()
	loadBatchSize := 5000
	// perTreeMetaNum := (operationCount / 5) / loadBatchSize // 5 trees
	// perTreeMetaNum := (operationCount) / loadBatchSize // 1 trees
	perTreeMetaNum := 60000000 / loadBatchSize // 60M per tree
	config := prepareConfig(uint64(perTreeMetaNum), uint64(loadBatchSize))
	config.BloomCap = 5000
	config.BloomRate = 0.01
	config.DataPath = dataPath
	config.VSize = uint32(valueSize)
	config.MaxCost = 0 // 1M, approx 1k 1024-byte records
	config.VlogSize = (1 << 30) * 4 // 1M, approx 1k 1024-byte records
	instance, err := vidbsvc.GetVIDBInstance(config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Total %d\n", n1)
	fmt.Printf("loadBatchSize %d\n", loadBatchSize)
	fmt.Printf("treeCap %d\n", 60000000)
	fmt.Printf("RotateThMetaNum %d\n", config.RotateThMetaNum)
	fmt.Printf("BloomCap %d\n", config.BloomCap)
	fmt.Printf("BloomRate %f\n", config.BloomRate)
	fmt.Printf("VlogSize %d\n", config.VlogSize)
	fmt.Printf("cacheCost %d\n", config.MaxCost)
	fmt.Printf("batchSize %d\n", n3)

	float64s := make([]float64, 0)
	// 预写数据量
	fmt.Println("Start Write Base data")
	writeDuration, err := vidbsvc.MicroWrite(instance, n1, loadBatchSize, keySize, valueSize)
	if err != nil {
		panic(err)
	}
	wTPS := float64(n1) / float64(writeDuration.Nanoseconds()) * 1e9
	float64s = append(float64s, wTPS, float64(writeDuration.Milliseconds()))
	printCmdSummary(cmd, zipfHeaderStringArr, float64s)

	var key_generator ycsb.Generator
	switch zipfType {
	case "zipfian":
		key_generator = generator.NewZipfianWithItems(int64(n1), p1)
	case "scrambled":
		key_generator = generator.NewScrambledZipfian(0, int64(n1), p1)
	default:
		key_generator = generator.NewZipfianWithItems(int64(n1), p1)
	}

	rTps_sum := 0.0
	rLat_sum := 0.0
	fmt.Printf(zipfHeader)
	for i := 0; i < batchCount; i++ {
		// 随机生成的 key
		zipfKeys := genRandomKeywGenerator(key_generator, n3, n1, p1)
		readBatchDuration, err := vidbsvc.ReadBatchKey(instance, zipfKeys)
		if err != nil {
			panic(err)
		}
		rTPS := float64(n3) / float64(readBatchDuration.Nanoseconds()) * 1e9
		rTps_sum += rTPS
		rLat_sum += float64(readBatchDuration.Seconds())
		fmt.Printf(zipfBodyR, n3, float64(readBatchDuration.Seconds()), rTPS)
	}

	wTps_sum := 0.0
	wLat_sum := 0.0
	for i := 0; i < batchCount; i++ {
		zipfKeys := genRandomKeywGenerator(key_generator, n3, n1, p1)
		// zipfKeys := genRandomKeyByNormal(batchSize)
		writeBatchDuration, err := vidbsvc.UpdateBatchKey(instance, zipfKeys, valueSize, false)
		if err != nil {
			panic(err)
		}
		wTps := float64(n3) / float64(writeBatchDuration.Nanoseconds()) * 1e9
		wTps_sum += wTps
		wLat_sum += float64(writeBatchDuration.Seconds())
		fmt.Printf(zipfBodyW, n3, float64(writeBatchDuration.Seconds()), wTps)
	}
	fmt.Println("Average")
	fmt.Printf(zipfBodyR, n3, rLat_sum/float64(batchCount), rTps_sum/float64(batchCount))
	fmt.Printf(zipfBodyW, n3, wLat_sum/float64(batchCount), wTps_sum/float64(batchCount))

}

func genRandomKeyByZipf(zType string, num int, total_rec int, p float64) [][]byte {
	data := make([][]byte, 0)
	var gen ycsb.Generator
	switch zType {
	case "zipfian":
		gen = generator.NewZipfianWithItems(int64(total_rec), p)
	case "scrambled":
		gen = generator.NewScrambledZipfian(0, int64(total_rec), p)
	default:
		gen = generator.NewZipfianWithItems(int64(total_rec), p)
	}

	// nums := make([]int, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		idx := gen.Next(r)
		data = append(data, getKey(total_rec-int(idx)))
		// nums = append(nums, int(idx))
	}
	// sort.Ints(nums)
	// fmt.Printf("排序后：[")
	// for _, n := range nums {
	// 	fmt.Printf("%d,", n)
	// }
	// fmt.Println("]")
	return data
}

func genRandomKeywGenerator(gen ycsb.Generator, num int, total_rec int, p float64) [][]byte {
	data := make([][]byte, 0)

	// nums := make([]int, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		idx := gen.Next(r)
		data = append(data, getKey(total_rec-int(idx)))
		// nums = append(nums, int(idx))
	}
	// sort.Ints(nums)
	// fmt.Printf("排序后：[")
	// for _, n := range nums {
	// 	fmt.Printf("%d,", n)
	// }
	// fmt.Println("]")
	return data
}

func genRandomKeywGeneratorMoreOld(gen ycsb.Generator, num int) [][]byte {
	data := make([][]byte, 0)

	// nums := make([]int, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		idx := gen.Next(r)
		data = append(data, getKey(int(idx)))
		// nums = append(nums, int(idx))
	}
	// sort.Ints(nums)
	// fmt.Printf("排序后：[")
	// for _, n := range nums {
	// 	fmt.Printf("%d,", n)
	// }
	// fmt.Println("]")
	return data
}

// printCmdSummary 打印 flags + 自定义指标
func printCmdSummary(cmd *cobra.Command, headers []string, values []float64) {
	// 先收集 flags
	var flagHeaders []string
	var flagValues []string
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		flagHeaders = append(flagHeaders, flag.Name)
		flagValues = append(flagValues, flag.Value.String())
	})
	// 打印表头
	allHeaders := append(flagHeaders, headers...)
	for _, h := range allHeaders {
		fmt.Printf("%s,", h)
	}
	fmt.Println()
	// 打印表内容
	for _, v := range flagValues {
		fmt.Printf("%s,", v)
	}
	for _, v := range values {
		fmt.Printf("%.2f,", v)
	}
	fmt.Println()
}
