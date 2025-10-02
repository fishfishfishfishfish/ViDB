/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	vidbconfig "github.com/bcds/go-hpc-vidb/config"
	vidb "github.com/bcds/go-hpc-vidb/database"
	"github.com/duke-git/lancet/v2/random"
	"github.com/spf13/cobra"
	"gitlab.bcds.org.cn/sunyang/letus-vidb/vidbsvc"
	"os"
	"path/filepath"
	"time"
)

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "rb",
	Long:  `数据回滚`,
	Run: func(cmd *cobra.Command, args []string) {
		executeRollback(cmd, args)
	},
}

var rollbacks []int

func init() {
	rootCmd.AddCommand(rollbackCmd)
	rollbackCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "存储路径")
	rollbackCmd.Flags().IntSliceVar(&rollbacks, "rollback", []int{5, 10, 20, 30}, "rollback不同数量版本")
}

func executeRollback(cmd *cobra.Command, args []string) {
	// 首先构建每一个会滚对象的 db
	configs := make([]*vidbconfig.VidbConfig, len(rollbacks))
	for idx, rollback := range rollbacks {
		rollbackPath := filepath.Join(dataPath, "rollback", fmt.Sprintf("%d", rollback))
		_ = os.RemoveAll(rollbackPath)
		config := vidbconfig.GetDefaultConfig()
		config.DataPath = rollbackPath
		configs[idx] = config
	}

	dbs := make([]*vidb.DB, len(configs))
	for idx, config := range configs {
		db, err := vidbsvc.GetVIDBInstance(config)
		if err != nil {
			panic(err)
		}
		dbs[idx] = db
	}
	value := random.RandBytes(1024)
	// 每个db 写入 1W 条数据 BatchSize = 500 Value 1024
	for _, db := range dbs {
		for i := 0; i <= 1000; i++ {
			tx, _ := db.NewBatch()
			for j := 0; j < 500; j++ {
				key := fmt.Sprintf("-account%20d%12d", i, j)
				_ = tx.Put([]byte(key), value)
			}
			_ = tx.Hash(uint64(i))
			_ = tx.Write(uint64(i))

			if i != 0 && i < 950 && i%10 == 0 {
				_ = db.Commit(uint64(i))
			}
		}
	}

	// 写入完成后，开始进行 rollback
	for idx, rollback := range rollbacks {
		db := dbs[idx]
		curSeq, _ := db.GetSeqNo()
		targetSeq := curSeq - uint64(rollback)
		now := time.Now()
		_ = db.Revert(targetSeq)
		duration := time.Since(now)

		afterRevertSeq, _ := db.GetSeqNo()
		if afterRevertSeq+1 != targetSeq {
			panic("afterRevertSeq != targetMeta")
		}
		fmt.Println(fmt.Sprintf("RollBack的版本数量: %d Lan: %d (us)", rollback, duration.Microseconds()))
	}
}
