package input

import (
	"bufio"
	"fmt"
	"os"
)

func ReadFromStdin() ([]byte, error) {
	var data []byte
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data = append(data, '\n')
		data = append(data, scanner.Bytes()...)
	}
	return data[1:], nil
}

func ReadFromFile(path string) ([]byte, error) {
	return nil, fmt.Errorf("Not implemented")
}
