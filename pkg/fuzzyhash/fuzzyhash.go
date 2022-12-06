package fuzzyhash

import (
	"image"
	"log"
	"os"

	"github.com/azr/phash"
)

// function to get fuzzy hash of an image
func GetFuzzyHash(filename string) uint64 {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// get image.Image from file
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	hash := phash.DTC(img)

	return hash
}
