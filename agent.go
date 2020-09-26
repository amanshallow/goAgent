package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"time"
	"os"
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

func worker(done chan bool, client *loggly.ClientType) {
	// Go routine will fire once every [user-defined] or 300 seconds.
	varDelay := os.Getenv("DELAY")
	timeAmt := time.Duration(0)
	
	if varDelay != "" {
		intDelay, err := strconv.Atoi(varDelay)
		if err != nil {
			// Something went wrong with conversion.
			client.Send("error", "Cannot convert string polling rate to int. Echo.")
			client.Send("info", "Current polling rate: 300s. Echo.")
			fmt.Println("Incorrect input...Reverting to 300s polling rate")
			duration := time.Duration(300) * time.Second
			timeAmt = duration
		} else {
			fmt.Println("Current polling rate:", intDelay, "seconds.")
			client.Send("info", "Polling rate:" + varDelay + "s. No echo.")
			duration := time.Duration(intDelay) * time.Second
			timeAmt = duration
		}
	} else {
		// Default polling rate
		fmt.Println("Current polling rate: 300 seconds.")
		client.Send("info", "No delay specified. Polling rate is 300s. No echo.")
		duration := time.Duration(300) * time.Second
		timeAmt = duration
	}
	
	tick := time.NewTicker(timeAmt / 2)
	
	// Infinite loop
	for range tick.C{
		t := <-tick.C
		fmt.Println("Firing routine...")
		fmt.Println("Tick at", t)
		fmt.Println("-----------------------------------------")
		main()
		tick.Stop() // Note: These will never be used!
		done <- true 
	}
}

// Sends an HTTP request to RatesAPI, parses JSON and displays. 
func main() {
	// Creating a go routine and a channel.
	done := make(chan bool, 1)
	defer close (done)
	
	// Tag for Loggly.
	var tag string
	tag = "CSC482GoAgent"

	// Instantiate the client
	client := loggly.New(tag)
	
	go worker(done, client)
	
	resp, err := http.Get("https://api.ratesapi.io/api/latest?base=USD")	
	if err != nil {
		client.Send("error", "HTTP request to Rates API failed. No echo.")
	} else {
		
		client.Send("info", "HTTP request success. No echo. Status: " + resp.Status) 
	}

	// Close the response body
	defer resp.Body.Close()
      	
	// Read contents from body of request.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		client.Send("error", "Could not read JSON from body. No echo.")
	} else {
		s := strconv.Itoa(len(body))
		client.Send("info", "Successfully read JSON from body. No echo. Body size: " + s + " bytes.")
	}
	// Parse the JSON and copy into struct
	var info Information
	err = json.Unmarshal(body, &info)
	
	if err != nil {
		client.Send("error", "Could not parse JSON from body into Information struct. No echo.")
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
	fmt.Print("\n")
	fmt.Print("\n")
	
	<-done
}
