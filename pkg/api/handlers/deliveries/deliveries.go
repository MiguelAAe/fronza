package deliveries

import (
	"encoding/json"
	"fmt"
	"net/http"
	"portal-backend/pkg/api/middleware"
	"portal-backend/pkg/deliverysvc"
	"portal-backend/pkg/usersvc"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type JobRequest struct {
	Origin      Recipient `json:"origin"`
	Destination Recipient `json:"destination"`
	WorkerNotes string    `json:"worker_notes"`
}

type JobResponse struct {
	ID               string    `json:"id"`
	ShortID          string    `json:"short_id"`
	CreateTime       time.Time `json:"create_time"`
	LastTimeModified time.Time `json:"last_time_modified"`
	TrackingURL      string    `json:"tracking_url"`
	Creator          int64     `json:"creator"`
	Worker           int64     `json:"worker"`
	Status           int       `json:"status"`
	OrderStatus      int       `json:"order_status"`
	Origin           Recipient `json:"origin"`
	Destination      Recipient `json:"destination"`
	WorkerNotes      string    `json:"worker_notes"`
}

type Recipient struct {
	CompanyName       string  `json:"company_name"`
	FirstName         string  `json:"first_name" validate:"required,alpha"`
	SecondName        string  `json:"second_name" validate:"alpha"`
	PhoneNumber       string  `json:"phone_number"`
	EmailAddress      string  `json:"email_address" validate:"email"`
	FirstLineAddress  string  `json:"first_line_address" validate:"required"`
	SecondLineAddress string  `json:"second_line_address"`
	ThirdLineAddress  string  `json:"third_line_address"`
	Town              string  `json:"town" validate:"alpha"`
	City              string  `json:"city" validate:"required,alpha"`
	PostCode          string  `json:"postcode" validate:"required,alphanum"`
	Latitude          float32 `json:"latitude"`
	Longitude         float32 `json:"longitude"`
	Notes             string  `json:"notes"`
}

// NewJob
// @Summary      Creates a new job with its status open
// @Description  creates a new job (its order status is open and needs to be closed to be processed)
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "JWT token"
// @Param        job            body      JobRequest     true  "Add account"
// @Success      201            {object}  JobResponse    "New job, replies with all the created details of a job"
// @Failure      400            {object}  ErrorResponse  "ErrorResponse"
// @Failure      500            {object}  ErrorResponse  "ErrorResponse"
// @Router       /job [post]
func NewJob(w http.ResponseWriter, r *http.Request) {
	claim, ok := middleware.GetClaim(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	req := JobRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	// data validation
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	jobSave := deliverysvc.Job{
		Creator:                      claim.ID,
		OriginCompanyName:            req.Origin.CompanyName,
		OriginFirstName:              req.Origin.FirstName,
		OriginSecondName:             req.Origin.SecondName,
		OriginPhoneNumber:            req.Origin.PhoneNumber,
		OriginEmailAddress:           req.Origin.EmailAddress,
		OriginFirstLineAddress:       req.Origin.FirstLineAddress,
		OriginSecondLineAddress:      req.Origin.SecondLineAddress,
		OriginThirdLineAddress:       req.Origin.ThirdLineAddress,
		OriginTown:                   req.Origin.Town,
		OriginCity:                   req.Origin.City,
		OriginPostcode:               req.Origin.PostCode,
		OriginLatitude:               req.Origin.Latitude,
		OriginLongitude:              req.Origin.Longitude,
		OriginNotes:                  req.Origin.Notes,
		DestinationCompanyName:       req.Destination.CompanyName,
		DestinationFirstName:         req.Destination.FirstName,
		DestinationSecondName:        req.Destination.SecondName,
		DestinationPhoneNumber:       req.Destination.PhoneNumber,
		DestinationEmailAddress:      req.Destination.EmailAddress,
		DestinationFirstLineAddress:  req.Destination.FirstLineAddress,
		DestinationSecondLineAddress: req.Destination.SecondLineAddress,
		DestinationThirdLineAddress:  req.Destination.ThirdLineAddress,
		DestinationTown:              req.Destination.Town,
		DestinationCity:              req.Destination.City,
		DestinationPostcode:          req.Destination.PostCode,
		DestinationLatitude:          req.Destination.Latitude,
		DestinationLongitude:         req.Destination.Longitude,
		DestinationNotes:             req.Destination.Notes,
	}

	job, err := deliverysvc.SaveJob(jobSave)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	resp := JobResponse{
		ID:               job.ID,
		ShortID:          job.ShortID,
		CreateTime:       job.CreateTime,
		LastTimeModified: job.LastTimeModified,
		TrackingURL:      job.TrackingURL,
		Creator:          job.Creator,
		Worker:           job.Worker,
		Status:           job.Status,
		OrderStatus:      job.OrderStatus,
		Origin: Recipient{
			CompanyName:       job.OriginCompanyName,
			FirstName:         job.OriginFirstName,
			SecondName:        job.OriginSecondName,
			PhoneNumber:       job.OriginPhoneNumber,
			EmailAddress:      job.OriginEmailAddress,
			FirstLineAddress:  job.OriginFirstLineAddress,
			SecondLineAddress: job.OriginSecondLineAddress,
			ThirdLineAddress:  job.OriginThirdLineAddress,
			Town:              job.OriginTown,
			City:              job.OriginCity,
			PostCode:          job.OriginPostcode,
			Latitude:          job.OriginLatitude,
			Longitude:         job.OriginLongitude,
			Notes:             job.OriginNotes,
		},
		Destination: Recipient{
			CompanyName:       job.DestinationCompanyName,
			FirstName:         job.DestinationFirstName,
			SecondName:        job.DestinationSecondName,
			PhoneNumber:       job.DestinationPhoneNumber,
			EmailAddress:      job.DestinationEmailAddress,
			FirstLineAddress:  job.DestinationFirstLineAddress,
			SecondLineAddress: job.DestinationSecondLineAddress,
			ThirdLineAddress:  job.DestinationThirdLineAddress,
			Town:              job.DestinationTown,
			City:              job.OriginCity,
			PostCode:          job.DestinationPostcode,
			Latitude:          job.DestinationLatitude,
			Longitude:         job.DestinationLongitude,
			Notes:             job.DestinationNotes,
		},
		WorkerNotes: job.WorkerNotes,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&resp)
}

// GetJob
// @Summary      Gets the information of a job
// @Description  Gets all the information and details of a job
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "JWT token"
// @Param        job-id         query     string         true  "job-id"
// @Success      200            {object}  JobResponse    "Get job, retrieves the details of a job"
// @Failure      400            {object}  ErrorResponse  "ErrorResponse"
// @Failure      500            {object}  ErrorResponse  "ErrorResponse"
// @Router       /job [get]
func GetJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job-id")

	_, err := uuid.Parse(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	job, err := deliverysvc.GetJob(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	resp := JobResponse{
		ID:               job.ID,
		ShortID:          job.ShortID,
		CreateTime:       job.CreateTime,
		LastTimeModified: job.LastTimeModified,
		TrackingURL:      job.TrackingURL,
		Creator:          job.Creator,
		Worker:           job.Worker,
		Status:           job.Status,
		OrderStatus:      job.OrderStatus,
		Origin: Recipient{
			CompanyName:       job.OriginCompanyName,
			FirstName:         job.OriginFirstName,
			SecondName:        job.OriginSecondName,
			PhoneNumber:       job.OriginPhoneNumber,
			EmailAddress:      job.OriginEmailAddress,
			FirstLineAddress:  job.OriginFirstLineAddress,
			SecondLineAddress: job.OriginSecondLineAddress,
			ThirdLineAddress:  job.OriginThirdLineAddress,
			Town:              job.OriginTown,
			City:              job.OriginCity,
			PostCode:          job.OriginPostcode,
			Latitude:          job.OriginLatitude,
			Longitude:         job.OriginLongitude,
			Notes:             job.OriginNotes,
		},
		Destination: Recipient{
			CompanyName:       job.DestinationCompanyName,
			FirstName:         job.DestinationFirstName,
			SecondName:        job.DestinationSecondName,
			PhoneNumber:       job.DestinationPhoneNumber,
			EmailAddress:      job.DestinationEmailAddress,
			FirstLineAddress:  job.DestinationFirstLineAddress,
			SecondLineAddress: job.DestinationSecondLineAddress,
			ThirdLineAddress:  job.DestinationThirdLineAddress,
			Town:              job.DestinationTown,
			City:              job.OriginCity,
			PostCode:          job.DestinationPostcode,
			Latitude:          job.DestinationLatitude,
			Longitude:         job.DestinationLongitude,
			Notes:             job.DestinationNotes,
		},
		WorkerNotes: job.WorkerNotes,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

// CloseJob
// @Summary      Updates the status of a job to close
// @Description  We normally close a job when its ready to be processed and its payment intend has gone through
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string             true  "JWT token"
// @Param        job-id         query   string  true  "job-id"
// @Success      200
// @Failure      400  {object}  ErrorResponse  "ErrorResponse"
// @Failure      500  {object}  ErrorResponse  "ErrorResponse"
// @Router       /job/close [post]
func CloseJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job-id")

	_, err := uuid.Parse(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	err = deliverysvc.UpdateJobOrderStatus(jobID, deliverysvc.OrderStatusClosed)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CancelJob
// @Summary      Updates the status of a job to cancel
// @Description  We normally cancel a job when its payment intend has not gone through
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "JWT token"
// @Param        job-id         query   string  true  "job-id"
// @Success      200
// @Failure      400  {object}  ErrorResponse  "ErrorResponse"
// @Failure      500  {object}  ErrorResponse  "ErrorResponse"
// @Router       /job/cancel [post]
func CancelJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job-id")

	_, err := uuid.Parse(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	err = deliverysvc.UpdateJobOrderStatus(jobID, deliverysvc.OrderStatusCancelled)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type JobStatusResponse struct {
	Status string `json:"status"`
}

// GetJobStatus
// @Summary      Retrieves the status of a job
// @Description  Retrieves the status of a job, this status refers to in which stage the delivery is at
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "JWT token"
// @Param        job-id         query     string             true  "job-id"
// @Success      200            {object}  JobStatusResponse  "The status of a job"
// @Failure      400            {object}  ErrorResponse      "ErrorResponse"
// @Failure      500            {object}  ErrorResponse      "ErrorResponse"
// @Router       /job/status [get]
func GetJobStatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job-id")

	_, err := uuid.Parse(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	status, err := deliverysvc.GetJobStatus(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"order_status": status,
	})
}

type JobOrderStatusResponse struct {
	OrderStatus string `json:"order_status"`
}

// GetJobOrderStatus
// @Summary      Retrieves the order status of a job
// @Description  Retrieves the order status of a job, this status refers to the payment stage open,close,cancel
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "JWT token"
// @Param        job-id         query     string                  true  "job-id"
// @Success      200            {object}  JobOrderStatusResponse  "GetJobOrderStatus, gets the status of a job"
// @Failure      400            {object}  ErrorResponse           "Get Job Status"
// @Failure      500            {object}  ErrorResponse           "Get Job Status"
// @Router       /job/order-status [get]
func GetJobOrderStatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job-id")

	_, err := uuid.Parse(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	status, err := deliverysvc.GetJobOrderStatus(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"order_status": status,
	})
}

// GetJobs
// @Summary      Gets all the jobs of a user
// @Description  retrieves all the jobs of a user
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "token: test"
// @Param        id             query     string         true  "id"
// @Success      200            {array}   JobResponse    "GetJobs, gets all the jobs of a user"
// @Failure      400            {object}  ErrorResponse  "GetJobs"
// @Failure      500            {object}  ErrorResponse  "GetJobs"
// @Router       /job/user [get]
func GetJobs(w http.ResponseWriter, r *http.Request) {
	urlValue := r.URL.Query().Get("id")

	userID, err := strconv.Atoi(urlValue)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	jobs, err := deliverysvc.GetUserJobs(userID)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	resp := make([]JobResponse, len(jobs))

	for i, job := range jobs {
		resp[i] = JobResponse{
			ID:               job.ID,
			ShortID:          job.ShortID,
			CreateTime:       job.CreateTime,
			LastTimeModified: job.LastTimeModified,
			TrackingURL:      job.TrackingURL,
			Creator:          job.Creator,
			Worker:           job.Worker,
			Status:           job.Status,
			OrderStatus:      job.OrderStatus,
			Origin: Recipient{
				CompanyName:       job.OriginCompanyName,
				FirstName:         job.OriginFirstName,
				SecondName:        job.OriginSecondName,
				PhoneNumber:       job.OriginPhoneNumber,
				EmailAddress:      job.OriginEmailAddress,
				FirstLineAddress:  job.OriginFirstLineAddress,
				SecondLineAddress: job.OriginSecondLineAddress,
				ThirdLineAddress:  job.OriginThirdLineAddress,
				Town:              job.OriginTown,
				City:              job.OriginCity,
				PostCode:          job.OriginPostcode,
				Latitude:          job.OriginLatitude,
				Longitude:         job.OriginLongitude,
				Notes:             job.OriginNotes,
			},
			Destination: Recipient{
				CompanyName:       job.DestinationCompanyName,
				FirstName:         job.DestinationFirstName,
				SecondName:        job.DestinationSecondName,
				PhoneNumber:       job.DestinationPhoneNumber,
				EmailAddress:      job.DestinationEmailAddress,
				FirstLineAddress:  job.DestinationFirstLineAddress,
				SecondLineAddress: job.DestinationSecondLineAddress,
				ThirdLineAddress:  job.DestinationThirdLineAddress,
				Town:              job.DestinationTown,
				City:              job.OriginCity,
				PostCode:          job.DestinationPostcode,
				Latitude:          job.DestinationLatitude,
				Longitude:         job.DestinationLongitude,
				Notes:             job.DestinationNotes,
			},
			WorkerNotes: job.WorkerNotes,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

// UpdateJobOnRouteToPickUpLocation
// @Summary      Updates a job to OnRouteToPickUpLocation
// @Description  retrieves all the jobs of a user
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "token: test"
// @Param        job-id         query     string         true  "job-id"
// @Success      200            {array}   JobResponse    "Updates a job to assigned"
// @Failure      400            {object}  ErrorResponse  "Updates a job to assigned"
// @Failure      500            {object}  ErrorResponse  "Updates a job to assigned"
// @Router       /driver/job/on-route-to-pick-up [post]
func UpdateJobOnRouteToPickUpLocation(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job-id")

	_, err := uuid.Parse(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	claim, ok := middleware.GetClaim(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if claim.Role != usersvc.Courier {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = deliverysvc.UpdateJobStatus(jobID, deliverysvc.JobStatusOnRouteToPickUpLocation)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateJobParcelCollected
// @Summary      Updates a job to ParcelCollected
// @Description  retrieves all the jobs of a user
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "token: test"
// @Param        job-id         query     string         true  "job-id"
// @Success      200            {array}   JobResponse    "Updates a job to assigned"
// @Failure      400            {object}  ErrorResponse  "Updates a job to assigned"
// @Failure      500            {object}  ErrorResponse  "Updates a job to assigned"
// @Router       /driver/job/parcel-collected [post]
func UpdateJobParcelCollected(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job-id")

	_, err := uuid.Parse(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	claim, ok := middleware.GetClaim(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if claim.Role != usersvc.Courier {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = deliverysvc.UpdateJobStatus(jobID, deliverysvc.JobStatusParcelCollected)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateJobOnRouteToDropOffLocation
// @Summary      Updates a job to OnRouteToDropOffLocation
// @Description  retrieves all the jobs of a user
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "token: test"
// @Param        job-id         query     string         true  "job-id"
// @Success      200            {array}   JobResponse    "Updates a job to assigned"
// @Failure      400            {object}  ErrorResponse  "Updates a job to assigned"
// @Failure      500            {object}  ErrorResponse  "Updates a job to assigned"
// @Router       /driver/job/on-route-to-drop-off [post]
func UpdateJobOnRouteToDropOffLocation(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job-id")

	_, err := uuid.Parse(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	claim, ok := middleware.GetClaim(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if claim.Role != usersvc.Courier {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = deliverysvc.UpdateJobStatus(jobID, deliverysvc.JobStatusOnRouteToDropOffLocation)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateJobStatusComplete
// @Summary      Updates a job to StatusComplete
// @Description  retrieves all the jobs of a user
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "token: test"
// @Param        job-id         query     string         true  "job-id"
// @Success      200            {array}   JobResponse    "Updates a job to assigned"
// @Failure      400            {object}  ErrorResponse  "Updates a job to assigned"
// @Failure      500            {object}  ErrorResponse  "Updates a job to assigned"
// @Router       /driver/job/complete [post]
func UpdateJobStatusComplete(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job-id")

	_, err := uuid.Parse(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	claim, ok := middleware.GetClaim(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if claim.Role != usersvc.Courier {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = deliverysvc.UpdateJobStatus(jobID, deliverysvc.JobStatusComplete)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetDriverJobs
// @Summary      Gets the jobs of a driver
// @Description  Gets the jobs of a driver
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "token: test"
// @Param        job-id         query     string         true  "driver-id"
// @Success      200            {array}   JobResponse    "Gets the jobs of a driver"
// @Failure      400            {object}  ErrorResponse  "Gets the jobs of a driver"
// @Failure      500            {object}  ErrorResponse  "Gets the jobs of a driver"
// @Router       /driver/jobs [get]
func GetDriverJobs(w http.ResponseWriter, r *http.Request) {
	queryValue := r.URL.Query().Get("driver-id")

	driverID, err := strconv.ParseInt(queryValue, 10, 64)

	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	jobs, err := deliverysvc.GetDriverJobs(driverID)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	resp := make([]JobResponse, len(jobs))

	for i, job := range jobs {
		resp[i] = JobResponse{
			ID:               job.ID,
			ShortID:          job.ShortID,
			CreateTime:       job.CreateTime,
			LastTimeModified: job.LastTimeModified,
			TrackingURL:      job.TrackingURL,
			Creator:          job.Creator,
			Worker:           job.Worker,
			Status:           job.Status,
			OrderStatus:      job.OrderStatus,
			Origin: Recipient{
				CompanyName:       job.OriginCompanyName,
				FirstName:         job.OriginFirstName,
				SecondName:        job.OriginSecondName,
				PhoneNumber:       job.OriginPhoneNumber,
				EmailAddress:      job.OriginEmailAddress,
				FirstLineAddress:  job.OriginFirstLineAddress,
				SecondLineAddress: job.OriginSecondLineAddress,
				ThirdLineAddress:  job.OriginThirdLineAddress,
				Town:              job.OriginTown,
				City:              job.OriginCity,
				PostCode:          job.OriginPostcode,
				Latitude:          job.OriginLatitude,
				Longitude:         job.OriginLongitude,
				Notes:             job.OriginNotes,
			},
			Destination: Recipient{
				CompanyName:       job.DestinationCompanyName,
				FirstName:         job.DestinationFirstName,
				SecondName:        job.DestinationSecondName,
				PhoneNumber:       job.DestinationPhoneNumber,
				EmailAddress:      job.DestinationEmailAddress,
				FirstLineAddress:  job.DestinationFirstLineAddress,
				SecondLineAddress: job.DestinationSecondLineAddress,
				ThirdLineAddress:  job.DestinationThirdLineAddress,
				Town:              job.DestinationTown,
				City:              job.OriginCity,
				PostCode:          job.DestinationPostcode,
				Latitude:          job.DestinationLatitude,
				Longitude:         job.DestinationLongitude,
				Notes:             job.DestinationNotes,
			},
			WorkerNotes: job.WorkerNotes,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

type ETAResponse struct {
	Duration int `json:"duration"`
}

// ETA
// @Summary      Gets the ETA of a job
// @Description  retrieves the estimated time for a driver to arrive at its next destination
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "token: test"
// @Param        job-id         query   string  true  "job-id"
// @Param        driver-id      query   string  true  "driver-id"
// @Success      200
// @Failure      400  {object}  ErrorResponse  "Get job"
// @Failure      500  {object}  ErrorResponse  "Get job"
// @Router       /job/eta [post]
func ETA(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job-id")
	driverID := r.URL.Query().Get("driver-id")

	driverIDParsed, err := strconv.ParseInt(driverID, 10, 64)

	_, err = uuid.Parse(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	jwt, _ := middleware.GetClaim(r)

	if jwt.Role != usersvc.Courier {
		middleware.WriteErr(w, http.StatusUnauthorized, "", fmt.Errorf("invalid role"))
		return
	}

	duration, err := deliverysvc.GetEstimatedJobJourneyDuration(jobID, driverIDParsed)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ETAResponse{
		Duration: duration,
	})
}

// PostProofOfDelivery
// @Summary      uploads an image for proof of delivery
// @Description  uploads an image for proof of delivery
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "token: test"
// @Param        job-id         query   string  true  "job-id"
// @Param        driver-id      query   string  true  "driver-id"
// @Success      200
// @Failure      400  {object}  ErrorResponse  "Get job"
// @Failure      500  {object}  ErrorResponse  "Get job"
// @Router       /job/poa [post]
func PostProofOfDelivery(w http.ResponseWriter, r *http.Request) {

}

// GetProofOfDelivery
// @Summary      retrieves an image for proof of delivery
// @Description  retrieves an image for proof of delivery
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "token: test"
// @Param        job-id         query   string  true  "job-id"
// @Param        driver-id      query   string  true  "driver-id"
// @Success      200
// @Failure      400  {object}  ErrorResponse  "Get job"
// @Failure      500  {object}  ErrorResponse  "Get job"
// @Router       /job/poa [get]
func GetProofOfDelivery(w http.ResponseWriter, r *http.Request) {

}
