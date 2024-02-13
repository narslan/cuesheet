//This code contains some parts from go source code.
// Mainly from src/text/template/parse/parse.go
// Copyright 2011 The Go Authors.

package cuesheet

import (
	"fmt"

	"github.com/narslan/tree"
)

// Parser
type Parser struct {
	lexer     *lexer
	tree      *tree.Tree
	token     [3]item // three-token lookahead for parser.
	peekCount int
}
type node struct {
	Type  itemType
	Value string
}

func (n node) String() string {
	return fmt.Sprintf("%s %s", n.Type, n.Value)
}
func NewParser(input string) *Parser {
	l := lex(input)
	p := &Parser{lexer: l}
	return p
}

func (p *Parser) Start() (*tree.Tree, error) {

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
			n := node{Type: item.typ, Value: s.Value}
			p.tree.AddTree(n)

		case itemTitle:
			s, err := p.matchTitle()
			if err != nil {
				return nil, err
			}
			n := node{Type: item.typ, Value: s.Value}
			p.tree.AddTree(n)
		case itemIndex:
			s, err := p.matchIndex()
			if err != nil {
				return nil, err
			}
			n := node{Type: item.typ, Value: s.Value}
			p.tree.AddTree(n)
		case itemFile:
			s, err := p.matchFile()
			if err != nil {
				return nil, err
			}
			n := node{Type: item.typ, Value: s.Value}
			ft := p.tree.AddTree(n)

			err = p.matchTrack(ft)
			if err != nil {
				return nil, err
			}

		}

	}
	return tree, nil
}
func (p *Parser) next() item {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.token[0] = p.lexer.nextItem()
	}
	return p.token[p.peekCount]
}

func (p *Parser) nextNonSpace() (token item) {
	for {
		token = p.next()
		if (token.typ != itemSpace) && (token.typ != itemNewline) {
			break
		}
	}
	return token
}

// backup backs the input stream up one token.
func (p *Parser) backup() {
	p.peekCount++
}
