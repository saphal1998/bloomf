package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/saphal1998/bloomf/bf"
)

type applicationMode string

const (
	create applicationMode = "create"
	load   applicationMode = "load"
)

var mode *string
var probability float64

func init() {
	mode = flag.String("mode", "create", "Mode")
	probability = *flag.Float64("probability", 95, "Probability of containing in the filter")
}

func main() {
	flag.Parse()

	appMode := applicationMode(*mode)

	fmt.Printf("App started in %v mode\n", appMode)

	if appMode != create && appMode != load {
		fmt.Fprintln(os.Stderr, "Usage: go run main.go -mode create|load < source.(txt|bf)")
		os.Exit(1)
	}

	if appMode == create {
		scanner := bufio.NewScanner(os.Stdin)
		var dataset [][]byte
		for {
			scanner.Scan()
			databyteLine := scanner.Bytes()
			if len(databyteLine) == 0 {
				break
			}
			dataset = append(dataset, databyteLine)
		}

		err := scanner.Err()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not read datasource from stdin")
			os.Exit(1)
		}

		bloomFilterFactory := bf.New(probability, uint64(len(dataset)))
		bloomFilterFactory.RunOver(dataset)

		bloomFilter := bloomFilterFactory.GetBloomFilter()

		// This is done to validate
		for _, val := range dataset {
			if bloomFilter.WillNotContain(val) != false {
				fmt.Fprintf(os.Stderr, "Something went wrong with our bloom filter, it was trained on %v, but it says that it will never contain it\n", string(val))
			}
		}

		bytesSaved := bloomFilter.Save()

		fmt.Printf("Size of filter: %d, %v\n", bloomFilter.FilterSize(), bloomFilter)

		if err := os.WriteFile("filter.txt", bytesSaved, 0666); err != nil {
			fmt.Fprintln(os.Stderr, "Could not save the bloom filter")
			os.Exit(1)
		}

	}

	if appMode == load {
		scanner := bf.BFScanner(os.Stdin)
		var bloomFilterBytes []byte
		for {
			scanner.Scan()
			databyteLine := scanner.Bytes()
			if len(databyteLine) == 0 {
				break
			}
			bloomFilterBytes = append(bloomFilterBytes, databyteLine[:]...)
		}

		err := scanner.Err()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Printf("Size of filter read: %d bytes\n", len(bloomFilterBytes))

		bloomFilter, err := bf.Load(bloomFilterBytes)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Printf("Size of filter: %d, %v\n", bloomFilter.FilterSize(), bloomFilter)
	}

}
