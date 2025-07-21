/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	vidb "github.com/bcds/go-hpc-vidb"
	vidbconfig "github.com/bcds/go-hpc-vidb"
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
	Long:  `rollback benchmark`,
	Run: func(cmd *cobra.Command, args []string) {
		executeRollback(cmd, args)
	},
}

var rollbacks []int

func init() {
	rootCmd.AddCommand(rollbackCmd)
	rollbackCmd.Flags().StringVar(&dataPath, "dataPath", filepath.Join("testdata", "letus"), "data path")
	rollbackCmd.Flags().IntSliceVar(&rollbacks, "rollback", []int{5, 10, 20, 30}, "rollback version count")
}

func executeRollback(cmd *cobra.Command, args []string) {
	// init vidb
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

	// data loading
	value := random.RandBytes(1024)
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

	// rollback
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
		fmt.Println(fmt.Sprintf("RollBack versions: %d Latency: %d (ns)", rollback, duration.Nanoseconds()))
	}
}
