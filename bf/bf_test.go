package bf_test

import (
	"testing"

	"github.com/saphal1998/bloomf/bf"
)

func generateBloomFilter(items [][]byte, probability float64) bf.BloomFilter {
	bloomFilterFactory := bf.New(probability, len(items))
	bloomFilterFactory.RunOver(items)

	return bloomFilterFactory.GetBloomFilter()
}

func TestBloomFilterSerialization(t *testing.T) {
	items := []string{
		"saphal",
	}

	itemsToBytes := make([][]byte, 0)

	for _, item := range items {
		itemByte := []byte(item)
		itemsToBytes = append(itemsToBytes, itemByte)
	}

	bloomFilter := generateBloomFilter(itemsToBytes, 99)
	serialized := bloomFilter.Save()

	bloomFilterBytes := []byte{
		byte(0x0A),
		byte(0x00),
	}

	if len(serialized) != len(bloomFilterBytes) {
		t.Fatalf("expected serialized bloom filter to have %v bytes, but it had %v bytes", len(bloomFilterBytes), len(serialized))
	}

	for idx, databyte := range bloomFilterBytes {
		if serialized[idx] != databyte {
			t.Fatalf("unexpected byte at idx[%d], expected=%b, got=%b\n", idx, databyte, serialized[idx])
		}
	}
}
