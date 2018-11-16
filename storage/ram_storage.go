package storage

import (
	"sync"

	"github.com/KyberNetwork/node-monitor/blockchain"
	"github.com/KyberNetwork/node-monitor/types"
)

type RamStorage struct {
	mu                 *sync.RWMutex
	totolRequest       int64
	statDataNormalTest map[string]map[string]*types.StatReturn
}

func NewRamStorage(listBlockchain map[string]*blockchain.Blockchain) *RamStorage {
	mu := &sync.RWMutex{}
	statDataNormalTest := make(map[string]map[string]*types.StatReturn)
	for _, b := range listBlockchain {
		// log.Println("endpoint: ", b.EndPointName)
		statDataNormalTest[b.EndPointName] = make(map[string]*types.StatReturn)
	}
	return &RamStorage{
		mu:                 mu,
		totolRequest:       0,
		statDataNormalTest: statDataNormalTest,
	}
}

func (self *RamStorage) GetTotalRequest() int64 {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.totolRequest
}

func (self *RamStorage) GetStatNormalData() map[string]map[string]*types.StatReturn {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.statDataNormalTest
}

func (self *RamStorage) UpdateStatNormalDataTest(endpointName string, arrayData *[]types.DataNormalTest) {
	self.mu.Lock()
	defer self.mu.Unlock()
	currentNodeData := self.statDataNormalTest[endpointName]
	arrayDataValue := *arrayData
	for _, dataTest := range arrayDataValue {
		typeRPC := dataTest.TypeRPC
		currentTypeData := currentNodeData[typeRPC]
		if currentTypeData == nil {
			var numRequestSuccess uint64
			var numRequestFailed uint64
			var successRate float64
			if dataTest.Success == true {
				numRequestSuccess = 1
				numRequestFailed = 0
				successRate = 1
			} else {
				numRequestSuccess = 0
				numRequestFailed = 1
				successRate = 0
			}
			currentTypeData = &types.StatReturn{
				TotalNumReQuest:   1,
				MinTimeRes:        dataTest.TimeResponse,
				MaxTimeRes:        dataTest.TimeResponse,
				AverageTimeRest:   dataTest.TimeResponse,
				NumRequestSuccess: numRequestSuccess,
				NumRequestFailed:  numRequestFailed,
				SuccessRate:       successRate,
			}
			currentNodeData[typeRPC] = currentTypeData
		} else {
			if dataTest.TimeResponse < currentTypeData.MinTimeRes {
				currentTypeData.MinTimeRes = dataTest.TimeResponse
			}
			if dataTest.TimeResponse > currentTypeData.MaxTimeRes {
				currentTypeData.MaxTimeRes = dataTest.TimeResponse
			}
			if dataTest.Success == true {
				currentTypeData.NumRequestSuccess = currentTypeData.NumRequestSuccess + 1
			} else {
				currentTypeData.NumRequestFailed = currentTypeData.NumRequestFailed + 1
			}
			currentTotalNumReQuest := currentTypeData.TotalNumReQuest + 1
			averageTimeRest := (currentTypeData.AverageTimeRest*float64(currentTypeData.TotalNumReQuest) + dataTest.TimeResponse) / float64(currentTotalNumReQuest)
			successRate := float64(currentTypeData.NumRequestSuccess) / float64(currentTotalNumReQuest)
			currentTypeData.AverageTimeRest = averageTimeRest
			currentTypeData.SuccessRate = successRate
			currentTypeData.TotalNumReQuest = currentTotalNumReQuest
		}
	}
	self.totolRequest = self.totolRequest + int64(len(arrayDataValue))
	self.statDataNormalTest[endpointName] = currentNodeData
}
