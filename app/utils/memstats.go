package utils

import "runtime"

func HeapAllocMb() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return float64(m.HeapAlloc) / 1024 / 1024
}
