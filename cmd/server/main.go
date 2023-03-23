package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"portal-backend/config"
	"portal-backend/migrations"
	"portal-backend/pkg/api"
	"portal-backend/pkg/api/handlers/payments"
	"portal-backend/pkg/db"
	"portal-backend/pkg/deliverysvc"
	_ "portal-backend/pkg/deliverysvc"
	"portal-backend/pkg/usersvc"

	_ "portal-backend/docs"
)

func main() {
	switch os.Getenv("ENV") {
	case "prod":
		config.StripePublishableKey = os.Getenv("Stripe_Publishable_Key")
		config.StripeSecretKey = os.Getenv("Stripe_Secret_Key")
		config.StripeWebhookSecret = os.Getenv("Stripe_Webhook_Secret")
		config.DatabaseURL = os.Getenv("DATABASE_URL")
		config.Environment = "prod"
		config.ListeningPort = os.Getenv("PORT")
		config.GoogleAPIKey = os.Getenv("GOOGLE_API_KEY")
	default:
		config.StripePublishableKey = "pk_test_51KqNvMKopqzDmLTZsKAW1dE3QPkyj5th97SevjBPIm5Kg9KWg1UkGmDkudcbJbIzMB28qnGp473K0lSTznUYWkoj00YhvgG3Ge"
		config.StripeSecretKey = "sk_test_51KqNvMKopqzDmLTZfEhOC8yHDCq7g4ISqy3pQSvhteNL8dKrd83BvlFYDPIDuA0lexHXQrfkJzYsoFjRkikgVcWt00nylZyZ7m"
		config.StripeWebhookSecret = "whsec_651618d564bd95bd5aacc02670ba7abcf801abd67a6e4276f3406749a9f426aa"
		config.Environment = "test"
		config.ListeningPort = "8080"
		config.GoogleAPIKey = "AIzaSyBrnQzp90T6jbBpRlTOLfRJAisDE11Q53E"
	}

	info := `Portal Deliveries backend is starting up
+++++++++++++++++++++
+++++++++++++++++++++
+++++++++++++++++++++
+++++++++++++++++++++`
	version := "Version:0.0.1"

	fmt.Printf("%s\n%s\n", info, version)

	dbMig, err := migrations.Init()
	if err != nil {
		log.Fatalf("fatal connecting to db: %v", err)
	}

	err = migrations.Run(dbMig)
	if err != nil {
		log.Fatalf("fatal running migration: %v", err)
	}

	db, err := db.Init()
	if err != nil {
		log.Fatalf("fatal initialising db conn: %v", err)
	}

	deliverysvc.Init(db)
	usersvc.Init(db)

	payments.Start()

	router := api.NewAPI()
	log.Printf("Listening on Port:%s\n", config.ListeningPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", config.ListeningPort), router)
	if err != nil {
		log.Printf("err from router: %v\n", err)
	}
}
