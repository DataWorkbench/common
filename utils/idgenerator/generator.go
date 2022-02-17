package idgenerator

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
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
func New(prefix string, opts ...Option) *IDGenerator {
	cfg := applyOptions(opts...)
	instanceId := getInstanceId(&cfg)
	worker, err := snowflake.New(instanceId)
	if err != nil {
		panic(fmt.Errorf("unexpected error %v", err))
	}

	g := &IDGenerator{
		worker: worker,
		prefix: prefix,
	}
	return g
}

// Take return a new unique id that format with `prefix` + `16 bytes string`.
func (g *IDGenerator) Take() (string, error) {
	id, err := g.worker.Next()
	if err != nil {
		log.Printf("IDGenerator: take new id from worker error: %v\n", err)
		return "", err
	}
	return g.encode(id), nil
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
