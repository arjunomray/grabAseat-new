//go:build integration

package handlers_test

import (
	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/models"
	"arjunomray/eventer/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func bookConcurrently(t *testing.T, eventID string, tokens []string) map[int]int {
	t.Helper()

	var wg sync.WaitGroup
	statuses := make(chan int, len(tokens))

	for _, token := range tokens {
		wg.Add(1)
		go func(token string) {
			defer wg.Done()
			body, _ := json.Marshal(map[string]string{"event_id": eventID})
			req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()
			engine.ServeHTTP(rec, req)
			statuses <- rec.Code
		}(token)
	}

	wg.Wait()
	close(statuses)

	counts := map[int]int{}
	for code := range statuses {
		counts[code]++
	}
	return counts
}

func mintTokens(t *testing.T, n int) []string {
	t.Helper()

	tokens := make([]string, n)
	for i := range tokens {
		user := models.User{
			ID:       uuid.New(),
			Name:     fmt.Sprintf("user-%d", i),
			Email:    fmt.Sprintf("race-%d-%d@example.com", i, uuid.New().ID()),
			Password: "x",
		}
		if err := dbConfig.Db.Create(&user).Error; err != nil {
			t.Fatalf("create user: %v", err)
		}
		token, err := utils.GenerateJWT(user.ID)
		if err != nil {
			t.Fatalf("generate token: %v", err)
		}
		tokens[i] = token
	}
	return tokens
}

func TestConcurrentBookingDoesNotOversell(t *testing.T) {
	resetDB(t)
	register(t, "Host", "host@example.com")
	host, _ := login(t, "host@example.com")

	const (
		seats    = 20
		bookers  = 60
		wantCode = http.StatusConflict
	)

	eventID := createEvent(t, host, "2027-01-01T10:00:00Z", seats)
	counts := bookConcurrently(t, eventID, mintTokens(t, bookers))

	if counts[http.StatusCreated] != seats {
		t.Fatalf("created=%d, want %d (all statuses: %v)", counts[http.StatusCreated], seats, counts)
	}
	if counts[http.StatusConflict] != bookers-seats {
		t.Fatalf("conflict=%d, want %d (all statuses: %v)", counts[wantCode], bookers-seats, counts)
	}

	var ticketCount int64
	dbConfig.Db.Model(&models.Ticket{}).Where("event_id = ?", eventID).Count(&ticketCount)
	if ticketCount != seats {
		t.Fatalf("tickets in database=%d, want %d", ticketCount, seats)
	}

	var event models.Event
	if err := dbConfig.Db.First(&event, "id = ?", eventID).Error; err != nil {
		t.Fatalf("load event: %v", err)
	}
	if event.SeatsAvailable != 0 {
		t.Fatalf("seats_available=%d, want 0", event.SeatsAvailable)
	}
}

func TestConcurrentDuplicateBooking(t *testing.T) {
	resetDB(t)
	register(t, "Host", "host@example.com")
	host, _ := login(t, "host@example.com")

	eventID := createEvent(t, host, "2027-01-01T10:00:00Z", 10)
	token := mintTokens(t, 1)[0]

	tokens := make([]string, 10)
	for i := range tokens {
		tokens[i] = token
	}
	counts := bookConcurrently(t, eventID, tokens)

	if counts[http.StatusCreated] != 1 {
		t.Fatalf("created=%d, want exactly 1 (all statuses: %v)", counts[http.StatusCreated], counts)
	}

	var ticketCount int64
	dbConfig.Db.Model(&models.Ticket{}).Where("event_id = ?", eventID).Count(&ticketCount)
	if ticketCount != 1 {
		t.Fatalf("tickets in database=%d, want 1", ticketCount)
	}
}

func verifyConcurrently(t *testing.T, token string, codes []string) map[int]int {
	t.Helper()

	var wg sync.WaitGroup
	statuses := make(chan int, len(codes))

	for _, code := range codes {
		wg.Add(1)
		go func(code string) {
			defer wg.Done()
			body, _ := json.Marshal(map[string]string{"code": code})
			req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets/verify", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()
			engine.ServeHTTP(rec, req)
			statuses <- rec.Code
		}(code)
	}

	wg.Wait()
	close(statuses)

	counts := map[int]int{}
	for code := range statuses {
		counts[code]++
	}
	return counts
}

func TestConcurrentVerificationOnlyOnce(t *testing.T) {
	resetDB(t)
	register(t, "Creator", "creator@example.com")
	register(t, "Booker", "booker@example.com")
	creator, _ := login(t, "creator@example.com")
	booker, _ := login(t, "booker@example.com")

	eventID := createEvent(t, creator, "2027-01-01T10:00:00Z", 5)
	rec, body := request(t, http.MethodPost, "/api/v1/tickets", booker, map[string]string{"event_id": eventID})
	if rec.Code != http.StatusCreated {
		t.Fatalf("book: got %d", rec.Code)
	}
	code, _ := body["code"].(string)

	const scanners = 20
	codes := make([]string, scanners)
	for i := range codes {
		codes[i] = code
	}

	counts := verifyConcurrently(t, creator, codes)

	if counts[http.StatusOK] != 1 {
		t.Fatalf("successful verifications=%d, want exactly 1 (all statuses: %v)", counts[http.StatusOK], counts)
	}
	if counts[http.StatusConflict] != scanners-1 {
		t.Fatalf("conflicts=%d, want %d (all statuses: %v)", counts[http.StatusConflict], scanners-1, counts)
	}
}
