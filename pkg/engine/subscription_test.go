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

// Sample Subscriptions to use for testing
var (

	// Subscriptions with only the required fields
	sub1 = SubscriptionParams{
		Plan: "plan1",
	}

	// Subscriptions with all fields, plus new Credit Card
	sub2 = SubscriptionParams{
		Plan:     "plan1",
		Coupon:   "test coupon 1",
		Prorate:  true,
		TrialEnd: time.Now().Unix() + 1000000,
		Quantity: 5,
		Card: &CardParams{
			Name:     "Pramila Kumari",
			Number:   "4242424242424242",
			ExpYear:  time.Now().Year() + 1,
			ExpMonth: 6,
		},
	}
)

func TestUpdateSubscription(t *testing.T) {
	// Create the customer, and defer its deletion
	cust, _ := Customers.Create(&cust1)
	defer Customers.Delete(cust.ID)

	// Create the plan, and defer its deletion
	Plans.Create(&p1)
	defer Customers.Delete(p1.ID)

	// Subscribe the Customer to the Plan
	resp, err := Subscriptions.Update(cust.ID, &sub1)
	if err != nil {
		t.Errorf("Expected Subscription, got error %s", err.Error())
	}
	if resp.Customer != cust.ID {
		t.Errorf("Expected Customer %s, got %s", cust.ID, resp.Customer)
	}
	if resp.Status != SubscriptionActive {
		t.Errorf("Expected Active Subscription, got %s", resp.Status)
	}
}

func TestUpdateSubscriptionCard(t *testing.T) {

	// Create the customer, and defer its deletion
	cust, _ := Customers.Create(&cust1)
	defer Customers.Delete(cust.ID)
	if cust.Cards.Count != 0 {
		t.Errorf("Expected Customer to be created with a nil card")
		return
	}

	// Create the plan, and defer its deletion
	Plans.Create(&p1)
	defer Customers.Delete(p1.ID)

	// Create the coupon, and defer its deletion
	Coupons.Create(&c1)
	defer Coupons.Delete(c1.ID)

	// Subscribe a Customer to a new plan, using a new Credit Card
	resp, err := Subscriptions.Update(cust.ID, &sub2)
	if err != nil {
		t.Errorf("Expected Subscription, got error %s", err.Error())
	}
	if resp.Quantity != sub2.Quantity {
		t.Errorf("Expected Quantity %d, got %d", sub2.Quantity, resp.Quantity)
	}

	// Check to see if the customer's card was added
	cust, _ = Customers.Retrieve(cust.ID)
	if cust.Cards.Count == 0 {
		t.Errorf("Expected Subscription to assign a new active customer card")
	}
}

func TestUpdateSubscriptionToken(t *testing.T) {
	// Create the customer, and defer its deletion
	cust, _ := Customers.Create(&cust1)
	defer Customers.Delete(cust.ID)
	if cust.Cards.Count != 0 {
		t.Errorf("Expected Customer to be created with a nil card")
		return
	}

	// Create the plan, and defer its deletion
	Plans.Create(&p1)
	defer Customers.Delete(p1.ID)

	// Create a Token for the credit card
	token, _ := Tokens.Create(&token1)

	// Subscribe the Customer to the Plan, using the Token
	params := SubscriptionParams{Plan: "plan1", Token: token.ID}
	_, err := Subscriptions.Update(cust.ID, &params)
	if err != nil {
		t.Errorf("Expected Subscription with Token, got error %s", err.Error())
	}

	// Check to see if the customer's card was added
	cust, _ = Customers.Retrieve(cust.ID)
	if cust.Cards.Count == 0 {
		t.Errorf("Expected Subscription to assign a new active customer card")
	}
}

func TestCancelSubscription(t *testing.T) {
	// Create the customer, and defer its deletion
	cust, _ := Customers.Create(&cust1)
	defer Customers.Delete(cust.ID)

	// Create the plan, and defer its deletion
	Plans.Create(&p1)
	defer Customers.Delete(p1.ID)

	// Subscribe the Customer to the Plan
	_, err := Subscriptions.Update(cust.ID, &sub1)
	if err != nil {
		t.Errorf("Expected Subscription, got error %s", err.Error())
	}

	// Now cancel the subscription
	subs, err := Subscriptions.Cancel(cust.ID)
	if err != nil {
		t.Errorf("Expected Subscription Cancellation, got error %s", err.Error())
	}

	if subs.Status != SubscriptionCanceled {
		t.Errorf("Expected Subscription Status %s, got %s", SubscriptionCanceled, subs.Status)
	}
}

func TestCancelSubscriptionAtPeriodEnd(t *testing.T) {
	// Create the customer, and defer its deletion
	cust, _ := Customers.Create(&cust1)
	defer Customers.Delete(cust.ID)

	// Create the plan, and defer its deletion
	Plans.Create(&p1)
	defer Customers.Delete(p1.ID)

	// Subscribe the Customer to the Plan
	_, err := Subscriptions.Update(cust.ID, &sub1)
	if err != nil {
		t.Errorf("Expected Subscription, got error %s", err.Error())
	}

	// Now cancel the subscription
	subs, err := Subscriptions.CancelAtPeriodEnd(cust.ID)
	if err != nil {
		t.Errorf("Expected Subscription Cancellation, got error %s", err.Error())
	}

	if subs.Status != SubscriptionActive {
		t.Errorf("Expected Subscription Status %s, got %s", SubscriptionCanceled, subs.Status)
	}

	if subs.CancelAtPeriodEnd != true {
		t.Errorf("Expected CancelAtPeriodEnd to be %s, got %s", true, subs.CancelAtPeriodEnd)
	}
}
