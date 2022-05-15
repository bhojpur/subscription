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
	"strconv"
)

// ISO 3-digit Currency Codes for major currencies (not the full list).
const (
	INR = "inr" // Indian Rupee (₹)
	USD = "usd" // US Dollar ($)
	EUR = "eur" // Euro (€)
	GBP = "gbp" // British Pound Sterling (UK£)
	JPY = "jpy" // Japanese Yen (¥)
	CAD = "cad" // Canadian Dollar (CA$)
	HKD = "hkd" // Hong Kong Dollar (HK$)
	CNY = "cny" // Chinese Yuan (CN¥)
	AUD = "aud" // Australian Dollar (A$)
)

// Charge represents details about a credit card charge in Bhojpur subscription.
type Charge struct {
	ID                   string        `json:"id"`
	Desc                 String        `json:"description"`
	Amount               float64       `json:"amount"`
	Card                 *Card         `json:"card"`
	Currency             string        `json:"currency"`
	Created              int64         `json:"created"`
	Customer             String        `json:"customer"`
	Invoice              String        `json:"invoice"`
	Fee                  float64       `json:"fee"`
	Paid                 bool          `json:"paid"`
	Details              []*FeeDetails `json:"fee_details"`
	Refunded             bool          `json:"refunded"`
	AmountRefunded       float64       `json:"amount_refunded"`
	FailureMessage       String        `json:"failure_message"`
	Disputed             bool          `json:"disputed"`
	Livemode             bool          `json:"livemode"`
	StatementDescription string        `json:"statement_description"`
}

// FeeDetails represents a single fee associated with a Charge.
type FeeDetails struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Type        string  `json:"type"`
	Application String  `json:"application"`
}

// ChargeParams encapsulates options for creating a new Charge.
type ChargeParams struct {
	// A positive integer in paisa representing how much to charge the card.
	// The minimum amount is 50 paisa.
	Amount float64

	// 3-letter ISO code for currency.
	Currency string

	// (Optional) Either customer or card is required, but not both The ID of an
	// existing customer that will be charged in this request.
	Customer string

	// (Optional) Credit Card that should be charged.
	Card *CardParams

	// (Optional) Credit Card token that should be charged.
	Token string

	// An arbitrary string which you can attach to a charge object. It is
	// displayed when in the web interface alongside the charge. It's often a
	// good idea to use an email address as a description for tracking later.
	Desc string

	// An arbitrary string to be displayed alongside your company name on your
	// customer's credit card statement. This may be up to 15 characters. As an
	// example, if your website is RunClub and you specify 5K Race Ticket, the
	// user will see:
	//     RUNCLUB 5K RACE TICKET.
	// The statement description may not include <>"' characters. While most
	// banks display this information consistently, some may display it
	// incorrectly or not at all.
	StatementDescription string
}

// ChargeClient encapsulates operations for creating, updating, deleting and
// querying charges using the Bhojpur Subscription REST API.
type ChargeClient struct{}

// Creates a new credit card Charge.
func (self *ChargeClient) Create(params *ChargeParams) (*Charge, error) {
	charge := Charge{}
	values := url.Values{
		"amount":      {strconv.FormatFloat(params.Amount, 'E', -1, 64)},
		"currency":    {params.Currency},
		"description": {params.Desc},
	}

	// add optional credit card details, if specified
	if params.Card != nil {
		appendCardParamsToValues(params.Card, &values)
	} else if len(params.Token) > 0 {
		values.Add("card", params.Token)
	} else {
		// if no credit card is provide we need to specify the customer
		values.Add("customer", params.Customer)
	}

	// add optional statment description, if specified
	if params.StatementDescription != "" {
		values.Add("statement_description", params.StatementDescription)
	}

	err := query("POST", "/v1/charges", values, &charge)
	return &charge, err
}

// Retrieves the details of a charge with the given ID.
func (self *ChargeClient) Retrieve(id string) (*Charge, error) {
	charge := Charge{}
	path := "/v1/charges/" + url.QueryEscape(id)
	err := query("GET", path, nil, &charge)
	return &charge, err
}

// Refunds a charge for the full amount.
func (self *ChargeClient) Refund(id string) (*Charge, error) {
	values := url.Values{}
	charge := Charge{}
	path := "/v1/charges/" + url.QueryEscape(id) + "/refund"
	err := query("POST", path, values, &charge)
	return &charge, err
}

// Refunds a charge for the specified amount.
func (self *ChargeClient) RefundAmount(id string, amt float64) (*Charge, error) {
	values := url.Values{
		"amount": {strconv.FormatFloat(amt, 'E', -1, 64)},
	}
	charge := Charge{}
	path := "/v1/charges/" + url.QueryEscape(id) + "/refund"
	err := query("POST", path, values, &charge)
	return &charge, err
}

// Returns a list of your Charges.
func (self *ChargeClient) List() ([]*Charge, error) {
	return self.list("", 10, 0)
}

// Returns a list of your Charges with the specified range.
func (self *ChargeClient) ListN(count int, offset int) ([]*Charge, error) {
	return self.list("", count, offset)
}

// Returns a list of your Charges with the given Customer ID.
func (self *ChargeClient) CustomerList(id string) ([]*Charge, error) {
	return self.list(id, 10, 0)
}

// Returns a list of your Charges with the given Customer ID and range.
func (self *ChargeClient) CustomerListN(id string, count int, offset int) ([]*Charge, error) {
	return self.list(id, count, offset)
}

func (self *ChargeClient) list(id string, count int, offset int) ([]*Charge, error) {
	// define a wrapper function for the Charge List, so that we can
	// cleanly parse the JSON
	type listChargesResp struct{ Data []*Charge }
	resp := listChargesResp{}

	// add the count and offset to the list of url values
	values := url.Values{
		"count":  {strconv.Itoa(count)},
		"offset": {strconv.Itoa(offset)},
	}

	// query for customer id, if provided
	if id != "" {
		values.Add("customer", id)
	}

	err := query("GET", "/v1/charges", values, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
