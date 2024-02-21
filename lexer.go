//This code contains some parts from go source code.
// From src/text/template/parse/lexer.go
// Copyright 2011 The Go Authors.

package cuesheet

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type itemType int

// Item represents a token or text string returned from the scanner.
type item struct {
	typ  itemType // The type of this Item.
	pos  int      // The starting position, in bytes, of this Item in the input string.
	val  string   // The value of this Item.
	line int      // The line number at the start of this Item.
}

const eof = -1

const (
	itemError itemType = iota // error occurred; value is text of error
	itemEOF
	itemSpace
	itemText
	itemNumber
	itemColon
	itemNewline

	itemKeyword
	//every keyword come after here
	itemFile
	itemTitle
	itemRem
	itemPerformer
	itemIndex
	itemTrack

	//keywords go above,The last argument of file command
	itemFileType
	itemBinary
	itemMotorola
	itemAiff
	itemWave
	itemMp3
)

var key = map[string]itemType{
	"FILE": itemFile,

	"TITLE":     itemTitle,
	"REM":       itemRem,
	"PERFORMER": itemPerformer,
	"INDEX":     itemIndex,
	"TRACK":     itemTrack,
	"BINARY":    itemBinary,
	"MOTOROLA":  itemMotorola,
	"AIFF":      itemAiff,
	"WAVE":      itemWave,
	"MP3":       itemMp3,
}

func (i item) String() string {

	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ > itemKeyword:
		return i.val
	}
	return fmt.Sprintf("%q", i.val)

}

// stateFn represents the state of the lexer as a function that returns the
// next state.

type stateFn func(*lexer) stateFn

type lexer struct {
	input string // the string being scanned
	pos   int    // current position in the input
	start int    // start position of this Item
	width int    // width of last rune read from input
	item  item

	line      int // 1+number of newlines seen
	startLine int // start line of this Item
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	if r == '\n' {
		l.line++

	}
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.item = item{itemError, l.start, fmt.Sprintf(format, args...), l.startLine}
	l.start = 0
	l.pos = 0
	l.input = l.input[:0]
	return nil
}

// // ignore skips over the pending input before this point.
// func (l *lexer) ignore() {
// 	l.line += strings.Count(l.input[l.start:l.pos], "\n")
// 	l.start = l.pos
// 	l.startLine = l.line

// }

// thisItem returns the item at the current input point with the specified type
// and advances the input.
func (l *lexer) thisItem(t itemType) item {
	i := item{t, l.start, l.input[l.start:l.pos], l.startLine}
	l.start = l.pos
	l.startLine = l.line
	//log.Print(i.typ, " ", i.val)
	return i
}

// emit passes an Item back to the client.
func (l *lexer) emit(t itemType) stateFn {
	return l.emitItem(l.thisItem(t))
}

// emitItem passes the specified item to the parser.
func (l *lexer) emitItem(i item) stateFn {
	l.item = i
	return nil
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	l.item = item{itemEOF, l.pos, "EOF", l.startLine}
	state := lexText
	for {
		state = state(l)
		if state == nil {
			return l.item
		}
	}
}

// lexText scans everything.
func lexText(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(itemEOF)
		return nil
	case isEndOfLine(r):
		return l.emit(itemNewline)
	case isSpace(r):
		return lexSpace
	case r == ':':
		return l.emit(itemColon)
	case r == '"':
		return lexQuote
	case unicode.IsNumber(r):
		return lexNumber
	case isIdentifierChar(r):
		return lexIdentifier
	default:
		return l.errorf("bad syntax: %q", l.input[l.start:l.pos])
	}

}

func lexNumber(l *lexer) stateFn {

	f := func(r rune) bool {
		return unicode.IsLetter(r)
	}

	for {
		switch r := l.next(); {
		case isIdentifierChar(r):
			// absorb.
		default:
			l.backup()
			word := strings.ToUpper(l.input[l.start:l.pos]) //Commands are case insensitive.
			if strings.ContainsFunc(word, f) {
				return l.emit(itemText)
			} else {
				return l.emit(itemNumber)
			}

		}
	}

}

// lexQuote scans until the end of a quote
func lexQuote(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	return l.emit(itemText)
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isIdentifierChar(r):
			// absorb.
		default:
			l.backup()
			word := strings.ToUpper(l.input[l.start:l.pos]) //Commands are case insensitive.
			switch {
			case key[word] > itemKeyword:
				item := key[word]
				return l.emit(item)
			default:
				return l.emit(itemText)
			}
		}
	}
}

func lexSpace(l *lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}
	l.emit(itemSpace)
	return lexText
}

// lex creates a new scanner for the input string.
func lex(input string) *lexer {
	l := &lexer{
		input:     input,
		line:      1,
		startLine: 1,
	}
	return l
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isIdentifierChar reports whether r is an alphabetic, digit, underscore, star, period or dash.
func isIdentifierChar(r rune) bool {
	return r == '*' || r == '_' || r == '-' || r == '.' || r == '/' || unicode.IsLetter(r) ||
		unicode.IsDigit(r)
}
