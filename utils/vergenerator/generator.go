package vergenerator

import (
	"fmt"

	"github.com/yu31/snowflake"
)

// VerGenerator implements an Version Generator uses to generate unique version.
type VerGenerator struct {
	worker *snowflake.Snowflake
}

// New return an new VerGenerator
func New(opts ...Option) *VerGenerator {
	cfg := applyOptions(opts...)
	worker, err := snowflake.New(cfg.instanceId)
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
