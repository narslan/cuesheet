package cuesheet

import (
	"fmt"
	"strings"
)

type titleCmd struct {
	Value string
}

// matchTitle captures the carguments of TITLE command.
func (p *Parser) matchTitle() (c titleCmd, err error) {

	var s strings.Builder

	item := p.lexer.nextItem()
	if item.typ == itemEOF || item.typ == itemNewline {
		return c, fmt.Errorf("met EOF or newline at %d %d", p.lexer.line, p.lexer.pos)
	}

	if item.typ == itemText || item.typ == itemNumber {
		//Put a space between items.
		if s.Len() > 0 {
			s.WriteRune(' ')
		}
		s.WriteString(item.val)
	} else {
		return c, fmt.Errorf("expected text at %d %d", p.lexer.line, p.lexer.pos)
	}

	rs := titleCmd{Value: s.String()}
	return rs, nil
}

type remCmd struct {
	Value string
}

// matchRem captures the comment of REM command.
func (p *Parser) matchRem() (c remCmd, err error) {

	var s strings.Builder

	for {
		item := p.lexer.nextItem()
		if item.typ == itemEOF || item.typ == itemNewline {
			break
		}

		if item.typ == itemText {

			if s.Len() > 0 {
				s.WriteRune(' ')
			}
			s.WriteString(item.val)
		}

		if item.typ == itemNumber {

			if s.Len() > 0 {
				s.WriteRune(' ')
			}
			s.WriteString(item.val)
		}

	}
	c = remCmd{Value: s.String()}
	return c, nil
}

type indexCmd struct {
	Value string
}

// matchIndex captures the arguments of INDEX command.
func (p *Parser) matchIndex() (c indexCmd, err error) {

	var s strings.Builder

	item := p.nextNonSpace()
	if item.typ == itemEOF || item.typ == itemNewline {
		return c, fmt.Errorf("met EOF or newline at %d %d", p.lexer.line, p.lexer.pos)
	}

	if item.typ == itemNumber {
		//We want to put a space between items.
		if s.Len() > 0 {
			s.WriteRune(' ')
		}
		s.WriteString(item.val)
	} else {
		return c, fmt.Errorf("expected text at %d %d but found %q ", p.lexer.line, p.lexer.pos, item.val)
	}
	s.WriteRune(' ')
	//parsing mm:ss:ff part of the INDEX command.
	for i := 0; i < 5; i++ {
		item = p.nextNonSpace()

		switch item.typ {
		case itemNumber, itemColon:
			s.WriteString(item.val)
		default:
			return c, fmt.Errorf("expected text at %d %d but found %q", p.lexer.line, p.lexer.pos, item.val)
		}

	}

	c = indexCmd{Value: s.String()}
	return c, nil
}

type fileCmd struct {
	Value string
}

// matchIndex captures the arguments of FILE command.
func (p *Parser) matchFile() (c fileCmd, err error) {

	var s strings.Builder

	item := p.nextNonSpace()
	if item.typ == itemEOF {
		return c, fmt.Errorf("met EOF or newline at %d %d", p.lexer.line, p.lexer.pos)
	}

	if item.typ == itemText {
		if s.Len() > 0 {
			s.WriteRune(' ')
		}
		s.WriteString(item.val)

	} else {
		return c, fmt.Errorf("expected text at %d %d", p.lexer.line, p.lexer.pos)
	}

	item = p.nextNonSpace()

	switch item.typ {

	case itemBinary, itemMotorola, itemAiff, itemWave, itemMp3:
		if s.Len() > 0 {
			s.WriteRune(' ')
		}
		s.WriteString(item.val)
	default:
		return c, fmt.Errorf("expected a file type at %d %d", p.lexer.line, p.lexer.pos)

	}

	c = fileCmd{Value: s.String()}
	return c, nil
}
