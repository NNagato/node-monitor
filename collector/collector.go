package collector

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/KyberNetwork/node-monitor/blockchain"
	"github.com/KyberNetwork/node-monitor/storage"
	"github.com/KyberNetwork/node-monitor/types"
)

const (
	HOUR      = 3600
	TIME_REST = 1800
)

type Collector struct {
	mu             *sync.RWMutex
	tokens         map[string]types.Token
	listNode       map[string]string
	listBlockChain map[string]*blockchain.Blockchain
	dataGetRate    string
	storage        BoltStorage
	configData     types.ConfigData
	ramStorage     RamStorage
	inStressTest   bool
}

func NewCollector(bolt *storage.BoltStorage) *Collector {
	mu := &sync.RWMutex{}
	file, err := ioutil.ReadFile("../env/production.json")
	if err != nil {
		log.Print(err)
		panic(err)
	}
	configData := types.ConfigData{}
	err = json.Unmarshal(file, &configData)
	if err != nil {
		log.Println(err)
	}

	dataGetRate := strings.Join(configData.DataGetRate[:], "")

	tokens := configData.Tokens
	listNode := make(map[string]string)
	listBlockChain := make(map[string]*blockchain.Blockchain)
	toAddress := configData.Wrapper
	for _, node := range configData.ListNode {
		listNode[node.Name] = node.EndPoint
		listBlockChain[node.Name], err = blockchain.NewBlockchain(node.Name, node.EndPoint, toAddress)
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}
	ramStorage := storage.NewRamStorage(listBlockChain)

	return &Collector{
		mu:             mu,
		storage:        bolt,
		tokens:         tokens,
		listNode:       listNode,
		listBlockChain: listBlockChain,
		dataGetRate:    dataGetRate,
		configData:     configData,
		ramStorage:     ramStorage,
		inStressTest:   false,
	}
}

func (self *Collector) CollectData() {

	// tickerStress := time.NewTicker(time.Duration(self.configData.StressTestInterval) * time.Second)
	// tickerNormal := time.NewTicker(time.Duration(self.configData.NormalTestInterval) * time.Second)
	// go func() {
	// 	time.Sleep(259200 * time.Second)
	// 	for {
	// 		self.RunStressTest()
	// 		<-tickerStress.C
	// 	}
	// }()

	go func() {
		// for {
		// 	inStressTest := self.IsInStressTest()
		// 	if inStressTest == false {
		// 		self.RunNormalTest()
		// 	}
		// 	<-tickerNormal.C
		// }
		self.RunNormalTest()
	}()
}

func (self *Collector) IsInStressTest() bool {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.inStressTest
}

func (self *Collector) SetInStressTest(b bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.inStressTest = b
}

func (self *Collector) RunStressTest() {
	self.SetInStressTest(true)

	self.RunTestGetRate(self.configData.NumberRequest)
	// time.Sleep(TIME_REST * time.Second)

	// self.mu.Lock()
	// self.inStressTest = false
	// self.mu.Unlock()

	time.Sleep(HOUR * time.Second)

	// self.mu.Lock()
	// self.inStressTest = true
	// self.mu.Unlock()

	self.RunTestEstimateGas(self.configData.NumberRequest)
	// time.Sleep(TIME_REST * time.Second)

	// self.mu.Lock()
	// self.inStressTest = false
	// self.mu.Unlock()

	time.Sleep(HOUR * time.Second)

	// self.mu.Lock()
	// self.inStressTest = true
	// self.mu.Unlock()
	self.RunTestGetKyberEnnable(self.configData.NumberRequest)
	// time.Sleep(TIME_REST * time.Second)

	// self.mu.Lock()
	// self.inStressTest = false
	// self.mu.Unlock()

	time.Sleep(HOUR * time.Second)

	// self.mu.Lock()
	// self.inStressTest = true
	// self.mu.Unlock()

	self.RunTestGetBlockNum(self.configData.NumberRequest)
	// time.Sleep(TIME_REST * time.Second)
	time.Sleep(86400 * time.Second)

	self.SetInStressTest(false)
}

func (self *Collector) RunNormalTest() {
	for _, b := range self.listBlockChain {
		go func(b1 *blockchain.Blockchain) {
			for i := 0; i < 10; i++ {
				go func(blockchain *blockchain.Blockchain) {
					for {
						self.RunTestNode(blockchain)
					}
				}(b1)
			}
		}(b)
	}
}

func (self *Collector) RunTestNode(blockchain *blockchain.Blockchain) {
	arrDataTest := &[]types.DataNormalTest{}
	// timeTest := time.Now().Unix()
	wg := &sync.WaitGroup{}
	// wg.Add(4)
	// go self.RunTestEstimateGasNormal(blockchain, wg, arrDataTest)
	// go self.RunTestGetRateNormal(self.configData.Wrapper, self.dataGetRate, blockchain, wg, arrDataTest)
	// go self.RunTestGetBlockNumber(blockchain, wg, arrDataTest)
	// go self.RunTestCheckKyberEnable(blockchain, wg, arrDataTest)

	self.RunTestEstimateGasNormal(blockchain, wg, arrDataTest)
	self.RunTestGetRateNormal(self.configData.Wrapper, self.dataGetRate, blockchain, wg, arrDataTest)
	self.RunTestGetBlockNumber(blockchain, wg, arrDataTest)
	self.RunTestCheckKyberEnable(blockchain, wg, arrDataTest)

	// wg.Wait()
	// self.storage.StoreDataNormalTest(blockchain.EndPointName, arrDataTest, timeTest)
	// log.Println("data: ", blockchain.EndPointName, arrDataTest, len(*arrDataTest))
	self.ramStorage.UpdateStatNormalDataTest(blockchain.EndPointName, arrDataTest)
}

