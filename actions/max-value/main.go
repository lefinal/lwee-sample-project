package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	outputFilename := flag.String("out", "", "Output filename.")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	currentMax := -math.MaxFloat64
	var currentMaxWithTimestamp string
	var value float64
	var err error
	for scanner.Scan() {
		value, err = strconv.ParseFloat(strings.Split(scanner.Text(), ";")[1], 32)
		if err != nil {
			log.Fatal(err)
		}
		if value > currentMax {
			currentMax = value
			currentMaxWithTimestamp = scanner.Text()
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Writer result.
	err = os.WriteFile(*outputFilename, []byte(currentMaxWithTimestamp), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
