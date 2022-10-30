package idgeneratorlockfree

import (
	"errors"
	"time"
)

const (
	TimestampLength = 41
	NodeIDLength    = 10
	SequenceLength  = 12
	MaxSequence     = 1<<SequenceLength - 1
	MaxTimestamp    = 1<<TimestampLength - 1
	MaxNodeID       = 1<<NodeIDLength - 1
)

type SequenceResolver func(ms int64) (uint16, error)

var (
	resolver      SequenceResolver
	nodeID              = PrivateIPToNodeID()
	lastTimestamp int64 = -1
)

func NextID() (uint64, error) {
	current := currentMillis()
	seqResolver := callSequenceResolver()
	seq, err := seqResolver(current)

	if current < lastTimestamp {
		return 0, errors.New("system clock moving backward. Refuse to generate new id")
	}

	if current == lastTimestamp {
		seq = (seq + 1) & MaxSequence
		if seq == 0 {
			current = tillNextMillis(current)
		}
	} else {
		seq = 0
	}

	lastTimestamp = current

	return uint64(current)<<(NodeIDLength+SequenceLength) |
		((uint64(nodeID)) << SequenceLength) |
		uint64(seq), err
}

func callSequenceResolver() SequenceResolver {
	if resolver == nil {
		return AtomicResolver
	}

	return resolver
}

func tillNextMillis(lastTimestamp int64) int64 {
	now := currentMillis()

	for now == lastTimestamp {
		now = currentMillis()
	}

	return now
}

func currentMillis() int64 {
	return time.Now().UnixMilli()
}
