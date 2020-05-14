package main

import (
	"fmt"
	"log"
	"os"

	"github.com/matiss/go-microservice-test/services"
)

const (
	feedURL = "https://www.bank.lv/vk/ecb_rss.xml"
)

// Setup database tables
func setup(mysqlService *services.MySQLService) error {
	fmt.Println("Running setup...")
	err := mysqlService.Setup()
	if err != nil {
		return err
	}

	fmt.Println("Done!")
	return nil
}

// Update currencies
func update(mysqlService *services.MySQLService, currencyFeedService *services.CurrencyFeedService) error {
	// Fetch and parse currency feed
	fmt.Println("Fetching latest currency feed")

	err := currencyFeedService.Fetch()
	if err != nil {
		log.Println(err)
	}

	ok, err := currencyFeedService.Check()
	if err != nil {
		log.Println(err)
	}

	// Check if feed data needs to be stored in database
	if ok {
		// Store feed data in database
		err := currencyFeedService.Persist()
		if err != nil {
			panic(err)
		}

		fmt.Println("Currencies updated successfully!")
	} else {
		fmt.Println("No need to update!")
	}

	return nil
}

func serve(mysqlService *services.MySQLService, currencyService *services.CurrencyService) error {
	fmt.Println("Starting HTTP server...")

	currencies, err := currencyService.BySymbol("PHP", 20)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Values (PHP): %+v\n", currencies)

	currencies, err = currencyService.Latest()
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Latest currencies: %+v\n", currencies)

	return nil
}

const infoText = `
Go Microservice Test

Usage:
  go-microservice [command]

Available Commands:
  update Fetch and update latest currencies
  serve Start HTTP server
  setup Setup database tables

`

func printInfo() {
	fmt.Println(infoText)
}

func main() {
	// Print info text
	if len(os.Args) < 2 {
		printInfo()
		return
	}

	// Setup MySQL service
	mysqlService, err := services.NewMySQLService("root", "12345678", "tcp(127.0.0.1:3306)", "go_microservice_test")
	if err != nil {
		panic(err)
	}

	// Terminate database service before exiting
	defer mysqlService.Close()

	// Setup currency feed service
	currencyFeedService := services.NewCurrencyFeedService(mysqlService, feedURL)

	// Setup currency service
	currencyService := services.NewCurrencyService(mysqlService)

	switch os.Args[1] {
	case "update":
		// Run update command
		update(mysqlService, currencyFeedService)
	case "serve":
		// Run serve command
		serve(mysqlService, currencyService)
	case "setup":
		// Run setup command
		setup(mysqlService)
	default:
		printInfo()
	}

}
