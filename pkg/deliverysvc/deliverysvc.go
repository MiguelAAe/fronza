package deliverysvc

import (
	"context"
	"encoding/base64"
	"fmt"
	"portal-backend/pkg/addresssvc"
	"portal-backend/pkg/location"
	"portal-backend/pkg/usersvc"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	JobStatusPending                  = 0
	JobStatusOnRouteToPickUpLocation  = 1 // on route to pick up location
	JobStatusParcelCollected          = 2 // parcel collected
	JobStatusOnRouteToDropOffLocation = 3 // on route to drop off location
	JobStatusComplete                 = 4 // delivered

	// order statuses
	OrderStatusOpen      = 0
	OrderStatusCancelled = 1
	OrderStatusClosed    = 2
)

func SaveJob(job Job) (Job, error) {
	// todo: set context timeouts for requests
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	id := uuid.New()

	// source: https://stackoverflow.com/questions/37934162/output-uuid-in-go-as-a-short-string
	var escaper = strings.NewReplacer("9", "99", "-", "90", "_", "91")
	shortID := escaper.Replace(base64.RawURLEncoding.EncodeToString(id.NodeID()))

	createTime := time.Now().UTC()
	trackingURL := "portaldeliveries-backend.herokuapp.com/job/" + shortID

	job.ID = id.String()
	job.ShortID = shortID
	job.CreateTime = createTime
	job.LastTimeModified = createTime
	job.TrackingURL = trackingURL
	job.Status = JobStatusPending
	job.OrderStatus = OrderStatusOpen

	_, err := db.NewInsert().Model(&job).Exec(context.TODO())
	if err != nil {
		return Job{}, fmt.Errorf("failed to save job, %v", err)
	}

	var jobResult Job

	err = db.NewSelect().Model(&jobResult).Where("id = ?", id).Scan(context.TODO())
	if err != nil {
		return Job{}, fmt.Errorf("failed to save job id:%s, %v", id.String(), err)
	}
	return jobResult, nil
}

func GetJob(jobID string) (Job, error) {
	var jobResult Job
	err := db.NewSelect().Model(&jobResult).Where("id = ?", jobID).Scan(context.TODO())
	if err != nil {
		return Job{}, fmt.Errorf("failed to retrieve job, %v", err)
	}
	return jobResult, nil
}

func GetJobStatus(jobID string) (string, error) {
	var jobResult Job
	err := db.NewSelect().Model(&jobResult).Where("id = ?", jobID).Scan(context.TODO())
	if err != nil {
		return "", err
	}

	switch jobResult.Status {
	case JobStatusPending:
		return "Pending", nil
	case JobStatusOnRouteToPickUpLocation:
		return "On route to pick up location", nil
	case JobStatusParcelCollected:
		return "Parcel collected", nil
	case JobStatusOnRouteToDropOffLocation:
		return "On route to drop off location", nil
	case JobStatusComplete:
		return "Job complete", nil
	}

	return "", nil
}

func GetDriver(jobID string) (int, error) {
	var jobResult Job
	err := db.NewSelect().Model(&jobResult).Where("id = ?", jobID).Scan(context.TODO())
	if err != nil {
		return 0, err
	}

	return int(jobResult.Worker), nil
}

// TODO: THIS NEEDS TO BE HANDLE PART IN WITH THE DELIVERY SVC AND THE USER SVC
//func GetDriverInformation(jobID string)(WorkerInfo,error){
//	var jobResult Job
//	err := db.NewSelect().Model(&jobResult).Where("id = ?", jobID).Scan(context.TODO())
//
//	return jobResults,err
//}

func GetUserJobs(userID int) ([]Job, error) {
	jobResults := &[]Job{}
	err := db.NewSelect().Model(jobResults).Where("creator = ?", userID).Scan(context.TODO())

	return *jobResults, err
}

func GetAllJobs() ([]Job, error) {
	jobResults := &[]Job{}
	err := db.NewSelect().Model(jobResults).Scan(context.TODO())

	if jobResults == nil {
		return []Job{}, nil
	}
	return *jobResults, err
}

