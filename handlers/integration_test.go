//go:build integration

package handlers_test

import (
	"arjunomray/eventer/config"
	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/router"
	"arjunomray/eventer/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var engine *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "test configuration error:", err)
		os.Exit(1)
	}
	utils.SetJWTSecret(cfg.JWTSecret)

	db, err := dbConfig.Connect(cfg.Database)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot connect to the test database (is Postgres running and migrated?):", err)
		os.Exit(1)
	}
	dbConfig.Db = db
	engine = router.New(cfg)

	os.Exit(m.Run())
}

func resetDB(t *testing.T) {
	t.Helper()
	if err := dbConfig.Db.Exec("TRUNCATE tickets, events, users RESTART IDENTITY CASCADE").Error; err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

func request(t *testing.T, method, path, token string, body any) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	decoded := map[string]any{}
	_ = json.Unmarshal(rec.Body.Bytes(), &decoded)
	return rec, decoded
}

func register(t *testing.T, name, email string) {
	t.Helper()
	rec, _ := request(t, http.MethodPost, "/api/v1/register", "", map[string]string{
		"name": name, "email": email, "password": "password123",
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("register %s: got %d, body %s", email, rec.Code, rec.Body.String())
	}
}

func login(t *testing.T, email string) (token, userID string) {
	t.Helper()
	rec, body := request(t, http.MethodPost, "/api/v1/login", "", map[string]string{
		"email": email, "password": "password123",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("login %s: got %d, body %s", email, rec.Code, rec.Body.String())
	}
	token, _ = body["token"].(string)
	userID, _ = body["userId"].(string)
	if token == "" {
		t.Fatalf("login returned no token: %v", body)
	}
	return token, userID
}

func createEvent(t *testing.T, token, when string, seats int) string {
	t.Helper()
	rec, body := request(t, http.MethodPost, "/api/v1/events", token, map[string]any{
		"name": "Test Event", "time": when, "location": "Hall",
		"priceCents": 500, "seatsAvailable": seats, "tags": []string{"test"},
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create event: got %d, body %s", rec.Code, rec.Body.String())
	}
	id, _ := body["id"].(string)
	return id
}

func TestRegisterLoginProtected(t *testing.T) {
	resetDB(t)

	register(t, "Alice", "alice@example.com")

	rec, _ := request(t, http.MethodPost, "/api/v1/register", "", map[string]string{
		"name": "Alice", "email": "ALICE@example.com", "password": "password123",
	})
	if rec.Code != http.StatusConflict {
		t.Fatalf("duplicate register: got %d, want 409", rec.Code)
	}

	token, userID := login(t, "alice@example.com")

	if rec, _ := request(t, http.MethodGet, "/api/v1/events", "", nil); rec.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated events: got %d, want 401", rec.Code)
	}
	if rec, _ := request(t, http.MethodGet, "/api/v1/events", token, nil); rec.Code != http.StatusOK {
		t.Fatalf("authenticated events: got %d, want 200", rec.Code)
	}

	rec, body := request(t, http.MethodGet, "/api/v1/users/"+userID, token, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get profile: got %d", rec.Code)
	}
	if _, leaked := body["password"]; leaked {
		t.Fatal("profile response leaked the password field")
	}
}

func TestEventValidation(t *testing.T) {
	resetDB(t)
	register(t, "Alice", "alice@example.com")
	token, _ := login(t, "alice@example.com")

	cases := []struct {
		name string
		body map[string]any
	}{
		{"past", map[string]any{"name": "E", "time": "2020-01-01T10:00:00Z", "location": "H", "priceCents": 0, "seatsAvailable": 1}},
		{"negative price", map[string]any{"name": "E", "time": "2027-01-01T10:00:00Z", "location": "H", "priceCents": -1, "seatsAvailable": 1}},
		{"zero seats", map[string]any{"name": "E", "time": "2027-01-01T10:00:00Z", "location": "H", "priceCents": 0, "seatsAvailable": 0}},
		{"missing name", map[string]any{"time": "2027-01-01T10:00:00Z", "location": "H", "priceCents": 0, "seatsAvailable": 1}},
	}
	for _, tc := range cases {
		if rec, _ := request(t, http.MethodPost, "/api/v1/events", token, tc.body); rec.Code != http.StatusBadRequest {
			t.Errorf("%s: got %d, want 400", tc.name, rec.Code)
		}
	}

	if id := createEvent(t, token, "2027-01-01T10:00:00Z", 1); id == "" {
		t.Fatal("valid event was not created")
	}
}

func TestBookingFlow(t *testing.T) {
	resetDB(t)
	register(t, "Alice", "alice@example.com")
	register(t, "Bob", "bob@example.com")
	alice, _ := login(t, "alice@example.com")
	bob, _ := login(t, "bob@example.com")

	eventID := createEvent(t, alice, "2027-01-01T10:00:00Z", 1)

	rec, _ := request(t, http.MethodPost, "/api/v1/tickets", bob, map[string]string{"event_id": eventID})
	if rec.Code != http.StatusCreated {
		t.Fatalf("bob booking: got %d, body %s", rec.Code, rec.Body.String())
	}

	rec, _ = request(t, http.MethodPost, "/api/v1/tickets", alice, map[string]string{"event_id": eventID})
	if rec.Code != http.StatusConflict {
		t.Fatalf("sold-out booking: got %d, want 409", rec.Code)
	}

	rec, _ = request(t, http.MethodPost, "/api/v1/tickets", bob, map[string]string{"event_id": eventID})
	if rec.Code != http.StatusConflict {
		t.Fatalf("duplicate booking: got %d, want 409", rec.Code)
	}

	rec, body := request(t, http.MethodGet, "/api/v1/events/"+eventID, alice, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get event: got %d", rec.Code)
	}
	if seats, _ := body["seatsAvailable"].(float64); seats != 0 {
		t.Fatalf("seats available: got %v, want 0", body["seatsAvailable"])
	}

	rec, _ = request(t, http.MethodGet, "/api/v1/tickets/me", bob, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get my tickets: got %d", rec.Code)
	}
}

func TestOwnership(t *testing.T) {
	resetDB(t)
	register(t, "Alice", "alice@example.com")
	register(t, "Bob", "bob@example.com")
	alice, aliceID := login(t, "alice@example.com")
	bob, _ := login(t, "bob@example.com")

	if rec, _ := request(t, http.MethodGet, "/api/v1/users/"+aliceID, bob, nil); rec.Code != http.StatusForbidden {
		t.Fatalf("bob reading alice: got %d, want 403", rec.Code)
	}
	if rec, _ := request(t, http.MethodGet, "/api/v1/users/"+aliceID+"/tickets", bob, nil); rec.Code != http.StatusForbidden {
		t.Fatalf("bob reading alice tickets: got %d, want 403", rec.Code)
	}
	if rec, _ := request(t, http.MethodGet, "/api/v1/users/"+aliceID, alice, nil); rec.Code != http.StatusOK {
		t.Fatalf("alice reading alice: got %d, want 200", rec.Code)
	}
}
