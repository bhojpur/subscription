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
	"time"
)

func init() {
	// In order to execute Unit Test, you must set your Bhojpur Subscription API Key
	// as environment variable, BHOJPUR_API_KEY=xxxx
	if err := SetKeyEnv(); err != nil {
		panic(err)
	}
}

var (
	// These cards will be successfully charged.
	goodCards = []string{
		"4242424242424242",
		"4012888888881881",
		"5555555555554444",
		"5105105105105100",
		"378282246310005",
		"371449635398431",
		"6011111111111117",
		"6011000990139424",
		"30569309025904",
		"38520000023237",
		"3530111333300000",
		"3566002020360505",
		"4000000000000010",
		"4000000000000028",
		"4000000000000036",
		"4000000000000044",
		"4000000000000101",
	}
	// "These cards will produce specific responses that are useful for testing different scenarios"
	badCardsAndErrorCodes = map[string]string{
		"4000000000000341": ErrCodeCardDeclined,
		"4000000000000002": ErrCodeCardDeclined,
		"4000000000000127": ErrCodeIncorrectCVC,
		"4000000000000069": ErrCodeExpiredCard,
		"4000000000000119": ErrCodeProcessingError,
	}
	// Charge with only the required fields
	charge = ChargeParams{
		Desc:     "Litti Chokha",
		Amount:   300,
		Currency: INR,
		Card: &CardParams{
			Name: "Pramila Kumari",
			//Number:   "", // This gets changed per-test
			ExpYear:  time.Now().Year() + 1,
			ExpMonth: 5,
		},
	}
)

// TestGoodCards ensures we can charge all of Bhojpur Subscription's "good" test cards.
func TestGoodCards(t *testing.T) {
	for _, cardNumber := range goodCards {
		charge.Card.Number = cardNumber
		if _, err := Charges.Create(&charge); err != nil {
			t.Errorf("Expected Successful Charge, got Error %s", err.Error())
		}
	}
}

// TestBadCards ensures we can't charge any of Bhojpur Subscription's "bad" test cards,
// and that the resulting error types and codes are correctly mapped.
func TestBadCards(t *testing.T) {
	for cardNumber, errCode := range badCardsAndErrorCodes {
		charge.Card.Number = cardNumber
		_, err := Charges.Create(&charge)
		bhojpurErr := err.(*Error)
		if bhojpurErr.Detail.Type != ErrTypeCard {
			t.Errorf("Expected Error Type %s, got %s", ErrTypeCard, bhojpurErr.Detail.Type)
		}
		if bhojpurErr.Detail.Code != errCode {
			t.Errorf("Expected Error Code %s, got %s", errCode, bhojpurErr.Detail.Code)
		}
	}
}
