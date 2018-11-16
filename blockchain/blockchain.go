package blockchain

import (
	"context"
	"log"
	"math/big"
	"strings"
	"time"

	// "github.com/KyberNetwork/server-go/ethereum"
	"github.com/KyberNetwork/node-monitor/types"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	WRAPPER_ABI = `[{"constant":true,"inputs":[{"name":"x","type":"bytes14"},{"name":"byteInd","type":"uint256"}],"name":"getInt8FromByte","outputs":[{"name":"","type":"int8"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"reserve","type":"address"},{"name":"tokens","type":"address[]"}],"name":"getBalances","outputs":[{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"pricingContract","type":"address"},{"name":"tokenList","type":"address[]"}],"name":"getTokenIndicies","outputs":[{"name":"","type":"uint256[]"},{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"x","type":"bytes14"},{"name":"byteInd","type":"uint256"}],"name":"getByteFromBytes14","outputs":[{"name":"","type":"bytes1"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"network","type":"address"},{"name":"sources","type":"address[]"},{"name":"dests","type":"address[]"},{"name":"qty","type":"uint256[]"}],"name":"getExpectedRates","outputs":[{"name":"expectedRate","type":"uint256[]"},{"name":"slippageRate","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"pricingContract","type":"address"},{"name":"tokenList","type":"address[]"}],"name":"getTokenRates","outputs":[{"name":"","type":"uint256[]"},{"name":"","type":"uint256[]"},{"name":"","type":"int8[]"},{"name":"","type":"int8[]"},{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"}]`
	NETWORK_ABI = `[{"constant":false,"inputs":[{"name":"alerter","type":"address"}],"name":"removeAlerter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"reserve","type":"address"},{"name":"src","type":"address"},{"name":"dest","type":"address"},{"name":"add","type":"bool"}],"name":"listPairForReserve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"bytes32"}],"name":"perReserveListedPairs","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"enabled","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"pendingAdmin","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getOperators","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"token","type":"address"},{"name":"amount","type":"uint256"},{"name":"sendTo","type":"address"}],"name":"withdrawToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"maxGasPrice","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newAlerter","type":"address"}],"name":"addAlerter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"negligibleRateDiff","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"feeBurnerContract","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"expectedRateContract","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"whiteListContract","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"user","type":"address"}],"name":"getUserCapInWei","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newAdmin","type":"address"}],"name":"transferAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_enable","type":"bool"}],"name":"setEnable","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"claimAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"isReserve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getAlerters","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"src","type":"address"},{"name":"dest","type":"address"},{"name":"srcQty","type":"uint256"}],"name":"getExpectedRate","outputs":[{"name":"expectedRate","type":"uint256"},{"name":"slippageRate","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"reserves","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newOperator","type":"address"}],"name":"addOperator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"reserve","type":"address"},{"name":"add","type":"bool"}],"name":"addReserve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"operator","type":"address"}],"name":"removeOperator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_whiteList","type":"address"},{"name":"_expectedRate","type":"address"},{"name":"_feeBurner","type":"address"},{"name":"_maxGasPrice","type":"uint256"},{"name":"_negligibleRateDiff","type":"uint256"}],"name":"setParams","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"src","type":"address"},{"name":"dest","type":"address"},{"name":"srcQty","type":"uint256"}],"name":"findBestRate","outputs":[{"name":"","type":"uint256"},{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"src","type":"address"},{"name":"srcAmount","type":"uint256"},{"name":"dest","type":"address"},{"name":"destAddress","type":"address"},{"name":"maxDestAmount","type":"uint256"},{"name":"minConversionRate","type":"uint256"},{"name":"walletId","type":"address"}],"name":"trade","outputs":[{"name":"","type":"uint256"}],"payable":true,"stateMutability":"payable","type":"function"},{"constant":false,"inputs":[{"name":"amount","type":"uint256"},{"name":"sendTo","type":"address"}],"name":"withdrawEther","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getNumReserves","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"token","type":"address"},{"name":"user","type":"address"}],"name":"getBalance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"admin","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_admin","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"sender","type":"address"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"EtherReceival","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"source","type":"address"},{"indexed":false,"name":"dest","type":"address"},{"indexed":false,"name":"actualSrcAmount","type":"uint256"},{"indexed":false,"name":"actualDestAmount","type":"uint256"}],"name":"ExecuteTrade","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"reserve","type":"address"},{"indexed":false,"name":"add","type":"bool"}],"name":"AddReserveToNetwork","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"reserve","type":"address"},{"indexed":false,"name":"src","type":"address"},{"indexed":false,"name":"dest","type":"address"},{"indexed":false,"name":"add","type":"bool"}],"name":"ListReservePairs","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"token","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"sendTo","type":"address"}],"name":"TokenWithdraw","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"sendTo","type":"address"}],"name":"EtherWithdraw","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"pendingAdmin","type":"address"}],"name":"TransferAdminPending","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newAdmin","type":"address"},{"indexed":false,"name":"previousAdmin","type":"address"}],"name":"AdminClaimed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newAlerter","type":"address"},{"indexed":false,"name":"isAdd","type":"bool"}],"name":"AlerterAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newOperator","type":"address"},{"indexed":false,"name":"isAdd","type":"bool"}],"name":"OperatorAdded","type":"event"}]`
)

