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
	"strings"
	"testing"
)

func TestSingleToken(t *testing.T) {
	type TestCase struct {
		input    string
		expected Token
		literal  string
	}

	tests := map[string]TestCase{
		// Valid inputs for lexer
		"EOF":              {input: string(rune(0)), expected: EOF},
		"Space( )":         {input: "   ", expected: WS},
		"Tab(\\t)":         {input: "\t\t", expected: WS},
		"Newline(\\n)":     {input: "\n\n", expected: WS},
		"Mixed Whitespace": {input: "\n\t ", expected: WS},
		"Variable 0":       {input: "x0", expected: VARIABLE, literal: "x0"},
		"Variable 1":       {input: "x1", expected: VARIABLE, literal: "x1"},
		"Variable 42":      {input: "x42", expected: VARIABLE, literal: "x42"},
		"Constant 0":       {input: "0", expected: CONSTANT, literal: "0"},
		"Constant 1":       {input: "1", expected: CONSTANT, literal: "1"},
		"Semicolon ;":      {input: ";", expected: SEMICOLON},
		"Assign :=":        {input: ":=", expected: ASSIGN},
		"Not Equal !=":     {input: "!=", expected: NOTEQUAL},
		"Operator +":       {input: "+", expected: PLUS},
		"Operator -":       {input: "-", expected: MINUS},
		"Keyword WHILE":    {input: "WHILE", expected: WHILE},
		"Keyword DO":       {input: "DO", expected: DO},
		"Keyword END":      {input: "END", expected: END},

		// Tests for invalid inputs
		"Invalid Assign":    {input: ":!", expected: ILLEGAL},
		"Invalid Not Equal": {input: "!!", expected: ILLEGAL},
		// TODO: Fix variable scanning
		// "Variable 042": {input: "x042", expected: ILLEGAL},
	}

	for caseName, testCase := range tests {
		reader := strings.NewReader(testCase.input)
		scanner := NewScanner(reader)

		tok, lit, err := scanner.Scan()
		if testCase.expected != ILLEGAL && err != nil {
			t.Fatalf("error scanning case for \"%s\": %s", caseName, err)
		}

		if tok != testCase.expected {
			t.Fatalf("%s: expected %s but got %s",
				caseName, testCase.expected, tok)

			if testCase.literal != "" && lit != testCase.literal {
				t.Fatalf("%s: expected literal \"%s\" but got \"%s\"",
					caseName, testCase.literal, lit)
			}
		}
	}
}

func TestMultipleTokens(t *testing.T) {
	type TestCase struct {
		input    string
		expected []Token
		literals []string
	}

	tests := map[string]TestCase{
		"Some Program": {input: "WHILE x1 != 0 DO x1 := x1 - 1 END", expected: []Token{
			WHILE, WS, VARIABLE, WS, NOTEQUAL, WS, CONSTANT, WS, DO, WS,
			VARIABLE, WS, ASSIGN, WS, VARIABLE, WS, MINUS, WS, CONSTANT, WS, END}},
	}

	for caseName, testCase := range tests {
		reader := strings.NewReader(testCase.input)
		scanner := NewScanner(reader)

		var tok Token
		for i := 0; i < len(testCase.expected); i++ {
			tok, _, _ = scanner.Scan()
			if tok != testCase.expected[i] {
				t.Fatalf("%s: expected token %s but was %s", caseName, testCase.expected[i], tok)

			}
		}
	}

}
