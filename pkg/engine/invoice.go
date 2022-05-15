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

// Invoice represents statements of what a customer owes for a particular
// billing period, including subscriptions, invoice items, and any automatic
// proration adjustments if necessary.
type Invoice struct {
	ID              string        `json:"id"`
	AmountDue       float64       `json:"amount_due"`
	AttemptCount    int           `json:"attempt_count"`
	Attempted       bool          `json:"attempted"`
	Closed          bool          `json:"closed"`
	Paid            bool          `json:"paid"`
	PeriodEnd       int64         `json:"period_end"`
	PeriodStart     int64         `json:"period_start"`
	Subtotal        float64       `json:"subtotal"`
	Total           float64       `json:"total"`
	Charge          String        `json:"charge"`
	Customer        string        `json:"customer"`
	Date            int64         `json:"date"`
	Discount        *Discount     `json:"discount"`
	Lines           *InvoiceLines `json:"lines"`
	StartingBalance float64       `json:"starting_balance"`
	EndingBalance   float64       `json:"ending_balance"`
	NextPayment     float64       `json:"next_payment_attempt"`
	Livemode        bool          `json:"livemode"`
}

// InvoiceLines represents an individual line items that is part of an invoice.
type InvoiceLines struct {
	InvoiceItems  []*InvoiceItem      `json:"invoiceitems"`
	Prorations    []*InvoiceItem      `json:"prorations"`
	Subscriptions []*SubscriptionItem `json:"subscriptions"`
}

type SubscriptionItem struct {
	Amount float64 `json:"amount"`
	Period *Period `json:"period"`
	Plan   *Plan   `json:"plan"`
}

type Period struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// InvoiceClient encapsulates operations for querying invoices using the
// Bhojpur Subscription REST API.
type InvoiceClient struct{}

// Retrieves the invoice with the given ID.
func (self *InvoiceClient) Retrieve(id string) (*Invoice, error) {
	invoice := Invoice{}
	path := "/v1/invoices/" + url.QueryEscape(id)
	err := query("GET", path, nil, &invoice)
	return &invoice, err
}

// Retrieves the upcoming invoice the given customer ID.
func (self *InvoiceClient) RetrieveCustomer(cid string) (*Invoice, error) {
	invoice := Invoice{}
	values := url.Values{"customer": {cid}}
	err := query("GET", "/v1/invoices/upcoming", values, &invoice)
	return &invoice, err
}

// Returns a list of Invoices.
func (self *InvoiceClient) List() ([]*Invoice, error) {
	return self.list("", 10, 0)
}

// Returns a list of Invoices at the specified range.
func (self *InvoiceClient) ListN(count int, offset int) ([]*Invoice, error) {
	return self.list("", count, offset)
}

// Returns a list of Invoices with the given Customer ID.
func (self *InvoiceClient) CustomerList(id string) ([]*Invoice, error) {
	return self.list(id, 10, 0)
}

// Returns a list of Invoices with the given Customer ID, at the specified range.
func (self *InvoiceClient) CustomerListN(id string, count int, offset int) ([]*Invoice, error) {
	return self.list(id, count, offset)
}

func (self *InvoiceClient) list(id string, count int, offset int) ([]*Invoice, error) {
	// define a wrapper function for the Invoice List, so that we can
	// cleanly parse the JSON
	type listInvoicesResp struct{ Data []*Invoice }
	resp := listInvoicesResp{}

	// add the count and offset to the list of url values
	values := url.Values{
		"count":  {strconv.Itoa(count)},
		"offset": {strconv.Itoa(offset)},
	}

	// query for customer id, if provided
	if id != "" {
		values.Add("customer", id)
	}

	err := query("GET", "/v1/invoices", values, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
