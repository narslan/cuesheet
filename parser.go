//This code contains some parts from go source code.
// Mainly from src/text/template/parse/parse.go
// Copyright 2011 The Go Authors.

package cuesheet

import (
	"fmt"

	"github.com/narslan/tree"
)

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

func (p *parser) Start() (*tree.Tree, error) {

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
				return nil, err
			}
			n := node{Type: item, Value: s}
			p.tree.AddTree(n)

		case itemTitle:
			s, err := p.matchTitle()
			if err != nil {
				return nil, err
			}
			n := node{Type: item, Value: s}
			p.tree.AddTree(n)
		case itemIndex:
			s, err := p.matchIndex()
			if err != nil {
				return nil, err
			}
			n := node{Type: item, Value: s}
			p.tree.AddTree(n)
		case itemFile:
			s, err := p.matchFile()
			if err != nil {
				return nil, err
			}
			n := node{Type: item, Value: s}
			ft := p.tree.AddTree(n)

			err = p.matchTrack(ft)
			if err != nil {
				return nil, err
			}

		case itemError:
			return nil, fmt.Errorf("reading error %s at pos %d of line %d", item, item.pos, item.line)

		}

	}
	return tree, nil
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
