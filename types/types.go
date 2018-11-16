package types

type Token struct {
	Symbol  string `json:"symbol"`
	Decimal uint64 `json:"decimal"`
	Address string `json:"address"`
}

type Node struct {
	EndPoint string `json:"endpoint"`
	Name     string `json:"name"`
}

type ConfigData struct {
	ListNode           []Node           `json:"listNode"`
	Tokens             map[string]Token `json:"tokens"`
	Reserve            string           `json:"reserve"`
	Network            string           `json:"network"`
	Wrapper            string           `json:"wrapper"`
	DataGetRate        []string         `json:"dataGetRate"`
	StressTestInterval int64            `json:"stressInterval"`
	NormalTestInterval int64            `json:"normalInterval"`
	NumberRequest      int              `json:"numberRequest"`
}

type RawStressData struct {
	TimeResponse float64
	Success      bool
}

type DataStressTest struct {
	MinTimeRes        float64 `json:"minTimeRes"`
	MaxTimeRes        float64 `json:"maxTimeRes"`
	AverageTimeRest   float64 `json:"averageTimeRes"`
	TotalNumReQuest   uint64  `json:"totalNumRequest"`
	NumRequestSuccess uint64  `json:"numRequestSuccess"`
	NumRequestFailed  uint64  `json:"numRequestFailed"`
	SuccessRate       float64 `json:"successRate"`
	TypeRPC           string  `json:"typeRPC"`
}

type DataNormalTest struct {
	Success      bool    `json:"success"`
	TimeResponse float64 `json:"timeResponse"`
	TypeRPC      string  `json:"typeRPC"`
}

func GetStatReturn(data map[string][]DataNormalTest) map[string]StatReturn {
	result := make(map[string]StatReturn)
	for key, listData := range data {
		var minTimeRes float64
		var maxTimeRes float64
		if len(listData) > 0 {
			minTimeRes = listData[0].TimeResponse
		}
		var averageTimeRes float64
		var totalRequest uint64
		var numRequestSuccess uint64
		var numRequestFailed uint64
		var successRate float64
		for _, ele := range listData {
			totalRequest++
			timeResponse := ele.TimeResponse
			if timeResponse > maxTimeRes {
				maxTimeRes = timeResponse
			}
			if timeResponse < minTimeRes {
				minTimeRes = timeResponse
			}
			if ele.Success == true {
				numRequestSuccess++
			} else {
				numRequestFailed++
			}
		}
		successRate = float64(numRequestSuccess) / float64(totalRequest)
		statReturn := StatReturn{
			MinTimeRes:        minTimeRes,
			MaxTimeRes:        maxTimeRes,
			AverageTimeRest:   averageTimeRes,
			TotalNumReQuest:   totalRequest,
			NumRequestSuccess: numRequestSuccess,
			NumRequestFailed:  numRequestFailed,
			SuccessRate:       successRate,
		}
		result[key] = statReturn
	}
	return result
}

type DataNormalReturn struct {
	TimeTest     uint64  `json:"timeExecute"`
	Success      bool    `json:"success"`
	TimeResponse float64 `json:"timeResponse"`
}

func GetDataNormalReturn(data DataNormalTest, timeExecute uint64) DataNormalReturn {
	return DataNormalReturn{
		TimeTest:     timeExecute,
		TimeResponse: data.TimeResponse,
		Success:      data.Success,
	}
}

type DataReturn struct {
	TimeTest          uint64  `json:"timeExecute"`
	MinTimeRes        float64 `json:"minTimeRes"`
	MaxTimeRes        float64 `json:"maxTimeRes"`
	AverageTimeRest   float64 `json:"averageTimeRes"`
	TotalNumReQuest   uint64  `json:"totalNumRequest"`
	NumRequestSuccess uint64  `json:"numRequestSuccess"`
	NumRequestFailed  uint64  `json:"numRequestFailed"`
	SuccessRate       float64 `json:"successRate"`
}

type StatReturn struct {
	MinTimeRes        float64 `json:"minTimeRes"`
	MaxTimeRes        float64 `json:"maxTimeRes"`
	AverageTimeRest   float64 `json:"averageTimeRes"`
	TotalNumReQuest   uint64  `json:"totalNumRequest"`
	NumRequestSuccess uint64  `json:"numRequestSuccess"`
	NumRequestFailed  uint64  `json:"numRequestFailed"`
	SuccessRate       float64 `json:"successRate"`
}

func GetDataReturn(data DataStressTest, timeExecute uint64) DataReturn {
	return DataReturn{
		TimeTest:          timeExecute,
		MinTimeRes:        data.MinTimeRes,
		MaxTimeRes:        data.MaxTimeRes,
		AverageTimeRest:   data.AverageTimeRest,
		TotalNumReQuest:   data.TotalNumReQuest,
		NumRequestSuccess: data.NumRequestSuccess,
		NumRequestFailed:  data.NumRequestFailed,
		SuccessRate:       data.SuccessRate,
	}
}
