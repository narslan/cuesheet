//This code contains some parts from go source code.
// Mainly from src/text/template/parse/parse.go
// Copyright 2011 The Go Authors.

package cuesheet

import (
	"fmt"

	"github.com/narslan/tree"
)

type Index struct {
	Minutes int
	Seconds int
	Frames  int
}

type Track struct {
	ID      int
	Indices []*Index
}

type File struct {
	Path   string
	Tracks []*Track
}

// parser
type parser struct {
	lexer     *lexer
	tree      *tree.Tree
	token     [3]item // three-token lookahead for parser.
	peekCount int
}
type Command interface {
	String() string
}
type node struct {
	Type  item
	Value Command
}

func (n node) String() string {
	return fmt.Sprintf("%v %s", n.Type, n.Value)
}
func newParser(input string) *parser {
	l := lex(input)
	p := &parser{lexer: l}
	return p
}

func (p *parser) Start() (*tree.Tree, []*File, error) {

	f := make([]*File, 0)
	tree := tree.New("root")
	p.tree = tree
	for {
		item := p.nextNonSpace()
		if item.typ == itemEOF {
			break
		}
		switch item.typ {
		case itemRem:
			s, err := p.matchRem()
			if err != nil {
				return nil, f, err
			}
			n := node{Type: item, Value: s}
			p.tree.AddTree(n)

		case itemTitle:
			s, err := p.matchTitle()
			if err != nil {
				return nil, f, err
			}
			n := node{Type: item, Value: s}
			p.tree.AddTree(n)
		case itemIndex:
			s, err := p.matchIndex()
			if err != nil {
				return nil, f, err
			}
			n := node{Type: item, Value: s}
			p.tree.AddTree(n)
		case itemFile:
			s, err := p.matchFile()
			if err != nil {
				return nil, f, err
			}
			n := node{Type: item, Value: s}
			ft := p.tree.AddTree(n)
			file := &File{Path: s.Path, Tracks: make([]*Track, 0)}
			err = p.matchTrack(ft, file)
			if err != nil {
				return nil, f, err
			}
			f = append(f, file)

		case itemError:
			return nil, f, fmt.Errorf("%s on line %d", item, item.line)

		}

	}
	return tree, f, nil
}
func (p *parser) next() item {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.token[0] = p.lexer.nextItem()
	}
	return p.token[p.peekCount]
}

func (p *parser) nextNonSpace() (token item) {
	for {
		token = p.next()
		if (token.typ != itemSpace) && (token.typ != itemNewline) {
			break
		}
	}
	return token
}

// backup backs the input stream up one token.
func (p *parser) backup() {
	p.peekCount++
}
