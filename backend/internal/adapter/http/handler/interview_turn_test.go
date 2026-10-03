package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/middleware"
	interviewapp "job-copilot-backend/internal/application/interview"
	interviewdomain "job-copilot-backend/internal/domain/interview"
	"job-copilot-backend/internal/port"
)

const testTurnSessionID = "11111111-1111-1111-1111-111111111111"

// 以下替身只模拟认证、模型和持久化，不访问真实账号或外部服务。
type turnAuthStub struct{}

func (turnAuthStub) VerifyAccessToken(context.Context, string) (port.AuthenticatedUser, error) {
	return port.AuthenticatedUser{UserID: "test-user"}, nil
}

type turnRepositoryStub struct {
	port.InterviewRepository
	session   interviewdomain.InterviewSession
	messages  []interviewdomain.TranscriptMessage
	saves     int
	saveError error
}

func (stub *turnRepositoryStub) FindTurnContext(context.Context, string, string) (interviewdomain.InterviewTurnContext, error) {
	return interviewdomain.InterviewTurnContext{Session: stub.session, Messages: stub.messages}, nil
}

func (stub *turnRepositoryStub) SaveTurn(_ context.Context, session interviewdomain.InterviewSession, answer interviewdomain.InterviewMessage, _ interviewdomain.InterviewTurnResult) error {
	if stub.saveError != nil {
		return stub.saveError
	}
	if answer.Round() != stub.session.CurrentRound() {
		return port.ErrRepositoryConflict
	}
	stub.session = session
	stub.messages = append(stub.messages, interviewdomain.TranscriptMessage{
		Role: interviewdomain.InterviewMessageRoleCandidate, Round: answer.Round(),
	})
	stub.saves++
	return nil
}

type turnAIStub struct {
	port.AIClient
	calls int
	err   error
}

func (stub *turnAIStub) EvaluateInterviewAnswer(_ context.Context, prompt interviewdomain.InterviewTurnPrompt) (interviewdomain.InterviewTurnResult, error) {
	stub.calls++
	if stub.err != nil {
		return interviewdomain.InterviewTurnResult{}, stub.err
	}
	feedback, err := interviewdomain.NewInterviewFeedback(80, "模拟反馈", []string{"模拟优势"}, []string{"模拟建议"})
	if err != nil {
		return interviewdomain.InterviewTurnResult{}, err
	}
	if !prompt.GenerateNextQuestion {
		return interviewdomain.NewFinalInterviewTurnResult(feedback), nil
	}
	return interviewdomain.NewInterviewTurnResult(feedback, "模拟下一题")
}

func newTurnTestRouter(t *testing.T, round int) (*gin.Engine, *turnRepositoryStub, *turnAIStub) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	session, err := interviewdomain.RestoreInterviewSession(interviewdomain.InterviewSessionParams{
		ID: testTurnSessionID, UserID: "test-user", AnalysisID: "test-analysis",
		Status: interviewdomain.InterviewStatusInProgress, CurrentRound: round, MaxRounds: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	repository := &turnRepositoryStub{session: session}
	ai := &turnAIStub{}
	service, err := interviewapp.NewTurnService(ai, repository)
	if err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.POST("/turn", middleware.Authenticate(turnAuthStub{}), NewInterviewTurnHandler(service).Handle)
	return router, repository, ai
}

func sendTestTurn(router *gin.Engine, roundField string) *httptest.ResponseRecorder {
	body := fmt.Sprintf(`{"sessionId":%q,"answer":"模拟回答"%s}`, testTurnSessionID, roundField)
	request := httptest.NewRequest(http.MethodPost, "/turn", strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer mock")
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestInterviewTurnRejectsInvalidOrStaleRoundBeforeAI(t *testing.T) {
	for _, test := range []struct {
		name       string
		roundField string
		wantStatus int
	}{
		{"missing", "", http.StatusBadRequest},
		{"zero", `,"expectedRound":0`, http.StatusBadRequest},
		{"negative", `,"expectedRound":-1`, http.StatusBadRequest},
		{"too_large", `,"expectedRound":6`, http.StatusBadRequest},
		{"string", `,"expectedRound":"2"`, http.StatusBadRequest},
		{"fraction", `,"expectedRound":2.5`, http.StatusBadRequest},
		{"null", `,"expectedRound":null`, http.StatusBadRequest},
		{"stale_tab", `,"expectedRound":1`, http.StatusConflict},
		{"future_round", `,"expectedRound":3`, http.StatusConflict},
	} {
		t.Run(test.name, func(t *testing.T) {
			router, repository, ai := newTurnTestRouter(t, 2)
			response := sendTestTurn(router, test.roundField)
			if response.Code != test.wantStatus || ai.calls != 0 || repository.saves != 0 || repository.session.CurrentRound() != 2 {
				t.Fatalf("status=%d, AI calls=%d, saves=%d, round=%d", response.Code, ai.calls, repository.saves, repository.session.CurrentRound())
			}
		})
	}
}

func TestInterviewTurnReplayDoesNotAdvanceTwice(t *testing.T) {
	for _, round := range []int{1, 5} {
		t.Run(fmt.Sprint(round), func(t *testing.T) {
			router, repository, ai := newTurnTestRouter(t, round)
			field := fmt.Sprintf(`,"expectedRound":%d`, round)
			first := sendTestTurn(router, field)
			var body struct {
				Data struct {
					CurrentRound int     `json:"currentRound"`
					NextQuestion *string `json:"nextQuestion"`
				} `json:"data"`
			}
			if err := json.Unmarshal(first.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			wantRound := min(round+1, 5)
			if first.Code != http.StatusOK || body.Data.CurrentRound != wantRound || (body.Data.NextQuestion == nil) != (round == 5) {
				t.Fatalf("first response status=%d, data=%+v", first.Code, body.Data)
			}
			second := sendTestTurn(router, field)
			if second.Code != http.StatusConflict || ai.calls != 1 || repository.saves != 1 || repository.session.CurrentRound() != wantRound {
				t.Fatalf("replay status=%d, AI calls=%d, saves=%d, round=%d", second.Code, ai.calls, repository.saves, repository.session.CurrentRound())
			}
			if repository.session.Status() != interviewdomain.InterviewStatusInProgress {
				t.Fatal("answer submission must not complete the session before a report exists")
			}
		})
	}
}

func TestInterviewTurnFailureDoesNotAdvanceRound(t *testing.T) {
	for _, conflict := range []bool{false, true} {
		t.Run(fmt.Sprint(conflict), func(t *testing.T) {
			router, repository, ai := newTurnTestRouter(t, 2)
			wantStatus := http.StatusBadGateway
			if conflict {
				// 模拟模型调用期间，另一请求已抢先完成数据库事务。
				repository.saveError = port.ErrRepositoryConflict
				wantStatus = http.StatusConflict
			} else {
				ai.err = port.ErrAIUpstream
			}
			response := sendTestTurn(router, `,"expectedRound":2`)
			if response.Code != wantStatus || ai.calls != 1 || repository.saves != 0 || repository.session.CurrentRound() != 2 {
				t.Fatalf("status=%d, AI calls=%d, saves=%d, round=%d", response.Code, ai.calls, repository.saves, repository.session.CurrentRound())
			}
		})
	}
}
