package depthparser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func newLexer(input string) *lexer {
	return &lexer{input: input}
}

type lexer struct {
	start     int
	pos       int
	input     string
	lastToken token
	atEof     bool
}

func (l *lexer) lex() token {
	for {
		if l.scanNumber() {
			return l.emit(tokenNumber)
		}
		if l.scanSpaceType() {
			return l.emit(tokenSpaceType)
		}
		if l.atEof {
			return l.eof()
		}
		return l.errorf("unsupported characters in input")
	}
}

const eof rune = -1

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.atEof = true
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += w
	return rune(r)
}

func (l *lexer) backup() {
	if !l.atEof && l.pos > 0 {
		_, w := utf8.DecodeLastRuneInString(l.input[:l.pos])
		l.pos -= w
	}
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

const digits = "0123456789"

func (l *lexer) scanSpaceType() bool {
	return l.accept("std")
}

func (l *lexer) scanNumber() bool {
	l.acceptRun(digits)
	return l.pos > l.start
}

func (l *lexer) emit(typ tokenType) token {
	t := token{
		pos:   l.start,
		value: l.input[l.start:l.pos],
		typ:   typ,
	}
	l.start = l.pos
	return t
}

func (l *lexer) eof() token {
	return token{typ: tokenEOF, value: "", pos: l.start}
}

func (l *lexer) errorf(msg string, args ...any) token {
	t := token{typ: tokenError, value: fmt.Sprintf(msg, args...), pos: l.start}
	l.start = 0
	l.pos = 0
	l.input = l.input[:0]
	return t
}

type token struct {
	typ   tokenType
	pos   int
	value string
}

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenSpaceType
	tokenNumber
)

func (t tokenType) String() string {
	switch t {
	case tokenError:
		return "error"
	case tokenEOF:
		return "EOF"
	case tokenSpaceType:
		return "spaceType"
	case tokenNumber:
		return "number"
	}

	panic("unknown token type")
}
