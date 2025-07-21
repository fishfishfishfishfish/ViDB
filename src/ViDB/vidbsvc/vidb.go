package vidbsvc

import (
	"fmt"
	"math"
	"strconv"
	"time"

	vidb "github.com/bcds/go-hpc-vidb"
	vidbinterface "github.com/bcds/go-hpc-vidb"
	"golang.org/x/exp/rand"
)

const (
	K = 1_000
	M = 1000 * K
	G = 1000 * M
)

var globalSeq = uint64(0)

func GetVIDBInstance(configInterface vidbinterface.VidbConfigInterface) (*vidb.DB, error) {
	db, err := vidb.Open(configInterface, Logger{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MicroWrite(db *vidb.DB, operationCount, batchSize int, keySize int, valueSize uint32) (time.Duration, error) {
	txCount := operationCount / batchSize // get batch count
	seq, err := db.GetSeqNo()
	if err != nil {
		panic(err)
	}
	if seq == math.MaxUint64{
		seq = 0
	} else{
		seq = seq + 1
	}
	
	value := make([]byte, valueSize)
	_, _ = rand.Read(value)
	num := 0
	now := time.Now()
	for i := seq; i < seq + uint64(txCount); i++ {
		tx, err := db.NewBatchWithEngine()
		if err != nil {
			return 0, err
		}
		for j := 0; j < batchSize; j++ {
			key := fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", num)
			if err := tx.Put([]byte(key), value); err != nil {
				return 0, err
			}
			num += 1
		}

		if err := tx.Hash(uint64(i)); err != nil {
			return 0, err
		}

		if err := tx.Write(uint64(i)); err != nil {
			return 0, err
		}
		if i != 0 && i%10 == 0 {
			if err := db.Commit(uint64(i)); err != nil {
				return 0, err
			}
		}
	}
	return time.Since(now), nil
}

func MicroRangeQuery(db *vidb.DB, operationCount, batchSize int, keySize int, r int) (time.Duration, error) {
	num := rand.Intn(operationCount - r)
	start := fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", num)
	end := fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", num+r)

	fmt.Printf("Start Key: %s, end Key: %s\n", start, end)
	iterator := db.NewIterator([]byte(start), []byte(end))
	defer iterator.Release()
	iteratorStartTime := time.Now()
	for iterator.Next() {}

	duration := time.Since(iteratorStartTime) 
	return duration, nil
}

func MicroRead(db *vidb.DB, operationCount, batchSize int) (time.Duration, error) {
	num := 0
	queryKeys := make([][]byte, operationCount)
	for i := 0; i < operationCount/batchSize; i++ {
		for j := 0; j < batchSize; j++ {
			key := fmt.Sprintf("-account%032d", num)
			queryKeys = append(queryKeys, []byte(key))
			num += 1
		}
	}

	startTime := time.Now()
	for _, key := range queryKeys {
		_, _ = db.Get(key)
	}
	return time.Since(startTime), nil
}

func MicroRRead(db *vidb.DB, operationCount, keySize, r int) (time.Duration, error) {
	num := rand.Intn(operationCount - r)
	keys := make([][]byte, 0)
	for i := 0; i < r; i++ {
		key_str := fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", num)
		keys = append(keys, []byte(key_str))
		num += 1
	}
	fmt.Printf("Start Key: %s, end Key: %s\n", keys[0], keys[r-1])
	startTime := time.Now()
	for _, key := range keys {
		_, _ = db.Get(key)
	}
	return time.Since(startTime), nil
}

func RandomWrite(db *vidb.DB, operationCount, batchSize int, keySize int, valueSize uint32) (time.Duration, error) {
	value := make([]byte, valueSize)
	_, _ = rand.Read(value)

	keys := make([][]byte, 0)
	for i := 0; i < batchSize; i++ {
		num := rand.Intn(operationCount)
		key_str := fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", num)
		keys = append(keys, []byte(key_str))
	}

	tx, _ := db.NewBatch()
	now := time.Now()

	for _, key := range keys {
		if err := tx.Put(key, value); err != nil {
			return 0, err
		}
	}

	seq, err := db.GetSeqNo()
	fmt.Println("current seq:", seq)
	if err != nil {
		return 0, err
	}
	seq += 1

	if err := tx.Hash(seq); err != nil {
		return 0, err
	}

	if err := tx.Write(seq); err != nil {
		return 0, err
	}

	if seq != 0 && seq%10 == 0 {
		if err := db.Commit(seq); err != nil {
			return 0, err
		}
	}

	return time.Since(now), nil
}

func RandomRead(db *vidb.DB, operationCount, batchSize int, keySize int, valueSize uint32) (time.Duration, error) {
	keys := make([][]byte, 0)
	for i := 0; i < batchSize; i++ {
		num := rand.Intn(operationCount)
		key_str := fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", num)
		keys = append(keys, []byte(key_str))
	}
	startTime := time.Now()
	for _, key := range keys {
		_, _ = db.Get(key)
	}
	return time.Since(startTime), nil
}
