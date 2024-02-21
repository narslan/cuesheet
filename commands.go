package cuesheet

import (
	"fmt"
	"strings"
)

type titleCmd struct {
	Value string
}

func (c titleCmd) String() string {
	return c.Value
}

// matchTitle captures the carguments of TITLE command.
func (p *parser) matchTitle() (c titleCmd, err error) {

	item := p.lexer.nextItem()
	if item.typ == itemEOF || item.typ == itemNewline {
		return c, fmt.Errorf("met EOF or newline at %d %d", p.lexer.line, p.lexer.pos)
	}

	if item.typ == itemText || item.typ == itemNumber {
		//Put a space between items.
		return titleCmd{Value: item.val}, nil
	} else {
		return c, fmt.Errorf("expected text at %d %d", p.lexer.line, p.lexer.pos)
	}

}

type remCmd struct {
	Values []string
}

func (c remCmd) String() string {
	var s strings.Builder

	for i, v := range c.Values {
		s.WriteString(v)
		if i < len(v)-1 {
			s.WriteRune(' ')
		}
	}

	return s.String()
}

// matchRem captures the comment of REM command.
func (p *parser) matchRem() (c remCmd, err error) {

	c.Values = make([]string, 0)

	for {
		item := p.lexer.nextItem()
		if item.typ == itemEOF || item.typ == itemNewline {
			break
		}

		if item.typ == itemText {
			c.Values = append(c.Values, item.val)
		}

		if item.typ == itemNumber {
			c.Values = append(c.Values, item.val)
		}

	}

	return c, nil
}

type indexCmd struct {
	ID     string
	Min    string
	Sec    string
	Frames string
}

func (c indexCmd) String() string {
	return fmt.Sprintf("%s %s %s %s", c.ID, c.Min, c.Sec, c.Frames)
}

// matchIndex captures the arguments of INDEX command.
func (p *parser) matchIndex() (c indexCmd, err error) {

	idx := indexCmd{}

	item := p.nextNonSpace()
	if item.typ == itemEOF || item.typ == itemNewline {
		return c, fmt.Errorf("met EOF or newline at %d %d", p.lexer.line, p.lexer.pos)
	}

	if item.typ == itemNumber {
		//We want to put a space between items.
		idx.ID = item.val

	} else {
		return c, fmt.Errorf("expected text at %d %d but found %q ", p.lexer.line, p.lexer.pos, item.val)
	}

	//parsing mm:ss:ff part of the INDEX command.
	for i := 0; i < 5; i++ {
		item = p.nextNonSpace()

		switch item.typ {
		case itemNumber:
			if i == 0 {
				idx.Min = item.val
			}
			if i == 2 {
				idx.Sec = item.val
			}
			if i == 4 {
				idx.Frames = item.val
			}
		case itemColon:
		default:
			return c, fmt.Errorf("expected text at %d %d but found %q", p.lexer.line, p.lexer.pos, item.val)
		}

	}

	return idx, nil
}

type fileCmd struct {
	Path   string
	Format string
}

func (f fileCmd) String() string {
	return fmt.Sprintf("%s %s", f.Path, f.Format)
}

// matchIndex captures the arguments of FILE command.
func (p *parser) matchFile() (c fileCmd, err error) {

	item := p.nextNonSpace()
	if item.typ == itemEOF {
		return c, fmt.Errorf("met EOF or newline at %d %d", p.lexer.line, p.lexer.pos)
	}

	if item.typ == itemText {
		c.Path = item.val
	} else {
		return c, fmt.Errorf("expected text at %d %d", p.lexer.line, p.lexer.pos)
	}

	item = p.nextNonSpace()

	switch item.typ {

	case itemBinary, itemMotorola, itemAiff, itemWave, itemMp3:
		c.Format = item.val
	default:
		return c, fmt.Errorf("expected a file type at %d %d", p.lexer.line, p.lexer.pos)

	}

	return c, nil
}
