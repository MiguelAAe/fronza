package useraccount

import (
	"encoding/json"
	"log"
	"net/http"
	"portal-backend/pkg/api/middleware"
	"portal-backend/pkg/auth"
	"portal-backend/pkg/deliverysvc"
	"portal-backend/pkg/usersvc"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Name        string `json:"name"  validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserResponse struct {
	ID          int64
	Name        string
	Email       string
	PhoneNumber string
	CreateTime  time.Time
}

// RegisterCustomer
// @Summary      Registers a new customer user
// @Description  Registers a new customer user
// @Accept       json
// @Produce      json
// @Param        job  body      UserRequest    true  "Registers a new user"
// @Success      201  {object}  UserResponse   "Registers a new user"
// @Header       201  {string}  Token          "authorisation token"
// @Failure      400  {object}  ErrorResponse  "Registers a new user"
// @Failure      500  {object}  ErrorResponse  "Registers a new user"
// @Router       /register/customer [post]
func RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	req := UserRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	// Generate "hashedPassword" to store from user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	err = usersvc.CreateUser(usersvc.User{
		Name:           req.Name,
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		HashedPassword: hashedPassword,
		Role:           usersvc.Customer,
	})
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	user, err := usersvc.GetUser(req.Email)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	token, err := auth.NewJWT(user.ID, user.Email, "customer", "")
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.Header().Add("Authorization", token)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CreateTime:  user.CreateTime,
	})
}

// RegisterCourier
// @Summary      Registers a new courier user
// @Description  Registers a new courier user
// @Accept       json
// @Produce      json
// @Param        job  body      UserRequest    true  "Registers a new user"
// @Success      201  {object}  UserResponse   "Registers a new user"
// @Header       201  {string}  Token          "authorisation token"
// @Failure      400  {object}  ErrorResponse  "Registers a new user"
// @Failure      500  {object}  ErrorResponse  "Registers a new user"
// @Router       /register/courier [post]
func RegisterCourier(w http.ResponseWriter, r *http.Request) {
	req := UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	// Generate "hashedPassword" to store from user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	err = usersvc.CreateUser(usersvc.User{
		Name:           req.Name,
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		HashedPassword: hashedPassword,
		Role:           usersvc.Courier,
	})
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	user, err := usersvc.GetUser(req.Email)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	token, err := auth.NewJWT(user.ID, user.Email, "courier", "")
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.Header().Add("Authorization", token)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CreateTime:  user.CreateTime,
	})
}

// RegisterAdmin
// @Summary      Registers a new admin user
// @Description  Registers a new admin user
// @Accept       json
// @Produce      json
// @Param        job  body      UserRequest    true  "Registers a new user"
// @Success      201  {object}  UserResponse   "Registers a new user"
// @Header       201  {string}  Token          "authorisation token"
// @Failure      400  {object}  ErrorResponse  "Registers a new user"
// @Failure      500  {object}  ErrorResponse  "Registers a new user"
// @Router       /register/admin [post]
func RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	req := UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	// Generate "hashedPassword" to store from user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	err = usersvc.CreateUser(usersvc.User{
		Name:           req.Name,
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		HashedPassword: hashedPassword,
		Role:           usersvc.Admin,
	})
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	user, err := usersvc.GetUser(req.Email)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	token, err := auth.NewJWT(user.ID, user.Email, "admin", "")
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.Header().Add("Authorization", token)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CreateTime:  user.CreateTime,
	})
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUsers
// @Summary      Logins a user
// @Description  Logins a user
// @Accept       json
// @Produce      json
// @Param        job  body  LoginUser  true  "Logins a user"
// @Success      200
// @Header       200  {string}  Token          "authorisation token"
// @Failure      400  {object}  ErrorResponse  "Logins a user"
// @Failure      500  {object}  ErrorResponse  "Logins a user"
// @Router       /login [post]
func LoginUsers(w http.ResponseWriter, r *http.Request) {
	req := LoginUser{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	user, err := usersvc.GetUser(req.Email)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(req.Password))
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	token, err := auth.NewJWT(user.ID, user.Email, user.Role, "")
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.Header().Add("Authorization", token)
	w.WriteHeader(http.StatusOK)
}

type UpdatePasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UpdatePassword
// @Summary      Updates a user password
// @Description  Updates a user password
// @Accept       json
// @Produce      json
// @Param        Authorization   header  string                 true  "token: test"
// @Param        updatePassword  body    UpdatePasswordRequest  true  "Logins a user"
// @Success      204
// @Failure      400  {object}  ErrorResponse  "Logins a user"
// @Failure      500  {object}  ErrorResponse  "Logins a user"
// @Router       /update-login [put]
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	req := UpdatePasswordRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	user, err := usersvc.GetUser(req.Email)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	//todo: security make sure the JWT matches with email

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(req.OldPassword))
	log.Println("after hash comparison")
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	// Generate "hashedPassword" to store from user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	// update password
	err = usersvc.UpdatePassword(req.Email, hashedPassword)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type NewAPIKeyResponse struct {
	Key string
}

// NewAPIKey
// @Summary      Issues an API key
// @Description  Issues a user with a new api key, reissuing a new one invalidates old api keys
// @Accept       json
// @Produce      json
// @Param        Authorization   header  string          true  "token: test"
// @Param        email          query     string             true  "email"
// @Success      201            {object}  NewAPIKeyResponse  "Issues an API key"
// @Failure      400            {object}  ErrorResponse      "Issues an API key"
// @Failure      500            {object}  ErrorResponse      "Issues an API key"
// @Router       /apikey [get]
func NewAPIKey(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	validate := validator.New()
	err := validate.Var(email, "email")
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	claim, ok := middleware.GetClaim(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if claim.Email != email {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	guid := uuid.New()

	id := strings.ReplaceAll(guid.String(), "-", "")

	//hash the id and save it to the db
	hashedID, err := bcrypt.GenerateFromPassword([]byte(id), bcrypt.DefaultCost)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	user, err := usersvc.GetUser(email)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	//save it to the db
	err = usersvc.UpsertApiKey(hashedID, user.ID)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	// return the plain id to the user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&NewAPIKeyResponse{Key: id})
}

type UserInfo struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

// UserInformation
// @Summary      Retrieves user information
// @Description  Retrieves user information from the JWT token
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string             true  "token: test"
// @Success      200            {object}  UserInfo       "Retrieves user information"
// @Failure      400            {object}  ErrorResponse  "Retrieves user information"
// @Failure      500            {object}  ErrorResponse  "Retrieves user information"
// @Router       /user/info [get]
func UserInformation(w http.ResponseWriter, r *http.Request) {
	claim, ok := middleware.GetClaim(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := usersvc.GetUser(claim.Email)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&UserInfo{
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	})
}

type DriverInfo struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// DriverInformation
// @Summary      Retrieves the information of a driver from an order
// @Description  Retrieves the information of a driver from an order
// @Accept       json
// @Produce      json
// @Param        Authorization   header  string        true  "token: test"
// @Param        job-id         query     string         true  "job-id"
// @Success      200            {object}  DriverInfo     "Gets the information of a driver from an order"
// @Failure      400            {object}  ErrorResponse  "Gets the information of a driver from an order"
// @Failure      500            {object}  ErrorResponse  "Gets the information of a driver from an order"
// @Router       /job/driver-info [get]
func DriverInformation(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job-id")

	_, err := uuid.Parse(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	jobResp, err := deliverysvc.GetJob(jobID)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	userResp, err := usersvc.GetUserByID(jobResp.Worker)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DriverInfo{
		Name:        userResp.Name,
		Email:       userResp.Email,
		PhoneNumber: userResp.PhoneNumber,
	})
}

type DriverLocation struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// SaveDriverLocation
// @Summary      saves the location of a driver
// @Description  saves the location of a driver
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "token: test"
// @Param        driverLocation  body    DriverLocation  true  "latitude and longitude"
// @Success      200
// @Failure      400  {object}  ErrorResponse  "Updates a job to assigned"
// @Failure      500  {object}  ErrorResponse  "Updates a job to assigned"
// @Router       /driver/save-location [post]
func SaveDriverLocation(w http.ResponseWriter, r *http.Request) {
	claim, ok := middleware.GetClaim(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if claim.Role != usersvc.Courier {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := DriverLocation{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	//todo: validate data

	coordinates := strings.ReplaceAll(req.Latitude, " ", "") + "," + strings.ReplaceAll(req.Longitude, " ", "")

	err = usersvc.SaveDriverLocation(claim.ID, coordinates)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type DriverStatus struct {
	Status bool `json:"status"`
}

// SaveDriverStatus
// @Summary      saves the status of a driver
// @Description  saves the status of a driver
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "token: test"
// @Param        driverLocation  body    DriverStatus  true  "the status"
// @Success      200
// @Failure      400  {object}  ErrorResponse  "Updates a job to assigned"
// @Failure      500  {object}  ErrorResponse  "Updates a job to assigned"
// @Router       /driver/status [post]
func SaveDriverStatus(w http.ResponseWriter, r *http.Request) {
	claim, ok := middleware.GetClaim(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if claim.Role != usersvc.Courier {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := DriverStatus{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middleware.WriteErr(w, http.StatusBadRequest, "", err)
		return
	}

	err = usersvc.SaveDriverStatus(claim.ID, req.Status)
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
