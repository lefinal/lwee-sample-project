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
	entriesInput, err := client.RequestInputStream("entries")
	if err != nil {
		log.Fatal(err)
	}
	entryCountOutput, err := client.ProvideOutputStream("entryCount")
	if err != nil {
		log.Fatal(err)
	}
	defer entryCountOutput.Close(nil)
	go func() {
		err := client.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// Read and count entries.
	err = entriesInput.WaitForOpen(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(entriesInput)
	entryCount := 0
	for scanner.Scan() {
		entryCount++
	}
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	// Write result.
	err = entryCountOutput.Open()
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(entryCountOutput, strings.NewReader(strconv.Itoa(entryCount)))
	if err != nil {
		log.Fatal(err)
	}
	entryCountOutput.Close(nil)
	<-client.Lifetime().Done()
}
