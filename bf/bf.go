package bf

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type BloomFilter interface {
	MayContain([]byte) bool
	Save() []byte
	Equal(BloomFilter) bool
	FilterSize() uint64
	SetIndices() []uint64
}

func (bF *bloomFilter) MayContain(data []byte) bool {
	for _, hashFn := range bF.hashFns {
		idx := hashFn(data)
		if bF.bitArray[idx] == true {
			return true
		}
	}

	return false
}

func (bF *bloomFilter) safeHash(_ string, fn func([]byte) uint64) func([]byte) uint64 {
	return func(b []byte) uint64 {
		original := fn(b)
		amended := original % uint64(len(bF.bitArray))

		return amended
	}
}

func (bF *bloomFilter) setupBloomFilter(numberOfBits uint64) {
	bF.bitArray = make([]bool, numberOfBits)

	bF.hashFns = []func([]byte) uint64{
		bF.safeHash("fnv1Hash", fnv1Hash),
		bF.safeHash("fnv1aHash", fnv1aHash),
	}
}

func Load(rawBytes []byte) (BloomFilter, error) {
	// Check file header
	if len(rawBytes) < 2 || (rawBytes[0] != 'B' && rawBytes[1] != 'F') {
		return nil, errors.New("Invalid bloomFilter")
	}
	rawBytes = rawBytes[2:]

	if len(rawBytes) < 2 {
		return nil, errors.New("Invalid bloomFilter version: No version found")
	}

	// TODO: A version number can be added here for backwards compatibility
	_, err := strconv.ParseInt(string(rawBytes[:2]), 10, 16)
	if err != nil {
		return nil, errors.New("Invalid bloomFilter version")
	}
	rawBytes = rawBytes[2:]

	if len(rawBytes) < 2 {
		return nil, errors.New("Invalid bloomFilter hash function count: No hash functions found")
	}

	// TODO: A maximum of 2^16 hash functions can be used
	_, err = strconv.ParseInt(string(rawBytes[:2]), 10, 16)
	if err != nil {
		return nil, errors.New("Invalid bloomFilter hash function count: Could not determine the hash function count")
	}
	rawBytes = rawBytes[2:]

	if len(rawBytes) < 8 {
		return nil, errors.New("Invalid bloomFilter bit array count: No bit array count found")
	}

	// TODO: A maximum of 2^64 bit array size is permitted
	bitArraySize, err := strconv.ParseUint(string(rawBytes[:8]), 10, 64)
	if err != nil {
		return nil, errors.New("Invalid bloomFilter bit array count: Could not determine the bit array count")
	}
	rawBytes = rawBytes[8:]

	// For the carriage return
	// Note: We do not need this since the scanner merely skips it for us, see the SplitFunc of the BFScanner
	// rawBytes = rawBytes[2:]

	filter := &bloomFilter{}
	filter.setupBloomFilter(bitArraySize)
	filter.bitArray = bytesToBools(rawBytes)[:bitArraySize]

	return filter, nil
}

func bytesToBools(b []byte) []bool {
	t := make([]bool, 8*len(b))
	for i, x := range b {
		for j := 0; j < 8; j++ {
			if (x<<uint(j))&0x80 == 0x80 {
				t[8*i+j] = true
			}
		}
	}
	return t
}

func (bF *bloomFilter) Save() []byte {
	rawBytes := make([]byte, 0)

	header := make([]byte, 2)
	header[0] = 'B'
	header[1] = 'F'
	rawBytes = append(rawBytes, header[:]...)

	version := make([]byte, 2)
	version[0] = '0'
	version[1] = '1'
	rawBytes = append(rawBytes, version[:]...)

	hashFnCount := make([]byte, 2)
	hashFnCount[0] = '0'
	hashFnCount[1] = '2'
	rawBytes = append(rawBytes, hashFnCount[:]...)

	bitArrayCount := make([]byte, 8)
	actualNumberAsString := fmt.Sprintf("%08d", len(bF.bitArray))
	for idx, ch := range actualNumberAsString {
		bitArrayCount[idx] = byte(ch)
	}
	rawBytes = append(rawBytes, bitArrayCount[:]...)

	carriageReturn := make([]byte, 2)
	carriageReturn[0] = '\r'
	carriageReturn[1] = '\n'
	rawBytes = append(rawBytes, carriageReturn[:]...)

	bitArrayAsByteArray := boolsToBytes(bF.bitArray)
	rawBytes = append(rawBytes, bitArrayAsByteArray[:]...)

	return rawBytes
}

func boolsToBytes(t []bool) []byte {
	b := make([]byte, (len(t)+7)/8)
	for i, x := range t {
		if x {
			b[i/8] |= 0x80 >> uint(i%8)
		}
	}
	return b
}

func (bF *bloomFilter) FilterSize() uint64 {
	return uint64(len(bF.bitArray))
}

func (bF *bloomFilter) SetIndices() []uint64 {
	indicies := make([]uint64, 0)

	for idx, val := range bF.bitArray {
		if val {
			indicies = append(indicies, uint64(idx))
		}
	}

	return indicies
}

func (bF *bloomFilter) Equal(otherBF BloomFilter) bool {
	return otherBF.FilterSize() == bF.FilterSize() && reflect.DeepEqual(bF.SetIndices(), otherBF.SetIndices())
}
