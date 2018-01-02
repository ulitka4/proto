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

import (
	"text/scanner"
)

// Enum definition consists of a name and an enum body.
type Enum struct {
	Position scanner.Position
	Comment  *Comment
	Name     string
	Elements []Visitee
}

// Accept dispatches the call to the visitor.
func (e *Enum) Accept(v Visitor) {
	v.VisitEnum(e)
}

// Doc is part of Documented
func (e *Enum) Doc() *Comment {
	return e.Comment
}

// addElement is part of elementContainer
func (e *Enum) addElement(v Visitee) {
	e.Elements = append(e.Elements, v)
}

// elements is part of elementContainer
func (e *Enum) elements() []Visitee {
	return e.Elements
}

// takeLastComment is part of elementContainer
// removes and returns the last element of the list if it is a Comment.
func (e *Enum) takeLastComment(expectedOnLine int) (last *Comment) {
	last, e.Elements = takeLastCommentIfOnLine(e.Elements, expectedOnLine)
	return
}

func (e *Enum) parse(p *Parser) error {
	pos, tok, lit := p.next()
	if tok != tIDENT {
		if !isKeyword(tok) {
			return p.unexpected(lit, "enum identifier", e)
		}
	}
	e.Name = lit
	_, tok, lit = p.next()
	if tok != tLEFTCURLY {
		return p.unexpected(lit, "enum opening {", e)
	}
	for {
		pos, tok, lit = p.next()
		switch tok {
		case tCOMMENT:
			if com := mergeOrReturnComment(e.elements(), lit, pos); com != nil { // not merged?
				e.Elements = append(e.Elements, com)
			}
		case tOPTION:
			v := new(Option)
			v.Position = pos
			v.Comment = e.takeLastComment(pos.Line)
			err := v.parse(p)
			if err != nil {
				return err
			}
			e.Elements = append(e.Elements, v)
		case tRIGHTCURLY, tEOF:
			goto done
		case tSEMICOLON:
			maybeScanInlineComment(p, e)
		default:
			p.nextPut(pos, tok, lit)
			f := new(EnumField)
			f.Position = pos
			f.Comment = e.takeLastComment(pos.Line - 1)
			err := f.parse(p)
			if err != nil {
				return err
			}
			e.Elements = append(e.Elements, f)
		}
	}
done:
	if tok != tRIGHTCURLY {
		return p.unexpected(lit, "enum closing }", e)
	}
	return nil
}

// EnumField is part of the body of an Enum.
type EnumField struct {
	Position      scanner.Position
	Comment       *Comment
	Name          string
	Integer       int
	ValueOption   *Option
	InlineComment *Comment
}

// Accept dispatches the call to the visitor.
func (f *EnumField) Accept(v Visitor) {
	v.VisitEnumField(f)
}

// inlineComment is part of commentInliner.
func (f *EnumField) inlineComment(c *Comment) {
	f.InlineComment = c
}

// Doc is part of Documented
func (f *EnumField) Doc() *Comment {
	return f.Comment
}

func (f *EnumField) parse(p *Parser) error {
	_, tok, lit := p.nextIdentifier()
	if tok != tIDENT {
		if !isKeyword(tok) {
			return p.unexpected(lit, "enum field identifier", f)
		}
	}
	f.Name = lit
	pos, tok, lit := p.next()
	if tok != tEQUALS {
		return p.unexpected(lit, "enum field =", f)
	}
	i, err := p.nextInteger()
	if err != nil {
		return p.unexpected(lit, "enum field integer", f)
	}
	f.Integer = i
	pos, tok, lit = p.next()
	if tok == tLEFTSQUARE {
		o := new(Option)
		o.Position = pos
		o.IsEmbedded = true
		err := o.parse(p)
		if err != nil {
			return err
		}
		f.ValueOption = o
		pos, tok, lit = p.next()
		if tok != tRIGHTSQUARE {
			return p.unexpected(lit, "option closing ]", f)
		}
	}
	if tSEMICOLON == tok {
		p.nextPut(pos, tok, lit) // put back this token for scanning inline comment
	}
	return nil
}
