/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	vidbconfig "github.com/bcds/go-hpc-vidb/config"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"

	"github.com/spf13/cobra"
)

// coldHotCmd represents the coldHot command
var coldHotCmd = &cobra.Command{
	Use:   "coldHot",
	Short: "ch",
	Long:  `冷热数据的写入测试`,
	Run: func(cmd *cobra.Command, args []string) {
		executeColdHot(cmd, args)
	},
}

var coldRate float64

func init() {
	rootCmd.AddCommand(coldHotCmd)

	coldHotCmd.Flags().IntVar(&operationCount, "operationCount", 1*vidbsvc.M, "操作次数")
	coldHotCmd.Flags().IntVar(&batchSize, "batchSize", 500, "批次大小")
	coldHotCmd.Flags().IntVar(&keySize, "keySize", 32, "键的大小")
	coldHotCmd.Flags().IntVar(&valueSize, "valueSize", 1024, "值的大小")
	coldHotCmd.Flags().Float64Var(&coldRate, "cr", 0.2, "冷数据占比 0.2")
	coldHotCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "存储路径")
}

func executeColdHot(cmd *cobra.Command, args []string) {

	minRate := min(coldRate, 1-coldRate)
	minCount := float64(operationCount) * minRate

	coldNum := float64(operationCount) * coldRate
	hotNum := float64(operationCount) * (1 - coldRate)

	// 先执行 VIDB Cold 的数据写入
	coldDataPath := filepath.Join(dataPath, "cold")
	coldConfig := vidbconfig.GetDefaultConfig()
	coldConfig.DataPath = coldDataPath
	_ = os.RemoveAll(coldConfig.DataPath)
	insCold, err := vidbsvc.GetVIDBInstance(coldConfig)
	if err != nil {
		panic(err)
	}
	// 先执行数据的写入
	duration, err := vidbsvc.MicroWrite(insCold, int(coldNum), batchSize, keySize, valueSize)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("Cold write done. Num: %d Lan: %d us", int(coldNum), duration.Microseconds()))
	fmt.Println("Now close cold db")
	_ = insCold.Close()
	insCold = nil
	// 开始执行 Cold 文件数据的 Copy
	hotDataPath := filepath.Join(dataPath, "hot")
	// command := exec.Command("cp", "-r", coldConfig.DataPath, hotDataPath)
	// // 执行命令
	// err = command.Run()
	// if err != nil {
	// 	log.Fatalf("执行 cp 命令出错: %v", err)
	// }

	hotConfig := vidbconfig.GetDefaultConfig()
	hotConfig.DataPath = hotDataPath
	_ = os.RemoveAll(hotConfig.DataPath)
	insHot, err := vidbsvc.GetVIDBInstance(hotConfig)
	if err != nil {
		panic(err)
	}
	// 先执行数据的写入
	duration, err = vidbsvc.MicroWrite(insHot, int(hotNum), batchSize, keySize, valueSize)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("Hot write done. Num: %d Lan: %d us", int(hotNum), duration.Microseconds()))
	fmt.Println("Now start cold db")

	insCold, err = vidbsvc.GetVIDBInstance(coldConfig)
	// cold 数据库读取数据
	coldDurtion, err := vidbsvc.MicroRead(insCold, int(minCount), batchSize)
	if err != nil {
		panic(err)
	}

	// Hot 数据库读取数据,读取一样的多的数据
	hotDurtion, err := vidbsvc.MicroRead(insHot, int(minCount), batchSize)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("冷数据占比: %.2f,读取数据总量:%d Cold 读取耗时: %d ms Hot 读取耗时: %d ms", coldRate, int(minCount), coldDurtion.Milliseconds(), hotDurtion.Milliseconds()))
}
