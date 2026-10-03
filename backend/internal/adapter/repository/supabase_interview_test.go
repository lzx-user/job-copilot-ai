package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	interviewdomain "job-copilot-backend/internal/domain/interview"
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

func TestSupabaseInterviewMessagesUseRoundOrderForEqualTimestamps(t *testing.T) {
	client := testHTTPClient(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		if request.URL.Path != "/rest/v1/interview_messages" ||
			query.Get("order") != "round.desc,role.asc" || query.Get("limit") != "8" ||
			query.Get("session_id") != "eq.session-1" || query.Get("user_id") != "eq.user-1" {
			t.Fatalf("unexpected message request: %s", request.URL)
		}
		// 模拟按请求顺序返回的记录：上一轮回答与下一题拥有相同时间。
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`[
			{"id":"q3","role":"interviewer","round":3,"content":"模拟问题三","created_at":"2026-10-02T00:02:00Z"},
			{"id":"a2","role":"candidate","round":2,"content":"模拟回答二","created_at":"2026-10-02T00:02:00Z"},
			{"id":"q2","role":"interviewer","round":2,"content":"模拟问题二","created_at":"2026-10-02T00:01:00Z"},
			{"id":"a1","role":"candidate","round":1,"content":"模拟回答一","created_at":"2026-10-02T00:01:00Z"},
			{"id":"q1","role":"interviewer","round":1,"content":"模拟问题一","created_at":"2026-10-02T00:00:00Z"}
		]`))}, nil
	})
	repository, err := NewSupabaseInterviewRepository("https://example.supabase.co", "anon-key", client)
	if err != nil {
		t.Fatal(err)
	}
	ctx := port.WithAuthenticatedAccessToken(context.Background(), "access-token")
	messages, err := repository.findMessages(ctx, "user-1", "session-1", 8)
	if err != nil {
		t.Fatal(err)
	}
	wantIDs := []string{"q1", "a1", "q2", "a2", "q3"}
	if len(messages) != len(wantIDs) {
		t.Fatalf("message count=%d", len(messages))
	}
	for index, want := range wantIDs {
		if messages[index].ID != want {
			t.Fatalf("message %d=%s, want %s", index, messages[index].ID, want)
		}
	}
}

func TestSupabaseSaveTurnPreservesAnsweredRound(t *testing.T) {
	client := testHTTPClient(func(request *http.Request) (*http.Response, error) {
		var body struct {
			ExpectedRound int `json:"p_expected_round"`
		}
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if request.URL.Path != "/rest/v1/rpc/submit_interview_turn" || body.ExpectedRound != 1 {
			t.Fatalf("RPC path=%s, expected round=%d", request.URL.Path, body.ExpectedRound)
		}
		return &http.Response{StatusCode: http.StatusNoContent, Body: io.NopCloser(strings.NewReader(""))}, nil
	})
	repository, err := NewSupabaseInterviewRepository("https://example.supabase.co", "anon-key", client)
	if err != nil {
		t.Fatal(err)
	}
	session, err := interviewdomain.RestoreInterviewSession(interviewdomain.InterviewSessionParams{
		ID: "session-1", UserID: "user-1", AnalysisID: "analysis-1",
		Status: interviewdomain.InterviewStatusInProgress, CurrentRound: 2, MaxRounds: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	answer, err := interviewdomain.NewInterviewMessage(interviewdomain.InterviewMessageParams{
		SessionID: "session-1", UserID: "user-1", Role: interviewdomain.InterviewMessageRoleCandidate,
		Round: 1, Content: "模拟回答",
	})
	if err != nil {
		t.Fatal(err)
	}
	feedback, err := interviewdomain.NewInterviewFeedback(80, "模拟反馈", []string{"模拟优势"}, []string{"模拟建议"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := interviewdomain.NewInterviewTurnResult(feedback, "模拟下一题")
	if err != nil {
		t.Fatal(err)
	}
	ctx := port.WithAuthenticatedAccessToken(context.Background(), "access-token")
	if err := repository.SaveTurn(ctx, session, answer, result); err != nil {
		t.Fatal(err)
	}
}
