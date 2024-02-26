package cuesheet

import (
	"github.com/narslan/tree"
)

type Cuefile struct {
	t     *tree.Tree
	files []*File
}

func New(input string) (*Cuefile, error) {

	parse := newParser(input)
	tree, f, err := parse.Start()
	if err != nil {
		return nil, err
	}
	c := &Cuefile{t: tree, files: f}
	return c, nil
}

func (c *Cuefile) Files() []*File {

	return c.files
}
