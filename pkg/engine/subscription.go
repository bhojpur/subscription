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

// Subscription Statuses
const (
	SubscriptionTrialing = "trialing"
	SubscriptionActive   = "active"
	SubscriptionPastDue  = "past_due"
	SubscriptionCanceled = "canceled"
	SubscriptionUnpaid   = "unpaid"
)

// Subscriptions represents a recurring charge a customer's card.
type Subscription struct {
	Customer           string `json:"customer"`
	Status             string `json:"status"`
	Plan               *Plan  `json:"plan"`
	Start              int64  `json:"start"`
	EndedAt            Int64  `json:"ended_at"`
	CurrentPeriodStart Int64  `json:"current_period_start"`
	CurrentPeriodEnd   Int64  `json:"current_period_end"`
	TrialStart         Int64  `json:"trial_start"`
	TrialEnd           Int64  `json:"trial_end"`
	CanceledAt         Int64  `json:"canceled_at"`
	CancelAtPeriodEnd  bool   `json:"cancel_at_period_end"`
	Quantity           int64  `json:"quantity"`
}

// SubscriptionClient encapsulates operations for updating and canceling
// customer subscriptions using the Bhojpur Subscription REST API.
type SubscriptionClient struct{}

// SubscriptionParams encapsulates options for updating a Customer's
// subscription.
type SubscriptionParams struct {
	// The identifier of the plan to subscribe the customer to.
	Plan string

	// (Optional) The code of the coupon to apply to the customer if you would
	// like to apply it at the same time as creating the subscription.
	Coupon string

	// (Optional) Flag telling us whether to prorate switching plans during a
	// billing cycle
	Prorate bool

	// (Optional) UTC integer timestamp representing the end of the trial period
	// the customer will get before being charged for the first time. If set,
	// trial_end will override the default trial period of the plan the customer
	// is being subscribed to.
	TrialEnd int64

	// (Optional) A new card to attach to the customer.
	Card *CardParams

	// (Optional) A new card Token to attach to the customer.
	Token string

	// (Optional) The quantity you'd like to apply to the subscription you're creating.
	Quantity int64
}

// Subscribes a customer to a new plan.
func (self *SubscriptionClient) Update(customerId string, params *SubscriptionParams) (*Subscription, error) {
	values := url.Values{"plan": {params.Plan}}

	// set optional parameters
	if len(params.Coupon) != 0 {
		values.Add("coupon", params.Coupon)
	}
	if params.Prorate {
		values.Add("prorate", "true")
	}
	if params.TrialEnd != 0 {
		values.Add("trial_end", strconv.FormatInt(params.TrialEnd, 10))
	}
	if params.Quantity != 0 {
		values.Add("quantity", strconv.FormatInt(params.Quantity, 10))
	}
	// attach a new card, if requested
	if len(params.Token) != 0 {
		values.Add("card", params.Token)
	} else if params.Card != nil {
		appendCardParamsToValues(params.Card, &values)
	}

	s := Subscription{}
	path := "/v1/customers/" + url.QueryEscape(customerId) + "/subscription"
	err := query("POST", path, values, &s)
	return &s, err
}

// Cancels the customer's subscription if it exists. It cancels the
// subscription immediately.
func (self *SubscriptionClient) Cancel(customerId string) (*Subscription, error) {
	s := Subscription{}
	path := "/v1/customers/" + url.QueryEscape(customerId) + "/subscription"
	err := query("DELETE", path, nil, &s)
	return &s, err
}

// Cancels the customer's subscription at the end of the billing period.
func (self *SubscriptionClient) CancelAtPeriodEnd(customerId string) (*Subscription, error) {
	values := url.Values{}
	values.Add("at_period_end", "true")

	s := Subscription{}
	path := "/v1/customers/" + url.QueryEscape(customerId) + "/subscription"
	err := query("DELETE", path, values, &s)
	return &s, err
}
