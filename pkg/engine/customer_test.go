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

// Sample Customers to use when creating, deleting, updating Customer data.
var (
	// Customer with only the required fields
	cust1 = CustomerParams{
		Email: "test1@bhojpur.net",
		Desc:  "a test customer",
	}

	// Customer with all required fields + required credit card fields.
	cust2 = CustomerParams{
		Email:  "test2@bhojpur.net",
		Desc:   "a 2nd test customer",
		Coupon: c1.ID,
		Plan:   p1.ID,
		Card: &CardParams{
			Name:     "Pramila Kumari",
			Number:   "4242424242424242",
			ExpYear:  time.Now().Year() + 1,
			ExpMonth: 1,
		},
	}

	// Another Customer with only the required fields
	cust3 = CustomerParams{
		Email: "test3@bhojpur.net",
		Desc:  "a 3rd test customer",
	}

	// A customer with the required fields + a credit card
	cust4 = CustomerParams{
		Email: "test3@bhojpur.net",
		Desc:  "a 3rd test customer",
		Card: &CardParams{
			Name:     "Sanjay Kumar",
			Number:   "4242424242424242",
			ExpYear:  time.Now().Year() + 1,
			ExpMonth: 1,
		},
	}
)

// TestCreateCustomer will test that we can successfully Create a Customer,
// parse the JSON reponse from Bhojpur subscription, and that all values are
// populated as expected.
func TestCreateCustomer(t *testing.T) {
	// Create the customer, and defer its deletion
	cust, err := Customers.Create(&cust1)
	defer Customers.Delete(cust.ID)

	if err != nil {
		t.Errorf("Expected Customer, got Error %s", err.Error())
	}
	if string(cust.Email) != cust1.Email {
		t.Errorf("Expected Customer Email %s, got %v", cust1.Email, cust.Email)
	}
	if string(cust.Desc) != cust1.Desc {
		t.Errorf("Expected Customer Desc %s, got %v", cust1.Desc, cust.Desc)
	}
}

// TestCreateCustomerToken will test that we can successfully Create a Customer
// using a credit card Token.
func TestCreateCustomerToken(t *testing.T) {

	// Create a Token for the credit card
	token, _ := Tokens.Create(&token1)

	// Create a Charge that uses a Token
	cust := CustomerParams{
		Token: token.ID,
		Desc:  "Customer for site@bhojpur.net",
	}

	// Create the charge
	resp, err := Customers.Create(&cust)
	defer Customers.Delete(resp.ID)
	if err != nil {
		t.Errorf("Expected Create Customer, got Error %s", err.Error())
	}
	if resp.Cards.Count == 0 {
		t.Errorf("Expected Customer Card from Token, got nil")
	}
	// Sanity check to make sure card was attached to customer
	if string(resp.Cards.Data[0].Name) != string(token.Card.Name) {
		t.Errorf("Expected Card Name %s, got %v", token.Card.Name, resp.Cards.Data[0].Name)
	}
}

// TestRetrieveCustomer will test that we can successfully Retrieve a Customer,
// parse the JSON response, and that all values are populated as expected.
func TestRetrieveCustomer(t *testing.T) {

	// setup default plans and coupons, defer deletion
	Plans.Create(&p1)
	Coupons.Create(&c1)
	defer Plans.Delete(p1.ID)
	defer Coupons.Delete(c1.ID)

	// Create the customer, and defer its deletion
	resp, err := Customers.Create(&cust2)
	defer Customers.Delete(resp.ID)
	if err != nil {
		t.Errorf("Expected Customer, got Error %s", err.Error())
		return
	}

	// Retrieve the Customer by Id
	cust, err := Customers.Retrieve(resp.ID)
	if err != nil {
		t.Errorf("Expected Customer, got Error %s", err.Error())
	}
	if string(cust.Email) != cust2.Email {
		t.Errorf("Expected Customer Email %s, got %v", cust2.Email, cust.Email)
	}
	if string(cust.Desc) != cust2.Desc {
		t.Errorf("Expected Customer Desc %s, got %v", cust2.Desc, cust.Desc)
	}
	if cust.Cards.Count == 0 {
		t.Errorf("Expected Credit Card %s, got nil", cust2.Card.Number)
		return
	}

	if string(cust.Cards.Data[0].Name) != cust2.Card.Name {
		t.Errorf("Expected Card Name %s, got %s", cust2.Card.Name, cust.Cards.Data[0].Name)
	}
	if cust.Cards.Data[0].Last4 != "4242" {
		t.Errorf("Expected Card Last4 %d, got %d", "4242", cust.Cards.Data[0].Last4)
	}
	if cust.Cards.Data[0].ExpYear != cust2.Card.ExpYear {
		t.Errorf("Expected Card ExpYear %d, got %d", cust2.Card.ExpYear, cust.Cards.Data[0].ExpYear)
	}
	if cust.Cards.Data[0].ExpMonth != cust2.Card.ExpMonth {
		t.Errorf("Expected Card ExpMonth %d, got %d", cust2.Card.ExpMonth, cust.Cards.Data[0].ExpMonth)
	}
	if cust.Cards.Data[0].Type != Visa {
		t.Errorf("Expected Card Type %s, got %s", Visa, cust.Cards.Data[0].Type)
	}
}

// TestUpdateCustomer will test that we can successfully update a Customer,
// parse the JSON reponse, and verify the updated name was returned.
func TestUpdateCustomer(t *testing.T) {
	// Create the Customer, and defer its deletion
	resp, _ := Customers.Create(&cust1)
	defer Customers.Delete(resp.ID)

	cust, err := Customers.Update(resp.ID, &CustomerParams{Email: "pramila@bhojpur.net"})
	if err != nil {
		t.Errorf("Expected Customer update, got Error %s", err.Error())
	}
	if cust.Email != "pramila@bhojpur.net" {
		t.Errorf("Expected Updated Customer Email")
	}
}

// TestDeleteCustomer will test that we can successfully remove a Customer,
// parse the JSON reponse, and that the deletion flag is captured as a boolean
// value.
func TestDeleteCustomer(t *testing.T) {
	// Create the Customer, and defer its deletion
	resp, _ := Customers.Create(&cust1)
	defer Customers.Delete(resp.ID)

	// let's try to delete the customer
	ok, err := Customers.Delete(resp.ID)
	if err != nil {
		t.Errorf("Expected Customer deletion, got Error %s", err.Error())
	}
	if !ok {
		t.Errorf("Expected Customer deleted true, got false")
	}
}

// TestListCustomers will test that we can successfully retrieve a list of
// Customers, parse the JSON reponse, and that the length of the coupon array
// matches our expectations.
func TestListCustomers(t *testing.T) {

	// create two dummy customers that we can retrieve
	resp1, _ := Customers.Create(&cust1)
	resp2, _ := Customers.Create(&cust3)
	defer Customers.Delete(resp1.ID)
	defer Customers.Delete(resp2.ID)

	// get the list from Bhojpur Subscription
	customers, err := Customers.ListN(2, 0)
	if err != nil {
		t.Errorf("Expected Customer List, got Error %s", err.Error())
	}

	// since we added two dummy customers, we expect the array to be a size of 2
	if len(customers) != 2 {
		t.Errorf("Expected two Customers, got %s", len(customers))
	}
}
