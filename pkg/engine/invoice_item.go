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

// InvoiceItem represents a charge (or credit) that should be applied to the
// customer at the end of a billing cycle.
type InvoiceItem struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Customer string  `json:"customer"`
	Date     int64   `json:"date"`
	Desc     String  `json:"description"`
	Invoice  String  `json:"invoice"`
	Livemode bool    `json:"livemode"`
}

// InvoiceItemParams encapsulates options for creating a new Invoice Items.
type InvoiceItemParams struct {
	// The ID of the customer who will be billed when this invoice item is
	// billed.
	Customer string

	// The decimal amount in paisa of the charge to be applied to the upcoming
	// invoice. If you want to apply a credit to the customer's account, pass a
	// negative amount.
	Amount float64

	// 3-letter ISO code for currency. Currently, only 'inr' is supported.
	Currency string

	// (Optional) An arbitrary string which you can attach to the invoice item.
	// The description is displayed in the invoice for easy tracking.
	Desc string

	// (Optional) The ID of an existing invoice to add this invoice item to.
	// When left blank, the invoice item will be added to the next upcoming
	// scheduled invoice.
	Invoice string
}

// InvoiceItemClient encapsulates operations for creating, updating, deleting
// and querying invoices using the Bhojpur Subscription REST API.
type InvoiceItemClient struct{}

// Create adds an arbitrary charge or credit to the customer's upcoming invoice.
func (self *InvoiceItemClient) Create(params *InvoiceItemParams) (*InvoiceItem, error) {
	item := InvoiceItem{}
	values := url.Values{
		"amount":   {strconv.FormatFloat(params.Amount, 'E', -1, 64)},
		"currency": {params.Currency},
		"customer": {params.Customer},
	}

	// add optional parameters
	if len(params.Desc) != 0 {
		values.Add("description", params.Desc)
	}
	if len(params.Invoice) != 0 {
		values.Add("invoice", params.Invoice)
	}

	err := query("POST", "/v1/invoiceitems", values, &item)
	return &item, err
}

// Retrieves the Invoice Item with the given ID.
func (self *InvoiceItemClient) Retrieve(id string) (*InvoiceItem, error) {
	item := InvoiceItem{}
	path := "/v1/invoiceitems/" + url.QueryEscape(id)
	err := query("GET", path, nil, &item)
	return &item, err
}

// Update changes the amount or description of an Invoice Item on an upcoming
// invoice, using the given Invoice Item ID.
func (self *InvoiceItemClient) Update(id string, params *InvoiceItemParams) (*InvoiceItem, error) {
	item := InvoiceItem{}
	values := url.Values{}

	if len(params.Desc) != 0 {
		values.Add("description", params.Desc)
	}
	if params.Amount != 0 {
		values.Add("invoice", strconv.FormatFloat(params.Amount, 'E', -1, 64))
	}

	err := query("POST", "/v1/invoiceitems/"+url.QueryEscape(id), values, &item)
	return &item, err
}

// Removes an Invoice Item with the given ID.
func (self *InvoiceItemClient) Delete(id string) (bool, error) {
	resp := DeleteResp{}
	path := "/v1/invoiceitems/" + url.QueryEscape(id)
	if err := query("DELETE", path, nil, &resp); err != nil {
		return false, err
	}
	return resp.Deleted, nil
}

// Returns a list of Invoice Items.
func (self *InvoiceItemClient) List() ([]*InvoiceItem, error) {
	return self.list("", 10, 0)
}

// Returns a list of Invoice Items at the specified range.
func (self *InvoiceItemClient) ListN(count int, offset int) ([]*InvoiceItem, error) {
	return self.list("", count, offset)
}

// Returns a list of Invoice Items for the specified Customer ID.
func (self *InvoiceItemClient) CustomerList(id string) ([]*InvoiceItem, error) {
	return self.list(id, 10, 0)
}

// Returns a list of Invoice Items for the specified Customer ID, at the
// specified range.
func (self *InvoiceItemClient) CustomerListN(id string, count int, offset int) ([]*InvoiceItem, error) {
	return self.list(id, count, offset)
}

func (self *InvoiceItemClient) list(id string, count int, offset int) ([]*InvoiceItem, error) {
	// define a wrapper function for the Invoice Items List, so that we can
	// cleanly parse the JSON
	type listInvoiceItemsResp struct{ Data []*InvoiceItem }
	resp := listInvoiceItemsResp{}

	// add the count and offset to the list of url values
	values := url.Values{
		"count":  {strconv.Itoa(count)},
		"offset": {strconv.Itoa(offset)},
	}

	// query for customer id, if provided
	if id != "" {
		values.Add("customer", id)
	}

	err := query("GET", "/v1/invoiceitems", values, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
