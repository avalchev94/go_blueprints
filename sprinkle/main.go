package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func getTransformations() []string {
	f, err := os.Open("transformations.txt")
	if err != nil {
		fmt.Print("Failed to open the transformation file!")
		return nil
	}

	var transformations []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		transformations = append(transformations, scanner.Text())
	}

	return transformations
}

func main() {

	transforms := getTransformations()

	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := transforms[rand.Intn(len(transforms))]
		fmt.Println(strings.Replace(t, "*", s.Text(), -1))
	}
}
