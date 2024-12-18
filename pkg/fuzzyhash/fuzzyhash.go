package fuzzyhash

import (
	"image"
	"log"
	"os"

	"github.com/azr/phash"
)

// GetFuzzyHash computes a pHash-based fuzzy hash of an image file.
func GetFuzzyHash(filename string) uint64 {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open image file %s: %v", filename, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("failed to decode image %s: %v", filename, err)
	}

	hash := phash.DTC(img)
	return hash
}