// *******---------------------- Run stress test

func (self *Collector) RunTestEstimateGas(numberRequest int) {

	for _, blockchain := range self.listBlockChain {
		timeTest := time.Now().Unix()
		// wGroup := &sync.WaitGroup{}
		listRawStressData := &[]types.RawStressData{}
		chanRawStressData := make(chan types.RawStressData)
		for i := 0; i < numberRequest; i++ {
			// wGroup.Add(1)
			go blockchain.EstimateGas("", "", chanRawStressData)
		}

		for i := 0; i < numberRequest; i++ {
			// wGroup.Add(1)
			*listRawStressData = append(*listRawStressData, <-chanRawStressData)
		}
		// wGroup.Wait()
		go self.storage.StoreDataStressTest(blockchain.EndPointName, listRawStressData, timeTest, "eth_estimateGas")
	}
}

func (self *Collector) RunTestGetRate(numberRequest int) {
	for _, blockchain := range self.listBlockChain {
		timeTest := time.Now().Unix()
		listRawStressData := &[]types.RawStressData{}
		chanRawStressData := make(chan types.RawStressData)
		for i := 0; i < numberRequest; i++ {
			go blockchain.GetRate("", self.dataGetRate, chanRawStressData)
		}
		for i := 0; i < numberRequest; i++ {
			// wGroup.Add(1)
			*listRawStressData = append(*listRawStressData, <-chanRawStressData)
		}
		go self.storage.StoreDataStressTest(blockchain.EndPointName, listRawStressData, timeTest, "getRate")
	}
}

func (self *Collector) RunTestGetBlockNum(numberRequest int) {
	for _, blockchain := range self.listBlockChain {
		timeTest := time.Now().Unix()
		listRawStressData := &[]types.RawStressData{}
		chanRawStressData := make(chan types.RawStressData)
		for i := 0; i < numberRequest; i++ {
			go blockchain.GetLatestBlock("", "", chanRawStressData)
		}
		for i := 0; i < numberRequest; i++ {
			// wGroup.Add(1)
			*listRawStressData = append(*listRawStressData, <-chanRawStressData)
		}
		go self.storage.StoreDataStressTest(blockchain.EndPointName, listRawStressData, timeTest, "eth_blockNumber")
	}
}

func (self *Collector) RunTestGetKyberEnnable(numberRequest int) {
	for _, blockchain := range self.listBlockChain {
		timeTest := time.Now().Unix()
		listRawStressData := &[]types.RawStressData{}
		chanRawStressData := make(chan types.RawStressData)
		for i := 0; i < numberRequest; i++ {
			go blockchain.CheckKyberEnable(self.configData.Network, "0x238dafe0", chanRawStressData)
		}
		for i := 0; i < numberRequest; i++ {
			// wGroup.Add(1)
			*listRawStressData = append(*listRawStressData, <-chanRawStressData)
		}
		go self.storage.StoreDataStressTest(blockchain.EndPointName, listRawStressData, timeTest, "kyberEnable")
	}
}

// *******---------------------- Run normal test

func (self *Collector) RunTestEstimateGasNormal(blockchain *blockchain.Blockchain, wg *sync.WaitGroup, arrDataTest *[]types.DataNormalTest) {
	// defer wg.Done()
	result := blockchain.EstimateGasNormal()
	*arrDataTest = append(*arrDataTest, result)
}

func (self *Collector) RunTestGetRateNormal(to string, data string, blockchain *blockchain.Blockchain, wg *sync.WaitGroup, arrDataTest *[]types.DataNormalTest) {
	// defer wg.Done()
	result := blockchain.GetRateNormal(to, data)
	*arrDataTest = append(*arrDataTest, result)
}

func (self *Collector) RunTestGetBlockNumber(blockchain *blockchain.Blockchain, wg *sync.WaitGroup, arrDataTest *[]types.DataNormalTest) {
	// defer wg.Done()
	result := blockchain.GetLatestBlockNormal()
	*arrDataTest = append(*arrDataTest, result)
}

func (self *Collector) RunTestCheckKyberEnable(blockchain *blockchain.Blockchain, wg *sync.WaitGroup, arrDataTest *[]types.DataNormalTest) {
	// defer wg.Done()
	result := blockchain.CheckKyberEnableNormal("0x238dafe0", self.configData.Network)
	*arrDataTest = append(*arrDataTest, result)
}

// read data

func (self *Collector) GetStatNormalData() map[string]map[string]*types.StatReturn {
	return self.ramStorage.GetStatNormalData()
}

func (self *Collector) GetTotalRequest() int64 {
	return self.ramStorage.GetTotalRequest()
}
