package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_-"

func getTLDs() []string {
	var fileName = flag.String("tlds", "", "Text file with tlds.")
	flag.Parse()

	var tlds = []string{"com", "net"}

	if *fileName != "" {
		f, err := os.Open(*fileName)
		if err != nil {
			fmt.Println("Failed to open ", *fileName)
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			tlds = append(tlds, scanner.Text())
		}
	}

	return tlds
}

func main() {
	tlds := getTLDs()

	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		text := strings.ToLower(s.Text())
		var newText []rune
		for _, r := range text {
			if unicode.IsSpace(r) {
				r = '-'
			}
			if !strings.ContainsRune(allowedChars, r) {
				continue
			}
			newText = append(newText, r)
		}
		fmt.Println(string(newText) + "." + tlds[rand.Intn(len(tlds))])
	}
}
