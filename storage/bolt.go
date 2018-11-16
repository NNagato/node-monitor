package storage

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/KyberNetwork/node-monitor/types"
	"github.com/boltdb/bolt"
)

const (
	PATH = "../db/gin.db"

	DAY = 86400
)

var listBucket = []string{"infura", "semi-node", "quik-node", "knstat-node"}

var listBucketNormal = []string{"normal_infura", "normal_semi-node", "normal_quik-node", "normal_knstat-node"}

type BoltStorage struct {
	db *bolt.DB
}

func NewStorage() *BoltStorage {
	var err error
	var db *bolt.DB
	db, err = bolt.Open(PATH, 0600, nil)
	if err != nil {
		panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		for _, value := range listBucket {
			tx.CreateBucket([]byte(value))
		}
		for _, value := range listBucketNormal {
			tx.CreateBucket([]byte(value))
		}
		return nil
	})
	storage := &BoltStorage{
		db: db,
	}
	return storage
}

func uint64ToBytes(u uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	return b
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func handeRawData(listRawData *[]types.RawStressData, typeRPC string) types.DataStressTest {
	localValue := *listRawData
	totalNumReQuest := len(localValue)
	var maxTimeRes float64
	var minTimeRes float64
	if len(localValue) > 0 {
		minTimeRes = localValue[0].TimeResponse
	}
	var numRequestSuccess uint64
	var numRequestFailed uint64
	var totalTimeRes float64
	for _, data := range *listRawData {
		totalTimeRes += data.TimeResponse
		if data.Success == true {
			numRequestSuccess++
		} else {
			numRequestFailed++
		}
		if data.TimeResponse > maxTimeRes {
			maxTimeRes = data.TimeResponse
		}
		if data.TimeResponse < minTimeRes {
			minTimeRes = data.TimeResponse
		}
	}
	dataStressTest := types.DataStressTest{
		MinTimeRes:        minTimeRes,
		MaxTimeRes:        maxTimeRes,
		NumRequestSuccess: numRequestSuccess,
		NumRequestFailed:  numRequestFailed,
		AverageTimeRest:   totalTimeRes / float64(totalNumReQuest),
		SuccessRate:       float64(numRequestSuccess) / float64(totalNumReQuest),
		TypeRPC:           typeRPC,
		TotalNumReQuest:   uint64(totalNumReQuest),
	}
	return dataStressTest
}

func (self *BoltStorage) StoreDataStressTest(endPointName string, listRawData *[]types.RawStressData, timeTest int64, typeRPC string) {
	var err error
	err = self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(endPointName))
		dataStressTest := handeRawData(listRawData, typeRPC)

		dataJson, err = json.Marshal(dataStressTest)
		if err != nil {
			log.Println(err)
			return err
		}
		err = b.Put(uint64ToBytes(uint64(timeTest)), dataJson)
		if err != nil {
			return err
			log.Println(err)
		}
		return err
	})
	if err != nil {
		log.Println(err)
	}
}

func (self *BoltStorage) GetData() (map[string]map[string][]types.DataReturn, error) {
	var err error
	result := make(map[string]map[string][]types.DataReturn)
	err = self.db.View(func(tx *bolt.Tx) error {
		for _, bucket := range listBucket {
			// allDataReturn := []types.DataReturn{}
			dataWithType := make(map[string][]types.DataReturn)
			b := tx.Bucket([]byte(bucket))
			b.ForEach(func(k, v []byte) error {
				timeExecute := bytesToUint64(k)
				data := types.DataStressTest{}
				err := json.Unmarshal(v, &data)
				if err != nil {
					return err
				}
				dataReturn := types.GetDataReturn(data, timeExecute)
				dataWithType[data.TypeRPC] = append(dataWithType[data.TypeRPC], dataReturn)
				return err
			})
			result[bucket] = dataWithType
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (self *BoltStorage) StoreDataNormalTest(endPointName string, data *[]types.DataNormalTest, timeTest int64) {
	nameBucket := "normal_" + endPointName
	arrData := *data
	var err error
	err = self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nameBucket))
		var dataJSON []byte
		for index, data := range arrData {
			dataJSON, err = json.Marshal(data)
			if err != nil {
				log.Println(err)
				return err
			}
			err = b.Put(uint64ToBytes(uint64(timeTest*1000+int64(index))), dataJSON)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		return err
	})
	if err != nil {
		log.Println(err)
	}
}

func getNameNormalBucket(normalBucket string) string {
	arrString := strings.Split(normalBucket, "_")
	return arrString[1]
}

func (self *BoltStorage) GetDataNormal(fromTime, toTime uint64) (map[string]map[string][]types.DataNormalReturn, error) {
	var err error
	result := make(map[string]map[string][]types.DataNormalReturn)
	if fromTime >= toTime || fromTime+DAY < toTime {
		return result, errors.New("time range should be smaller than one day (86400s)")
	}
	fromTimeCompare := fromTime * 1000
	toTimeCompare := (toTime + 1) * 1000
	err = self.db.View(func(tx *bolt.Tx) error {
		for _, bucket := range listBucketNormal {
			dataWithType := make(map[string][]types.DataNormalReturn)
			b := tx.Bucket([]byte(bucket))
			b.ForEach(func(k, v []byte) error {
				timeExecute := bytesToUint64(k)
				if timeExecute >= fromTimeCompare && timeExecute <= toTimeCompare {
					data := types.DataNormalTest{}
					err := json.Unmarshal(v, &data)
					if err != nil {
						return err
					}
					dataNormalReturn := types.GetDataNormalReturn(data, timeExecute)
					dataWithType[data.TypeRPC] = append(dataWithType[data.TypeRPC], dataNormalReturn)
					return err
				}
				return err
			})
			result[getNameNormalBucket(bucket)] = dataWithType
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (self *BoltStorage) GetDataStatNormal(timeStamp uint64) (map[string]map[string]types.StatReturn, error) {
	// timeQuery := time.Unix(int64(timeStamp), 0)
	// year, month, day := timeQuery.Date()
	// beginOfDay := time.Date(year, month, day, 0, 0, 0, 0, nil).Unix()
	// endOfDay := time.Date(year, month, day, 23, 59, 59, 0, nil).Unix()
	timeCompare := (timeStamp + 1) * 1000
	var err error
	result := make(map[string]map[string]types.StatReturn)
	err = self.db.View(func(tx *bolt.Tx) error {
		for _, bucket := range listBucketNormal {
			dataWithType := make(map[string][]types.DataNormalTest)
			b := tx.Bucket([]byte(bucket))
			b.ForEach(func(k, v []byte) error {
				timeExecute := bytesToUint64(k)
				// if timeExecute >= beginOfDay && timeExecute <= endOfDay {
				if timeExecute <= timeCompare {
					data := types.DataNormalTest{}
					err := json.Unmarshal(v, &data)
					if err != nil {
						return err
					}
					dataWithType[data.TypeRPC] = append(dataWithType[data.TypeRPC], data)
				}
				return err
			})
			statReturn := types.GetStatReturn(dataWithType)
			result[getNameNormalBucket(bucket)] = statReturn
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}
