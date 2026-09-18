//go:build integration

package handlers_test

import (
	"net/http"
	"testing"
)

func TestEventAdminManagement(t *testing.T) {
	resetDB(t)
	register(t, "Creator", "creator@example.com")
	register(t, "Helper", "helper@example.com")
	register(t, "Other", "other@example.com")

	creator, creatorID := login(t, "creator@example.com")
	helper, helperID := login(t, "helper@example.com")
	other, _ := login(t, "other@example.com")

	eventID := createEvent(t, creator, "2027-01-01T10:00:00Z", 5)

	if rec, _ := request(t, http.MethodGet, "/api/v1/events/"+eventID+"/admins", other, nil); rec.Code != http.StatusForbidden {
		t.Fatalf("non-admin list: got %d, want 403", rec.Code)
	}
	if rec, _ := request(t, http.MethodGet, "/api/v1/events/"+eventID+"/admins", creator, nil); rec.Code != http.StatusOK {
		t.Fatalf("creator list: got %d, want 200", rec.Code)
	}

	// Only the creator can add.
	if rec, _ := request(t, http.MethodPost, "/api/v1/events/"+eventID+"/admins", other, map[string]string{"email": "helper@example.com"}); rec.Code != http.StatusForbidden {
		t.Fatalf("non-creator add: got %d, want 403", rec.Code)
	}
	if rec, _ := request(t, http.MethodPost, "/api/v1/events/"+eventID+"/admins", creator, map[string]string{"email": "nobody@example.com"}); rec.Code != http.StatusNotFound {
		t.Fatalf("unknown email: got %d, want 404", rec.Code)
	}
	if rec, _ := request(t, http.MethodPost, "/api/v1/events/"+eventID+"/admins", creator, map[string]string{"email": "creator@example.com"}); rec.Code != http.StatusConflict {
		t.Fatalf("add creator: got %d, want 409", rec.Code)
	}

	// Email is normalised: mixed case still matches.
	if rec, _ := request(t, http.MethodPost, "/api/v1/events/"+eventID+"/admins", creator, map[string]string{"email": "HELPER@example.com"}); rec.Code != http.StatusCreated {
		t.Fatalf("add helper: got %d, want 201", rec.Code)
	}
	if rec, _ := request(t, http.MethodPost, "/api/v1/events/"+eventID+"/admins", creator, map[string]string{"email": "helper@example.com"}); rec.Code != http.StatusConflict {
		t.Fatalf("duplicate add: got %d, want 409", rec.Code)
	}

	if rec, _ := request(t, http.MethodGet, "/api/v1/events/"+eventID+"/admins", helper, nil); rec.Code != http.StatusOK {
		t.Fatalf("added admin list: got %d, want 200", rec.Code)
	}

	// An added admin cannot remove anyone.
	if rec, _ := request(t, http.MethodDelete, "/api/v1/events/"+eventID+"/admins/"+helperID, helper, nil); rec.Code != http.StatusForbidden {
		t.Fatalf("added admin remove: got %d, want 403", rec.Code)
	}
	// The creator cannot be removed.
	if rec, _ := request(t, http.MethodDelete, "/api/v1/events/"+eventID+"/admins/"+creatorID, creator, nil); rec.Code != http.StatusBadRequest {
		t.Fatalf("remove creator: got %d, want 400", rec.Code)
	}

	if rec, _ := request(t, http.MethodDelete, "/api/v1/events/"+eventID+"/admins/"+helperID, creator, nil); rec.Code != http.StatusOK {
		t.Fatalf("creator remove: got %d, want 200", rec.Code)
	}
	if rec, _ := request(t, http.MethodDelete, "/api/v1/events/"+eventID+"/admins/"+helperID, creator, nil); rec.Code != http.StatusNotFound {
		t.Fatalf("remove again: got %d, want 404", rec.Code)
	}
}

func TestEventAdminVerification(t *testing.T) {
	resetDB(t)
	register(t, "Creator", "creator@example.com")
	register(t, "Helper", "helper@example.com")
	register(t, "Other", "other@example.com")

	creator, _ := login(t, "creator@example.com")
	helper, _ := login(t, "helper@example.com")
	other, _ := login(t, "other@example.com")

	eventID := createEvent(t, creator, "2027-01-01T10:00:00Z", 5)

	book := func(token string) string {
		rec, body := request(t, http.MethodPost, "/api/v1/tickets", token, map[string]string{"event_id": eventID})
		if rec.Code != http.StatusCreated {
			t.Fatalf("book: got %d", rec.Code)
		}
		code, _ := body["code"].(string)
		return code
	}
	verify := func(token, ticketCode string) int {
		rec, _ := request(t, http.MethodPost, "/api/v1/tickets/verify", token, map[string]string{"code": ticketCode})
		return rec.Code
	}

	code := book(other)

	if got := verify(helper, code); got != http.StatusForbidden {
		t.Fatalf("non-admin verify: got %d, want 403", got)
	}
	if got := verify(creator, code); got != http.StatusOK {
		t.Fatalf("creator verify: got %d, want 200", got)
	}
	if got := verify(creator, code); got != http.StatusConflict {
		t.Fatalf("second verify: got %d, want 409", got)
	}

	// Once granted admin, the helper can verify another event ticket.
	register(t, "Fourth", "fourth@example.com")
	fourth, _ := login(t, "fourth@example.com")
	code2 := book(fourth)

	if rec, _ := request(t, http.MethodPost, "/api/v1/events/"+eventID+"/admins", creator, map[string]string{"email": "helper@example.com"}); rec.Code != http.StatusCreated {
		t.Fatalf("add helper: got %d", rec.Code)
	}
	if got := verify(helper, code2); got != http.StatusOK {
		t.Fatalf("added admin verify: got %d, want 200", got)
	}
}
