package cuesheet

import (
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
			switch f := v.Value.(type) {
			case fileCmd:
				result = append(result, f.Path)
			}
		}
	}

	return result
}
