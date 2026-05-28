package handler

import (
	"net/http"

	"assignment-platform/internal/dto"
	"assignment-platform/internal/service"
	"assignment-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

// SubmissionHandler menangani submit, resubmit, list submission, dan progress.
type SubmissionHandler struct {
	submissionService *service.SubmissionService
}

func NewSubmissionHandler(submissionService *service.SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{submissionService: submissionService}
}

// SubmitAssignment dipakai talent untuk mengirim link tugas pertama kali.
func (h *SubmissionHandler) SubmitAssignment(c *gin.Context) {
	assignmentID, err := parseIDParam(c, "assignmentId")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid assignment id", err.Error())
		return
	}

	var req dto.SubmitAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	talentID := c.GetUint64("user_id")
	submissionID, err := h.submissionService.SubmitAssignment(c.Request.Context(), assignmentID, talentID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to submit assignment", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Assignment submitted", gin.H{"submission_id": submissionID})
}

// ResubmitSubmission dipakai talent untuk resubmit saat status revision_required.
func (h *SubmissionHandler) ResubmitSubmission(c *gin.Context) {
	submissionID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid submission id", err.Error())
		return
	}

	var req dto.SubmitAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	talentID := c.GetUint64("user_id")
	if err := h.submissionService.ResubmitSubmission(c.Request.Context(), submissionID, talentID, req); err != nil {
		utils.ErrorResponse(c, http.StatusConflict, "Failed to resubmit assignment", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Assignment resubmitted", nil)
}

// ListMySubmissions menampilkan submission milik talent yang login.
func (h *SubmissionHandler) ListMySubmissions(c *gin.Context) {
	talentID := c.GetUint64("user_id")

	items, err := h.submissionService.ListMySubmissions(c.Request.Context(), talentID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to list submissions", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "My submissions", items)
}

// ListSubmissionsByClass dipakai trainer/admin untuk melihat submission di class.
func (h *SubmissionHandler) ListSubmissionsByClass(c *gin.Context) {
	classID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid class id", err.Error())
		return
	}

	userID := c.GetUint64("user_id")
	role := c.GetString("role")

	items, err := h.submissionService.ListSubmissionsByClass(c.Request.Context(), classID, userID, role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to list class submissions", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Class submissions", items)
}

// GetClassProgress menghitung summary progress class.
func (h *SubmissionHandler) GetClassProgress(c *gin.Context) {
	classID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid class id", err.Error())
		return
	}

	userID := c.GetUint64("user_id")
	role := c.GetString("role")

	progress, err := h.submissionService.GetClassProgress(c.Request.Context(), classID, userID, role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to get class progress", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Class progress summary", progress)
}

// GetTalentProgress menghitung summary progress talent.
func (h *SubmissionHandler) GetTalentProgress(c *gin.Context) {
	talentID, err := parseIDParam(c, "talentId")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid talent id", err.Error())
		return
	}

	requesterID := c.GetUint64("user_id")
	role := c.GetString("role")

	progress, err := h.submissionService.GetTalentProgress(c.Request.Context(), talentID, requesterID, role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to get talent progress", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Talent progress summary", progress)
}
