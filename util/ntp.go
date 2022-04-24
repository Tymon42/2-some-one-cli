package util

import (
	"time"

	"github.com/lixiangyun/go_ntp"
)

const (
	ServerAddr = ":8080"
)

func GetTimeOffset() int64 {
	ntpc := ntp.NewNTPC(ServerAddr, time.Second*1)

	results := ntpc.SyncBatch(10)

	result := ntp.ResultAverage(results)

	return (result.Offset.NanoSecond / int64(time.Millisecond))
}
