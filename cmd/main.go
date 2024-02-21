package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/narslan/cuesheet"
)

// main ...
func main() {

	args := os.Args
	if len(args) == 1 {
		log.Fatal("expected a file input")
	}

	path := args[1]
	fd, err := os.Open(path)
	if err != nil {
		log.Fatal("error reading source file:", err)
	}

	//this is a hcak to deal with BOM (Byte order Mark).
	//https://stackoverflow.com/questions/21371673/reading-files-with-a-bom-in-go
	br := bufio.NewReader(fd)
	r, _, err := br.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	if r != '\uFEFF' {
		br.UnreadRune() // Not a BOM -- put the rune back
	}

	source, err := io.ReadAll(br)
	if err != nil {
		log.Fatal(err)
	}

	c, err := cuesheet.NewCueFile(string(source))

	if err != nil {
		log.Fatal(err)
	}
	log.Println(c.Files())

}
