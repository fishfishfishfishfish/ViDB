/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	vidbconfig "github.com/bcds/go-hpc-vidb"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"

	"github.com/spf13/cobra"
)

// coldHotCmd represents the coldHot command
var coldHotCmd = &cobra.Command{
	Use:   "coldHot",
	Short: "ch",
	Long:  `cold-hot test under storage tiering`,
	Run: func(cmd *cobra.Command, args []string) {
		executeColdHot(cmd, args)
	},
}

var coldRate float64

func init() {
	rootCmd.AddCommand(coldHotCmd)

	coldHotCmd.Flags().IntVar(&operationCount, "operationCount", 1*vidbsvc.M, "total number of records")
	coldHotCmd.Flags().IntVar(&batchSize, "batchSize", 500, "batch size")
	coldHotCmd.Flags().IntVar(&keySize, "keySize", 32, "key size")
	coldHotCmd.Flags().Uint32Var(&valueSize, "valueSize", 1024, "value szie")
	coldHotCmd.Flags().Float64Var(&coldRate, "cr", 0.2, "cold/hot data ratio")
	coldHotCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "data path")
}

func executeColdHot(cmd *cobra.Command, args []string) {

	minRate := min(coldRate, 1-coldRate)
	minCount := float64(operationCount) * minRate

	coldNum := float64(operationCount) * coldRate
	hotNum := float64(operationCount) * (1 - coldRate)

	// init vidb
	coldDataPath := filepath.Join(dataPath, "cold")
	coldConfig := vidbconfig.GetDefaultConfig()
	coldConfig.DataPath = coldDataPath
	coldConfig.VSize = valueSize
	_ = os.RemoveAll(coldConfig.DataPath)
	insCold, err := vidbsvc.GetVIDBInstance(coldConfig)
	if err != nil {
		panic(err)
	}
	// data loading
	_, err = vidbsvc.MicroWrite(insCold, int(coldNum), batchSize, keySize, valueSize)
	if err != nil {
		panic(err)
	}
	_ = insCold.Close()
	insCold = nil

	// migrate data
	hotDataPath := filepath.Join(dataPath, "hot")
	hotConfig := vidbconfig.GetDefaultConfig()
	hotConfig.DataPath = hotDataPath
	hotConfig.VSize = valueSize
	_ = os.RemoveAll(hotConfig.DataPath)
	insHot, err := vidbsvc.GetVIDBInstance(hotConfig)
	if err != nil {
		panic(err)
	}
	_, err = vidbsvc.MicroWrite(insHot, int(hotNum), batchSize, keySize, valueSize)
	if err != nil {
		panic(err)
	}

	insCold, err = vidbsvc.GetVIDBInstance(coldConfig)
	// query cold data
	coldDurtion, err := vidbsvc.MicroRead(insCold, int(minCount), batchSize)
	if err != nil {
		panic(err)
	}

	// query hot data
	hotDurtion, err := vidbsvc.MicroRead(insHot, int(minCount), batchSize)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("Cold data ratio: %.2f, cold query latency: %d ns, hot query latency: %d ns", coldRate, coldDurtion.Nanoseconds(), hotDurtion.Nanoseconds()))
}
