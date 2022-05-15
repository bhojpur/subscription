package engine

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"strconv"
)

// Int is a special type of integer that can unmarshall a JSON value of
// "null", which cannot be parsed by the Go JSON parser as of Go v1.
type Int int

func (self *Int) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		return nil
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		return err
	}

	*self = Int(i)
	return nil
}

// Int64 is a special type of int64 that can unmarshall a JSON value of
// "null", which cannot be parsed by the Go JSON parser as of Go v1.
type Int64 int64

func (self *Int64) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		return nil
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}

	*self = Int64(i)
	return nil
}

// Bool is a special type of bool that can unmarshall a JSON value of
// "null", which cannot be parsed by the Go JSON parser as of Go v1.
type Bool bool

func (self *Bool) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		return nil
	}

	b, err := strconv.ParseBool(str)
	if err != nil {
		return err
	}

	*self = Bool(b)
	return nil
}

// String is a special type of string that can unmarshall a JSON value of
// "null", which cannot be parsed by the Go JSON parser as of Go v1.
type String string

func (self *String) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		return nil
	}

	str, err := strconv.Unquote(str)
	if err != nil {
		return err
	}

	*self = String(str)
	return nil
}
