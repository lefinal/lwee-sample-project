package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	humidityFilename := flag.String("humidity-out", "", "Humidity output file.")
	flag.Parse()

	if *humidityFilename == "" {
		log.Fatal("missing humidity out flag")
	}
	humidityFile, err := os.OpenFile(*humidityFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = humidityFile.Close() }()
	bufferedStdout := bufio.NewWriterSize(os.Stdout, 1024*1024)
	bufferedHumidityFile := bufio.NewWriterSize(humidityFile, 1024*1024)
	extractAndPrint(os.Stdin, bufferedStdout, bufferedHumidityFile)
	_ = bufferedStdout.Flush()
	_ = bufferedHumidityFile.Flush()
}

func extractAndPrint(r io.Reader, temperatureWriter io.Writer, humidityWriter io.Writer) {
	scanner := bufio.NewScanner(r)
	first := true

	var timestampStr, temperatureStr, humidityStr string
	var timestamp, lastTimestamp time.Time
	var year, month, day, hour int
	var temperature, humidity float64
	var err error
	for scanner.Scan() {
		// Skip first.
		if first {
			first = false
			continue
		}
		// Format of a line:        '5155;1948010110;    5;  -1.6;  88.0;eor'
		segments := strings.Split(scanner.Text(), ";")
		// Extract timestamp.
		timestampStr = segments[1]
		timestampStr = strings.TrimSpace(timestampStr)
		year, err = strconv.Atoi(timestampStr[:4])
		if err != nil {
			log.Fatal(err)
		}
		month, err = strconv.Atoi(timestampStr[4:6])
		if err != nil {
			log.Fatal(err)
		}
		day, err = strconv.Atoi(timestampStr[6:8])
		if err != nil {
			log.Fatal(err)
		}
		hour, err = strconv.Atoi(timestampStr[8:10])
		if err != nil {
			log.Fatal(err)
		}
		timestamp = time.Time{}
		timestamp = timestamp.AddDate(year, month, day)
		timestamp = timestamp.Add(time.Duration(hour) * time.Hour)
		timestamp = timestamp.UTC()
		// Extract temperature.
		temperatureStr = strings.TrimSpace(segments[3])
		temperature, err = strconv.ParseFloat(temperatureStr, 64)
		if err != nil {
			log.Fatal(err)
		}
		// Extract humidity.
		humidityStr = strings.TrimSpace(segments[4])
		humidity, err = strconv.ParseFloat(humidityStr, 64)
		if err != nil {
			log.Fatal(err)
		}
		// Assure larger than previous timestamp in order to allow higher efficiency in
		// the steps after.
		if timestamp.Before(lastTimestamp) {
			// Some timestamps are not ordered because of leap years. We skip this one.
			continue
		}
		lastTimestamp = timestamp
		_, _ = fmt.Fprintf(temperatureWriter, "%s;%f\n", timestamp.Format(time.RFC3339), temperature)
		_, _ = fmt.Fprintf(humidityWriter, "%s;%f\n", timestamp.Format(time.RFC3339), humidity)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
