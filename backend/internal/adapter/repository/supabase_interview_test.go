package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"job-copilot-backend/internal/port"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func testHTTPClient(function roundTripFunc) *http.Client {
	return &http.Client{Transport: function}
}

func TestSupabaseInterviewRepositoryListSessions(t *testing.T) {
	client := testHTTPClient(func(request *http.Request) (*http.Response, error) {
		if request.URL.Path != "/rest/v1/rpc/list_interview_history" || request.Method != http.MethodPost {
			t.Fatalf("request = %s %s", request.Method, request.URL.Path)
		}
		if request.Header.Get("Authorization") != "Bearer access-token" {
			t.Fatalf("authorization = %q", request.Header.Get("Authorization"))
		}
		var body struct {
			Limit int `json:"p_limit"`
		}
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil || body.Limit != 20 {
			t.Fatalf("body = %+v, error = %v", body, err)
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`[{
			"session_id":"11111111-1111-1111-1111-111111111111",
			"company_name":"示例公司","job_title":"Go 开发",
			"status":"completed","current_round":5,"max_rounds":5,"overall_score":88,
			"created_at":"2026-09-27T01:00:00Z","updated_at":"2026-09-27T02:00:00Z",
			"completed_at":"2026-09-27T02:00:00Z"
		}]`))}, nil
	})

	repository, err := NewSupabaseInterviewRepository("https://example.supabase.co", "anon-key", client)
	if err != nil {
		t.Fatal(err)
	}
	ctx := port.WithAuthenticatedAccessToken(context.Background(), "access-token")
	records, err := repository.ListSessions(ctx, "user-1", 20)
	if err != nil {
		t.Fatalf("ListSessions() error = %v", err)
	}
	if len(records) != 1 || records[0].OverallScore() == nil || *records[0].OverallScore() != 88 {
		t.Fatalf("records = %+v", records)
	}
}

func TestSupabaseInterviewRepositoryDeleteSession(t *testing.T) {
	client := testHTTPClient(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodDelete || request.URL.Path != "/rest/v1/interview_sessions" {
			t.Fatalf("request = %s %s", request.Method, request.URL.Path)
		}
		if !strings.Contains(request.URL.RawQuery, "user_id=eq.user-1") || request.Header.Get("Prefer") != "return=representation" {
			t.Fatalf("query = %q, prefer = %q", request.URL.RawQuery, request.Header.Get("Prefer"))
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`[{"id":"11111111-1111-1111-1111-111111111111"}]`))}, nil
	})

	repository, err := NewSupabaseInterviewRepository("https://example.supabase.co", "anon-key", client)
	if err != nil {
		t.Fatal(err)
	}
	ctx := port.WithAuthenticatedAccessToken(context.Background(), "access-token")
	if err := repository.DeleteSession(ctx, "user-1", "11111111-1111-1111-1111-111111111111"); err != nil {
		t.Fatalf("DeleteSession() error = %v", err)
	}
}
