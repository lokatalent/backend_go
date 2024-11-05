package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	// "strings"
	"time"
)

var baseURL = "https://maps.googleapis.com/maps/api/distancematrix/json?"

const timeout = 10 * time.Second // timeout for api request

// Google distance matrix api response structure
type DistanceMatrixResponse struct {
	DestinationAddr []string `json:"destination_addresses"`
	OriginAddr      []string `json:"origin_addresses"`
	Rows            []DistanceMatrixRow
	Status          string
	ErrMsg          string `json:"error_message,omitempty"`
}

type DistanceMatrixRow struct {
	Elements []DistanceMatrixElement
}

type DistanceMatrixElement struct {
	Status          string
	Distance        DistanceTextValueObject `json:"distance,omitempty"`
	Duration        DistanceTextValueObject `json:"duration,omitempty"`
	TrafficDuration DistanceTextValueObject `json:"duration_in_traffic,omitempty"`
	Fare            DistanceFare            `json:"fare,omitempty"`
}

type DistanceTextValueObject struct {
	Text  string
	Value float32
}

type DistanceFare struct {
	Currency string
	Text     string
	Value    float32
}

// Fetch distance information using the origin and destination.
func fetchDistanceInfo(origin, destination, apiKey string) (DistanceMatrixResponse, error) {
	escapedOrigin := url.QueryEscape(origin)
	escapedDestination := url.QueryEscape(destination)
	URL := fmt.Sprintf(
		"%slanguage=en&key=%s&origins=%s&destinations=%s",
		baseURL, apiKey, escapedOrigin, escapedDestination)

	// Retry for 10 seconds in case of any error while fetching traffic
	// informatin from Google API.
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		response, err := http.Get(URL)
		if err != nil {
			continue // retry
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			var distResponse DistanceMatrixResponse
			err := json.NewDecoder(response.Body).Decode(&distResponse)
			if err != nil {
				return DistanceMatrixResponse{}, fmt.Errorf(
					"Error decoding response. %v", err)
			}
			return distResponse, nil
		}
	}
	return DistanceMatrixResponse{}, fmt.Errorf(
		"Error getting traffic information after %s.",
		timeout)
}
