package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
	"job-copilot-backend/internal/port"
)

type OpenAICompatibleAdapter struct {
	endpoint   string
	apiKey     string
	model      string
	httpClient *http.Client
}

var _ port.AIClient = (*OpenAICompatibleAdapter)(nil)

func NewOpenAICompatibleAdapter(
	baseURL string,
	apiKey string,
	model string,
	httpClient *http.Client,
) (*OpenAICompatibleAdapter, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	apiKey = strings.TrimSpace(apiKey)
	model = strings.TrimSpace(model)
	if baseURL == "" || apiKey == "" || model == "" {
		return nil, port.ErrAIUnavailable
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	return &OpenAICompatibleAdapter{
		endpoint:   baseURL + "/chat/completions",
		apiKey:     apiKey,
		model:      model,
		httpClient: httpClient,
	}, nil
}

type chatCompletionRequest struct {
	Model          string         `json:"model"`
	Messages       []chatMessage  `json:"messages"`
	Temperature    float64        `json:"temperature"`
	ResponseFormat responseFormat `json:"response_format"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func (adapter *OpenAICompatibleAdapter) AnalyzeJD(
	ctx context.Context,
	request analysisdomain.AnalysisRequest,
) (analysisdomain.AnalysisResult, error) {
	payload := chatCompletionRequest{
		Model: adapter.model,
		Messages: []chatMessage{
			{
				Role:    "system",
				Content: "你是校招岗位匹配分析助手。对比候选人经历与岗位要求，只返回 JSON：{\"matchScore\": 0到100的整数}。",
			},
			{
				Role:    "user",
				Content: buildAnalysisPrompt(request),
			},
		},
		Temperature:    0.2,
		ResponseFormat: responseFormat{Type: "json_object"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return analysisdomain.AnalysisResult{}, fmt.Errorf("encode AI request: %w", err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, adapter.endpoint, bytes.NewReader(body))
	if err != nil {
		return analysisdomain.AnalysisResult{}, fmt.Errorf("create AI request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+adapter.apiKey)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := adapter.httpClient.Do(httpRequest)
	if err != nil {
		return analysisdomain.AnalysisResult{}, errors.Join(port.ErrAIUpstream, err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		// 不读取或记录上游响应体，避免泄露请求上下文。
		_, _ = io.Copy(io.Discard, httpResponse.Body)
		return analysisdomain.AnalysisResult{}, fmt.Errorf("%w: status %d", port.ErrAIUpstream, httpResponse.StatusCode)
	}

	var completion chatCompletionResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&completion); err != nil ||
		len(completion.Choices) == 0 || strings.TrimSpace(completion.Choices[0].Message.Content) == "" {
		return analysisdomain.AnalysisResult{}, port.ErrAIInvalidResponse
	}

	return decodeAnalysisResult(completion.Choices[0].Message.Content)
}

func buildAnalysisPrompt(request analysisdomain.AnalysisRequest) string {
	return fmt.Sprintf(
		"公司：%s\n岗位：%s\n岗位 JD：\n%s\n\n候选人经历摘要：\n%s\n\n候选人技能：%s",
		request.CompanyName(),
		request.JobTitle(),
		request.Description().Content(),
		request.ResumeSummary(),
		strings.Join(request.Skills(), "、"),
	)
}

func decodeAnalysisResult(content string) (analysisdomain.AnalysisResult, error) {
	var output struct {
		MatchScore *int `json:"matchScore"`
	}
	decoder := json.NewDecoder(strings.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&output); err != nil {
		return analysisdomain.AnalysisResult{}, errors.Join(port.ErrAIInvalidResponse, err)
	}
	if err := ensureJSONEnded(decoder); err != nil {
		return analysisdomain.AnalysisResult{}, errors.Join(port.ErrAIInvalidResponse, err)
	}
	if output.MatchScore == nil {
		return analysisdomain.AnalysisResult{}, port.ErrAIInvalidResponse
	}

	result, err := analysisdomain.NewAnalysisResult(*output.MatchScore)
	if err != nil {
		return analysisdomain.AnalysisResult{}, errors.Join(port.ErrAIInvalidResponse, err)
	}
	return result, nil
}

func ensureJSONEnded(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}
	return nil
}
