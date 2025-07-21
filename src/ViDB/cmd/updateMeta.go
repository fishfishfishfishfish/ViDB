/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	vidbconfig "github.com/bcds/go-hpc-vidb"
	"github.com/duke-git/lancet/v2/random"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"

	"github.com/spf13/cobra"
)

// updateMetaCmd represents the updateMeta command
var updateMetaCmd = &cobra.Command{
	Use:   "updateMeta",
	Short: "um",
	Long:  `版本更新`,
	Run: func(cmd *cobra.Command, args []string) {
		executeUpdateMeta(cmd, args)
	},
	
}
var metas []int

func init() {
	rootCmd.AddCommand(updateMetaCmd)
	updateMetaCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "data path")
	updateMetaCmd.Flags().IntSliceVar(&metas, "metas", []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, "versions")
	updateMetaCmd.Flags().IntVar(&batchSize, "batchSize", 500, "records updated in one version")
}

func executeUpdateMeta(cmd *cobra.Command, args []string) {
	path := filepath.Join(dataPath, "updateMeta")
	_ = os.RemoveAll(path)

	config := vidbconfig.GetDefaultConfig()
	config.DataPath = path

	ins, err := vidbsvc.GetVIDBInstance(config)
	if err != nil {
		panic(err)
	}
	value := random.RandBytes(1024)

	for _, meta := range metas {
		batch, _ := ins.NewBatchWithEngine()
		now := time.Now()
		for i := 0; i < batchSize; i++ {
			_ = batch.Put([]byte(fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", i)), value)
		}
		if err := batch.Hash(uint64(meta - 1)); err != nil {
			fmt.Printf("%d meta, Hash error: %v\n", meta, err)
			continue
		}
		if err := batch.Write(uint64(meta - 1)); err != nil {
			fmt.Printf("%d meta, Write error: %v\n", meta, err)
			continue
		}
		since := time.Since(now).Nanoseconds() // nanoseconds

		cmd := exec.Command("bash", "-c", fmt.Sprintf("du -sk %s", filepath.Join(path, "index")))
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("command fail:: %v\n", err)
			continue
		}
		cmd = exec.Command("bash", "-c", fmt.Sprintf("du -sk %s", filepath.Join(path, "data")))
		output1, err := cmd.Output()
		if err != nil {
			fmt.Printf("command fail:: %v\n", err)
			continue
		}
		cmd = exec.Command("bash", "-c", fmt.Sprintf("du -sk %s", path))
		output2, err := cmd.Output()
		if err != nil {
			fmt.Printf("command fail: %v\n", err)
			continue
		}
		fmt.Printf("Version:[%d],Latency:[%d](ns),Index size:[%s],Data size:[%s],Total size:[%s]\n", meta, since, bytes.TrimSpace(output), bytes.TrimSpace(output1), bytes.TrimSpace(output2))
	}
}
