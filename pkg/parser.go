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
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ExprType denotes the type of a WHILE expression.
type ExprType int

const (
	// INVALID_EXPR indicates an invalid expression
	INVALID_EXPR ExprType = iota

	// INCR_EXPR indicates an expression of the from `xN := xN +/- 1`
	INCR_EXPR

	// SEQ_EXPR indicates a sequence of two expressions, e.g. `P1;P2`
	SEQ_EXPR
	// WHILE_EXPR indicates an expression of the from `WHILE xN != 0 DO P END`
	WHILE_EXPR
)

// Expr is an expression of the WHILE language.
// TODO: Maybe Refactor? (Composition?, Reflection?, Inheritance?)
type Expr struct {
	Type ExprType

	IncrExpr  *IncrExpr
	SeqExpr   *SeqExpr
	WhileExpr *WhileExpr
}

// String returns a simple string representation of the expression.
// It is meant for debugging, not for pretty printing.
func (e Expr) String() string {
	s := "{"
	switch e.Type {
	case INVALID_EXPR:
		s += "Invalid: "
	case INCR_EXPR:
		s += "IncrExpr: "
		s += fmt.Sprint(e.IncrExpr)
	case SEQ_EXPR:
		s += "SeqExpr: "
		s += fmt.Sprintf("P1: %s, P2: %s", e.SeqExpr.P1, e.SeqExpr.P2)
	case WHILE_EXPR:
		s += "WhileExpr: "
	default:
		s += "Unknown: "
	}

	s += "}"
	return s
}

// IncrExpr represents a expression in the form `xN := xN +/- 1`
type IncrExpr struct {
	// The number of the variable in range {0, ...}
	Variable int
	// true means decrement, false increment
	Decrement bool
}

// SeqExpr represents a sequence of two expressions, e.g. `P1;P2`
type SeqExpr struct {
	// P1 is the first program to run.
	P1 *Expr
	// P1 is the second program to run after P1.
	P2 *Expr
}

// WhileExpr represents an expression of the from `WHILE xN != 0 DO P END`
type WhileExpr struct {
	// Variable is the variable N to check `xN != 0` for.
	Variable int
	// P is the program to run while `xN != 0` is true.
	P *Expr
}

// Parser represents a parser for the WHILE language.
type Parser struct {
	s *Scanner
	// Buffer for lookahead
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size(max=1)
	}
}

// NewParser creates a new instance of a WHILE parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the scanner.
// If a token was unscanned, it will return the buffered one instead.
// In case of an error it will also return the error as the third value.
func (p *Parser) scan() (tok Token, lit string, err error) {
	// Take token from buffer if available
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit, nil
	}

	// Write token into buffer in case we unscan later
	tok, lit, err = p.s.Scan()
	p.buf.tok, p.buf.lit = tok, lit

	// This returns the values we have written to tok, lit and err
	return
}

// unscan "pushes" the previously read token back onto the buffer.
func (p *Parser) unscan() {
	p.buf.n = 1
}

// scanIgnoreWhitespace scans the next non-whitespace token.
// If there was an error during scanning it will also return it.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string, err error) {
	tok, lit, err = p.scan()
	if err != nil {
		return
	}
	// Scan next token, if a whitespace was read.
	if tok == WS {
		tok, lit, err = p.scan()
	}
	return
}

// Parse parses the input, given to the parser using the reader.
func (p *Parser) Parse() (*Expr, error) {
	ex1 := &Expr{}
	ex2 := &Expr{}
	var expr *Expr

	tok, _, err := p.scanIgnoreWhitespace()
	if err != nil {
		return expr, fmt.Errorf("error tokenizing: %s", err)
	}

	// TODO: WhileExpr

	// Base case: assignment
	if tok == VARIABLE {
		p.unscan()
		incExpr, err := p.parseIncr()
		if err != nil {
			return nil, err
		}
		ex1.Type = INCR_EXPR
		ex1.IncrExpr = incExpr
		expr = ex1
	}

	tok, _, err = p.scanIgnoreWhitespace()
	// TODO: Check if following condition is sane
	if tok == ILLEGAL || tok == EOF {
		return expr, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error tokenizing after ws: %s", err)
	}

	// If there is a semicolon following the assignment
	if tok == SEMICOLON {
		// Try to parse the following expression
		ex2, err = p.Parse()
		if err != nil {
			return nil,
				fmt.Errorf("no valid expression after semicolon: %s", err)
		}
		expr = &Expr{}
		expr.Type = SEQ_EXPR
		expr.SeqExpr = &SeqExpr{ex1, ex2}
	}

	return expr, nil
}

// parseIncr parses the increment expression of the WHILE language.
func (p *Parser) parseIncr() (*IncrExpr, error) {
	incrExpr := &IncrExpr{}

	// Read left side variable.
	// TODO: Introduce helpers for error reporting
	tok, lit, err := p.scanIgnoreWhitespace()
	if err != nil {
		return nil, fmt.Errorf("error parsing left side variable: %s", err)
	}
	if tok != VARIABLE {
		return nil, errors.New("initial token of increment has to be a variable")
	}
	firstVarNum, err := strconv.Atoi(strings.TrimPrefix(lit, "x"))
	if err != nil {
		return nil, fmt.Errorf("error parsing variable number: %s", err)
	}
	incrExpr.Variable = firstVarNum

	// Check if a assignment token follows
	tok, lit, err = p.scanIgnoreWhitespace()
	if err != nil {
		return nil, fmt.Errorf("error parsing equal sign: %s", err)
	}
	if tok != ASSIGN {
		return nil, fmt.Errorf("expected assignment operator after variable")
	}

	// Read right side variable.
	tok, lit, err = p.scanIgnoreWhitespace()
	if err != nil {
		return nil, fmt.Errorf("error parsing left side variable: %s", err)
	}
	if tok != VARIABLE {
		return nil, errors.New("initial token of increment has to be a variable")
	}
	secondVarNum, err := strconv.Atoi(strings.TrimPrefix(lit, "x"))
	if err != nil {
		return nil, fmt.Errorf("error parsing variable number: %s", err)
	}
	if firstVarNum != secondVarNum {
		return nil,
			fmt.Errorf("second variable index %d has to match the first one which is %d",
				firstVarNum, secondVarNum)
	}

	// Determine increment or decrement
	tok, lit, err = p.scanIgnoreWhitespace()
	if err != nil {
		return nil, fmt.Errorf("error parsing increment/decrement: %s", err)
	}
	switch tok {
	case PLUS:
		incrExpr.Decrement = false
	case MINUS:
		incrExpr.Decrement = true
	default:
		return nil, fmt.Errorf("token \"%s\" has to be - or + sign", lit)
	}

	// Make sure there is a 1 following the +/- sign
	tok, lit, err = p.scanIgnoreWhitespace()
	if err != nil {
		return nil, fmt.Errorf("error parsing number after increment/decrement: %s", err)
	}
	if tok != CONSTANT || lit != "1" {
		return nil, fmt.Errorf("there has to follow a 1 after +/-, got \"%s\"", lit)
	}

	return incrExpr, nil
}
