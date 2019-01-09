package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	//"github.com/lib/pq"
	"github.com/satori/go.uuid"
	"gopkg.in/inf.v0"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
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

var (
	theBusiness         = make([]Business, 0)
	theCheck            = make([]Check, 0)
	theOrderedItem      = make([]OrderedItem, 0)
	theMenuItem         = make([]MenuItem, 0)
	theEmployee         = make([]Employee, 0)
	theLaborEntry       = make([]LaborEntry, 0)
)

// MockPOSService provides mock POS operations.
type MockPOSService interface {
	Businesses(BusinessesRequest) (Business, error)
	MenuItems(string) int
	
	//Ping() (string, error)
	//CreateDB() ([]interface{}, error)
	//CreateBusinessBusiness) (string, error)
	//GetBusinessByID(string) Business, error)
	//GetAllBusinesses() ([]Business, error)
	//CreateEmployeeEmployee) (string, error)
	//GetEmployeeByID(string) Employee, error)
	//GetEmployeeByName(string) Employee, error)
	//GetAllEmployees() ([]Employee, error)
	//CreateMenuItemMenuItem) (string, error)
	//GetMenuItemByID(string) MenuItem, error)
	//GetMenuItemByName(string) MenuItem, error)
	//GetMenuItemsByBusinessID(string) ([]MenuItem, error)
	//GetAllMenuItems() ([]MenuItem, error)
	//DeleteMenuItemByID(string) error
	//UpdateMenuItemByIDMenuItem, string) (string, error)
	//GetMenuItemsByStartEndDate(start, end time.Time) ([]MenuItem, error)
	//GetMenuItemsByBusinessIDStartEndDate(businessID string, start, end time.Time) ([]MenuItem, error)
}

// mockPOSService is a concrete implementation of MockPOSService
type mockPOSService struct{}

func (s mockPOSService) Businesses(req BusinessesRequest) (Business, error) {
	if req.BusinessID == "" {
		return Business{}, ErrEmpty
	}

	bus, _ := GetBusinessByID(req.BusinessID)

	return bus, nil
}

func (mockPOSService) MenuItems(s string) int {
	return len(s)
}

//func (s reportingService) GetMenuItemsByBusinessID(id string) ([]MenuItem, error) {
//	entity, err := s.dao.GetMenuItemsByBusinessID(id)
//	if err != nil {
//		return nil, errFailedConnPostgre
//	}
//	return entity, nil
//}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

