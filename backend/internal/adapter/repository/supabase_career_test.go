package repository

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"job-copilot-backend/internal/port"
)

func TestDecodeCareerRepositoryErrorRecognizesMissingSchema(t *testing.T) {
	response := &http.Response{
		StatusCode: http.StatusNotFound,
		Body: io.NopCloser(strings.NewReader(`{
			"code":"PGRST205",
			"message":"Could not find the table 'public.resumes' in the schema cache"
		}`)),
	}

	err := decodeCareerRepositoryError(response)
	if !errors.Is(err, port.ErrRepositorySchema) {
		t.Fatalf("decodeCareerRepositoryError() error = %v, want schema error", err)
	}
}
