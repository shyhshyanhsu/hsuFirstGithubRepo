package main

import (
	"gopkg.in/inf.v0"
	"time"
)

// POS Data Model
type Business struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Hours []int `json:"hours"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type EncryptedBusiness struct {
	Business
	NameHash           string `json:"name_hash"`
	EncryptedEnvelopeKey string   `json:"encrypted_envelope_key"`
	EnvelopeKeyID        string   `json:"envelope_key_id"`
	ServiceKeyID         string   `json:"service_key_id"`
	InitializationVector string   `json:"initialization_vector"`
}

type Check struct {
	ID string `json:"id"`
	BusinessID string `json:"business_id"`
	EmployeeID string `json:"employee_id"`
	Name string `json:"name"`
	Closed bool `json:"closed"`
	ClosedAt time.Time `json:"closed_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderedItem struct {
	ID string `json:"id"`
	BusinessID string `json:"business_id"`
	EmployeeID string `json:"employee_id"`
	CheckID string `json:"check_id"`
	ItemID string `json:"item_id"`
	Name string `json:"name"`
	Cost *inf.Dec `json:"cost"`
	Price *inf.Dec `json:"price"`
	Voided bool `json:"voided"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type MenuItem struct {
	ID string `json:"id"`
	BusinessID string `json:"business_id"`
	Name string `json:"name"`
	Cost *inf.Dec `json:"cost"`
	Price *inf.Dec `json:"price"`
	//Cost json.Number `json:"cost"`
	//Price json.Number `json:"price"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Employee struct {
	ID string `json:"id"`
	BusinessID string `json:"business_id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	PayRate *inf.Dec `json:"pay_rate"`
	//PayRate json.Number `json:"pay_rate"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type LaborEntry struct {
	ID string `json:"id"`
	BusinessID string `json:"business_id"`
	EmployeeID string `json:"employee_id"`
	Name string `json:"name"`
	ClockIn time.Time `json:"clock_in"`
	ClockOut time.Time `json:"clock_out"`
	PayRate *inf.Dec `json:"pay_rate"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
