package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/middleware"
	dashboardapp "job-copilot-backend/internal/application/dashboard"
	"job-copilot-backend/internal/port"
	"job-copilot-backend/pkg/response"
)

type DashboardHandler struct{ service *dashboardapp.Service }

func NewDashboardHandler(service *dashboardapp.Service) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (handler *DashboardHandler) Get(ctx *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return
	}
	summary, err := handler.service.Get(ctx.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, dashboardapp.ErrInvalidQuery):
			response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "用户身份无效")
		case errors.Is(err, port.ErrRepositoryUnavailable):
			response.Error(ctx, http.StatusServiceUnavailable, "DATABASE_NOT_CONFIGURED", "数据存储服务未配置")
		default:
			response.Error(ctx, http.StatusBadGateway, "DATABASE_ERROR", "读取仪表盘失败，请稍后重试")
		}
		return
	}
	recentRecords := make([]gin.H, 0, len(summary.RecentRecords))
	for _, record := range summary.RecentRecords {
		recentRecords = append(recentRecords, gin.H{
			"id": record.ID, "kind": record.Kind, "companyName": record.CompanyName,
			"jobTitle": record.JobTitle, "description": record.Description, "createdAt": record.CreatedAt,
		})
	}
	response.Success(ctx, gin.H{
		"jdCount": summary.JDCount, "interviewCount": summary.InterviewCount,
		"applicationCount":           summary.ApplicationCount,
		"activeApplicationCount":     summary.ActiveApplicationCount,
		"realInterviewCount":         summary.RealInterviewCount,
		"offerCount":                 summary.OfferCount,
		"applicationToInterviewRate": summary.ApplicationToInterviewRate,
		"interviewToOfferRate":       summary.InterviewToOfferRate,
		"weeklyRecordCount":          summary.WeeklyRecordCount,
		"profileCompleteness":        summary.ProfileCompleteness,
		"recentRecords":              recentRecords,
	})
}
