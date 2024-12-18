package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"fuzzycap/pkg/fuzzyhash"
	"fuzzycap/pkg/screenshot"

	"github.com/alexflint/go-arg"
)

type args struct {
	Input            string `arg:"-i,--input" help:"input file with list of URLs" required:"true"`
	Output           string `arg:"--out" help:"output JSON file" default:"fuzzyhash-results.json"`
	Quality          int    `arg:"--quality" help:"screenshot quality" default:"100"`
	HammingThreshold int    `arg:"--threshold" help:"Hamming distance threshold for considering images different" default:"5"`
}

type Snapshot struct {
	Time     string `json:"time"`
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
}

type SiteHistory struct {
	Url       string     `json:"url"`
	Snapshots []Snapshot `json:"snapshots"`
}

func main() {
	var argVals args
	arg.MustParse(&argVals)

	urls := readUrls(argVals.Input)
	urls = normalizeUrls(urls)
	createDirIfNotExist("screenshots")

	history := loadHistory(argVals.Output)

	historyMap := make(map[string]*SiteHistory)
	for i, h := range history {
		historyMap[h.Url] = &history[i]
	}

	for _, url := range urls {
		filename, err := screenshot.Screenshot(url, argVals.Quality)
		fatalIfErr(err, "failed to take screenshot")

		hash := fuzzyhash.GetFuzzyHash(filename)
		hashString := fmt.Sprintf("%d", hash)

		siteHist, exists := historyMap[url]
		if !exists {
			// New URL
			newEntry := SiteHistory{
				Url: url,
				Snapshots: []Snapshot{
					{
						Time:     time.Now().Format(time.RFC3339),
						Filename: filename,
						Hash:     hashString,
					},
				},
			}
			history = append(history, newEntry)
			historyMap[url] = &history[len(history)-1]
			fmt.Printf("[NEW] %s (hash: %s)\n", url, hashString)
			continue
		}

		// Compare with the most recent snapshot
		lastSnapshot := siteHist.Snapshots[len(siteHist.Snapshots)-1]

		oldHashVal, err := strconv.ParseUint(lastSnapshot.Hash, 10, 64)
		fatalIfErr(err, "failed to parse old hash")
		newHashVal, err := strconv.ParseUint(hashString, 10, 64)
		fatalIfErr(err, "failed to parse new hash")

		dist := hammingDistance(oldHashVal, newHashVal)
		if dist > argVals.HammingThreshold {
			fmt.Printf("[CHANGED] %s: old hash=%s, new hash=%s (Hamming distance: %d)\n",
				url, lastSnapshot.Hash, hashString, dist)
		} else {
			fmt.Printf("[NO SIGNIFICANT CHANGE] %s (Hamming distance: %d)\n", url, dist)
		}

		newSnap := Snapshot{
			Time:     time.Now().Format(time.RFC3339),
			Filename: filename,
			Hash:     hashString,
		}
		siteHist.Snapshots = append(siteHist.Snapshots, newSnap)
	}

	saveHistory(history, argVals.Output)
	fmt.Println("Results saved to", argVals.Output)
}

func readUrls(filename string) []string {
	data, err := ioutil.ReadFile(filename)
	fatalIfErr(err, "failed to read input file")
	text := string(data)
	lines := strings.Split(text, "\n")
	return lines
}

func normalizeUrls(urls []string) []string {
	var fullUrls []string
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
			u = "https://" + u
		}
		fullUrls = append(fullUrls, u)
	}
	return fullUrls
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		fatalIfErr(err, "failed to create directory "+dir)
	}
}

func loadHistory(file string) []SiteHistory {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return []SiteHistory{}
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Could not read previous file: %v", err)
		return []SiteHistory{}
	}
	var results []SiteHistory
	err = json.Unmarshal(data, &results)
	if err != nil {
		log.Printf("Could not parse previous results: %v", err)
		return []SiteHistory{}
	}
	return results
}

func saveHistory(results []SiteHistory, outFile string) {
	data, err := json.MarshalIndent(results, "", "  ")
	fatalIfErr(err, "failed to marshal results to JSON")

	err = ioutil.WriteFile(outFile, data, 0644)
	fatalIfErr(err, "failed to write output file")
}

func fatalIfErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

// hammingDistance calculates the Hamming distance between two 64-bit integers.
// It counts the number of differing bits between the two.
func hammingDistance(a, b uint64) int {
	x := a ^ b
	dist := 0
	for x != 0 {
		dist++
		x = x & (x - 1) // remove the lowest set bit
	}
	return dist
}
