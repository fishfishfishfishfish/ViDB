/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	// "bytes"
	// "os/exec"
	"fmt"
	"regexp"
	vidbconfig "github.com/bcds/go-hpc-vidb/config"
	"github.com/duke-git/lancet/v2/random"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"
	"os"
	"path/filepath"
	"time"

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

func getFileSize(filename string)(int64,error){
	fi, err := os.Stat(filename)
	if err != nil{
		return 0,err
	}
	return fi.Size(), nil
}

func init() {
	rootCmd.AddCommand(updateMetaCmd)
	updateMetaCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "存储路径")
	updateMetaCmd.Flags().IntSliceVar(&metas, "metas", []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, "版本范围")
	updateMetaCmd.Flags().IntVar(&batchSize, "batchSize", 500, "版本范围")
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

	// 写入需要从 0 开始
	randString := random.RandString(31)
	for _, meta := range metas {
		batch, _ := ins.NewBatchWithEngine()
		now := time.Now()
		for i := 0; i < batchSize; i++ {
			_ = batch.Put([]byte("-account"+randString+fmt.Sprintf("%d", i)), value)
		}
		_ = batch.Hash(uint64(meta - 1))
		_ = batch.Write(uint64(meta - 1))
		since := time.Since(now).Milliseconds() // ms


		files, err := os.ReadDir(path) 
		if err != nil {
			fmt.Printf("error reading directory:%s", err)
			continue
		}
		partition_cnt := 0
		regex := regexp.MustCompile(`partition-\d+`)
		for _, file := range files {
			if file.IsDir() && regex.MatchString(file.Name()) {
				partition_cnt += 1
			}
		}

		var index_size = int64(0)
		for i := 0; i < partition_cnt; i++ {
			size, err := getFileSize(filepath.Join(path, fmt.Sprintf("partition-%d/index", i)))
			if err != nil {
				fmt.Printf("error get file size: %v\n", err)
				continue
			}
			index_size += size
		}

		// cmd = exec.Command("bash", "-c", fmt.Sprintf("du -sh %s", path))
		// output2, err := cmd.Output()
		// if err != nil {
		// 	fmt.Printf("执行命令出错: %v\n", err)
		// 	continue
		// }
		fmt.Printf("Version: %d, time: %d, Index size: %d\n", meta, since, index_size)
		// time.Sleep(5 * time.Second)
	}
}
