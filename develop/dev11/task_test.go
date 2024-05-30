package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func resetEvents() {
	events = make(map[int]Event)
	nextID = 1
}

func TestCreateEventHandler(t *testing.T) {
	resetEvents()

	form := url.Values{}
	form.Add("user_id", "1")
	form.Add("date", "2024-05-30")
	form.Add("title", "Test Event")

	req, err := http.NewRequest("POST", "/create_event", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createEventHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"result":"Event created"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdateEventHandler(t *testing.T) {
	resetEvents()

	form := url.Values{}
	form.Add("user_id", "1")
	form.Add("date", "2024-05-30")
	form.Add("title", "Test Event")

	req, err := http.NewRequest("POST", "/create_event", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createEventHandler)
	handler.ServeHTTP(rr, req)

	form = url.Values{}
	form.Add("id", "1")
	form.Add("user_id", "1")
	form.Add("date", "2024-06-01")
	form.Add("title", "Updated Event")

	req, err = http.NewRequest("POST", "/update_event", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(updateEventHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"result":"Event updated"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteEventHandler(t *testing.T) {
	resetEvents()

	form := url.Values{}
	form.Add("user_id", "1")
	form.Add("date", "2024-05-30")
	form.Add("title", "Test Event")

	req, err := http.NewRequest("POST", "/create_event", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createEventHandler)
	handler.ServeHTTP(rr, req)

	form = url.Values{}
	form.Add("id", "1")

	req, err = http.NewRequest("POST", "/delete_event", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(deleteEventHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"result":"Event deleted"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetEventsForDayHandler(t *testing.T) {
	resetEvents()

	form := url.Values{}
	form.Add("user_id", "1")
	form.Add("date", "2024-05-30")
	form.Add("title", "Test Event")

	req, err := http.NewRequest("POST", "/create_event", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createEventHandler)
	handler.ServeHTTP(rr, req)

	req, err = http.NewRequest("GET", "/events_for_day?date=2024-05-30", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(getEventsForDayHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"user_id":1,"date":"2024-05-30T00:00:00Z","title":"Test Event"}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetEventsForWeekHandler(t *testing.T) {
	resetEvents()

	form := url.Values{}
	form.Add("user_id", "1")
	form.Add("date", "2024-05-30")
	form.Add("title", "Test Event")

	req, err := http.NewRequest("POST", "/create_event", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createEventHandler)
	handler.ServeHTTP(rr, req)

	req, err = http.NewRequest("GET", "/events_for_week?date=2024-05-30", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(getEventsForWeekHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"user_id":1,"date":"2024-05-30T00:00:00Z","title":"Test Event"}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetEventsForMonthHandler(t *testing.T) {
	resetEvents()

	form := url.Values{}
	form.Add("user_id", "1")
	form.Add("date", "2024-05-30")
	form.Add("title", "Test Event")

	req, err := http.NewRequest("POST", "/create_event", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createEventHandler)
	handler.ServeHTTP(rr, req)

	req, err = http.NewRequest("GET", "/events_for_month?date=2024-05-30", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(getEventsForMonthHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"user_id":1,"date":"2024-05-30T00:00:00Z","title":"Test Event"}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
