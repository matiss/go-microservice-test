package main

import (
	"log"

	"github.com/matiss/go-microservice-test/services"
)

const (
	feedURL = "https://www.bank.lv/vk/ecb_rss.xml"
)

func init() {

}

func main() {
	// Setup MySQL service
	mysqlService, err := services.NewMySQLService("root", "12345678", "tcp(127.0.0.1:3306)", "go_microservice_test")
	if err != nil {
		panic(err)
	}

	// Run setup
	// TODO: Move this into seperate command
	err = mysqlService.Setup()
	if err != nil {
		panic(err)
	}

	// Setup currency feed service
	currencyFeedService := services.NewCurrencyFeedService(mysqlService, feedURL)

	// Fetch and parse currency feed
	log.Println("Fetching latest currency feed")

	err = currencyFeedService.Fetch()
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

		log.Println("Currencies updated successfully!")
	} else {
		log.Println("No need to update the currencies")
	}

	// Terminate database service before exiting
	mysqlService.Close()
}
