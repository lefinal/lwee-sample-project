package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	outputFilename := flag.String("out", "", "Output filename.")
	flag.Parse()

	f, err := os.OpenFile(*outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(os.Stdin)
	var timestamp time.Time
	var currentYear int
	var currentMax float64
	value := -math.MaxFloat64
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic while handling scanner line", scanner.Text())
			panic(r)
		}
	}()
	for scanner.Scan() {
		if strings.Count(scanner.Text(), ";") > 2 {
			log.Fatalf("that's not right: %s", scanner.Text())
		}
		value, err = strconv.ParseFloat(strings.Split(scanner.Text(), ";")[1], 32)
		if err != nil {
			log.Fatal(err)
		}
		timestamp, err = time.Parse(time.RFC3339, strings.Split(scanner.Text(), ";")[0])
		if err != nil {
			log.Fatal(err)
		}
		if timestamp.Year() == 0 {
			log.Fatalf("why is year 0 of %s", scanner.Text())
		}
		if timestamp.Year() != currentYear && currentYear != 0 {
			printResult(f, currentYear, currentMax)
			currentMax = -math.MaxFloat64
		}
		currentYear = timestamp.Year()
		currentMax = max(value, currentMax)
	}

	printResult(f, currentYear, currentMax)
}

func printResult(w io.Writer, year int, max float64) {
	_, err := fmt.Fprintf(w, "%d;%f\n", year, max)
	if err != nil {
		log.Fatal(err)
	}
}
