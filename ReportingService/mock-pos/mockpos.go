package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// MockPOSService provides mock POS operations.
type MockPOSService interface {
	Businesses(string) (string, error)
	MenuItems(string) int
}

// mockPOSService is a concrete implementation of MockPOSService
type mockPOSService struct{}

func (mockPOSService) Businesses(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return "/businesses : " + strings.ToUpper(s), nil //TODO : provides real POS operation
}

func (mockPOSService) MenuItems(s string) int {
	return len(s) //TODO : provides real POS operation
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

// For each method, we define request and response structs
type businessesRequest struct {
	S string `json:"s"`
}

type businessesResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

type menuItemsRequest struct {
	S string `json:"s"`
}

type menuItemsResponse struct {
	V int `json:"v"`
}

func makeBusinessesEndpoint(svc MockPOSService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(businessesRequest)
		v, err := svc.Businesses(req.S)
		if err != nil {
			return businessesResponse{v, err.Error()}, nil
		}
		return businessesResponse{v, ""}, nil
	}
}

func makeMenuItemsEndpoint(svc MockPOSService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(menuItemsRequest)
		v := svc.MenuItems(req.S)
		return menuItemsResponse{v}, nil
	}
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	var svc MockPOSService
	svc = mockPOSService{}
	svc = loggingMiddleware{logger, svc}

	businessesHandler := httptransport.NewServer(
		makeBusinessesEndpoint(svc),
		decodeBusinessesRequest,
		encodeResponse,
	)

	menuItemsHandler := httptransport.NewServer(
		makeMenuItemsEndpoint(svc),
		decodeMenuItemsRequest,
		encodeResponse,
	)

	http.Handle("/businesses", businessesHandler)
	http.Handle("/menuItems", menuItemsHandler)
	http.ListenAndServe(":8091", nil)
}

func decodeBusinessesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request businessesRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeMenuItemsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request menuItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

const TimeFormat  = "2006-01-02T15:04:05Z07:00"

type loggingMiddleware struct { //decorator
	logger log.Logger
	next   MockPOSService //Since our MockPOSService is defined as an interface, we just need to make a new type which wraps an existing MockPOSService, and performs the extra logging duties.
}

func (mw loggingMiddleware) Businesses (s string) (output string, err error) {
	defer func(begin time.Time) {
		callTime := time.Now().Format(TimeFormat)
		mw.logger.Log(
			"callTime", callTime,
			"method", "businesses",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.Businesses(s)
	return
}

func (mw loggingMiddleware) MenuItems(s string) (n int) {
	defer func(begin time.Time) {
		callTime := time.Now().Format(TimeFormat)
		mw.logger.Log(
			"callTime", callTime,
			"method", "menuItems",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = mw.next.MenuItems(s)
	return
}
