package cuesheet

import (
	"fmt"

	"github.com/narslan/tree"
)

type trackCmd struct {
	ID     string
	Format string
}

func (c trackCmd) String() string {
	return fmt.Sprintf("%s %s", c.ID, c.Format)
}

// matchTrack captures the arguments of TRACK command and its childeren.
func (p *parser) matchTrack(ft *tree.Tree) (err error) {

	item := p.nextNonSpace()
	if item.typ != itemTrack {
		return fmt.Errorf("[matchTrack] expected track at %d %d", p.lexer.line, p.lexer.pos)
	}
	//we save the first match of "TRACK" item to append into tree later
	firstTrackItem := item
	command := trackCmd{}

	item = p.nextNonSpace()
	if item.typ == itemNumber {
		command.ID = item.val
	} else {
		return fmt.Errorf("[matchTrack] expected text at %d %d", p.lexer.line, p.lexer.pos)
	}

	item = p.nextNonSpace()
	if item.typ == itemText {
		command.Format = item.val
	} else {
		return fmt.Errorf("[matchTrack] expected text at %d %d", p.lexer.line, p.lexer.pos)
	}

	n := node{Type: firstTrackItem, Value: command}
	tt := ft.AddTree(n)
	//parse children of track
OUT:
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
			n := node{Type: item, Value: s}
			tt.AddTree(n)
		case itemIndex:
			s, err := p.matchIndex()
			if err != nil {
				return err
			}
			n := node{Type: item, Value: s}
			tt.AddTree(n)
		case itemPerformer:
			s, err := p.matchTitle()
			if err != nil {
				return err
			}
			n := node{Type: item, Value: s}
			tt.AddTree(n)
		case itemRem:
			s, err := p.matchRem()
			if err != nil {
				return err
			}
			n := node{Type: item, Value: s}
			tt.AddTree(n)
		case itemTrack:
			p.backup()
			err := p.matchTrack(ft)
			if err != nil {
				return err
			}
		case itemFile:
			p.backup()
			break OUT

		case itemError:
			return fmt.Errorf("reading error %s at pos %d of line %d", item, item.pos, item.line)

		}
	}

	return nil
}
