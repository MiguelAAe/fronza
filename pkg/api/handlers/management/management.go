package management

import (
	"encoding/json"
	"net/http"
	"portal-backend/pkg/api/middleware"
	"portal-backend/pkg/deliverysvc"
	"portal-backend/pkg/usersvc"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// GetAllJobs
// @Summary      Retrieves all the jobs
// @Description  Retrieves all the jobs
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string           true       "token: test"
// @Success      200            {object}  deliverysvc.Job  Retrieves  all  the  jobs
// @Failure      400            {object}  ErrorResponse    "job"
// @Failure      500            {object}  ErrorResponse    "job"
// @Router       /management/jobs [get]
func GetAllJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := deliverysvc.GetAllJobs()
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&jobs)
}

// GetDriverStates
// @Summary      Retrieves all the current driver statuses
// @Description  Retrieves all the current driver statuses
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string               true       "token: test"
// @Success      200            {object}  usersvc.DriverState  Retrieves  all  the  jobs
// @Failure      400            {object}  ErrorResponse        "driver state"
// @Failure      500            {object}  ErrorResponse        "driver state"
// @Router       /management/driver-state [get]
func GetDriverStates(w http.ResponseWriter, r *http.Request) {
	driverStates, err := usersvc.GetDriverStates()
	if err != nil {
		middleware.WriteErr(w, http.StatusInternalServerError, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&driverStates)
}
