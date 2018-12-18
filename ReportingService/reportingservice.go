package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"strings"
)

// ReportingService retreives the source data from POS APIs,
// and implement a reporting API to calculate and deliver a number of common metrics.
type ReportingService interface {
	reporting(string) (string, error)
}

type reportingService struct{}

func (reportingService) reporting(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil //TODO : provides real POS operation
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

func makePOSEndpoint(svc ReportingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(businessesRequest)
		v, err := svc.reporting(req.S)
		if err != nil {
			return businessesResponse{v, err.Error()}, nil
		}
		return businessesResponse{v, ""}, nil
	}
}

type businessesRequest struct {
	S string `json:"s"`
}

type businessesResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

func decodeBusinessesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request businessesRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// ServiceMiddleware is a chainable behavior modifier for ReportingService.
type ServiceMiddleware func(ReportingService) ReportingService