type Blockchain struct {
	client       *rpc.Client
	endpoint     string
	EndPointName string
	toAddress    string
	wrapperAbi   abi.ABI
	networkAbi   abi.ABI
	ethclient    *ethclient.Client
}

type RateWrapper struct {
	ExpectedRate []*big.Int `json:"expectedRate"`
	SlippageRate []*big.Int `json:"slippageRate"`
}

type ResultEstimateGas struct {
	ID     int    `json:"id"`
	Result string `json:"result"`
}

func NewBlockchain(endPointName string, endpoint string, toAddress string) (*Blockchain, error) {
	client, err := rpc.DialHTTP(endpoint)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	wrapperAbi, err := abi.JSON(strings.NewReader(WRAPPER_ABI))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	networkAbi, err := abi.JSON(strings.NewReader(NETWORK_ABI))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	blockchain := Blockchain{
		client:       client,
		endpoint:     endpoint,
		EndPointName: endPointName,
		toAddress:    toAddress,
		wrapperAbi:   wrapperAbi,
		networkAbi:   networkAbi,
		ethclient:    ethclient.NewClient(client),
	}
	return &blockchain, nil
}

// ***------------------- ethercall

func (self *Blockchain) EthCall(from, to, data, method string) (string, error, time.Duration) {
	params := make(map[string]string)
	if data != "" {
		params["data"] = data
		if to == "" {
			to = self.toAddress
		}
		params["to"] = to
		if from != "" {
			params["from"] = from
		}
	}
	var result string
	timeStart := time.Now()
	err := self.client.Call(&result, method, params, "latest")
	timeExecute := time.Since(timeStart)
	if err != nil {
		log.Println(err)
		return "", err, timeExecute
	}
	return result, nil, timeExecute
}

func (self *Blockchain) EthCallNoParams(method string) (string, error, time.Duration) {
	var result string
	timeStart := time.Now()
	err := self.client.Call(&result, method)
	timeExecute := time.Since(timeStart)
	if err != nil {
		log.Println(err)
		return "", err, timeExecute
	}
	return result, nil, timeExecute
}

// *** -------------------------- Stress test

func (self *Blockchain) GetRate(to string, data string, chanRawStressData chan types.RawStressData) {
	result, err, timeExecute := self.EthCall("", to, data, "eth_call")

	rawStressData := types.RawStressData{
		TimeResponse: timeExecute.Seconds(),
		Success:      false,
	}
	if err == nil {
		rateByte, err := hexutil.Decode(result)
		if err == nil {
			rateWapper := RateWrapper{}
			err = self.wrapperAbi.Unpack(&rateWapper, "getExpectedRates", rateByte)
			if err == nil {
				rawStressData.Success = true
			}
		}
	}
	chanRawStressData <- rawStressData
}

func (self *Blockchain) EstimateGas(to string, data string, chanRawStressData chan types.RawStressData) {
	from := "0x2cc72d8857ac57ba058eddd36b2f14adc2"
	data = "0xcb3c28c7000000000000000000000000eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee00000000000000000000000000000000000000000000000000038d7ea4c680000000000000000000000000004e470dc7321e84ca96fcaedd0c8abcebbaeb68c60000000000000000000000002cc72d8857ac57ba058eddd36b2f14adc2a058bd800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001431b0294fa51ef00000000000000000000000000000000000000000000000000000000000003e1e4b"
	to = "0x9E67f627a17EDed3Fb7C71417DFe4aa7bFb4CaB7"
	toAddress := common.HexToAddress(to)
	value := big.NewInt(1000000000)
	dataByte, err := hexutil.Decode(data)
	if err != nil {
		log.Println("cant decode data: ", err)
	}

	msg := ethereum.CallMsg{
		From:  common.HexToAddress(from),
		To:    &toAddress,
		Value: value,
		Data:  dataByte,
	}

	timeStart := time.Now()
	result, err := self.ethclient.EstimateGas(context.Background(), msg)
	timeExecute := time.Since(timeStart)
	rawStressData := types.RawStressData{
		TimeResponse: timeExecute.Seconds(),
		Success:      false,
	}
	if err == nil && result.Cmp(big.NewInt(0)) != 0 {
		rawStressData.Success = true
	}
	chanRawStressData <- rawStressData
}

