package fuzzycap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jonathanwalker/fuzzycap/pkg/fuzzyhash"
	"github.com/jonathanwalker/fuzzycap/pkg/screenshot"

	"github.com/alexflint/go-arg"
)

var args struct {
	Input string `arg:"-i,--input" help:"input file with list of urls"`
}

// create struct with url, filename, and hash
type outputJson struct {
	Url      string
	Filename string
	Hash     string
}

func main() {
	arg.MustParse(&args)

	urls := readUrls(args.Input)

	var fullurls []string
	for _, url := range urls {
		if url != "" {
			if strings.HasPrefix(url, "http") {
				fullurls = append(fullurls, url)
			} else {
				fullurls = append(fullurls, "https://"+url)
			}
		}
	}

	// loop through urls and get fuzzy hash of each screenshot
	var outputData []outputJson
	for _, url := range fullurls {
		filename := screenshot.Screenshot(url, 100)
		hash := fuzzyhash.GetFuzzyHash(filename)
		// convert hash to string
		hashString := fmt.Sprintf("%d", hash)
		outputData = append(outputData, outputJson{url, filename, hashString})
	}

	// output the data as json
	outputJson, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(outputJson))
}

// function to read a file and return urls in a list
func readUrls(filename string) []string {
	// Open a txt file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the data to a string
	text := string(data)

	// Split the string into a list of urls
	urls := strings.Split(text, "\n")

	return urls
}
