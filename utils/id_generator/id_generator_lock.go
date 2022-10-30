package idgenerator

import (
	"errors"
	"sync"
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

var (
	nodeID        = PrivateIPToNodeID()
	sequence      uint64
	m             sync.Mutex
	lastTimestamp int64 = -1
)

func NextID() (uint64, error) {
	m.Lock()
	defer m.Unlock()
	current := currentMillis()

	if current < lastTimestamp {
		return 0, errors.New("system clock moving backward. Refuse to generate new id")
	}

	if current == lastTimestamp {
		sequence = (sequence + 1) & MaxSequence
		if sequence == 0 {
			current = tillNextMillis(current)
		}
	} else {
		sequence = 0
	}

	lastTimestamp = current

	return uint64(current)<<(NodeIDLength+SequenceLength) |
		((uint64(nodeID)) << SequenceLength) |
		sequence, nil
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
