package idgenerator

import (
	"fmt"

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
	worker, err := snowflake.New(cfg.instanceId)
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
		return "", err
	}
	idStr, err := g.hashId.EncodeInt64([]int64{id})
	if err != nil {
		return "", err
	}
	if len(idStr) != 16 {
		return "", fmt.Errorf("generate new id failed, exepcted length 16 bytes, get: %s", idStr)
	}
	return g.prefix + idStr, nil
}

//// TakeV1 return a new unique id that format with `prefix` + `16 bytes string`.
//func (g *IDGenerator) TakeV1() (string, error) {
//	id, err := g.worker.Next()
//	if err != nil {
//		return "", err
//	}
//	// Encode id to hex string.
//	buf := make([]byte, 8)
//	binary.BigEndian.PutUint64(buf, uint64(id))
//
//	lp := len(g.prefix)
//	dst := make([]byte, lp+hex.EncodedLen(len(buf)))
//
//	copy(dst[:lp], g.prefix)
//	hex.Encode(dst[lp:], buf)
//
//	return *(*string)(unsafe.Pointer(&dst)), nil
//}
