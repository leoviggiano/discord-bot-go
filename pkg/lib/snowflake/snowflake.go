package snowflake

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	// Epoch representa o timestamp de início em milissegundos
	Epoch = 1609459200000 // 2021-01-01 00:00:00 UTC
	// NodeID é o ID do nó (máquina)
	NodeID = 1
	// SequenceBits define o número de bits para o sequencial
	SequenceBits = 12
)

var (
	ErrInvalidSnowflake = errors.New("snowflake: invalid snowflake")
)

type Snowflake uint64

func (s Snowflake) String() string {
	return fmt.Sprintf("%d", s)
}

func FromString(s string) (Snowflake, error) {
	snowflake, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, ErrInvalidSnowflake
	}

	return Snowflake(snowflake), nil
}

type SnowflakeGenerator interface {
	Generate() Snowflake
}

type snowflakeGenerator struct {
	lastTimestamp uint64
	sequence      uint64
}

func NewSnowflakeGenerator() SnowflakeGenerator {
	return &snowflakeGenerator{
		lastTimestamp: 0,
		sequence:      0,
	}
}

func (s *snowflakeGenerator) Generate() Snowflake {
	timestamp := uint64(time.Now().UnixNano() / 1000000)
	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & ((1 << SequenceBits) - 1)
		if s.sequence == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = uint64(time.Now().UnixNano() / 1000000)
			}
		}
	} else {
		s.sequence = 0
	}
	s.lastTimestamp = timestamp
	snowflake := ((timestamp - Epoch) << (64 - 41)) | (uint64(NodeID) << (64 - 41 - 10)) | s.sequence
	return Snowflake(snowflake)
}
