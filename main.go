package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const maxGoroutines = 10000

type Input struct {
	A int `json:"a"`
	B int `json:"b"`
}

func worker(inputs []Input, wg *sync.WaitGroup, sumChan chan int) {
	defer wg.Done()
	sum := 0
	for _, input := range inputs {
		sum += input.A + input.B
	}
	sumChan <- sum
}

func findJSONFile(dir string) (string, error) {
	var jsonFile string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			jsonFile = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if jsonFile == "" {
		return "", fmt.Errorf("no JSON file found")
	}

	return jsonFile, nil
}


func readJSONFile(jsonFile string) ([]Input, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	t, err := decoder.Token()
	if err != nil || t != json.Delim('[') {
		return nil, fmt.Errorf("invalid JSON file")
	}

	var numbers []Input
	for decoder.More() {
		var input Input
		err := decoder.Decode(&input)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, input)
	}
	t, err = decoder.Token()
	if err != nil || t != json.Delim(']') {
		return nil, fmt.Errorf("invalid JSON file")
	}

	if len(numbers) == 0 {
		return nil, fmt.Errorf("no data found in JSON file")
	}
	return numbers, nil
}

func calcSum(numbers []Input, workers int) (int, error) {
	chunkSize := (len(numbers) + workers - 1) / workers
	sumChan := make(chan int, workers)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(numbers) {
			end = len(numbers)
		}
		if start >= end {
			break
		}

		wg.Add(1)
		go worker(numbers[start:end], &wg, sumChan)
	}

	wg.Wait()
	close(sumChan)

	totalSum := 0
	for sum := range sumChan {
		totalSum += sum
	}
	return totalSum, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <dir> <numWorkers>")
		return
	}

	jsonDir := os.Args[1]
	numWorkers, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if numWorkers > maxGoroutines {
		fmt.Printf("Error: numWorkers should be less than %d\n", maxGoroutines)
		return
	} else if numWorkers <= 0 {
		fmt.Printf("Error: numWorkers should be greater than 0\n")
		return
	}

	jsonFile, err := findJSONFile(jsonDir)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	numbers, err := readJSONFile(jsonFile)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return
	}

	startTime := time.Now()
	totalSum, err := calcSum(numbers, numWorkers)

	if err != nil {
		fmt.Printf("Error calculating sum: %v\n", err)
		return

	}
	endTime := time.Now()

	duration := endTime.Sub(startTime)

	fmt.Printf("Total sum: %d\n", totalSum)
	fmt.Printf("Time taken: %v\n", duration)
}
