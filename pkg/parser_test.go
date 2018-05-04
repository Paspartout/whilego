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
	"reflect"
	"strings"
	"testing"
)

func makeIncrExpr(v int, dec bool) Expr {
	incrExpr := &IncrExpr{v, dec}
	return Expr{Type: INCR_EXPR, IncrExpr: incrExpr}
}

func makeSeqExpr(p1, p2 *Expr) Expr {
	incrExpr := &SeqExpr{p1, p2}
	return Expr{Type: SEQ_EXPR, SeqExpr: incrExpr}
}

func TestParse(t *testing.T) {
	type TestCase struct {
		input    string
		expected Expr
	}

	incrX1 := makeIncrExpr(1, false)
	decrX1 := makeIncrExpr(1, true)
	tests := map[string]TestCase{
		"Increment x1":  {"x1 := x1 + 1", incrX1},
		"Decrement x1":  {"x1 := x1 - 1", decrX1},
		"Increment x42": {"x42 := x42 + 1", makeIncrExpr(42, false)},
		"Decrement x42": {"x42 := x42 - 1", makeIncrExpr(42, true)},
		"x1++;x1--":     {"x1 := x1 + 1 ; x1 := x1 - 1", makeSeqExpr(&incrX1, &decrX1)},
	}

	for caseName, testCase := range tests {
		reader := strings.NewReader(testCase.input)
		parser := NewParser(reader)

		expr, err := parser.Parse()
		if err != nil {
			t.Errorf("%s: %s", caseName, err)
		}
		gotExpr := *expr
		if !reflect.DeepEqual(gotExpr, testCase.expected) {
			// TODO: Implement Stringer for Expr
			t.Errorf("%s: expected %s, got %s", caseName, testCase.expected, gotExpr)
		}
	}
}
