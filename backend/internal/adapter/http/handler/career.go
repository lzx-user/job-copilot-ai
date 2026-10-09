package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/middleware"
	careerapp "job-copilot-backend/internal/application/career"
	careerdomain "job-copilot-backend/internal/domain/career"
	"job-copilot-backend/internal/port"
	"job-copilot-backend/pkg/response"
)

const maxCareerRequestBytes = 32 << 10

type CareerHandler struct{ service *careerapp.Service }

func NewCareerHandler(service *careerapp.Service) *CareerHandler {
	return &CareerHandler{service: service}
}

func (handler *CareerHandler) GetWorkspace(ctx *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return
	}
	workspace, err := handler.service.GetWorkspace(ctx.Request.Context(), userID)
	if err != nil {
		handleCareerLoadError(ctx, err)
		return
	}
	response.Success(ctx, workspace)
}

func handleCareerLoadError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, port.ErrRepositorySchema):
		response.Error(ctx, http.StatusServiceUnavailable, "DATABASE_MIGRATION_REQUIRED", "求职数据表尚未初始化，请先执行 202610060001_create_career_pipeline.sql")
	case errors.Is(err, port.ErrRepositoryUnavailable):
		response.Error(ctx, http.StatusServiceUnavailable, "DATABASE_NOT_CONFIGURED", "数据存储服务未配置")
	case errors.Is(err, port.ErrRepositorySchema):
		response.Error(ctx, http.StatusServiceUnavailable, "DATABASE_MIGRATION_REQUIRED", "求职数据表尚未初始化，请先执行 202610060001_create_career_pipeline.sql")
	default:
		response.Error(ctx, http.StatusBadGateway, "DATABASE_ERROR", "求职工作台加载失败，请稍后重试")
	}
}

func (handler *CareerHandler) CreateResume(ctx *gin.Context) {
	var input careerdomain.CreateResumeInput
	userID, ok := handler.bind(ctx, &input)
	if !ok {
		return
	}
	resumeID, versionID, err := handler.service.CreateResume(ctx.Request.Context(), userID, input)
	if err != nil {
		handleCareerError(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"resumeId": resumeID, "versionId": versionID})
}

func (handler *CareerHandler) CreateResumeVersion(ctx *gin.Context) {
	var input careerdomain.CreateResumeVersionInput
	userID, ok := handler.bind(ctx, &input)
	if !ok {
		return
	}
	input.ResumeID = ctx.Param("id")
	id, err := handler.service.CreateResumeVersion(ctx.Request.Context(), userID, input)
	if err != nil {
		handleCareerError(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"versionId": id})
}

func (handler *CareerHandler) CreateJob(ctx *gin.Context) {
	var input careerdomain.CreateJobInput
	userID, ok := handler.bind(ctx, &input)
	if !ok {
		return
	}
	jobID, versionID, err := handler.service.CreateJob(ctx.Request.Context(), userID, input)
	if err != nil {
		handleCareerError(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"jobId": jobID, "jdVersionId": versionID})
}

func (handler *CareerHandler) CreateJDVersion(ctx *gin.Context) {
	var input careerdomain.CreateJDVersionInput
	userID, ok := handler.bind(ctx, &input)
	if !ok {
		return
	}
	input.JobID = ctx.Param("id")
	id, err := handler.service.CreateJDVersion(ctx.Request.Context(), userID, input)
	if err != nil {
		handleCareerError(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"jdVersionId": id})
}

func (handler *CareerHandler) CreateApplication(ctx *gin.Context) {
	var input careerdomain.CreateApplicationInput
	userID, ok := handler.bind(ctx, &input)
	if !ok {
		return
	}
	id, err := handler.service.CreateApplication(ctx.Request.Context(), userID, input)
	if err != nil {
		handleCareerError(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"applicationId": id})
}

func (handler *CareerHandler) AddEvent(ctx *gin.Context) {
	var input careerdomain.CreateEventInput
	userID, ok := handler.bind(ctx, &input)
	if !ok {
		return
	}
	id, err := handler.service.AddEvent(ctx.Request.Context(), userID, input)
	if err != nil {
		handleCareerError(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"eventId": id})
}

func (handler *CareerHandler) CreateInterview(ctx *gin.Context) {
	var input careerdomain.CreateInterviewInput
	userID, ok := handler.bind(ctx, &input)
	if !ok {
		return
	}
	id, err := handler.service.CreateInterview(ctx.Request.Context(), userID, input)
	if err != nil {
		handleCareerError(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"interviewId": id})
}

func (handler *CareerHandler) UpdateInterview(ctx *gin.Context) {
	var input careerdomain.UpdateInterviewInput
	userID, ok := handler.bind(ctx, &input)
	if !ok {
		return
	}
	input.InterviewID = ctx.Param("id")
	if err := handler.service.UpdateInterview(ctx.Request.Context(), userID, input); err != nil {
		handleCareerError(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"saved": true})
}

func (handler *CareerHandler) SaveRetrospective(ctx *gin.Context) {
	var input careerdomain.SaveRetrospectiveInput
	userID, ok := handler.bind(ctx, &input)
	if !ok {
		return
	}
	if err := handler.service.SaveRetrospective(ctx.Request.Context(), userID, input); err != nil {
		handleCareerError(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"saved": true})
}

func (handler *CareerHandler) SaveOffer(ctx *gin.Context) {
	var input careerdomain.SaveOfferInput
	userID, ok := handler.bind(ctx, &input)
	if !ok {
		return
	}
	if err := handler.service.SaveOffer(ctx.Request.Context(), userID, input); err != nil {
		handleCareerError(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"saved": true})
}

func (handler *CareerHandler) bind(ctx *gin.Context, target any) (string, bool) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return "", false
	}
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxCareerRequestBytes)
	if err := ctx.ShouldBindJSON(target); err != nil {
		response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "请求参数格式不正确")
		return "", false
	}
	return userID, true
}

func handleCareerError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, careerapp.ErrInvalidInput):
		response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "请检查必填项、长度和状态值")
	case errors.Is(err, port.ErrRepositoryUnavailable):
		response.Error(ctx, http.StatusServiceUnavailable, "DATABASE_NOT_CONFIGURED", "数据存储服务未配置")
	case errors.Is(err, port.ErrRepositoryNotFound):
		response.Error(ctx, http.StatusNotFound, "NOT_FOUND", "关联记录不存在")
	case errors.Is(err, port.ErrRepositoryConflict):
		response.Error(ctx, http.StatusConflict, "CONFLICT", "记录状态已变化，请刷新后重试")
	default:
		response.Error(ctx, http.StatusBadGateway, "DATABASE_ERROR", "求职记录保存失败，请稍后重试")
	}
}
