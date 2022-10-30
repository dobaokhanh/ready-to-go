package idgeneratorlockfree

import "testing"

func BenchmarkTestIdGeneratorLockFree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextID()
	}
}
