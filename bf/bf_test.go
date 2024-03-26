package bf_test

import (
	"bloomfilter/bf"
	"os"
	"testing"
)

func generateBloomFilter(items [][]byte, probability float64) bf.BloomFilter {
	bloomFilterFactory := bf.New(probability, uint64(len(items)))
	bloomFilterFactory.RunOver(items)

	return bloomFilterFactory.GetBloomFilter()
}

func TestPresenceOfByteStream(t *testing.T) {
	items := []string{
		"saphal",
		"patro",
		"engineer",
		"student",
		"software",
	}

	itemsToBytes := make([][]byte, 0)

	for _, item := range items {
		itemByte := []byte(item)
		itemsToBytes = append(itemsToBytes, itemByte)
	}

	bloomFilter := generateBloomFilter(itemsToBytes, 99)

	t.Logf("bloomFilter %v\n", bloomFilter)

	for idx, databytes := range itemsToBytes {
		if bloomFilter.MayContain(databytes) != true {
			t.Fatalf("expected %s to be in the bloom filter, but it was not", items[idx])
		}
	}
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

	t.Logf("bloomFilter %v\n", bloomFilter)

	serialized := bloomFilter.Save()

	bloomFilterBytes := []byte{
		// File type
		byte('B'),
		byte('F'),
		// Version
		byte('0'),
		byte('1'),
		// # of hash functions
		byte('0'),
		byte('2'),
		// Length of bit array
		byte('0'),
		byte('0'),
		byte('0'),
		byte('0'),
		byte('0'),
		byte('0'),
		byte('0'),
		byte('9'),
		// Carriage Return
		byte('\r'),
		byte('\n'),
		byte(0x0A),
		byte(0x00),
	}

	if len(serialized) != len(bloomFilterBytes) {
		t.Fatalf("expected serialized bloom filter to have %v bytes, but it had %v bytes", len(bloomFilterBytes), len(serialized))
	}

	for idx, databyte := range bloomFilterBytes {
		if serialized[idx] != databyte {
			t.Fatalf("unexpected byte at idx[%d], expected=%s, got=%s\n", idx, string(databyte), string(serialized[idx]))
		}
	}
}

func TestBloomFilterDeserialization(t *testing.T) {
	items := []string{
		"saphal",
	}

	itemsToBytes := make([][]byte, 0)

	for _, item := range items {
		itemByte := []byte(item)
		itemsToBytes = append(itemsToBytes, itemByte)
	}

	originalBloomFilter := generateBloomFilter(itemsToBytes, 99)

	t.Logf("bloomFilter %v\n", originalBloomFilter)

	serialized := originalBloomFilter.Save()

	f, err := os.CreateTemp("", "tmpfile-*.bf")
	if err != nil {
		t.Fatalf("Cannot perform test, could not create a temporary file to test deserialization")
	}
	defer f.Close()
	defer os.Remove(f.Name())

	if _, err := f.Write(serialized); err != nil {
		t.Fatalf("Could not write the temporary bloom filter file")
	}

	f.Close()

	f, err = os.Open(f.Name())
	if err != nil {
		t.Fatalf("Cannot perform test, could not open the temporary file to test deserialization")
	}

	oldStdin := os.Stdin
	defer func() {
		os.Stdin = oldStdin
	}()

	os.Stdin = f

	scanner := bf.BFScanner(os.Stdin)
	bloomFilterBytes := make([]byte, 0)
	for scanner.Scan() {
		databytes := scanner.Bytes()
		bloomFilterBytes = append(bloomFilterBytes, databytes...)
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Cannot perform test, could not scan the temporary file correctly")
	}

	reconstructedBloomFilter, err := bf.Load(bloomFilterBytes)
	if err != nil {
		t.Fatalf("Could not reconstruct bloom filter")
	}

	if !reconstructedBloomFilter.Equal(originalBloomFilter) {
		t.Fatalf("Failed serialization, expected=%v, got=%v", originalBloomFilter, reconstructedBloomFilter)
	}
}

func TestBFScanner(t *testing.T) {
	data := "abcd\r\nefg"

	f, err := os.CreateTemp("", "tmpfile-*.bf")
	if err != nil {
		t.Fatalf("Cannot perform test, could not create a temporary file to test deserialization")
	}
	defer f.Close()
	defer os.Remove(f.Name())

	if _, err := f.Write([]byte(data)); err != nil {
		t.Fatalf("Could not write the temporary bloom filter file")
	}

	f.Close()

	f, err = os.Open(f.Name())
	if err != nil {
		t.Fatalf("Cannot perform test, could not open the temporary file to test deserialization")
	}

	scanner := bf.BFScanner(f)
	var dataBytes []byte

	for scanner.Scan() {
		fBytes := scanner.Bytes()
		dataBytes = append(dataBytes, fBytes...)
	}

	expectedResult := "abcdefg"

	if string(dataBytes) != expectedResult {
		t.Fatalf("expected=%v, got=%v", expectedResult, string(dataBytes))
	}

}
