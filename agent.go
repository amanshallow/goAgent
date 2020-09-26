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


func worker(done chan bool) {
	// Go routine will fire once every 60 seconds.
	varDelay := os.Getenv("DELAY")
	timeAmt := time.Duration(0)
	
	if varDelay != "" {
		intDelay, err := strconv.Atoi(varDelay)
		if err != nil {
			// Something went wrong with conversion.
			fmt.Println("Error:", err)
			fmt.Println("An error has occured...agent will default to polliing every 60 seconds")
			duration := time.Duration(60) * time.Second
			timeAmt = duration
		} else {
			fmt.Println("Current polling rate:", intDelay, "seconds")
			duration := time.Duration(intDelay) * time.Second
			timeAmt = duration
		}
	} else {
		fmt.Println("No time interval specified...polling rate is 60 seconds.")
		duration := time.Duration(60) * time.Second
		timeAmt = duration
	}
	
	tick := time.NewTicker(timeAmt)
	// Infinite loop
	for range tick.C{
		fmt.Println("Firing routine...")
		t := <-tick.C
		fmt.Println("Tick at", t)
		fmt.Println("-----------------------------------------")
		main()
		tick.Stop() // Note: These will never be used!
		done <- true 
	}
}

func main() {
	// Creating a go routine and a channel.
	done := make(chan bool, 1)
	defer close (done)
	go worker(done)
	
	// Tag for Loggly.
	var tag string
	tag = "CSC482GoAgent"

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
	fmt.Print("\n")
	fmt.Print("\n")
	
	<-done
}
