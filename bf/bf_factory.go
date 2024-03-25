package bf

import (
	"fmt"
	"math"
)

type BloomFilterFactory interface {
	RunOver([][]byte)
	GetBloomFilter() BloomFilter
}

type bloomFilter struct {
	bitArray []bool
	hashFns  []func([]byte) uint64
}

func (bF *bloomFilter) String() string {
	setBits := setCount(bF.bitArray)
	return fmt.Sprintf("bloomFilter[bitArray=(length=%d, setbitcount=%d[%v])]", len(bF.bitArray), len(setBits), setBits)
}

type bloomFilterFactory struct {
	bloomFilter
	probability   float64
	numberOfItems uint64
}

func (bF *bloomFilterFactory) String() string {
	return fmt.Sprintf("bloomFilterFactory[probability=%v, numberOfItems=%v, bloomFilter=(%v)]", bF.probability, bF.numberOfItems, bF.GetBloomFilter())
}

func setCount(bitArray []bool) (indices []uint64) {
	setBits := make([]uint64, 0)
	for idx, bit := range bitArray {
		if bit == true {
			setBits = append(setBits, uint64(idx))
		}
	}
	return setBits
}

func (bF *bloomFilterFactory) setup() {
	numberOfBits := uint64(float64(bF.numberOfItems) * math.Log(bF.probability) / math.Pow(math.Log(2), 2))
	bF.setupBloomFilter(numberOfBits)

}

func (bF *bloomFilterFactory) applyObject(data []byte) {
	for _, hashFn := range bF.hashFns {
		offset := hashFn(data)
		bF.bitArray[offset] = true
	}
}

func New(probability float64, numberOfitems uint64) BloomFilterFactory {
	bf := &bloomFilterFactory{
		probability:   probability,
		numberOfItems: numberOfitems,
	}
	bf.setup()
	return bf
}

func (bF *bloomFilterFactory) RunOver(dataset [][]byte) {
	for _, obj := range dataset {
		bF.applyObject(obj)
	}
}

func (bF *bloomFilterFactory) GetBloomFilter() BloomFilter {
	return &bloomFilter{
		bitArray: bF.bitArray,
		hashFns:  bF.hashFns,
	}
}
