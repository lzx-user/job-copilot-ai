package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
	"job-copilot-backend/internal/port"
)

func TestSupabaseAnalysisRepositorySaveUsesDatabaseDefaults(t *testing.T) {
	client := testHTTPClient(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodPost || request.URL.Path != "/rest/v1/jd_analyses" ||
			request.URL.Query().Get("select") != "id" {
			t.Fatalf("request = %s %s", request.Method, request.URL)
		}
		if request.Header.Get("Authorization") != "Bearer access-token" ||
			request.Header.Get("Prefer") != "return=representation" {
			t.Fatalf("unexpected request headers")
		}
		var payload map[string]json.RawMessage
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		for _, field := range []string{"id", "created_at"} {
			if _, exists := payload[field]; exists {
				t.Errorf("insert payload must omit %s", field)
			}
		}
		var userID string
		if err := json.Unmarshal(payload["user_id"], &userID); err != nil || userID != "user-1" {
			t.Errorf("user_id = %q, error = %v", userID, err)
		}
		for _, field := range []string{"company_name", "job_title", "jd_content", "resume_summary", "skills", "match_score", "job_summary", "core_requirements", "matched_skills", "missing_skills", "resume_suggestions", "preparation_topics", "greeting_message"} {
			if _, exists := payload[field]; !exists {
				t.Errorf("insert payload missing %s", field)
			}
		}
		return &http.Response{
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"id":"11111111-1111-1111-1111-111111111111"}]`)),
		}, nil
	})
	repository, err := NewSupabaseAnalysisRepository("https://example.supabase.co", "anon-key", client)
	if err != nil {
		t.Fatal(err)
	}
	request, err := analysisdomain.NewAnalysisRequest("示例公司", "开发实习生", strings.Repeat("岗位职责与要求。", 30), "做过 Web 项目", []string{"Go"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := analysisdomain.NewAnalysisResult(analysisdomain.AnalysisResultParams{
		MatchScore: 80, JobSummary: "开发业务功能", CoreRequirements: []string{"Go"},
		MatchedSkills: []string{"Go"}, MissingSkills: []string{}, ResumeSuggestions: []string{"补充项目成果"},
		PreparationTopics: []string{"接口设计"}, GreetingMessage: "继续准备",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := port.WithAuthenticatedAccessToken(context.Background(), "access-token")
	id, err := repository.Save(ctx, "user-1", request, result)
	if err != nil || id != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("Save() id = %q, error = %v", id, err)
	}
}
