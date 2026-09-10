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
				Content: analysisSystemPrompt,
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

const analysisSystemPrompt = `你是校招岗位匹配分析助手。请根据候选人信息与岗位要求给出客观、具体、可执行的分析。
只返回一个 JSON 对象，不要返回 Markdown、代码块或解释。必须且只能包含以下字段：
{
  "matchScore": 0到100的整数,
  "jobSummary": "岗位核心职责概述",
  "coreRequirements": ["岗位核心要求"],
  "matchedSkills": ["候选人已匹配的技能或经历"],
  "missingSkills": ["候选人欠缺或材料中未体现的技能"],
  "resumeSuggestions": ["具体的简历修改建议"],
  "preparationTopics": ["面试前应准备的主题"],
  "greetingMessage": "面向候选人的简短总结"
}
coreRequirements、resumeSuggestions、preparationTopics 至少包含一项。matchedSkills 和 missingSkills 没有内容时返回空数组。不要臆造候选人经历。`

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
		MatchScore        *int      `json:"matchScore"`
		JobSummary        *string   `json:"jobSummary"`
		CoreRequirements  *[]string `json:"coreRequirements"`
		MatchedSkills     *[]string `json:"matchedSkills"`
		MissingSkills     *[]string `json:"missingSkills"`
		ResumeSuggestions *[]string `json:"resumeSuggestions"`
		PreparationTopics *[]string `json:"preparationTopics"`
		GreetingMessage   *string   `json:"greetingMessage"`
	}
	decoder := json.NewDecoder(strings.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&output); err != nil {
		return analysisdomain.AnalysisResult{}, errors.Join(port.ErrAIInvalidResponse, err)
	}
	if err := ensureJSONEnded(decoder); err != nil {
		return analysisdomain.AnalysisResult{}, errors.Join(port.ErrAIInvalidResponse, err)
	}
	if output.MatchScore == nil || output.JobSummary == nil || output.CoreRequirements == nil ||
		output.MatchedSkills == nil || output.MissingSkills == nil || output.ResumeSuggestions == nil ||
		output.PreparationTopics == nil || output.GreetingMessage == nil {
		return analysisdomain.AnalysisResult{}, port.ErrAIInvalidResponse
	}

	result, err := analysisdomain.NewAnalysisResult(analysisdomain.AnalysisResultParams{
		MatchScore:        *output.MatchScore,
		JobSummary:        *output.JobSummary,
		CoreRequirements:  *output.CoreRequirements,
		MatchedSkills:     *output.MatchedSkills,
		MissingSkills:     *output.MissingSkills,
		ResumeSuggestions: *output.ResumeSuggestions,
		PreparationTopics: *output.PreparationTopics,
		GreetingMessage:   *output.GreetingMessage,
	})
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