func (self *Blockchain) GetLatestBlock(to string, data string, chanRawStressData chan types.RawStressData) {
	result, err, timeExecute := self.EthCallNoParams("eth_blockNumber")

	rawStressData := types.RawStressData{
		TimeResponse: timeExecute.Seconds(),
		Success:      false,
	}
	if err == nil && result != "0x" {
		rawStressData.Success = true
	} else {
		log.Println(err)
	}
	chanRawStressData <- rawStressData
}

func (self *Blockchain) CheckKyberEnable(to string, data string, chanRawStressData chan types.RawStressData) {
	result, err, timeExecute := self.EthCall("", to, data, "eth_call")

	rawStressData := types.RawStressData{
		TimeResponse: timeExecute.Seconds(),
		Success:      false,
	}
	if err == nil {
		enabledByte, err := hexutil.Decode(result)
		if err == nil {
			var enabled bool
			err = self.networkAbi.Unpack(&enabled, "enabled", enabledByte)
			if err == nil {
				rawStressData.Success = true
			} else {
				log.Println(err)
			}
		}
	} else {
		log.Println(err)
	}
	chanRawStressData <- rawStressData
}

// ***----------------------- normal test

func (self *Blockchain) GetRateNormal(to, data string) types.DataNormalTest {
	result, err, timeExecute := self.EthCall("", to, data, "eth_call")

	dataNormalTest := types.DataNormalTest{
		TimeResponse: timeExecute.Seconds(),
		Success:      false,
		TypeRPC:      "get-rate",
	}
	if err == nil {
		rateByte, err := hexutil.Decode(result)
		if err == nil {
			rateWapper := RateWrapper{}
			err = self.wrapperAbi.Unpack(&rateWapper, "getExpectedRates", rateByte)
			if err == nil {
				dataNormalTest.Success = true
			} else {
				log.Println(err)
				return dataNormalTest
			}
		} else {
			log.Println(err)
			return dataNormalTest
		}
	} else {
		log.Println(err)
		return dataNormalTest
	}
	return dataNormalTest
}

func (self *Blockchain) GetLatestBlockNormal() types.DataNormalTest {

	result, err, timeExecute := self.EthCallNoParams("eth_blockNumber")

	dataNormalTest := types.DataNormalTest{
		TimeResponse: timeExecute.Seconds(),
		Success:      false,
		TypeRPC:      "eth_blockNuber",
	}
	if err == nil && result != "0x" {
		dataNormalTest.Success = true
	} else {
		log.Println(err)
		return dataNormalTest
	}
	return dataNormalTest
}

func (self *Blockchain) CheckKyberEnableNormal(data string, to string) types.DataNormalTest {
	result, err, timeExecute := self.EthCall("", to, data, "eth_call")

	dataNormalTest := types.DataNormalTest{
		TimeResponse: timeExecute.Seconds(),
		Success:      false,
		TypeRPC:      "check-kyber-enable",
	}
	if err == nil {
		enabledByte, err := hexutil.Decode(result)
		if err == nil {
			var enabled bool
			err = self.networkAbi.Unpack(&enabled, "enabled", enabledByte)
			if err == nil {
				dataNormalTest.Success = true
			} else {
				log.Println(err)
				return dataNormalTest
			}
		} else {
			log.Println(err)
			return dataNormalTest
		}
	} else {
		log.Println(err)
		return dataNormalTest
	}
	return dataNormalTest
}

func (self *Blockchain) EstimateGasNormal() types.DataNormalTest {
	from := "0x2cc72d8857ac57ba058eddd36b2f14adc2"
	data := "0xcb3c28c7000000000000000000000000eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee00000000000000000000000000000000000000000000000000038d7ea4c680000000000000000000000000004e470dc7321e84ca96fcaedd0c8abcebbaeb68c60000000000000000000000002cc72d8857ac57ba058eddd36b2f14adc2a058bd800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001431b0294fa51ef00000000000000000000000000000000000000000000000000000000000003e1e4b"
	to := "0x9E67f627a17EDed3Fb7C71417DFe4aa7bFb4CaB7"
	toAddress := common.HexToAddress(to)
	value := big.NewInt(1000000000)
	dataByte, err := hexutil.Decode(data)
	if err != nil {
		log.Println("cant decode data: ", err)
	}

	msg := ethereum.CallMsg{
		From:  common.HexToAddress(from),
		To:    &toAddress,
		Value: value,
		Data:  dataByte,
	}

	timeStart := time.Now()
	result, err := self.ethclient.EstimateGas(context.Background(), msg)
	timeExecute := time.Since(timeStart)

	dataNormalTest := types.DataNormalTest{
		TimeResponse: timeExecute.Seconds(),
		Success:      false,
		TypeRPC:      "eth_estimateGas",
	}
	if err == nil && result.Cmp(big.NewInt(0)) > 0 {
		dataNormalTest.Success = true
	} else {
		log.Println(err)
		return dataNormalTest
	}
	return dataNormalTest
}
