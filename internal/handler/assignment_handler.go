package handler

import (
	"net/http"

	"assignment-platform/internal/dto"
	"assignment-platform/internal/service"
	"assignment-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

// AssignmentHandler menangani endpoint assignment.
type AssignmentHandler struct {
	assignmentService *service.AssignmentService
}

func NewAssignmentHandler(assignmentService *service.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{assignmentService: assignmentService}
}

// CreateAssignment membuat tugas baru untuk suatu class.
func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {
	classID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid class id", err.Error())
		return
	}

	var req dto.CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	userID := c.GetUint64("user_id")
	role := c.GetString("role")

	id, err := h.assignmentService.CreateAssignment(c.Request.Context(), classID, userID, role, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to create assignment", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Assignment created", gin.H{"id": id})
}

// ListAssignmentsByClass mengambil assignment berdasarkan class.
func (h *AssignmentHandler) ListAssignmentsByClass(c *gin.Context) {
	classID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid class id", err.Error())
		return
	}

	userID := c.GetUint64("user_id")
	role := c.GetString("role")

	assignments, err := h.assignmentService.ListByClass(c.Request.Context(), classID, userID, role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to list assignments", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Assignment list", assignments)
}

// GetAssignmentByID mengambil detail assignment.
func (h *AssignmentHandler) GetAssignmentByID(c *gin.Context) {
	assignmentID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid assignment id", err.Error())
		return
	}

	userID := c.GetUint64("user_id")
	role := c.GetString("role")

	assignment, err := h.assignmentService.GetAssignmentByID(c.Request.Context(), assignmentID, userID, role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to get assignment", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Assignment detail", assignment)
}

// UpdateAssignment mengubah assignment.
func (h *AssignmentHandler) UpdateAssignment(c *gin.Context) {
	assignmentID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid assignment id", err.Error())
		return
	}

	var req dto.UpdateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	userID := c.GetUint64("user_id")
	role := c.GetString("role")

	if err := h.assignmentService.UpdateAssignment(c.Request.Context(), assignmentID, userID, role, req); err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to update assignment", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Assignment updated", nil)
}

// CloseAssignment menutup assignment.
func (h *AssignmentHandler) CloseAssignment(c *gin.Context) {
	assignmentID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid assignment id", err.Error())
		return
	}

	userID := c.GetUint64("user_id")
	role := c.GetString("role")

	if err := h.assignmentService.CloseAssignment(c.Request.Context(), assignmentID, userID, role); err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to close assignment", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Assignment closed", nil)
}
