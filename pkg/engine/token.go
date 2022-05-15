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
	"net/url"
)

// Token represents a unique identifier for a credit card that can be safely
// stored without having to hold sensitive card information on your own servers.
type Token struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Created  int64   `json:"created"`
	Used     bool    `json:"used"`
	Livemode bool    `json:"livemode"`
	Type     string  `json:"type"`
	Card     *Card   `json:"card"`
}

// TokenClient encapsulates operations for creating and querying tokens using
// the Bhojpur Subscription REST API.
type TokenClient struct{}

// TokenParams encapsulates options for creating a new Card Token.
type TokenParams struct {
	//Currency string REMOVED! no longer part of the API
	Card *CardParams
}

// Creates a single use token that wraps the details of a credit card.
// This token can be used in place of a credit card hash with any API method.
// These tokens can only be used once: by creating a new charge object, or
// attaching them to a customer.
func (self *TokenClient) Create(params *TokenParams) (*Token, error) {
	token := Token{}
	values := url.Values{} // REMOVED "currency": {params.Currency}}
	appendCardParamsToValues(params.Card, &values)

	err := query("POST", "/v1/tokens", values, &token)
	return &token, err
}

// Retrieves the card token with the given Id.
func (self *TokenClient) Retrieve(id string) (*Token, error) {
	token := Token{}
	path := "/v1/tokens/" + url.QueryEscape(id)
	err := query("GET", path, nil, &token)
	return &token, err
}
