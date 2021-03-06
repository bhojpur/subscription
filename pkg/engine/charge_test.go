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
	// In order to execute Unit Test, you must set your Bhojpur API Key as
	// environment variable, BHOJPUR_API_KEY=xxxx
	if err := SetKeyEnv(); err != nil {
		panic(err)
	}
}

// Sample Charges to use when creating, deleting, updating Charge data.
var (

	// Charge with only the required fields
	charge1 = ChargeParams{
		Desc:     "Litti Chokha",
		Amount:   200,
		Currency: INR,
		Card: &CardParams{
			Name:     "Pramila Kumari",
			Number:   "4242424242424242",
			ExpYear:  time.Now().Year() + 1,
			ExpMonth: 5,
		},
	}
)

// TestCreateCharge will test that we can successfully Charge a credit card,
// parse the JSON reponse from Bhojpur Subscription, and that all values are
// populated as expected.
func TestCreateCharge(t *testing.T) {

	// Create the charge
	resp, err := Charges.Create(&charge1)

	if err != nil {
		t.Errorf("Expected Successful Charge, got Error %s", err.Error())
	}
	if string(resp.Desc) != charge1.Desc {
		t.Errorf("Expected Charge Desc %s, got %s", charge1.Desc, resp.Desc)
	}
	if resp.Amount != charge1.Amount {
		t.Errorf("Expected Charge Amount %v, got %v", charge1.Amount, resp.Amount)
	}
	if resp.Card == nil {
		t.Errorf("Expected Charge Response to include the Charged Credit Card")
		return
	}
	if resp.Paid != true {
		t.Errorf("Expected Charge was paid, got %v", resp.Paid)
	}
}

// TestCreateChargeToken attempts to charge using a Card Token.
func TestCreateChargeToken(t *testing.T) {

	// Create a Token for the credit card
	token, err := Tokens.Create(&token1)
	if err != nil {
		t.Errorf("Expected Token Creation, got Error %s", err.Error())
	}

	// Create a Charge that uses a Token
	charge := ChargeParams{
		Desc:     "Litti Chokha",
		Amount:   400,
		Currency: INR,
		Token:    token.ID,
	}

	// Create the charge
	_, err = Charges.Create(&charge)
	if err != nil {
		t.Errorf("Expected Successful Charge, got Error %s", err.Error())
	}
}

// TestCreateChargeCustomer attempts to charge a pre-defined customer, meaning
// we don't specify the credit card or token when Creating the charge.
func TestCreateChargeCustomer(t *testing.T) {

	// Create a Customer and defer deletion
	// This customer should have a credit card setup
	cust, _ := Customers.Create(&cust4)
	defer Customers.Delete(cust.ID)
	if cust.Cards.Count == 0 {
		t.Errorf("Cannot test charging a customer with no pre-defined Card")
		return
	}

	// Create a Charge that uses a Token
	charge := ChargeParams{
		Desc:     "Litti Chokha",
		Amount:   200,
		Currency: INR,
		Customer: cust.ID,
	}

	// Create the charge
	_, err := Charges.Create(&charge)
	if err != nil {
		t.Errorf("Expected Successful Charge, got Error %s", err.Error())
	}
}

func TestRetrieveCharge(t *testing.T) {
	// Create the charge
	resp, err := Charges.Create(&charge1)
	if err != nil {
		t.Errorf("Expected Successful Charge, got Error %s", err.Error())
		return
	}

	// Retrieve the charge from the database
	_, err = Charges.Retrieve(resp.ID)
	if err != nil {
		t.Errorf("Expected to retrieve Charge by Id, got Error %s", err.Error())
		return
	}
}

func TestRefundCharge(t *testing.T) {
	// Create the charge
	resp, err := Charges.Create(&charge1)
	if err != nil {
		t.Errorf("Expected Successful Charge, got Error %s", err.Error())
		return
	}

	// Refund the full amount
	charge, err := Charges.Refund(resp.ID)
	if err != nil {
		t.Errorf("Expected Refund, got Error %s", err.Error())
		return
	}

	if charge.Refunded == false {
		t.Errorf("Expected Refund, however Refund flag was set to false")
	}
	if float64(charge.AmountRefunded) != charge1.Amount {
		t.Errorf("Expected AmountRefunded %v, but got %v", charge1.Amount, int64(charge.AmountRefunded))
		return
	}
}
