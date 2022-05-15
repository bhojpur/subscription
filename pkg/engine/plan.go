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

// Subscription Plan Intervals
const (
	IntervalSecond  = "second"
	IntervalMinute  = "minute"
	IntervalHour    = "hour"
	IntervalDay     = "day"
	IntervalWeek    = "week"
	IntervalMonth   = "month"
	IntervalQuarter = "quarter"
	IntervalYear    = "year"
	IntervalDecade  = "decade"
	IntervalCentury = "century"
)

// Plan holds details about pricing information for different products and
// feature levels on your site. For example, you might have a INR 10/month plan
// for basic features and a different INR 20/month plan for premium features.
type Plan struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Amount          float64 `json:"amount"`
	Interval        string  `json:"interval"`
	IntervalCount   int     `json:"interval_count"`
	Currency        string  `json:"currency"`
	TrialPeriodDays Int     `json:"trial_period_days"`
	Livemode        bool    `json:"livemode"`
}

// PlanClient encapsulates operations for creating, updating, deleting and
// querying plans using the Bhojpur Subscription REST API.
type PlanClient struct{}

// PlanParams encapsulates options for creating a new Plan.
type PlanParams struct {
	// Unique string of your choice that will be used to identify this plan
	// when subscribing a customer.
	ID string

	// A positive integer in paisa (or 0 for a free plan) representing how much
	// to charge (on a recurring basis)
	Amount float64

	// 3-letter ISO code for currency. Currently, only 'inr' is supported.
	Currency string

	// Specifies billing frequency. Either month or year.
	Interval string

	// Name of the plan, to be displayed on invoices and in the web interface.
	Name string

	// (Optional) Specifies a trial period in (an integer number of) days. If
	// you include a trial period, the customer won't be billed for the first
	// time until the trial period ends. If the customer cancels before the
	// trial period is over, she'll never be billed at all.
	TrialPeriodDays int
}

// Creates a new Plan.
func (self *PlanClient) Create(params *PlanParams) (*Plan, error) {
	plan := Plan{}
	values := url.Values{
		"id":       {params.ID},
		"name":     {params.Name},
		"amount":   {strconv.FormatFloat(params.Amount, 'E', -1, 64)},
		"interval": {params.Interval},
		"currency": {params.Currency},
	}

	// trial_period_days is optional, add if specified
	if params.TrialPeriodDays != 0 {
		values.Add("trial_period_days", strconv.Itoa(params.TrialPeriodDays))
	}

	err := query("POST", "/v1/plans", values, &plan)
	return &plan, err
}

// Retrieves the plan with the given ID.
func (self *PlanClient) Retrieve(id string) (*Plan, error) {
	plan := Plan{}
	path := "/v1/plans/" + url.QueryEscape(id)
	err := query("GET", path, nil, &plan)
	return &plan, err
}

// Updates the name of a plan. Other plan details (price, interval, etc) are,
// by design, not editable.
func (self *PlanClient) Update(id string, newName string) (*Plan, error) {
	values := url.Values{"name": {newName}}
	plan := Plan{}
	path := "/v1/plans/" + url.QueryEscape(id)
	err := query("POST", path, values, &plan)
	return &plan, err
}

// Deletes a plan with the given ID.
func (self *PlanClient) Delete(id string) (bool, error) {
	resp := DeleteResp{}
	path := "/v1/plans/" + url.QueryEscape(id)
	if err := query("DELETE", path, nil, &resp); err != nil {
		return false, err
	}
	return resp.Deleted, nil
}

// Returns a list of your Plans.
func (self *PlanClient) List() ([]*Plan, error) {
	return self.ListN(10, 0)
}

// Returns a list of your Plans at the specified range.
func (self *PlanClient) ListN(count int, offset int) ([]*Plan, error) {
	// define a wrapper function for the Plan List, so that we can
	// cleanly parse the JSON
	type listPlanResp struct{ Data []*Plan }
	resp := listPlanResp{}

	// add the count and offset to the list of url values
	values := url.Values{
		"count":  {strconv.Itoa(count)},
		"offset": {strconv.Itoa(offset)},
	}

	err := query("GET", "/v1/plans", values, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
