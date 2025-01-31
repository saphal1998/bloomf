package bf

import (
	"github.com/saphal1998/bitarray"
)

type BloomFilter interface {
	Buffer() []byte
	Check([]byte) bool
}

type bloomFilter struct {
	bitArray bitarray.BitArray
	hashFns  []func([]byte) int
}

func safeHash(_ string, size int, fn func([]byte) uint64) func([]byte) int {
	return func(b []byte) int {
		original := fn(b)
		amended := original % uint64(size)
		return int(amended)
	}
}

func (filter *bloomFilter) Buffer() []byte {
	buffer := filter.bitArray.UnsafeRawBuffer()
	return buffer
}

func New(probability float64, numberOfItems int) BloomFilter {
	return &bloomFilter{
		bitArray: bitarray.NewBitArray(numberOfBits),
		hashFns: []func([]byte) int{
			safeHash("fnv1Hash", numberOfItems, fnv1Hash),
			safeHash("fnv1aHash", numberOfItems, fnv1aHash),
		},
	}
}

func (filter *bloomFilter) applyObject(data []byte) {
	for _, hashFn := range filter.hashFns {
		offset := hashFn(data)
		filter.bitArray.Set(offset)
	}
}
func Apply(filter bloomFilter, dataset [][]byte) {
	for _, obj := range dataset {
		filter.applyObject(obj)
	}
}
