// Copyright (c) 2017 Ernest Micklei
// 
// MIT License
// 
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
// 
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package proto

import "testing"

func TestSyntax(t *testing.T) {
	proto := `syntax = "proto";`
	p := newParserOn(proto)
	p.scanIgnoreWhitespace() // consume first token
	s := new(Syntax)
	err := s.parse(p)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := s.Value, "proto"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestCommentAroundSyntax(t *testing.T) {
	proto := `
	// comment1
	// comment2
	syntax = 'proto'; // comment3
	// comment4
`
	p := newParserOn(proto)
	r, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	comments := collect(r).Comments()
	if got, want := len(comments), 3; got != want {
		t.Fatalf("got [%v] want [%v]", got, want)
	}
}
