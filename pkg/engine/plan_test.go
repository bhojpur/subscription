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
	"strings"
	"testing"
)

func init() {
	// In order to execute Unit Test, you must set your Bhojpur Subscription API Key
	// as environment variable, BHOJPUR_API_KEY=xxxx
	if err := SetKeyEnv(); err != nil {
		panic(err)
	}
}

// Sample Plans to use when creating, deleting, updating Plan data.
var (
	// Plan with only the required fields
	p1 = PlanParams{
		ID:       "plan1",
		Name:     "plan 1",
		Amount:   1,
		Currency: INR,
		Interval: IntervalMonth,
	}

	// Plan with all required + optional fields.
	p2 = PlanParams{
		ID:              "plan9",
		Name:            "plan 9",
		Amount:          9,
		Currency:        INR,
		Interval:        IntervalMonth,
		TrialPeriodDays: 365,
	}
)

// TestCreatePlan will test that we can successfully Create a plan, parse the
// JSON reponse from Bhojpur Subscription, and that all values are populated
// as expected.
//
// Second, we will test that error handling works correctly by attempting to
// create a duplicate Plan, create a Plan with invalid currency, which should
// throw exceptions.
func TestCreatePlan(t *testing.T) {

	// Create the plan, and defer its deletion
	plan, err := Plans.Create(&p1)
	defer Plans.Delete(p1.ID)

	if err != nil {
		t.Errorf("Expected Plan %s, got Error %s", p1.ID, err.Error())
	}
	if plan.ID != plan.ID {
		t.Errorf("Expected Plan ID %s, got %s", p1.ID, plan.ID)
	}
	if plan.Name != p1.Name {
		t.Errorf("Expected Plan Name %v, got %v", p1.Name, plan.Name)
	}
	if plan.Amount != p1.Amount {
		t.Errorf("Expected Plan Amount %v, got %v", p1.Amount, plan.Amount)
	}
	if plan.Currency != p1.Currency {
		t.Errorf("Expected Plan Currency %v, got %v", p1.Currency, plan.Currency)
	}

	// Now try to re-create the existing plan, which should throw an exception
	_, err = Plans.Create(&p1)
	if err == nil {
		t.Error("Expected non-null Error when creating a duplicate Plan.")
	} else if err.Error() != "Plan already exists." {
		t.Errorf("Expected %s, got %s", "Plan already exists.", err.Error())
	}

	// Now use an invalid currency, which should throw an exception
	var p3 PlanParams
	p3 = p1
	p3.Currency = "XXX"
	_, err = Plans.Create(&p3)
	if err == nil {
		t.Error("Expected non-null Error when using an Invalid Currency.")
	} else if strings.HasPrefix(err.Error(), "Invalid currency: xxx.") == false {
		t.Errorf("Expected %s, got %s", "Invalid currency: xxx.", err.Error())
	}
}

// TestRetrievePlan will test that we can successfully Retrieve a Plan,
// parse the JSON response, and that all values are populated as expected.
func TestRetrievePlan(t *testing.T) {
	// Create the plan, and defer its deletion
	Plans.Create(&p2)
	defer Plans.Delete(p2.ID)

	// Retrieve the Plan by Id
	plan, err := Plans.Retrieve(p2.ID)
	if err != nil {
		t.Errorf("Expected Plan %s, got Error %s", p2.ID, err.Error())
	}
	if plan.ID != plan.ID {
		t.Errorf("Expected Plan Id %s, got %s", p2.ID, plan.ID)
	}
	if plan.Name != p2.Name {
		t.Errorf("Expected Plan Name %v, got %v", p2.Name, plan.Name)
	}
	if plan.Amount != p2.Amount {
		t.Errorf("Expected Plan Amount %v, got %v", p2.Amount, plan.Amount)
	}
	if plan.Currency != p2.Currency {
		t.Errorf("Expected Plan Currency %v, got %v", p2.Currency, plan.Currency)
	}
	if plan.TrialPeriodDays != Int(p2.TrialPeriodDays) {
		t.Errorf("Expected Plan Trial Period %v, got %v",
			p2.TrialPeriodDays, plan.TrialPeriodDays)
	}
}

// TestUpdatePlan will test that we can successfully update a Plan's name, parse
// the JSON reponse, and verify the updated name was returned.
func TestUpdatePlan(t *testing.T) {
	// Create the plan, and defer its deletion
	Plans.Create(&p1)
	defer Plans.Delete(p1.ID)

	plan, err := Plans.Update(p1.ID, "New Name")
	if err != nil {
		t.Errorf("Expected Plan update, got Error %s", err.Error())
	}
	if plan.Name != "New Name" {
		t.Errorf("Expected Updated Plan Name %v, got %v", p1.Name, plan.Name)
	}
}

// TestDeletePlan will test that we can successfully remove a Plan, parse
// the JSON reponse, and that the deletion flag is captured as a boolean value.
func TestDeletePlan(t *testing.T) {
	// create a Plan that we can delete
	Plans.Create(&p1)

	// let's try to delete the plan
	ok, err := Plans.Delete(p1.ID)
	if err != nil {
		t.Errorf("Expected Plan deletion, got Error %s", err.Error())
	}
	if !ok {
		t.Errorf("Expected Plan deletion true, got false")
	}
}

// TestListPlan will test that we can successfully retrieve a list of Plans,
// parse the JSON reponse, and that the length of the coupon array matches our
// expectations.
func TestListPlan(t *testing.T) {

	// create two dummy plans that we can retrieve
	Plans.Create(&p1)
	Plans.Create(&p2)
	defer Plans.Delete(p1.ID)
	defer Plans.Delete(p2.ID)

	// get the list from Bhojpur Subscription
	plans, err := Plans.List()
	if err != nil {
		t.Errorf("Expected Plan List, got Error %s", err.Error())
	}

	// since we added two dummy plans, we expect the array to be a size of 2
	if len(plans) != 2 {
		t.Errorf("Expected two Plans, got %d", len(plans))
	}
}
