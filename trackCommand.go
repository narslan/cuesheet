package cuesheet

import (
	"fmt"
	"strconv"

	"github.com/narslan/tree"
)

type trackCmd struct {
	ID     string
	Format string
}

func (c trackCmd) String() string {
	return fmt.Sprintf("%s %s", c.ID, c.Format)
}

func StrToInt(s string) (int, error) {

	if s[0] == '0' {
		s = s[1:]
	}
	return strconv.Atoi(s)
}

func ParseTimeIndex(m, s, f string) (id Index, err error) {

	m1, err := StrToInt(m)
	if err != nil {
		return
	}

	s1, err := StrToInt(s)
	if err != nil {
		return
	}

	f1, err := StrToInt(f)
	if err != nil {
		return
	}

	id.Minutes = m1
	id.Seconds = s1
	id.Frames = f1
	return
}

// matchTrack captures the arguments of TRACK command and its childeren.
func (p *parser) matchTrack(ft *tree.Tree, fc *File) (err error) {

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
	tid, err := StrToInt(command.ID)
	if err != nil {
		return err
	}
	tr := &Track{ID: tid, Indices: make([]*Index, 0)}
	fc.Tracks = append(fc.Tracks, tr)

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

			in, err := ParseTimeIndex(s.Min, s.Sec, s.Frames)
			if err != nil {
				return err
			}
			tr.Indices = append(tr.Indices, &in)
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
			err := p.matchTrack(ft, fc)
			if err != nil {
				return err
			}
		case itemFile:
			p.backup()
			break OUT

		case itemError:
			return fmt.Errorf("%s on line %d", item, item.line)

		}
	}

	return nil
}
