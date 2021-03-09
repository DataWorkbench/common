package idgenerator

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"math/rand"
	"net"
	"time"
	"unsafe"

	"github.com/yu31/snowflake"
)

// IDGenerator implements an ID Generator uses to generate unique ID
// ID format: "prefix" | 16 bytes string
type IDGenerator struct {
	worker *snowflake.Snowflake
	prefix string
}

// New return an new IDGenerator
func New(prefix string) *IDGenerator {
	worker, err := snowflake.New(defaultInstanceID())
	if err != nil {
		panic(fmt.Errorf("unexpected error %v", err))
	}

	g := &IDGenerator{
		worker: worker,
		prefix: prefix,
	}
	return g
}

func (g *IDGenerator) encode(x int64) string {
	buf := make([]byte, 8)

	binary.BigEndian.PutUint64(buf, uint64(x))

	lp := len(g.prefix)
	dst := make([]byte, lp+hex.EncodedLen(len(buf)))

	copy(dst[:lp], g.prefix)
	hex.Encode(dst[lp:], buf)

	return *(*string)(unsafe.Pointer(&dst))
}

// Take return a new unique id
func (g *IDGenerator) Take() (string, error) {
	id, err := g.worker.Next()
	if err != nil {
		return "", err
	}
	return g.encode(id), nil
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
