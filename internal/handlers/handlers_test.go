package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ol-ilyassov/booking-app/internal/models"
)

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
}

// {"post-search-availability", "/search-availability", "POST", []postData{
// 	{key: "start", value: "2020-01-01"},
// 	{key: "end", value: "2020-01-02"},
// }, http.StatusOK},

// {"post-search-availability-json", "/search-availability-json", "POST", []postData{
// 	{key: "start", value: "2020-01-01"},
// 	{key: "end", value: "2020-01-02"},
// }, http.StatusOK},

func TestHandlers(t *testing.T) {
	var response *http.Response
	var err error
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		response, err = ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if response.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d, got %d", e.name, e.expectedStatusCode, response.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	// #1
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, err := http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		log.Println(err)
	}
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// #2 test case: reservation is not in session
	req, err = http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		log.Println(err)
	}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// #3 test case: non-existent room
	req, err = http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		log.Println(err)
	}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 10
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

// {"post-make-reservation", "/make-reservation", "POST", []postData{
// 	{key: "first_name", value: "John"},
// 	{key: "last_name", value: "Doe"},
// 	{key: "email", value: "john.doe@gmail.com"},
// 	{key: "phone", value: "87777777777"},
// }, http.StatusOK},

func TestRepository_PostReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Test room",
		},
	}

	// #1 test case.
	strBody := "first_name=Rusaln&last_name=Jora&email=ll@ll.ru&phone=88005553535&start_date=2023-01-01&end_date=2024-01-01"
	body := strings.NewReader(strBody)

	req, _ := http.NewRequest("POST", "/make-reservation", body)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// #2 test case.
	strBody = "first_name=Rusaln&last_name=Jora&email=ll@ll.ru&phone=88005553535&start_date=2023-01-01&end_date=2024-01-01"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", nil)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code for missed reservation data in session: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// #3 test case.
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code for missed body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// #4 test case.
	strBody = "first_name=1&email=meEmail&phone=88005553535&start_date=2023-01-01&end_date=2024-01-01"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Handler PostReservation returned wrong response code when invalid name: got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// #5 test case.
	strBody = "first_name=Ruslan&last_name=Jora&email=ll@ll.ru&phone=88005553535&start_date=2023-01-01&end_date=2024-01-01"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	reservation.RoomID = 2
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostReservation returned wrong response code when insert reservation fails: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// #6 test case.
	strBody = "first_name=Ruslan&last_name=Jora&email=ll@ll.ru&phone=88005553535&start_date=2023-01-01&end_date=2024-01-01"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	reservation.RoomID = 3
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostReservation returned wrong response code when insert RoomRestriction fails: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostSearchAvailabilityJSON(t *testing.T) {
	// #1 test case
	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("room_id", "1")
	body := strings.NewReader(postedData.Encode())

	req, _ := http.NewRequest("POST", "/search-availability-json", body)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostSearchAvailabilityJSON)
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	// #2 test case
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostSearchAvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	// #3 test case
	strBody := "start_date=invalid&end_date=2050-01-02&room_id=1"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/search-availability-json", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostSearchAvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	// #4 test case
	strBody = "start_date=2050-01-01&end_date=invalid&room_id=1"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/search-availability-json", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostSearchAvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	// #5 test case
	strBody = "start_date=2050-01-01&end_date=2050-01-02&room_id=invalid"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/search-availability-json", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostSearchAvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	// #6 test case
	strBody = "start_date=2050-01-01&end_date=2050-01-02&room_id=10"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/search-availability-json", body)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostSearchAvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
}
