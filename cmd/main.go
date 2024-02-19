package main

import (
	"github.com/narslan/cuesheet"
	"log"
	"os"
)

// main ...
func main() {

	args := os.Args
	if len(args) == 1 {
		log.Fatal("expected a file input")
	}

	path := args[1]
	source, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("error reading source file:", err)
	}

	parse := cuesheet.NewParser(string(source))
	tree, err := parse.Start()
	if err != nil {
		log.Fatal("error constructing parse tree:", err)
	}
	log.Println(tree)

}
