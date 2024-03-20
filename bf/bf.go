package bf

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type BloomFilter interface {
	MayContain([]byte) bool
	Save() []byte
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

func (bF *bloomFilter) safeHash(fn func([]byte) uint64) func([]byte) uint64 {
	return func(b []byte) uint64 {
		return fn(b) % uint64(len(bF.bitArray))
	}
}

func (bF *bloomFilter) setupBloomFilter(numberOfBits uint64) {
	bF.bitArray = make([]bool, numberOfBits)

	bF.hashFns = []func([]byte) uint64{
		bF.safeHash(fnv1Hash),
		bF.safeHash(fnv1aHash),
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

	// TODO: A maximum of 2^16 hash functions can be used
	bitArraySize, err := strconv.ParseUint(string(rawBytes[:8]), 10, 64)
	if err != nil {
		return nil, errors.New("Invalid bloomFilter bit array count: Could not determine the bit array count")
	}
	rawBytes = rawBytes[8:]

	filter := &bloomFilter{}
	filter.setupBloomFilter(bitArraySize)

	for i := uint64(0); i <= bitArraySize; i++ {
		dataByte := rawBytes[i]
		go filter.memCopyByte(dataByte, i*8)
	}

	return filter, nil
}

func (bF *bloomFilter) memCopyByte(rawByte byte, startIdx uint64) {
	firstBit := rawByte&0b1 == 0b1
	secondBit := (rawByte>>1)&0b1 == 0b1
	thirdBit := (rawByte>>2)&0b1 == 0b1
	fourthBit := (rawByte>>3)&0b1 == 0b1
	fifthBit := (rawByte>>4)&0b1 == 0b1
	sixthBit := (rawByte>>5)&0b1 == 0b1
	seventhBit := (rawByte>>6)&0b1 == 0b1
	eighthBit := (rawByte>>7)&0b1 == 0b1

	for idx, value := range []bool{
		firstBit,
		secondBit,
		thirdBit,
		fourthBit,
		fifthBit,
		sixthBit,
		seventhBit,
		eighthBit,
	} {

		offset := startIdx + uint64(idx)
		if offset >= uint64(len(bF.bitArray)) {
			return
		}

		bF.bitArray[offset] = value

	}
}

func (bF *bloomFilter) Save() []byte {
	rawBytes := make([]byte, 2+2+2+8+len(bF.bitArray))

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
	rawBytes = append(rawBytes, version[:]...)

	bitArrayCount := make([]byte, 8)
	actualCount := len(bF.bitArray)
	actualNumberAsString := fmt.Sprintf("%08d", actualCount)
	for idx, ch := range actualNumberAsString {
		bitArrayCount[idx] = byte(ch)
	}
	rawBytes = append(rawBytes, bitArrayCount[:]...)

	bitArrayAsByteArray := bF.constructByteArray()
	rawBytes = append(rawBytes, bitArrayAsByteArray[:]...)

	return rawBytes
}

func (bF *bloomFilter) constructByteArray() []byte {
	byteArrayLength := len(bF.bitArray) / 8
	if len(bF.bitArray)%8 != 0 {
		byteArrayLength += 1
	}
	returnBytes := make([]byte, byteArrayLength)

	for i := 0; i < byteArrayLength*8; i += 8 {

		constructedByte := 0

		for power := 0; power < 8; power++ {
			if i+power < len(bF.bitArray) {
				bit := bF.bitArray[i]
				if bit {
					constructedByte += int(math.Pow(2, float64(power)))
				}
			}

		}

		returnBytes = append(returnBytes, byte(constructedByte))
	}

	return returnBytes
}
