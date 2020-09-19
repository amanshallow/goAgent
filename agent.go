package main

import (
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
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
	tag = "CSC482GoAgent"
	
	// Loggly Customer Token
	os.Setenv("LOGGLY_TOKEN", "7530acc8-7398-4ebd-9d2c-ce0730d87190")
	
	// Instantiate the client
	client := loggly.New(tag)
	
	resp, err := http.Get("https://api.ratesapi.io/api/latest?base=USD")	
	if err != nil {
		client.Send("Error", "HTTP request to Rates API failed. No echo.")
	} else {
		
		client.Send("info", "HTTP request success. No echo. Status: " + resp.Status) 
	}

	// Close the response body
	defer resp.Body.Close()
      	
	// Read contents from body of request.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		client.Send("Error", "Could not read JSON from body. No echo.")
	} else {
		s := strconv.Itoa(len(body))
		client.Send("info", "Successfully read JSON from body. No echo. Body size: " + s + " bytes.")
	}
	// Parse the JSON and copy into struct
	var info Information
	err = json.Unmarshal(body, &info)
	
	if err != nil {
		client.Send("Error", "Could not parse JSON from body into Information struct. No echo.")
	} else {
		client.Send("info", "Sucessfull unmarshal of JSON from response body. No echo.")
	}
	
	
	// Display the formatted output to the console.
	fmt.Println("Currency exchange rates for: " + info.Date)
	fmt.Println("Base currency: " + info.Base)
	fmt.Println("Exchange rates follow:")
	fmt.Print("USD: ")
	fmt.Println(info.Rates.USD)
	fmt.Print("GBP: ")
	fmt.Println(info.Rates.GBP)
	fmt.Print("CAD: ")
	fmt.Println(info.Rates.CAD)
	fmt.Print("INR: ")
	fmt.Println(info.Rates.INR)
	fmt.Print("AUD: ")
	fmt.Println(info.Rates.AUD)
}
