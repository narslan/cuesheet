package cuesheet

import (
	"fmt"

	"github.com/narslan/tree"
)

type Cuefile struct {
	t *tree.Tree
}

func NewCueFile(input string) (*Cuefile, error) {

	parse := newParser(input)
	tree, err := parse.Start()
	if err != nil {
		return nil, err
	}
	c := &Cuefile{t: tree}
	return c, nil
}

func (c *Cuefile) Files() []string {

	result := make([]string, 0)
	for _, item := range c.t.Traverse() {
		switch v := item.(type) {
		case node:
			switch v.Type.typ {
			case itemFile:
				if v.Type.typ == itemFile {
					result = append(result, v.Value.String())
				}
			}

		}
	}

	fmt.Println(c.t)
	return result
}
