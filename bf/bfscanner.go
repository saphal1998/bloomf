package bf

import (
	"bufio"
	"bytes"
	"os"
)

func BFScanner(f *os.File) *bufio.Scanner {
	scanner := bufio.NewScanner(f)

	// The first line needs to be parsed, after that everything is one chunk
	statLineDone := false

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			// If at end of file and no more data is left, return 0 to indicate no more tokens
			return 0, nil, nil
		}

		// Search for '\r\n' in the data
		if i := bytes.Index(data, []byte{'\r', '\n'}); i >= 0 && !statLineDone {
			statLineDone = true
			return i + 2, data[0:i], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		if len(data) >= bufio.MaxScanTokenSize {
			return bufio.MaxScanTokenSize, data[:bufio.MaxScanTokenSize], nil
		}

		return 0, nil, nil
	})
	return scanner

}
