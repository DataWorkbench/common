package idgeneratorv2

import (
	"fmt"
	"log"

	"github.com/speps/go-hashids/v2"
	"github.com/yu31/snowflake"
)

// IDGenerator implements an ID Generator uses to generate unique ID
// ID format: "prefix" | 16 bytes string
type IDGenerator struct {
	prefix string
	worker *snowflake.Snowflake
	hashId *hashids.HashID
}

// New return an new IDGenerator
func New(prefix string, opts ...Option) *IDGenerator {
	cfg := applyOptions(opts...)
	worker, err := snowflake.New(*cfg.instanceId)
	if err != nil {
		panic(fmt.Errorf("IDGenerator: unexpected error %v", err))
	}

	hd := hashids.NewData()
	hd.Alphabet = cfg.hashAlphabet
	hd.MinLength = cfg.hashMinLength
	hd.Salt = cfg.hashSalt
	hashId, err := hashids.NewWithData(hd)
	if err != nil {
		panic(fmt.Errorf("IDGenerator: unexpected error: %v", err))
	}

	g := &IDGenerator{
		prefix: prefix,
		worker: worker,
		hashId: hashId,
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
	return g.encode(id)
}

func (g *IDGenerator) encode(x int64) (string, error) {
	id, err := g.hashId.EncodeInt64([]int64{x})
	if err != nil {
		return "", err
	}
	return g.prefix + id, nil
}
