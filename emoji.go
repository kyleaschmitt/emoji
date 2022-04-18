package main

import (
	"fmt"
	"github.com/hbollon/go-edlib"
	"os"
	"strings"
)

// some happy globals
var (
	descriptions []string
	searchAlgo   edlib.Algorithm
)

func getKeys(some map[string][]string) []string {
	keys := make([]string, len(some))
	i := 0
	for k := range some {
		keys[i] = k
		i++
	}
	return keys
}

func getter(some string, algo edlib.Algorithm) []string {
	result, err := edlib.FuzzySearch(some, descriptions, algo)
	if err != nil {
		fmt.Println(err)
	} else {
	}
	return db[result]
}

func setup() {
	descriptions = getKeys(db)
	searchAlgo = edlib.Qgram // informal tests showed good results
}

func main() {
	setup()
	if len(os.Args[1:]) > 0 {
		targetEmoji := strings.Join(os.Args[1:], " ")
		results := getter(targetEmoji, searchAlgo)
		fmt.Println(results[3])
	}
}
