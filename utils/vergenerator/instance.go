package vergenerator

import (
	crand "crypto/rand"
	"hash/crc32"
	"math/big"
	mrand "math/rand"
	"net"
	"time"
)

func getInstanceId(cfg *config) int64 {
	if cfg.instanceId != nil {
		return *cfg.instanceId
	}
	return defaultInstanceID()
}

func defaultInstanceID() int64 {
	var ret int64

	itf, err := net.Interfaces()
	if err == nil {
		h := crc32.NewIEEE()
		for i := range itf {
			_, _ = h.Write(itf[i].HardwareAddr)
		}
		ret = int64(h.Sum32())
	} else {
		var n *big.Int
		n, err = crand.Int(crand.Reader, big.NewInt(1024))
		if err != nil {
			ret = mrand.New(mrand.NewSource(time.Now().UnixNano())).Int63()
		} else {
			ret = n.Int64()
		}
	}
	return ret % 1024
}
