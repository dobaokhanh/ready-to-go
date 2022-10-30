package idgenerator

import "testing"

func BenchmarkTestIdGeneratorLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextID()
	}
}