func UpdateJobStatus(jobID string, status int) error {
	_, err := db.NewUpdate().Model(&Job{}).Set("status = ?", status).Where("id = ?", jobID).Exec(context.TODO())
	return err
}

func GetDriverJobs(driverID int64) ([]Job, error) {
	jobResults := &[]Job{}
	err := db.NewSelect().Model(jobResults).Where("worker = ?", driverID).Scan(context.TODO())
	if err != nil {
		return []Job{}, err
	}

	if jobResults == nil {
		return []Job{}, fmt.Errorf("no driver jobs for driver: %d", driverID)
	}
	return *jobResults, nil
}

func GetJobOrderStatus(jobID string) (string, error) {
	jobResult := &Job{}

	err := db.NewSelect().Model(jobResult).Where("id = ?", jobID).Scan(context.TODO())
	if err != nil {
		return "", err
	}

	switch jobResult.OrderStatus {
	case OrderStatusOpen:
		return "open", nil
	case OrderStatusCancelled:
		return "cancelled", nil
	case OrderStatusClosed:
		return "closed", nil
	}

	return "", nil
}

func UpdateJobOrderStatus(jobID string, status int) error {
	_, err := db.NewUpdate().Model(&Job{}).Set("order_status = ?", status).Where("id = ?", jobID).Exec(context.TODO())
	return err
}

func GetOriginPostcode(jobID string) (string, error) {
	jobResult := &Job{}

	err := db.NewSelect().Model(jobResult).Where("id = ?", jobID).Scan(context.TODO())
	if err != nil {
		return "", err
	}

	if jobResult.OriginPostcode == "" {
		return "", fmt.Errorf("origin postcode was not set")
	}

	return jobResult.OriginPostcode, nil
}

func GetDestinationPostcode(jobID string) (string, error) {
	jobResult := &Job{}

	err := db.NewSelect().Model(jobResult).Where("id = ?", jobID).Scan(context.TODO())
	if err != nil {
		return "", err
	}

	if jobResult.DestinationPostcode == "" {
		return "", fmt.Errorf("destination postcode was not set")
	}

	return jobResult.DestinationPostcode, nil
}

func getRelevantJourneyPostcode(jobID string) (string, error) {
	jobResult := &Job{}

	err := db.NewSelect().Model(jobResult).Where("id = ?", jobID).Scan(context.TODO())
	if err != nil {
		return "", err
	}

	switch jobResult.Status {
	case JobStatusPending:
		return "", fmt.Errorf("%s pending", jobID)
	case JobStatusOnRouteToPickUpLocation:
		return jobResult.OriginPostcode, nil
	case JobStatusParcelCollected:
		return "", fmt.Errorf("%s parcel collected", jobID)
	case JobStatusOnRouteToDropOffLocation:
		return jobResult.DestinationPostcode, nil
	case JobStatusComplete:
		return "", fmt.Errorf("%s job completed", jobID)
	}
	return "", nil
}

func GetEstimatedJobJourneyDuration(jobID string, driverID int64) (int, error) {
	postcode, err := getRelevantJourneyPostcode(jobID)
	if err != nil {
		return 0, err
	}

	if postcode == "" {
		return 0, fmt.Errorf("empty postcode")
	}

	destinationLatitude, destinationLongitude := addresssvc.GetPostcodeCoordinates(postcode)

	coordinates, err := usersvc.GetDriverLocation(driverID)
	if err != nil {
		return 0, err
	}

	splittedCoordinates := strings.SplitAfter(coordinates, ",")
	if len(splittedCoordinates) != 2 {
		return 0, fmt.Errorf("failed to read coordinates of driver id: %s", coordinates)
	}

	driverLatitude := strings.ReplaceAll(splittedCoordinates[0], ",", "")
	driverLongitude := strings.ReplaceAll(splittedCoordinates[1], ",", "")

	duration, err := location.GetJourneyDuration(location.Point{
		Latitude:  driverLatitude,
		Longitude: driverLongitude,
	}, location.Point{
		Latitude:  destinationLatitude,
		Longitude: destinationLongitude,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to get journey duration: %v", err)
	}

	return duration, nil
}
