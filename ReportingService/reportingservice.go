package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"time"
)

var (
	errMissingParam 	= errors.New("missing param")
	errBadInput        	= errors.New("client: bad input")
	errBadCredential    = errors.New("client: HTTP Basic Auth: bad credential")
)

// ReportingService retreives the source data from POS APIs,
// and implement a reporting API to calculate and deliver a number of common metrics.
type ReportingService interface {
	reporting(BusinessesRequest) (Business, error)
}

type reportingService struct{}

func (reportingService) reporting(req BusinessesRequest) (Business, error) {
	if req.BusinessID == "" {
		return Business{}, ErrEmpty
	}
	return Business{}, nil //TODO : provides real POS operation
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

func makePOSEndpoint(svc ReportingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(BusinessesRequest)
		v, err := svc.reporting(req)
		if err != nil {
			return Business{}, nil
		}
		return v, nil
	}
}

//Deliver an http-based reporting api that implements a single /reporting endpoint
//A consumer should be able to provide a date range, a bucketing time interval, a business id, and a report type (defined below).
//The api should return a report, calculated from the source POS data.

type BusinessesRequest struct {
	Limit int `json:"limit"`
	Offset int `json:"offset"`
	BusinessID string `json:"business_id"`
}

//A consumer should be able to provide a date range, a bucketing time interval, a business id, and a report type (defined below).
type BusinessIDStartDateDaysAfterStartRequest struct {
	BusinessID 	string `json:"business_id"`
	StartDate           time.Time `json:"start_date"`
	DaysAfter 	int		`json:"days_after_start_date"`
	ReportType 	string	`json:"report_type"`
}

//limit (number) - The amount of results to return is 100 by default and and the max is 500.
//offset (number) - The amount of results to skip the default is 0.
//business_id (uuid) - The business_id of record used to constrain the results.

type businessesResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

func decodeBusinessesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request BusinessesRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeBusinessIDStartDateDaysAfterStartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	err := checkHTTPBasicAuth(r)
	if err != nil {
		return nil, err
	}
	var request BusinessIDStartDateDaysAfterStartRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	proxyRequest := BusinessesRequest {
		Limit : 500,
		Offset: 0,
		BusinessID: request.BusinessID,
	}
	return proxyRequest, nil
}

//base64 encoded "posUser:posPassword" = cG9zVXNlcjpwb3NQYXNzd29yZA==
func checkHTTPBasicAuth(r *http.Request) error {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "posUser" || pass != "posPassword" {
		return errBadCredential
	}
	return nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// ServiceMiddleware is a chainable behavior modifier for ReportingService.
type ServiceMiddleware func(ReportingService) ReportingService




