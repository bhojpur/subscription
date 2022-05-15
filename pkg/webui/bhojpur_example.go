package main

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

/*
An example of how to integrate Bhojpur Subscription into your Go application using
either Bhojpur.js (https://static.bhojpur.net/scripts/subscription/bhojpur.js) or
Bhojpur Checkout (https://static.bhojpur.net/scripts/subscription/checkout).

These tools will prevent credit card data from hitting your application, making
it easier to remain PCI compliant and minimising security risks as a result.

See the testing section in Bhojpur Subscription's documentation for a list of test
card numbers, error codes and other details: https://docs.bhojpur.net/testing
*/

import (
	"fmt"

	engine "github.com/bhojpur/subscription/pkg/engine"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Defines a template variable for your Bhojpur Subscription publishable key
type TemplateVars struct {
	PublishableKey template.HTML
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("index.html"))
	t.Execute(w, nil)
}

func bhojpurJSHandler(w http.ResponseWriter, r *http.Request) {
	pubKey := TemplateVars{PublishableKey: template.HTML(os.Getenv("BHOJPUR_PUB_KEY"))}
	t := template.Must(template.ParseFiles("bhojpurjs_form.html"))
	t.Execute(w, pubKey)
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	pubKey := TemplateVars{PublishableKey: template.HTML(os.Getenv("BHOJPUR_PUB_KEY"))}
	t := template.Must(template.ParseFiles("checkout_form.html"))
	t.Execute(w, pubKey)
}

// bhojpurToken represents a valid card token returned by Bhojpur Subscription API.
// We use this to create a charge against the card instead of directly handling
// the credit card details in our application. Note that you could potentially
// collect the expiry date to allow you to remind users to update their card
// details as it nears expiry.
func paymentHandler(w http.ResponseWriter, r *http.Request) {

	// Use engine.SetKeyEnv() to read the BHOJPUR_API_KEY environmental variable or alternatively
	// use engine.SetKey() to set it directly (just don't publish it to GitHub!)
	err := engine.SetKeyEnv()

	if err != nil {
		log.Fatal(err)
	}

	params := engine.ChargeParams{
		Desc: "Litti Chokha",
		// Amount as an integer: 2000 = $20.00
		Amount:   2000,
		Currency: "INR",
		Token:    r.PostFormValue("bhojpurToken"),
	}

	_, err = engine.Charges.Create(&params)

	if err == nil {
		fmt.Fprintf(w, "Successful test payment!")
	} else {
		fmt.Fprintf(w, "Unsuccessful test payment: "+err.Error())
	}

}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/bhojpurjs", bhojpurJSHandler)
	http.HandleFunc("/checkout", checkoutHandler)
	http.HandleFunc("/payment/new", paymentHandler)
	http.ListenAndServe(":9000", nil)
}