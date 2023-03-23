package api

import (
	"net/http"
	"portal-backend/pkg/api/handlers/deliveries"
	"portal-backend/pkg/api/handlers/management"
	"portal-backend/pkg/api/handlers/payments"
	"portal-backend/pkg/api/handlers/useraccount"
	"portal-backend/pkg/api/handlers/verification"
	portalMiddleware "portal-backend/pkg/api/middleware"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/httprate"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

func NewAPI() *chi.Mux {
	// set up router
	r := chi.NewRouter()

	// RESTy routes for "articles" resource
	// check mount for splitting big routes in apis 'example'
	r.Use(middleware.Logger)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use()

	//http://localhost:8080/swagger/index.html
	// todo: review all the documentation
	r.Get("/swagger/*", httpSwagger.Handler())

	// public routes
	r.Group(func(r chi.Router) {
		r.Route("/register", func(r chi.Router) {
			r.Post("/courier", useraccount.RegisterCourier)
			r.Post("/customer", useraccount.RegisterCustomer)
			r.Post("/admin", useraccount.RegisterAdmin)
		})

		r.Route("/login", func(r chi.Router) {
			r.Post("/", useraccount.LoginUsers)
		})

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello world, welcome to portal deliveries backend"))
		})

		r.Route("/stripe", func(r chi.Router) {
			r.Post("/webhook", payments.Webhook)
		})

		r.Route("/verify", func(r chi.Router) {
			r.Post("/postcode", verification.Postcode)
		})
	})

	// private common routes
	r.Group(func(r chi.Router) {
		r.Use(portalMiddleware.Auth)
		r.Route("/update-login", func(r chi.Router) {
			r.Put("/", useraccount.UpdatePassword)
		})
		r.Route("/apikey", func(r chi.Router) {
			r.Get("/", useraccount.NewAPIKey) // {email} creates a new api key
		})

		r.Route("/job", func(r chi.Router) {
			r.Post("/", deliveries.NewJob)                       // create job
			r.Post("/close", deliveries.CloseJob)                // close job
			r.Post("/cancel", deliveries.CancelJob)              // cancel job
			r.Get("/", deliveries.GetJob)                        // {job-id} get job job information
			r.Get("/status", deliveries.GetJobStatus)            // {job-id} get job status
			r.Get("/order-status", deliveries.GetJobOrderStatus) // order-status
			r.Get("/driver-info", useraccount.DriverInformation)
			r.Get("/user", deliveries.GetJobs) // {id} gets all jobs from a user
			r.Post("/on-route-to-pick-up", deliveries.UpdateJobOnRouteToPickUpLocation)
			r.Post("/parcel-collected", deliveries.UpdateJobParcelCollected)
			r.Post("/on-route-to-drop-off", deliveries.UpdateJobOnRouteToDropOffLocation)
			r.Post("/complete", deliveries.UpdateJobStatusComplete)
			r.Get("/eta", deliveries.ETA)
			r.Post("/poa", deliveries.PostProofOfDelivery)
			r.Get("/poa", deliveries.GetProofOfDelivery)
		})

		r.Route("/driver", func(r chi.Router) {
			r.Post("/jobs", deliveries.GetDriverJobs)
			// todo: retrieve a collection (a collection is a list of jobs to do in order)
			r.Post("/save-location", useraccount.SaveDriverLocation)
			r.Post("/status", useraccount.SaveDriverStatus)
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/info", useraccount.UserInformation)
		})

		r.Route("/management", func(r chi.Router) {
			r.Get("/jobs", management.GetAllJobs)
			r.Get("/driver-state", management.GetDriverStates)
			//r.Post("/assign-driver-job") // todo: assign driver to job
			//r.Get("")
			// todo: retrieve driver information
		})

		r.Route("/payment", func(r chi.Router) {
			r.Post("/create-intent", payments.CreatePaymentIntent)
			r.Get("/config", payments.Config)
		})
	})

	return r

}
