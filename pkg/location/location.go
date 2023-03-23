package location

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var GoogleAPIKey string

type DistanceMatrix struct {
	DestinationAddresses []string `json:"destination_addresses"`
	OriginAddresses      []string `json:"origin_addresses"`
	Rows                 []struct {
		Elements []struct {
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

type Point struct {
	Latitude  string
	Longitude string
}

var client = &http.Client{
	Timeout: 30 * time.Second,
}

// GetJourneyDuration retrives the duration of a journey in seconds
func GetJourneyDuration(origin, destination Point) (int, error) {
	GoogleAPIKey = "AIzaSyBrnQzp90T6jbBpRlTOLfRJAisDE11Q53E"
	url := urlBuilder(origin, destination)

	req, err := http.NewRequest("GET", url, nil)

	res, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to do request: %v", err)
	}
	defer res.Body.Close()

	var resp DistanceMatrix
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return 0, fmt.Errorf("failed to decode request: %v", err)
	}

	if len(resp.Rows) > 0 {
		if len(resp.Rows[0].Elements) > 0 {
			return resp.Rows[0].Elements[0].Duration.Value, nil
		}
	}

	return 0, fmt.Errorf("no rows returned")
}

func urlBuilder(origin, destination Point) string {
	url := "https://maps.googleapis.com/maps/api/distancematrix/json?"
	origins := "origins=" + origin.Latitude + "," + origin.Longitude
	destinations := "&destinations=" + destination.Latitude + "," + destination.Longitude
	key := "&key=" + GoogleAPIKey

	return url + origins + destinations + key
}
