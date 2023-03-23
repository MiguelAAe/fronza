package deliverysvc

import (
	"time"

	"github.com/uptrace/bun"
)

func Init(dbconn *bun.DB) {
	db = dbconn
}

var db *bun.DB

type Job struct {
	ID               string    `bun:"id,pk"`
	ShortID          string    `bun:"short_id"`
	CreateTime       time.Time `bun:"create_time"`
	LastTimeModified time.Time `bun:"last_time_modified"`
	TrackingURL      string    `bun:"tracking_url"`
	Creator          int64     `bun:"creator"`
	Worker           int64     `bun:"worker,nullzero"`
	Status           int       `bun:"status"`
	OrderStatus      int       `bun:"order_status"`
	// Destination Details
	OriginCompanyName       string  `bun:"origin_company_name"`
	OriginFirstName         string  `bun:"origin_first_name"`
	OriginSecondName        string  `bun:"origin_second_name"`
	OriginPhoneNumber       string  `bun:"origin_phone_number"`
	OriginEmailAddress      string  `bun:"origin_email_address"`
	OriginFirstLineAddress  string  `bun:"origin_first_line_address"`
	OriginSecondLineAddress string  `bun:"origin_second_line_address"`
	OriginThirdLineAddress  string  `bun:"origin_third_line_address"`
	OriginTown              string  `bun:"origin_town"`
	OriginCity              string  `bun:"origin_city"`
	OriginPostcode          string  `bun:"origin_postcode"`
	OriginLatitude          float32 `bun:"origin_latitude"`
	OriginLongitude         float32 `bun:"origin_longitude"`
	OriginNotes             string  `bun:"origin_notes"`
	//Origin Details
	DestinationCompanyName       string  `bun:"destination_company_name"`
	DestinationFirstName         string  `bun:"destination_first_name"`
	DestinationSecondName        string  `bun:"destination_second_name"`
	DestinationPhoneNumber       string  `bun:"destination_phone_number"`
	DestinationEmailAddress      string  `bun:"destination_email_address"`
	DestinationFirstLineAddress  string  `bun:"destination_first_line_address"`
	DestinationSecondLineAddress string  `bun:"destination_second_line_address"`
	DestinationThirdLineAddress  string  `bun:"destination_third_line_address"`
	DestinationTown              string  `bun:"destination_town"`
	DestinationCity              string  `bun:"destination_city"`
	DestinationPostcode          string  `bun:"destination_postcode"`
	DestinationLatitude          float32 `bun:"destination_latitude"`
	DestinationLongitude         float32 `bun:"destination_longitude"`
	DestinationNotes             string  `bun:"destination_notes"`

	//courier notes
	WorkerNotes string `bun:"worker_notes"`
}
