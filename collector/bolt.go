package collector

import "github.com/KyberNetwork/node-monitor/types"

type BoltStorage interface {
	StoreDataStressTest(string, *[]types.RawStressData, int64, string)
	StoreDataNormalTest(endPointName string, data *[]types.DataNormalTest, timeTest int64)
}