// For each method, we define request and response structs
type BusinessesRequest struct {
	Limit int `json:"limit"`
	Offset int `json:"offset"`
	BusinessID string `json:"business_id"`
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
		req := request.(BusinessesRequest)
		v, err := svc.Businesses(req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func makeMenuItemsEndpoint(svc MockPOSService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(menuItemsRequest)
		v := svc.MenuItems(req.S)
		return menuItemsResponse{v}, nil
	}
}

func init() {
	CreateDB()
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
	var request BusinessesRequest
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

func (mw loggingMiddleware) Businesses (req BusinessesRequest) (output Business, err error) {
	defer func(begin time.Time) {
		callTime := time.Now().Format(TimeFormat)
		mw.logger.Log(
			"callTime", callTime,
			"method", "businesses",
			"Limit", req.Limit,
			"Offset", req.Offset,
			"BusinessID", req.BusinessID,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.Businesses(req)
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

func CreateDB() (error) {

	businessID, err := AddBusiness()
	if err != nil {
		return err
	}
	err = AddEmployeesByBusinessID(businessID)
	if err != nil {
		return err
	}
	err = AddMenuItemsByBusinessID(businessID)
	if err != nil {
		return err
	}

	//AddLaborEntries
	employee1, err := GetEmployeeByName("JohnWayne")
	if err != nil {
		return err
	}
	err = AddChecksByBusinessIDEmployeeID(businessID, employee1.ID)
	if err != nil {
		return err
	}
	err = AddLaborEntriesByBusinessIDEmployeeID(businessID, employee1.ID)
	if err != nil {
		return err
	}
	employee2, err := GetEmployeeByName("MaryPoppins")
	if err != nil {
		return err
	}
	err = AddChecksByBusinessIDEmployeeID(businessID, employee2.ID)
	if err != nil {
		return err
	}
	err = AddLaborEntriesByBusinessIDEmployeeID(businessID, employee2.ID)
	if err != nil {
		return err
	}

	//AddOrderedItems
	item1, err := GetMenuItemByName("Buffalo Wing")
	check1, err := GetCheckByName("check1")
	AddOrderedItemsByBusinessIDEmployeeIDCheckIDItemID(businessID, employee1.ID, check1.ID, item1.ID)
	item2, err := GetMenuItemByName("Origional Recipe")
	AddOrderedItemsByBusinessIDEmployeeIDCheckIDItemID(businessID, employee2.ID, check1.ID, item2.ID)
	
	return nil
}

//func AddBusiness() (string, error) {
//	if db != nil {
//		txn, err := db.Begin()
//		if err != nil {
//			loggeFatal(erError())
//			return "", err
//		}
//
//		stmt, err := txn.Prepare(pq.CopyIn("business", "id", "name", "hours", "updated_at", "created_at")) //*database/sql.Stmt //COPY "link_test" ("url", "name") FROM STDIN
//		if err != nil {
//			loggeFatal(erError())
//			return "", err
//		}
//
//		business := Business{}
//		businessID, err := uuid.NewV4()
//		if err != nil {
//			return "", err
//		}
//		business.ID = businessID.String()
//		business.Name = "World's Best Fried Chicken"
//		business.Hours = []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22}
//		business.UpdatedAt = time.Now()
//		business.CreatedAt = time.Date(1900, 01, 02, 03, 04, 05, 123456789, time.UTC)
//
//		theBusiness = append(theBusiness, business)
//
//
//		//DB
//		_, err = stmt.Exec(business.ID, business.Name, business.Hours, business.UpdatedAt, business.CreatedAt)
//		if err != nil {
//			loggeFatal(erError())
//			return "", err
//		}
//
//		err = stmt.Close()
//		if err != nil {
//			loggeFatal(erError())
//			return "", err
//		}
//
//		err = txn.Commit()
//		if err != nil {
//			loggeFatal(erError())
//			return "", err
//		}
//
//		return business.ID, nil
//	}
//	return "", nil
//}

func AddBusiness() (string, error) {
	business := Business{}
	//businessID, err := uuid.NewV4()
	//if err != nil {
	//	return "", err
	//}
	//business.ID = businessID.String()
	business.ID = "businessID1"
	business.Name = "World's Best Fried Chicken"
	business.Hours = []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22}
	business.UpdatedAt = time.Now()
	business.CreatedAt = time.Date(1900, 01, 02, 03, 04, 05, 123456789, time.UTC)

	theBusiness = append(theBusiness, business)
	return "", nil
}

func GetBusinessByID(id string) (Business, error) {
	for _, entity := range theBusiness {
		if entity.ID == id {
			return entity, nil
		}
	}
	return Business{}, nil
}

func AddEmployeesByBusinessID(businessID string) (err error) {
	employee1 := Employee{}
	uuid1, err := uuid.NewV4()
	if err == nil {
		employee1.ID = uuid1.String()
	}
	employee1.BusinessID = businessID
	employee1.FirstName = "John"
	employee1.LastName = "Wayne"
	employee1.PayRate = inf.NewDec(2500, 2)
	employee1.UpdatedAt = time.Now()
	employee1.CreatedAt = time.Date(2015, 12, 21, 05, 34, 58, 651387237, time.UTC)

	theEmployee = append(theEmployee, employee1)
	AddChecksByBusinessIDEmployeeID(businessID, employee1.ID)
	AddLaborEntriesByBusinessIDEmployeeID(businessID, employee1.ID)

	employee2 := Employee{}
	uuid2, err := uuid.NewV4()
	if err == nil {
		employee2.ID = uuid2.String()
	}
	employee2.BusinessID = businessID
	employee2.FirstName = "Mary"
	employee2.LastName = "Poppins"
	employee2.PayRate = inf.NewDec(2100, 2)
	employee2.UpdatedAt = time.Now()
	employee2.CreatedAt = time.Date(2017, 01, 01, 05, 34, 58, 651387237, time.UTC)

	theEmployee = append(theEmployee, employee2)
	AddChecksByBusinessIDEmployeeID(businessID, employee2.ID)
	AddLaborEntriesByBusinessIDEmployeeID(businessID, employee2.ID)

	return nil
}

func GetEmployeeByID (employeeID string) (Employee, error) {
	for _, emp := range theEmployee {
		if emp.ID == employeeID {
			return emp, nil
		}
	}
	return Employee{}, nil
}

func GetEmployeeByName(name string) (Employee, error) {
	for _, entity := range theEmployee {
		if entity.FirstName + entity.LastName == name {
			return entity, nil
		}
	}
	return Employee{}, nil
}

func GetEmployeesByBusinessID(businessID string) ([]Employee, error) {
	items := make([]Employee, 0)
	for _, entity := range theEmployee {
		if entity.BusinessID == businessID {
			items = append(items, entity)
		}
	}
	return items, nil
}

func GetAllEmployees() ([]Employee, error) {
	return theEmployee, nil
}

func CreateMenuItem(m MenuItem) (string, error) {
	theMenuItem = append(theMenuItem, m)
	return m.ID, nil
}

func AddMenuItemsByBusinessID(businessID string) (err error) {
	menuitem1 := MenuItem{}
	uuid1, err := uuid.NewV4()
	if err == nil {
		menuitem1.ID = uuid1.String()
	}
	menuitem1.BusinessID = businessID
	menuitem1.Name = "Buffalo Wing"
	menuitem1.Cost = inf.NewDec(1000, 2)
	menuitem1.Price = inf.NewDec(1500, 2)
	menuitem1.UpdatedAt = time.Now()
	menuitem1.CreatedAt = time.Date(2015, 12, 21, 05, 34, 58, 651387237, time.UTC)

	theMenuItem = append(theMenuItem, menuitem1)

	menuitem2 := MenuItem{}
	uuid2, err := uuid.NewV4()
	if err == nil {
		menuitem2.ID = uuid2.String()
	}
	menuitem2.BusinessID = businessID
	menuitem2.Name = "Origional Recipe"
	menuitem2.Cost = inf.NewDec(1000, 2)
	menuitem2.Price = inf.NewDec(1500, 2)
	menuitem2.UpdatedAt = time.Now()
	menuitem2.CreatedAt = time.Date(2017, 01, 01, 05, 34, 58, 651387237, time.UTC)

	theMenuItem = append(theMenuItem, menuitem2)

	return nil
}

func GetMenuItemByID(id string) (MenuItem, error) {
	for _, entity := range theMenuItem {
		if entity.ID == id {
			return entity, nil
		}
	}
	return MenuItem{}, nil
}

func GetMenuItemByName(name string) (MenuItem, error) {
	for _, entity := range theMenuItem {
		if entity.Name == name {
			return entity, nil
		}
	}
	return MenuItem{}, nil
}

func GetMenuItemsByBusinessID(businessID string) ([]MenuItem, error) {
	items := make([]MenuItem, 0)
	for _, entity := range theMenuItem {
		if entity.BusinessID == businessID {
			items = append(items, entity)
		}
	}
	return items, nil
}

func GetAllMenuItems() ([]MenuItem, error) {
	return theMenuItem, nil
}

func DeleteMenuItemByID(id string) error {
	index := -1
	for i, entity := range theMenuItem {
		if entity.ID == id {
			index = i
			break
		}
	}
	if index > -1 {
		theMenuItem = append(theMenuItem[:index], theMenuItem[index+1:]...)
	}
	return nil
}

func UpdateMenuItemByID(m MenuItem, id string) (string, error) {
	index := -1
	for i, entity := range theMenuItem {
		if entity.ID == id {
			index = i
			break
		}
	}
	if index > -1 {
		theMenuItem = append(theMenuItem[:index], theMenuItem[index+1:]...)
	}
	theMenuItem = append(theMenuItem, m)
	return m.ID, nil
}

func GetMenuItemsByStartEndDate(start, end time.Time) ([]MenuItem, error) {
	items := make([]MenuItem, 0)
	for _, entity := range theMenuItem {
		if entity.UpdatedAt.After(start) && entity.UpdatedAt.Before(end) {
			items = append(items, entity)
		}
	}
	return items, nil
}

func GetMenuItemsByBusinessIDStartEndDate(businessID string, start, end time.Time) ([]MenuItem, error) {
	items := make([]MenuItem, 0)
	for _, entity := range theMenuItem {
		if entity.BusinessID == businessID && entity.UpdatedAt.After(start) && entity.UpdatedAt.Before(end) {
			items = append(items, entity)
		}
	}
	return items, nil
}

func AddChecksByBusinessIDEmployeeID(businessID, employeeID string) (err error) {
	check1 := Check{}
	uudi1, err := uuid.NewV4()
	if err == nil {
		check1.ID = uudi1.String()
	}
	check1.BusinessID = businessID
	check1.EmployeeID = employeeID
	check1.Name = "check1"
	check1.Closed = true
	check1.ClosedAt = time.Date(2018, 12, 23, 05, 06, 07, 123456789, time.UTC)
	check1.UpdatedAt = time.Now()
	check1.CreatedAt = time.Date(2015, 12, 21, 05, 34, 58, 651387237, time.UTC)

	theCheck = append(theCheck, check1)

	check2 := Check{}
	uudi2, err := uuid.NewV4()
	if err == nil {
		check2.ID = uudi2.String()
	}
	check2.BusinessID = businessID
	check2.EmployeeID = employeeID
	check2.Name = "check2"
	check2.Closed = false
	check2.UpdatedAt = time.Now()
	check2.CreatedAt = time.Date(2017, 01, 01, 05, 34, 58, 651387237, time.UTC)

	theCheck = append(theCheck, check2)

	return nil
}

func GetCheckByName(name string) (Check, error) {
	for _, entity := range theCheck {
		if entity.Name == name {
			return entity, nil
		}
	}
	return Check{}, nil
}

func GetAllChecks() ([]Check, error) {
	entities := []Check{}
	for _, entity := range theCheck {
		entities = append(entities, entity)
	}
	return entities, nil
}

func AddLaborEntriesByBusinessIDEmployeeID(businessID, employeeID string) (err error) {
	laborEntry1 := LaborEntry{}
	uudi1, err := uuid.NewV4()
	if err == nil {
		laborEntry1.ID = uudi1.String()
	}
	laborEntry1.BusinessID = businessID
	laborEntry1.EmployeeID = employeeID
	laborEntry1.Name = "laborEntry1"
	laborEntry1.ClockIn = time.Now()
	laborEntry1.ClockOut = time.Now().Add(time.Hour * 8)
	laborEntry1.PayRate = inf.NewDec(2500, 2)
	laborEntry1.UpdatedAt = time.Now()
	laborEntry1.CreatedAt = time.Now()

	theLaborEntry = append(theLaborEntry, laborEntry1)

	laborEntry2 := LaborEntry{}
	uudi2, err := uuid.NewV4()
	if err == nil {
		laborEntry2.ID = uudi2.String()
	}
	laborEntry2.BusinessID = businessID
	laborEntry2.EmployeeID = employeeID
	laborEntry2.Name = "laborEntry2"
	laborEntry2.ClockIn = time.Now()
	laborEntry2.ClockOut = time.Now().Add(time.Hour * 8)
	laborEntry2.PayRate = inf.NewDec(2100, 2)
	laborEntry2.UpdatedAt = time.Now()
	laborEntry2.CreatedAt = time.Now()

	theLaborEntry = append(theLaborEntry, laborEntry2)

	return nil
}

func GetAllLaborEntries() ([]LaborEntry, error) {
	entities := []LaborEntry{}
	for _, entity := range theLaborEntry {
		entities = append(entities, entity)
	}
	return entities, nil
}

func AddOrderedItemsByBusinessIDEmployeeIDCheckIDItemID(businessID, employeeID,
checkID, itemID string) (err error) {
	orderedItem1 := OrderedItem{}
	uudi1, err := uuid.NewV4()
	if err == nil {
		orderedItem1.ID = uudi1.String()
	}
	orderedItem1.BusinessID = businessID
	orderedItem1.EmployeeID = employeeID
	orderedItem1.CheckID = checkID
	orderedItem1.ItemID = itemID
	orderedItem1.Name = "orderedItem1"
	orderedItem1.Cost = inf.NewDec(1000, 2)
	orderedItem1.Price = inf.NewDec(1500, 2)
	orderedItem1.Voided = false
	orderedItem1.UpdatedAt = time.Now()
	orderedItem1.CreatedAt = time.Now()

	theOrderedItem = append(theOrderedItem, orderedItem1)

	orderedItem2 := OrderedItem{}
	uudi2, err := uuid.NewV4()
	if err == nil {
		orderedItem2.ID = uudi2.String()
	}
	orderedItem2.BusinessID = businessID
	orderedItem2.EmployeeID = employeeID
	orderedItem2.CheckID = checkID
	orderedItem2.ItemID = itemID
	orderedItem2.Name = "orderedItem2"
	orderedItem2.Cost = inf.NewDec(1000, 2)
	orderedItem2.Price = inf.NewDec(1500, 2)
	orderedItem2.Voided = true
	orderedItem2.UpdatedAt = time.Now()
	orderedItem2.CreatedAt = time.Now()

	theOrderedItem = append(theOrderedItem, orderedItem2)

	return nil
}

func GetAllOrderedItems() ([]OrderedItem, error) {
	entities := []OrderedItem{}
	for _, entity := range theOrderedItem {
		entities = append(entities, entity)
	}
	return entities, nil
}

func GetAllEntities() ([]interface{}, error) {

	var entities []interface{}

	bus, err := GetAllBusinesses()
	if err != nil {
		return nil, err
	}
	entities = append(entities, bus)

	emp, err := GetAllEmployees()
	if err != nil {
		return nil, err
	}
	entities = append(entities, emp)

	menu, err := GetAllMenuItems()
	if err != nil {
		return nil, err
	}
	entities = append(entities, menu)

	check, err := GetAllChecks()
	if err != nil {
		return nil, err
	}
	entities = append(entities, check)

	labor, err := GetAllLaborEntries()
	if err != nil {
		return nil, err
	}
	entities = append(entities, labor)

	order, err := GetAllOrderedItems()
	if err != nil {
		return nil, err
	}
	entities = append(entities, order)

	return entities, nil
}

func GetAllBusinesses() ([]Business, error) {
	return theBusiness, nil
}