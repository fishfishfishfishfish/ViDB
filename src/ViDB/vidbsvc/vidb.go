package vidbsvc

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	commonmetrics "github.com/bcds/go-hpc-common/metrics"
	commonprom "github.com/bcds/go-hpc-common/metrics/prometheus"
	vidb "github.com/bcds/go-hpc-vidb/database"
	vidbinterface "github.com/bcds/go-hpc-vidb/interfaces"
	// mathrand "math/rand/v2"
	"golang.org/x/exp/rand"
)

const (
	K = 1_000
	M = 1000 * K
	G = 1000 * M
)

var globalSeq = uint64(0)

func initMetricsProvider() commonmetrics.Provider {
	initMetrics()
	provider := commonprom.Provider{
		Name:      "flato",
		Namespace: "global",
	}
	return provider.SubProvider("vidb")
}

func initMetrics() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Starting metrics server on :8080/metrics")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()
}

func GetVIDBInstance(configInterface vidbinterface.VidbConfigInterface) (*vidb.DB, error) {
	provider := initMetricsProvider()
	db, err := vidb.Open(configInterface, Logger{}, provider)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetVIDBInstance4Rollback(configInterface vidbinterface.VidbConfigInterface) (*vidb.DB, error) {
	db, err := vidb.Open(configInterface, Logger{}, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MicroWrite(db *vidb.DB, operationCount, batchSize int, keySize int, valueSize int) (time.Duration, error) {
	txCount := operationCount / batchSize // 获取到当前写入的总的 批次

	value := make([]byte, valueSize)
	_, _ = rand.Read(value)
	num := 0
	now := time.Now()
	for i := 0; i < txCount; i++ {
		tx, err := db.NewBatchWithEngine()
		if err != nil {
			return 0, err
		}
		for j := 0; j < batchSize; j++ {
			// 根据keySize设置key的长度
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
		if (i != 0 && i%10 == 0) || i == txCount-1 {
			if i%1000 == 0 {
				fmt.Printf("commit seq: %d\n", uint64(i))
			}
			if err := db.Commit(uint64(i)); err != nil {
				return 0, err
			}
		}
	}
	return time.Since(now), nil
}

func MicroWrite4Growth(db *vidb.DB, startNum, operationCount, batchSize int, keySize int, valueSize int, dataPath string) (time.Duration, error) {
	txCount := operationCount / batchSize // 获取到当前写入的总的 批次
	value := make([]byte, valueSize)
	_, _ = rand.Read(value)
	num := startNum
	now := time.Now()
	for i := 0; i < txCount; i++ {
		tx, err := db.NewBatchWithEngine()
		if err != nil {
			return 0, err
		}
		for j := 0; j < batchSize; j++ {
			// 根据keySize设置key的长度
			key := fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", num)
			if err := tx.Put([]byte(key), value); err != nil {
				return 0, err
			}
			num += 1
		}

		seq, err := db.GetSeqNo()
		if err != nil {
			return 0, err
		}
		seq += 1
		if err := tx.Hash(uint64(seq)); err != nil {
			return 0, err
		}

		if err := tx.Write(uint64(seq)); err != nil {
			return 0, err
		}
		if seq != 0 && seq%1000 == 0 {
			if seq%1000 == 0 {
				fmt.Printf("commit seq: %d\n", uint64(seq))
			}
			if err := db.Commit(uint64(seq)); err != nil {
				return 0, err
			}
		}
		if num%(operationCount/10) == 0 {
			fmt.Printf("data volume: %d, time: %s,", uint64(num), time.Now().Format("2006-01-02 15:04:05"))
			// 读取文件获取partition数量
			files, err := os.ReadDir(dataPath) // 指定目录
			if err != nil {
				return 0, fmt.Errorf("error reading directory:%s", err)
			}
			partition_cnt := 0
			regex := regexp.MustCompile(`partition-\d+`)
			for _, file := range files {
				if file.IsDir() && regex.MatchString(file.Name()) {
					partition_cnt += 1
				}
			}
			fmt.Printf("partition count: %d\n", partition_cnt)
		}
	}
	return time.Since(now), nil
}

func MicroRangeQuery(db *vidb.DB, operationCount, batchSize int, keySize int, r int) (time.Duration, error) {
	// txCount := operationCount / batchSize
	// seq := rand.Intn(txCount)
	// seq := 0
	// pos := 0
	// num := rand.Intn(operationCount - r)
	num := 0
	start := fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", num)
	end := fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", num+r)

	fmt.Printf("测试起始Key: %s, 结束Key: %s\n", start, end)
	// fmt.Println(fmt.Sprintf("查询的区块为: %d, 在写入Batch中为位置: %d", seq, pos))
	// end := ""
	// for _, r := range rs {
	// if pos+r > batchSize {
	// 	sum := pos + r - batchSize
	// 	end = fmt.Sprintf("-account%020d%012d", seq+1, sum)
	// } else {
	// 	end = fmt.Sprintf("-account%020d%012d", seq, pos+r)
	// }

	// fmt.Println("Start= ", start, "End= ", end)
	// now := time.Now()
	// iterator := db.NewIterator([]byte(start), []byte(end))
	// defer iterator.Release()
	// iteratorStartTime := time.Now()
	// for iterator.Next() {
	// 	// 这里可以输出当前 迭代器的数据，这里不需要这个值
	// 	// key := iterator.Key().([]byte)
	// 	// fmt.Println(fmt.Sprintf("key: %v", string(key)))
	// }
	// duration := time.Since(iteratorStartTime) // 构建迭代器到迭代器迭代完成的总时间
	duration, err := ReadIteratorKey(db, []byte(start), []byte(end))

	// fmt.Printf("总耗时: %d 迭代耗时为: %d us , 查询区间数目为: %d TPS: %.2f\n", time.Since(now).Nanoseconds(), duration.Nanoseconds(), r, throughput)
	// fmt.Printf("总耗时: %.6f 迭代耗时为: %.6f s , 查询区间数目为: %d TPS: %.2f\n", time.Since(now).Seconds(), duration.Seconds(), r, throughput)
	return duration, err
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
	fmt.Printf("测试起始Key: %s, 结束Key: %s\n", keys[0], keys[r-1])
	startTime := time.Now()
	for _, key := range keys {
		_, _ = db.Get(key)
	}
	return time.Since(startTime), nil
}

func RandomWrite(db *vidb.DB, operationCount, batchSize int, keySize int, valueSize int) (time.Duration, error) {
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
	dur := time.Since(now)
	if seq != 0 && seq%10 == 0 {
		if seq%1000 == 0 {
			fmt.Printf("commit seq: %d\n", uint64(seq))
		}
		if err := db.Commit(seq); err != nil {
			return 0, err
		}
	}
	return dur, nil
}

func RandomRead(db *vidb.DB, operationCount, batchSize int, keySize int, valueSize int) (time.Duration, error) {
	keys := make([][]byte, 0)
	for i := 0; i < batchSize; i++ {
		num := rand.Intn(operationCount)
		key_str := fmt.Sprintf("-account%0"+strconv.Itoa(keySize)+"d", num)
		keys = append(keys, []byte(key_str))
	}
	startTime := time.Now()
	for _, key := range keys {
		_, _ = db.Get(key)
		// time.Sleep(500 * time.Nanosecond)
	}
	return time.Since(startTime), nil
}

func UpdateBatchKey(db *vidb.DB, keys [][]byte, valueSize int, delete bool) (time.Duration, error) {
	value := make([]byte, valueSize)
	if delete {
		value = nil
	} else {
		_, _ = rand.Read(value)
	}
	tx, _ := db.NewBatch()
	now := time.Now()

	for _, key := range keys {
		if delete {
			if err := tx.Delete(key); err != nil {
				return 0, err
			}
		} else {
			if err := tx.Put(key, value); err != nil {
				return 0, err
			}
		}

	}
	seq, err := db.GetSeqNo()
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
		if seq%1000 == 0 {
			fmt.Printf("commit seq: %d\n", uint64(seq))
		}
		if err := db.Commit(seq); err != nil {
			return 0, err
		}
	}
	return time.Since(now), nil
}

func ReadBatchKey(db *vidb.DB, keys [][]byte) (time.Duration, error) {
	startTime := time.Now()
	for _, key := range keys {
		_, _ = db.Get(key)
		// if err != nil {
		// 	fmt.Printf("key not found: %s, %s\n", string(key), err)
		// }
	}
	return time.Since(startTime), nil
}
func ReadIteratorKey(db *vidb.DB, startKey, endKey []byte) (time.Duration, error) {
	startTime := time.Now()
	manager := db.NewIteratorManager(startKey, endKey)
	defer manager.Release()
	for manager.Next() {
	}
	return time.Since(startTime), nil
}

func LatencyGet(db *vidb.DB, key []byte) (time.Duration, time.Duration, time.Duration, error) {
	_, t, t2, t3, err := db.GetDuration(key)
	return t, t2, t3, err
}

func LatencyIterator(db *vidb.DB, start, end []byte) (int64, int64, error) {
	manager := db.NewIteratorManager(start, end)
	defer manager.Release()
	var totalNextDuration time.Duration
	var totalTreeDuration time.Duration
	j := int64(0)
	for manager.Next() {
		next, tree := manager.Duration()
		totalTreeDuration += next
		totalNextDuration += tree
		j++
	}
	return totalNextDuration.Nanoseconds() / j, totalTreeDuration.Nanoseconds() / j, nil
}
