package main

import (
	"bloomfilter/bf"
	"bufio"
	"flag"
	"fmt"
	"os"
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

		bloomFilter := bf.New(probability, uint64(len(dataset)))
		fmt.Printf("Created %v\n", bloomFilter)
		bloomFilter.RunOver(dataset)
		fmt.Printf("After runover - %v\n", bloomFilter)

	}

	if appMode == load {
		fmt.Fprintln(os.Stderr, "Have not implemented load (yet)")
	}

}
