package main

import (
	"bufio"
	"context"
	"fmt"
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
		fmt.Println("got our entries. now we can output the result.")
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
		fmt.Println("and we are done")
		return nil
	})

	fmt.Println("serving")
	err = client.Serve()
	fmt.Println("serving done")
	if err != nil {
		log.Fatal(err)
	}
}
