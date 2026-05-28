package handler

import (
	"net/http"

	"assignment-platform/internal/dto"
	"assignment-platform/internal/service"
	"assignment-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

// ReviewHandler menangani endpoint review submission.
type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// ReviewSubmission dipakai trainer/admin untuk memberi score dan feedback.
func (h *ReviewHandler) ReviewSubmission(c *gin.Context) {
	submissionID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid submission id", err.Error())
		return
	}

	var req dto.ReviewSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	trainerID := c.GetUint64("user_id")
	role := c.GetString("role")

	if err := h.reviewService.ReviewSubmission(c.Request.Context(), submissionID, trainerID, role, req); err != nil {
		utils.ErrorResponse(c, http.StatusConflict, "Failed to review submission", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Submission reviewed", nil)
}

// RequestRevision adalah shortcut untuk status revision_required.
func (h *ReviewHandler) RequestRevision(c *gin.Context) {
	submissionID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid submission id", err.Error())
		return
	}

	var req dto.ReviewSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	req.Status = "revision_required"
	trainerID := c.GetUint64("user_id")
	role := c.GetString("role")

	if err := h.reviewService.ReviewSubmission(c.Request.Context(), submissionID, trainerID, role, req); err != nil {
		utils.ErrorResponse(c, http.StatusConflict, "Failed to request revision", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Revision requested", nil)
}
