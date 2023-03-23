package payments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"portal-backend/config"
	"portal-backend/pkg/api/middleware"
	"portal-backend/pkg/deliverysvc"

	"github.com/google/uuid"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/webhook"
)

func Start() {
	stripe.Key = config.StripeSecretKey

	// For sample support and debugging, not required for production:
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "portaldeliveries/backend",
		Version: "0.0.1",
		URL:     "https://portaldeliveries.co.uk",
	})
}

// ErrorResponseMessage represents the structure of the error
// object sent in failed responses.
type ErrorResponseMessage struct {
	Message string `json:"message"`
}

// ErrorResponse represents the structure of the error object sent
// in failed responses.
type ErrorResponse struct {
	Error *ErrorResponseMessage `json:"error"`
}

type PublicKey struct {
	PublishableKey string `json:"publishableKey"`
}

// Config
// @Summary      Gets the publishable key for stripe
// @Description  retrieves the publishable key for stripe
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "token: test"
// @Success      200            {object}  PublicKey      "Public key"
// @Failure      400            {object}  ErrorResponse  "Config"
// @Failure      500            {object}  ErrorResponse  "Config"
// @Router       /payment/config [get]
func Config(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, PublicKey{
		PublishableKey: config.StripePublishableKey,
	})
}

type ClientSecret struct {
	ClientSecret string `json:"clientSecret"`
}

// CreatePaymentIntent
// @Summary      Gets the publishable key for stripe
// @Description  retrieves the publishable key for stripe
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "token: test"
// @Param        job-id         query     string         true  "job-id"
// @Success      200            {object}  ClientSecret   "Public key"
// @Failure      400            {object}  ErrorResponse  "Client Secret"
// @Failure      500            {object}  ErrorResponse  "Client Secret"
// @Router       /payment/create-intent [post]
func CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(850),
		Currency: stripe.String(string(stripe.CurrencyGBP)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	jobID := r.URL.Query().Get("job-id")

	_, err := uuid.Parse(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	params.Metadata = map[string]string{
		"order_id": jobID,
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		// Try to safely cast a generic error to a stripe.Error so that we can get at
		// some additional Stripe-specific information about what went wrong.
		if stripeErr, ok := err.(*stripe.Error); ok {
			fmt.Printf("Other Stripe error occurred: %v\n", stripeErr.Error())
			writeJSONErrorMessage(w, stripeErr.Error(), 400)
		} else {
			fmt.Printf("Other error occurred: %v\n", err.Error())
			writeJSONErrorMessage(w, "Unknown server error", 500)
		}

		return
	}

	writeJSON(w, ClientSecret{
		ClientSecret: pi.ClientSecret,
	})
}

func Webhook(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Pass the request body & Stripe-Signature header to ConstructEvent, along with the webhook signing key
	// You can find your endpoint's secret in your webhook settings
	endpointSecret := config.StripeWebhookSecret
	event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), endpointSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.canceled":
		// Then define and call a function to handle the event payment_intent.canceled
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Sprintf("Failed: %v, %v\n", paymentIntent, paymentIntent.LastPaymentError)
		// Notify the customer that payment failed
	case "payment_intent.payment_failed":
		// Then define and call a function to handle the event payment_intent.payment_failed
	case "payment_intent.succeeded":
		// Then define and call a function to handle the event payment_intent.succeeded
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Sprintf("Succeeded: %v\n", paymentIntent)
		// Fulfil the customer's purchase

		jobID := paymentIntent.Metadata["order_id"]

		err = deliverysvc.UpdateJobOrderStatus(jobID, deliverysvc.OrderStatusClosed)
		if err != nil {
			//todo: this is very bad would need to try to update the job again or something
			fmt.Println("failed to update job order after payment" + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	// ... handle other event types
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
	return
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}

func writeJSONError(w http.ResponseWriter, v interface{}, code int) {
	w.WriteHeader(code)
	writeJSON(w, v)
	return
}

func writeJSONErrorMessage(w http.ResponseWriter, message string, code int) {
	resp := &ErrorResponse{
		Error: &ErrorResponseMessage{
			Message: message,
		},
	}
	writeJSONError(w, resp, code)
}
