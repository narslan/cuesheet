package cuesheet

import (
	"fmt"
	"strings"

	"github.com/narslan/tree"
)

// matchTrack captures the arguments of TRACK command and its childeren.
func (p *Parser) matchTrack(ft *tree.Tree) (err error) {

	var s strings.Builder

	item := p.nextNonSpace()
	if item.typ != itemTrack {
		return fmt.Errorf("[matchTrack] expected track at %d %d", p.lexer.line, p.lexer.pos)
	}

	item = p.nextNonSpace()
	if item.typ == itemNumber {
		if s.Len() > 0 {
			s.WriteRune(' ')
		}
		s.WriteString(item.val)

	} else {
		return fmt.Errorf("[matchTrack] expected text at %d %d", p.lexer.line, p.lexer.pos)
	}

	item = p.nextNonSpace()
	if item.typ == itemText {
		if s.Len() > 0 {
			s.WriteRune(' ')
		}
		s.WriteString(item.val)

	} else {
		return fmt.Errorf("[matchTrack] expected text at %d %d", p.lexer.line, p.lexer.pos)
	}

	n := node{Type: "TRACK", Value: s.String()}
	tt := ft.AddTree(n)
	//parse children of track

	for {
		item = p.nextNonSpace()
		if item.typ == itemEOF {
			break
		}
		switch item.typ {
		case itemTitle:
			s, err := p.matchTitle()
			if err != nil {
				return err
			}
			n := node{Type: item.String(), Value: s.Value}
			tt.AddTree(n)
		case itemIndex:
			s, err := p.matchIndex()
			if err != nil {
				return err
			}
			n := node{Type: item.String(), Value: s.Value}
			tt.AddTree(n)
		case itemPerformer:
			s, err := p.matchTitle()
			if err != nil {
				return err
			}
			n := node{Type: item.String(), Value: s.Value}
			tt.AddTree(n)
		case itemRem:
			s, err := p.matchRem()
			if err != nil {
				return err
			}
			n := node{Type: item.String(), Value: s.Value}
			tt.AddTree(n)
		case itemTrack:
			p.backup()
			err := p.matchTrack(ft)
			if err != nil {
				return err
			}
		case itemError:
			return fmt.Errorf("reading error %s at pos %d of line %d", item, item.pos, item.line)

		}
	}

	return nil
}
