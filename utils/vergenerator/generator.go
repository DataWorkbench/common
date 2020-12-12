package vergenerator

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"net"
	"time"

	"github.com/yu31/snowflake"
)

// VerGenerator implements an Version Generator uses to generate unique version.
type VerGenerator struct {
	worker *snowflake.Snowflake
}

// New return an new VerGenerator
func New() *VerGenerator {
	worker, err := snowflake.New(defaultInstanceID())
	if err != nil {
		panic(fmt.Errorf("unexpected error %v", err))
	}
	g := &VerGenerator{
		worker: worker,
	}
	return g
}

// Take return a new unique id
func (g *VerGenerator) Take() (int64, error) {
	return g.worker.Next()
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
		ret = rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
	}

	return ret % 1024
}
