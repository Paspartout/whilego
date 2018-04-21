// Copyright © 2018 Phileas Vöcking <paspartout@fogglabs.de>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package whilego

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// The following line enables go generate together with the tool
// enumer(https://github.com/alvaroloes/enumer) to generate e.g. Stringer for
// our type Token.
//go:generate enumer -type=Token

// Token is a lexical token for the while language.
type Token int

const (
	// Special tokens

	// ILLEGAL represents an illegal token.
	ILLEGAL Token = iota
	// EOF stands for End Of File.
	EOF
	// WS represents a White Space.
	WS

	// Literals

	// VARIABLE represents the name of a variable like x0, x1, x2, ...
	VARIABLE

	// CONSTANT is a numerical constant that is either 0 or 1
	CONSTANT

	// Symbols

	// SEMICOLON represents a ;
	SEMICOLON
	// ASSIGN is represented by :=
	ASSIGN
	// NOTEQUAL is represented by !=
	NOTEQUAL

	// Operators

	// PLUS is represented by +
	PLUS

	// MINUS is represented by -
	MINUS

	// Keywords

	// WHILE is the keyword represented by the string "WHILE"
	WHILE

	// DO is the keyword represented by the string "DO"
	DO

	// END is the keyword represented by the string "END"
	END
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

var eof = rune(0)

// Scanner is the lexical scanner for the WHILE language.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner creates and returns a new instance of a WHILE scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the buffered reader.
func (s *Scanner) read() (rune, error) {
	ch, _, err := s.r.ReadRune() // _ ignores the rune size
	if err != nil {
		return eof, err
	}
	return ch, nil
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() error { return s.r.UnreadRune() }

// Scan returns the next token and literal value.
// If an error occurs during reading it returns an error.
func (s *Scanner) Scan() (tok Token, lit string, err error) {
	ch, err := s.read()
	if err != nil {
		return scanError(fmt.Errorf("error reading next character: %s", err))
	}

	// Consume contiguous whitespace if we see one.
	if isWhitespace(ch) {
		err = s.unread()
		if err != nil {
			return scanError(fmt.Errorf("error unreading character: %s", err))
		}
		return s.scanWhitespace()
	}

	// Scan single character tokens
	switch ch {
	case eof:
		return EOF, "", nil
	case 'x':
		return s.scanVariable()
	case '0':
		return CONSTANT, string(ch), nil
	case '1':
		return CONSTANT, string(ch), nil
	case ':':
		return s.scanString(ASSIGN, ":=")
	case '!':
		return s.scanString(NOTEQUAL, "!=")
	case ';':
		return SEMICOLON, string(ch), nil
	case '+':
		return PLUS, string(ch), nil
	case '-':
		return MINUS, string(ch), nil
	case 'W':
		return s.scanString(WHILE, "WHILE")
	case 'D':
		return s.scanString(DO, "DO")
	case 'E':
		return s.scanString(END, "END")
	}

	return ILLEGAL, string(ch), nil
}

// scanWhitespace consumes the current rune and all following whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string, err error) {
	// Create buffer and read current character into it.
	var buf bytes.Buffer
	ch, err := s.read()
	if err != nil {
		return scanError(fmt.Errorf("error tokenizing whitespace: %s", err))
	}
	buf.WriteRune(ch)

	// Read every following whitespace character until the next EOF or
	// non whitespace character.
	for {
		ch, err = s.read()
		if err != nil && ch != eof {
			return scanError(fmt.Errorf("error tokenizing whitespace: %s", err))
		}
		if ch == eof {
			break
		} else if !isWhitespace(ch) {
			err = s.unread()
			if err != nil {
				return scanError(fmt.Errorf("error tokenizing whitespace: %s", err))
			}
			break
		} else {
			_, err = buf.WriteRune(ch)
			if err != nil {
				return scanError(fmt.Errorf("error tokenizing whitespace: %s", err))
			}
		}
	}

	return WS, buf.String(), nil
}

// scanVariable unreads the last character and reads in a variable name
// in the form of x0, x1, ...
func (s *Scanner) scanVariable() (tok Token, lit string, err error) {
	err = s.unread()
	if err != nil {
		return scanError(err)
	}

	var buf bytes.Buffer

	// Read initial x
	ch, err := s.read()
	if err != nil {
		return scanError(err)
	}
	if ch != 'x' {
		return scanError(fmt.Errorf("variable does not start with x"))
	}
	buf.WriteRune(ch)

	// Read every following digit until non digit character.
	for {
		ch, err = s.read()
		if err != nil && ch != eof {
			return scanError(fmt.Errorf("error tokenizing whitespace: %s", err))
		}
		if ch == eof {
			break
		} else if !isDigit(ch) {
			err = s.unread()
			if err != nil {
				return scanError(fmt.Errorf("error tokenizing digit: %s", err))
			}
			break
		} else {
			_, err = buf.WriteRune(ch)
			if err != nil {
				return scanError(fmt.Errorf("error tokenizing whitespace: %s", err))
			}
		}
	}

	return VARIABLE, buf.String(), nil
}

// scanString unreads the last rune and matches for expectedToken.
// If the string is not matched, it returns the ILLEGAL token.
func (s *Scanner) scanString(expectedToken Token, str string) (Token, string, error) {
	err := s.unread()
	if err != nil {
		return scanError(err)
	}
	var buf bytes.Buffer

	for _, strCh := range str {
		ch, err := s.read()
		if err != nil {
			return scanError(err)
		}
		if ch != strCh {
			return scanError(fmt.Errorf("unexpected string: %s", buf.String()))
		}
		_, err = buf.WriteRune(ch)
		if err != nil {
			return scanError(fmt.Errorf("error tokenizing whitespace: %s", err))
		}
	}

	return expectedToken, buf.String(), nil
}

// scanError is helper for returning an error in scanner functions.
func scanError(err error) (Token, string, error) {
	return ILLEGAL, "", err
}

// TODO: unreadAnd(scanMethod) instead of unreading in every method
