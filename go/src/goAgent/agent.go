package main

import (
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	loggly "github.com/jamespearly/loggly"
)

type Information struct {
	Base string `json:"base"`
	Date string `json:"date"`
	
	Rates struct {
		USD float64 `json:"USD"`
		GBP float64 `json:"GBP"`
		INR float64 `json:"INR"`
		CAD float64 `json:"CAD"`
		AUD float64 `json:"AUD"`
	} `json:"rates"`
}

func main() {
	var tag string
	tag = "My-Go-Agent"
	
	// Loggly Customer Token
	os.Setenv("LOGGLY_TOKEN", "7530acc8-7398-4ebd-9d2c-ce0730d87190")
	
	// Instantiate the client
	client := loggly.New(tag)
	
	resp, err := http.Get("https://api.ratesapi.io/api/latest")	
	if err != nil {
		client.Send("Error", "HTTP request to Rates API failed. No echo.")
	}
	else {
		client.Send("info", "HTTP request success. No echo.")
	}
	defer resp.Body.Close()
	
	// Read contents from body of request.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		client.Send("Error", "Could not read JSON from body. No echo.")
	}
	
	// Parse the JSON and copy into struct
	var info Information
	err = json.Unmarshal(body, &info)
	
	if err != nil {
		client.Send("Error", "Could not parse JSON from body into Information struct.")
	}
	
	
	// Display the formatted output to the console.
	fmt.Println("Currency exchange rates for date: " + info.Date)
	fmt.Println("Base currency: " + info.Base)
	fmt.Println("Exchange rates follow:")
	fmt.Print("US Dollar: ")
	fmt.Println(info.Rates.USD)
	fmt.Print("Great Britain Pound: ")
	fmt.Println(info.Rates.GBP)
	fmt.Print("Canadian Dollar: ")
	fmt.Println(info.Rates.CAD)
	fmt.Print("Indian Rupee: ")
	fmt.Println(info.Rates.INR)
	fmt.Print("Austrailian Dollar: ")
	fmt.Println(info.Rates.AUD)
}
