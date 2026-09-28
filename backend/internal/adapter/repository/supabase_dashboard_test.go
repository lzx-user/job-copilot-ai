package repository

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"job-copilot-backend/internal/port"
)

func TestSupabaseDashboardRepositoryGetSummary(t *testing.T) {
	client := testHTTPClient(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodPost || request.URL.Path != "/rest/v1/rpc/get_dashboard_summary" {
			t.Fatalf("request = %s %s", request.Method, request.URL.Path)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"jdCount":2,"interviewCount":1,"weeklyRecordCount":3,"profileCompleteness":75,
				"recentRecords":[{"id":"record-1","kind":"jd","companyName":"示例公司","jobTitle":"Go 开发","description":"匹配度 80 分","createdAt":"2026-09-27T01:00:00Z"}]
			}`)),
		}, nil
	})
	repository, err := NewSupabaseDashboardRepository("https://example.supabase.co", "anon-key", client)
	if err != nil {
		t.Fatal(err)
	}
	ctx := port.WithAuthenticatedAccessToken(context.Background(), "access-token")
	summary, err := repository.GetSummary(ctx, "user-1")
	if err != nil {
		t.Fatalf("GetSummary() error = %v", err)
	}
	if summary.JDCount != 2 || len(summary.RecentRecords) != 1 {
		t.Fatalf("summary = %+v", summary)
	}
}
