package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/hbollon/go-edlib"
	"github.com/mattn/go-sixel"
	"golang.org/x/image/draw"
	"image"
	_ "image/gif" // Yes.  Everything in full-emoji-list.html says it's a png. They lie
	_ "image/png"
	"os"
	"strconv"
	"strings"
)

// some happy globals
var (
	descriptions []string
	searchAlgo   edlib.Algorithm
	xSize        int
	ySize        int
	dbMap        map[string]int
	//flags
	charFlag        bool
	codeFlag        bool
	debugFlag       bool
	detailFlag      bool
	helpFlag        bool
	listFlag        bool
	renderFlag      bool
	renderWhichFlag string
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

func ratios() {
	for k, _ := range db {
		ratiosOf(k)
		/*for i:=2;i<len(v);i++ {
			bound:=makeImage(v[i]).Bounds()
			ratio:=float64(bound.Max.X)/float64(bound.Max.Y)
			if ratio!=1.0 {
				fmt.Printf("%s image[%d] has a ratio of %f\n",k,i,ratio)
			}
		}*/
	}
}

func ratiosOf(target string) {
	v := db[target]
	for i := 2; i < len(v); i++ {
		bound := makeImage(v[i]).Bounds()
		targetX := bound.Max.X
		targetY := bound.Max.Y
		ratio := float64(bound.Max.X) / float64(bound.Max.Y)
		if ratio != 1.0 {
			fmt.Printf("%s image[%d] is %dx%d pixels, and has a ratio of %f\n", target, i, targetX, targetY, ratio)
		}
	}
}

func setup() {
	descriptions = getKeys(db)
	searchAlgo = edlib.Qgram // informal tests showed good results
	flag.BoolVar(&helpFlag, "help", false, "displays help text")
	flag.BoolVar(&debugFlag, "debug", false, "Turns on debugging input.  Default off")
	flag.BoolVar(&detailFlag, "detail", false, "Gives detailed information. Default off")
	flag.BoolVar(&charFlag, "char", true, "Print the character. Default on (use --char=false to disable)")
	flag.BoolVar(&codeFlag, "code", true, "Print the code point. Default on (use --code=false to disable)")
	flag.BoolVar(&listFlag, "list", false, "Print the full description of every emoji in the database")
	flag.BoolVar(&renderFlag, "render", false, "Render the character using sixel graphics.  Default is off.")
	flag.StringVar(&renderWhichFlag, "image", "1", "Choose which image to render.  Defaults to 1.  Can be an integer or all")
	flag.Parse()
	dbMap = make(map[string]int)
	dbMap["code"] = 0
	dbMap["char"] = 1
	dbMap["image"] = 2
}

func makeImage(data string) image.Image {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	//img, format, err := image.Decode(reader)
	img, _, err := image.Decode(reader)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("Decoded image of format: %s\n", format)
	return img
}

func printImage(picture image.Image) {
	fmt.Println()
	sixel.NewEncoder(os.Stdout).Encode(picture)
}

func scaleImageTo(picture image.Image, x int, y int) image.Image {
	if picture.Bounds().Max.X <= x && picture.Bounds().Max.Y <= y {
		return picture
	}
	picX := float64(picture.Bounds().Max.X)
	picY := float64(picture.Bounds().Max.Y)
	larger := picX
	if picY > larger {
		larger = picY
	}
	s := float64(x) / larger
	scaledPicX := int(picX * s)
	scaledPicY := int(picY * s)
	picRatio := picX / picY
	//targetRatio := float64(x) / float64(y)
	if debugFlag {
		fmt.Printf("scaling a %00f,%00f %02f ratio image to %d,%d.  Good luck.\n", picX, picY, picRatio, scaledPicX, scaledPicY)
	}
	scaled := image.NewRGBA(image.Rect(0, 0, scaledPicX, scaledPicY))
	//draw.NearestNeighbor.Scale(scaled, scaled.Rect, picture, picture.Bounds(), draw.Over, nil)
	// nn is fast but ugly, try CatmullRom
	draw.CatmullRom.Scale(scaled, scaled.Rect, picture, picture.Bounds(), draw.Over, nil)
	return scaled
}

func main() {
	setup()
	xSize = 22
	ySize = 22
	if listFlag {
		for descrip, _ := range db {
			fmt.Println(descrip)
		}
	} else {
		if len(flag.Args()) > 0 {
			targetEmoji := strings.Join(flag.Args(), " ")
			realDescrip, results := getter(targetEmoji, searchAlgo)
			if len(results) == 0 {
				fmt.Printf("No emoji matched %s.  Maybe look through the list of descriptions with --list?\n", targetEmoji)
				os.Exit(1)
			}
			//fmt.Println(results[1])
			if debugFlag {
				fmt.Printf("our flag %s is currently at %v\n", "helpFlag", helpFlag)
				fmt.Printf("%s, codepoint %s, description %s\n", results[1], results[0], realDescrip)
				for i := dbMap["image"]; i < len(results); i++ {
					fmt.Printf("len(results[%d]) = %d\n", i, len(results[i]))
					//sixel.NewEncoder(os.Stdout).Encode(makeImage(results[i]))
					//printImage(makeImage(results[i]))
					//printImage(scaleImageTo(makeImage(results[i]),11,22))
					// that's the size on my current term :P, so double width:
					printImage(scaleImageTo(makeImage(results[i]), xSize, ySize))
				}
			}
			if detailFlag {
				fmt.Printf("%s translates to \"%s\".  Code point is %s.\n", targetEmoji, realDescrip, results[dbMap["code"]])
				fmt.Printf("This entry contains %d example images to render\n", len(results)-dbMap["image"])
				m := len(results) - dbMap["image"]
				for i := 0; i < m; i++ {
					fmt.Println(i + 1)
					printImage(scaleImageTo(makeImage(results[i+dbMap["image"]]), xSize, ySize))
				}
			}
			if renderFlag {
				images := results[dbMap["image"]:]
				if renderWhichFlag == "1" {
					// this is the default
					printImage(scaleImageTo(makeImage(images[0]), xSize, ySize))
				} else if renderWhichFlag == "all" {
				} else {
					index, _ := strconv.Atoi(renderWhichFlag)
					realindex := index - 1
					if realindex >= len(images) || realindex < 0 {
						fmt.Printf("there are %d images for %s, and you asked for %d\n", len(images), realDescrip, index)
						os.Exit(1)
					} else {
						printImage(scaleImageTo(makeImage(images[realindex]), xSize, ySize))
					}
				}
			}
			if charFlag {
				if dbMap["char"] > len(results) {
					fmt.Printf("there are only %d results?\n", len(results))
				}
				fmt.Println(results[dbMap["char"]])
			}
			if codeFlag {
				fmt.Println(results[dbMap["code"]])
			}
		}
	}
	//this abomination of a version has the key of description, codepoint, glyph, and
	//the rest are the raw image data
	//"grinning face":                           []string{"U+1F600", "😀"
	//ratios()
}

/*
const data = `
R0lGODlhDwANALMPAK6slY6OfmWz/vz8/XFxdKwAAENWjqFRTO9lTYjE/jiR5lkcRzk3aP43GTs7OwAAACH5BAEAAA8ALAAAAAAPAA0AAARK8Mn5wKmX6lYe1xpyiaAEnOgJBkRLsK1DAUFt2zLzBmlPSIzEYEgsDiaGA2LJZF4SUMNiSqUyHoKEQKDoer+KrXhMFjyg6HT6EQEAOw==
`*/
