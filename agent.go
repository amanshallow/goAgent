package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws" // AWS Dynamo DB
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	loggly "github.com/jamespearly/loggly" // Loggly
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Information struct {
	Date string `json:"date"`
	Base string `json:"base"`

	Currency struct {
		USD float64 `json:"USD"`
		GBP float64 `json:"GBP"`
		INR float64 `json:"INR"`
		CAD float64 `json:"CAD"`
		AUD float64 `json:"AUD"`
	} `json:"rates"`
}

// Necessary to historical data.
/*
var globalIndex = 1
var currentMonth int
var currentDate  int
var currentYear  int
var fullDate 	  string*/

// Decides the agent's polling interval.
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
			client.Send("info", "Polling rate:"+varDelay+"s. No echo.")
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
	for range tick.C {
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
	defer close(done)

	// Tag for Loggly.
	var tag string
	tag = "CSC482GoAgent"

	// Instantiate the client
	client := loggly.New(tag)

	// Algorithm to fetch historical data when needed.
	/*   if globalIndex == 1 {
		currentMonth = 10
		currentDate  = 21
		currentYear  = 2019
		fullDate     = ""
		globalIndex++
	}

	if currentMonth == 12 && currentDate == 31{
		currentMonth = 1
	} else if currentDate == 31 {
		currentMonth++;
	}
	if currentDate == 31 {
		currentDate = 1
	} else {
		currentDate++
	}
	if currentYear == 2019 && currentDate == 31 && currentMonth == 12 {
		currentYear = 2020
	}
	fullDate = strconv.Itoa(currentYear) + "-" + strconv.Itoa(currentMonth) + "-" + strconv.Itoa(currentDate)
	var apiString string = "https://api.ratesapi.io/api/" + fullDate + "?base=USD"
	resp, err := http.Get(apiString)*/

	// Run worker as a Go Routine.
	go worker(done, client)

	// Create an AWS session for US East 1.
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("https://dynamodb.us-east-1.amazonaws.com"),
	}))

	// Create a DynamoDB instance
	db := dynamodb.New(sess)

	resp, err := http.Get("https://api.ratesapi.io/api/latest?base=USD")
	if err != nil {
		client.Send("error", "HTTP request to Rates API failed. No echo.")
	} else {

		client.Send("info", "HTTP request success. No echo. Status: "+resp.Status)
	}

	// Close the response body
	defer resp.Body.Close()

	// Read contents from body of request.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		client.Send("error", "Could not read JSON from body. No echo.")
	} else {
		s := strconv.Itoa(len(body))
		client.Send("info", "Successfully read JSON from body. No echo. Body size: "+s+" bytes.")
	}

	// Parse the JSON and copy into struct
	var info Information
	err = json.Unmarshal(body, &info)

	if err != nil {
		client.Send("error", "Could not parse JSON from body into Information struct. No echo.")
	} else {
		client.Send("info", "Sucessfull unmarshal of JSON from response body. No echo.")
	}

	// Marshal data from Information struct into AWS attribute value.
	infoAVmap, err := dynamodbattribute.MarshalMap(info)
	if err != nil {
		client.Send("error", "Could not marshal information struct into attribute value map. No echo.")
	}

	// Create the api parameters
	params := &dynamodb.PutItemInput{
		TableName: aws.String("asingh2-rates"),
		Item:      infoAVmap,
	}

	// Push or Put the item into the table, no error checking here!
	db.PutItem(params)

	// Display the formatted output to the console.
	fmt.Println("Currency exchange rates for: " + info.Date)
	fmt.Println("Base currency: " + info.Base)
	fmt.Println("Exchange rates follow:")
	fmt.Print("USD: ")
	fmt.Println(info.Currency.USD)
	fmt.Print("GBP: ")
	fmt.Println(info.Currency.GBP)
	fmt.Print("CAD: ")
	fmt.Println(info.Currency.CAD)
	fmt.Print("INR: ")
	fmt.Println(info.Currency.INR)
	fmt.Print("AUD: ")
	fmt.Println(info.Currency.AUD)
	fmt.Print("\n")
	fmt.Print("\n")

	<-done
}
