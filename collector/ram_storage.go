package collector

import "github.com/KyberNetwork/node-monitor/types"

type RamStorage interface {
	UpdateStatNormalDataTest(endpointName string, arrayData *[]types.DataNormalTest)
	GetStatNormalData() map[string]map[string]*types.StatReturn
	GetTotalRequest() int64
}
