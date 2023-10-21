package main

import (
	"bufio"
	"context"
	"github.com/lefinal/lwee/go-sdk/lweeclient"
	"io"
	"log"
	"strconv"
	"strings"
)

func main() {
	client := lweeclient.New(lweeclient.Options{})

	// Request input streams. // TODO: Replace with your own implementation.
	entriesInput, err := client.RequestInputStream("entries")
	if err != nil {
		log.Fatal(err)
	}

	// Request output streams. // TODO: Replace with your own implementation.
	entryCountOutput, err := client.ProvideOutputStream("entryCount")
	if err != nil {
		log.Fatal(err)
	}
	defer entryCountOutput.Close(nil)

	// Process your data. // TODO: Replace with your own implementation.
	client.Do(func(ctx context.Context) error {
		// Read and count entries.
		err = entriesInput.WaitForOpen(ctx)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(entriesInput)
		entryCount := 0
		for scanner.Scan() {
			entryCount++
		}
		err = scanner.Err()
		if err != nil {
			return err
		}
		// Write result.
		err = entryCountOutput.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(entryCountOutput, strings.NewReader(strconv.Itoa(entryCount)))
		if err != nil {
			return err
		}
		entryCountOutput.Close(nil)
		return nil
	})

	// Serve the client and run registered functions until all inputs and outputs
	// have been processed.
	err = client.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
