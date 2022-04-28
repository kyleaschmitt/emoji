package main

import (
	"encoding/base64"
	"fmt"
	"github.com/hbollon/go-edlib"
	"github.com/mattn/go-sixel"
	"image"
	_ "image/png"
	_ "image/gif" // Yes.  Everything in full-emoji-list.html says it's a png. They lie
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

func getter(some string, algo edlib.Algorithm) (string, []string) {
	result, err := edlib.FuzzySearch(some, descriptions, algo)
	if err != nil {
		fmt.Println(err)
	} else {
	}
	return result, db[result]
}

func setup() {
	descriptions = getKeys(db)
	searchAlgo = edlib.Qgram // informal tests showed good results
}

func makeImage(data string) image.Image {
	reader := base64.NewDecoder(base64.StdEncoding,strings.NewReader(data))
	img, format, err := image.Decode(reader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Decoded image of format: %s\n",format)
	return img
}

func main() {
	setup()
	if len(os.Args[1:]) > 0 {
		targetEmoji := strings.Join(os.Args[1:], " ")
		realDescrip, results := getter(targetEmoji, searchAlgo)
		//fmt.Println(results[1])
		fmt.Println(len(results))
		fmt.Printf("%s, codepoint %s, description %s\n",results[1],results[0],realDescrip)
		for i:=2;i<len(results);i++ {
			fmt.Printf("len(results[%d]) = %d\n",i,len(results[i]))
			sixel.NewEncoder(os.Stdout).Encode(makeImage(results[i]))
		}
	}
	//this abomination of a version has the key of description, codepoint, glyph, and
	//the rest are the raw image data
	//"grinning face":                           []string{"U+1F600", "😀"
}
/*
const data = `
R0lGODlhDwANALMPAK6slY6OfmWz/vz8/XFxdKwAAENWjqFRTO9lTYjE/jiR5lkcRzk3aP43GTs7OwAAACH5BAEAAA8ALAAAAAAPAA0AAARK8Mn5wKmX6lYe1xpyiaAEnOgJBkRLsK1DAUFt2zLzBmlPSIzEYEgsDiaGA2LJZF4SUMNiSqUyHoKEQKDoer+KrXhMFjyg6HT6EQEAOw==
`*/
