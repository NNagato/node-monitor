package server

import "github.com/KyberNetwork/node-monitor/types"

type Storage interface {
	GetData() (map[string]map[string][]types.DataReturn, error)
	GetDataNormal(fromTime, toTime uint64) (map[string]map[string][]types.DataNormalReturn, error)
	GetDataStatNormal(timeStamp uint64) (map[string]map[string]types.StatReturn, error)
}
