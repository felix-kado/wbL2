package main

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// событие
type Event struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Date   time.Time `json:"date"`
	Title  string    `json:"title"`
}

// In-memory хранилище для events
var events = make(map[int]Event)
var nextID = 1

// Middleware для логов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func parseIntParam(r *http.Request, key string) (int, error) {
	return strconv.Atoi(r.FormValue(key))
}

func parseDateParam(r *http.Request, key string) (time.Time, error) {
	return time.Parse("2006-01-02", r.FormValue(key))
}

// Обработчики
func createEventHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := parseIntParam(r, "user_id")
	if err != nil {
		http.Error(w, `{"error": "Invalid user_id"}`, http.StatusBadRequest)
		return
	}

	date, err := parseDateParam(r, "date")
	if err != nil {
		http.Error(w, `{"error": "Invalid date"}`, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
		return
	}

	event := Event{
		ID:     nextID,
		UserID: userID,
		Date:   date,
		Title:  title,
	}
	events[nextID] = event
	nextID++

	writeJSON(w, http.StatusOK, map[string]string{"result": "Event created"})
}

func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseIntParam(r, "id")
	if err != nil {
		http.Error(w, `{"error": "Invalid event id"}`, http.StatusBadRequest)
		return
	}

	event, exists := events[eventID]
	if !exists {
		http.Error(w, `{"error": "Event not found"}`, http.StatusServiceUnavailable)
		return
	}

	userID, err := parseIntParam(r, "user_id")
	if err != nil {
		http.Error(w, `{"error": "Invalid user_id"}`, http.StatusBadRequest)
		return
	}

	date, err := parseDateParam(r, "date")
	if err != nil {
		http.Error(w, `{"error": "Invalid date"}`, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
		return
	}

	event.UserID = userID
	event.Date = date
	event.Title = title
	events[eventID] = event

	writeJSON(w, http.StatusOK, map[string]string{"result": "Event updated"})
}

func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseIntParam(r, "id")
	if err != nil {
		http.Error(w, `{"error": "Invalid event id"}`, http.StatusBadRequest)
		return
	}

	if _, exists := events[eventID]; !exists {
		http.Error(w, `{"error": "Event not found"}`, http.StatusServiceUnavailable)
		return
	}

	delete(events, eventID)
	writeJSON(w, http.StatusOK, map[string]string{"result": "Event deleted"})
}

func getEventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseDateParam(r, "date")
	if err != nil {
		http.Error(w, `{"error": "Invalid date"}`, http.StatusBadRequest)
		return
	}

	var result []Event
	for _, event := range events {
		if event.Date.Format("2006-01-02") == date.Format("2006-01-02") {
			result = append(result, event)
		}
	}

	writeJSON(w, http.StatusOK, result)
}

func getEventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseDateParam(r, "date")
	if err != nil {
		http.Error(w, `{"error": "Invalid date"}`, http.StatusBadRequest)
		return
	}

	startOfWeek := date.AddDate(0, 0, -int(date.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	var result []Event
	for _, event := range events {
		if event.Date.After(startOfWeek) && event.Date.Before(endOfWeek) {
			result = append(result, event)
		}
	}

	writeJSON(w, http.StatusOK, result)
}

func getEventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseDateParam(r, "date")
	if err != nil {
		http.Error(w, `{"error": "Invalid date"}`, http.StatusBadRequest)
		return
	}

	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	var result []Event
	for _, event := range events {
		if event.Date.After(startOfMonth) && event.Date.Before(endOfMonth) {
			result = append(result, event)
		}
	}

	writeJSON(w, http.StatusOK, result)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", createEventHandler)
	mux.HandleFunc("/update_event", updateEventHandler)
	mux.HandleFunc("/delete_event", deleteEventHandler)
	mux.HandleFunc("/events_for_day", getEventsForDayHandler)
	mux.HandleFunc("/events_for_week", getEventsForWeekHandler)
	mux.HandleFunc("/events_for_month", getEventsForMonthHandler)

	loggedMux := loggingMiddleware(mux)

	port := ":8080"
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(port, loggedMux))
}
