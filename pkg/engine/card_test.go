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
	"testing"
)

type card struct {
	Number string
	Type   string
	Valid  bool
}

var cards = []*card{
	&card{"4242424242424242", Visa, true},            // should pass
	&card{"4213729238347292", Visa, false},           // should fail
	&card{"79927398713", UnknownCard, true},          // should pass
	&card{"79927398710", UnknownCard, false},         // should fail
	&card{"601134239348202", Discover, false},        // should fail
	&card{"344347386473833", AmericanExpress, false}, // should fail
	&card{"374347386473833", AmericanExpress, false}, // should fail
	&card{"361134239348202", DinersClub, false},      // should fail
	&card{"300134239348202", DinersClub, false},      // should fail
	&card{"521134239348202", MasterCard, false},      // should fail
	&card{"380134239348202", JCB, false},             // should fail
	&card{"180034239348202", JCB, false},             // should fail
}

func TestLuhn(t *testing.T) {
	for _, card := range cards {
		valid, _ := IsLuhnValid(card.Number)
		cardType := GetCardType(card.Number)

		if valid != card.Valid {
			t.Errorf("card validation [%v]; want [%v]", valid, card.Valid)
		}
		if cardType != card.Type {
			t.Errorf("card type [%s]; want [%s]", cardType, card.Type)
		}
	}
}
