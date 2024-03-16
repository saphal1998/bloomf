package bf

import (
	"fmt"
	"math"
)

type BloomFilter interface {
	MayContain([]byte) bool
}

type BloomFilterFactory interface {
	RunOver([][]byte)
}

type bloomFilter struct {
	probability   float64
	numberOfItems uint64
	bitArray      []bool
	hashFns       []func([]byte) uint64
}

func (bF *bloomFilter) String() string {
	return fmt.Sprintf("bloomFilter[probability=%v, numberOfItems=%v, bitArray=(length=%d, setbitcount=%d)]", bF.probability, bF.numberOfItems, len(bF.bitArray), bF.setCount())
}

func (bF *bloomFilter) setCount() uint64 {
	count := uint64(0)
	for _, bit := range bF.bitArray {
		if bit == true {
			count += 1
		}
	}
	return count
}

func (bF *bloomFilter) MayContain([]byte) bool {
	return false
}

func (bF *bloomFilter) setup() {
	numberofBits := uint64(float64(bF.numberOfItems) * math.Log(bF.probability) / math.Pow(math.Log(2), 2))
	bF.bitArray = make([]bool, numberofBits)

	bF.hashFns = []func([]byte) uint64{
		bF.safeHash(fnv1Hash),
		bF.safeHash(fnv1aHash),
	}
}

func (bF *bloomFilter) safeHash(fn func([]byte) uint64) func([]byte) uint64 {
	return func(b []byte) uint64 {
		return fn(b) % uint64(len(bF.bitArray))
	}
}

func (bF *bloomFilter) applyObject(data []byte) {
	for _, hashFn := range bF.hashFns {
		offset := hashFn(data)
		bF.bitArray[offset] = true
	}
}

func New(probability float64, numberOfitems uint64) BloomFilterFactory {
	bf := &bloomFilter{
		probability:   probability,
		numberOfItems: numberOfitems,
	}
	bf.setup()
	return bf
}

func (bF *bloomFilter) RunOver(dataset [][]byte) {
	for _, obj := range dataset {
		bF.applyObject(obj)
	}
}
