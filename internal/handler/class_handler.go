package handler

import (
	"net/http"
	"strconv"

	"assignment-platform/internal/dto"
	"assignment-platform/internal/service"
	"assignment-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

// ClassHandler menangani endpoint class/batch.
type ClassHandler struct {
	classService *service.ClassService
}

func NewClassHandler(classService *service.ClassService) *ClassHandler {
	return &ClassHandler{classService: classService}
}

// CreateClass hanya untuk admin.
func (h *ClassHandler) CreateClass(c *gin.Context) {
	var req dto.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	id, err := h.classService.CreateClass(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create class", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Class created", gin.H{"id": id})
}

// ListClasses mengambil daftar class dengan pagination.
func (h *ClassHandler) ListClasses(c *gin.Context) {
	page := parseQueryInt(c, "page", 1)
	limit := parseQueryInt(c, "limit", 10)

	result, err := h.classService.ListClasses(c.Request.Context(), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to list classes", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Class list", result)
}

// GetClassByID mengambil detail class.
func (h *ClassHandler) GetClassByID(c *gin.Context) {
	classID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid class id", err.Error())
		return
	}

	class, err := h.classService.GetClassByID(c.Request.Context(), classID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Class not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Class detail", class)
}

// AssignTrainer menghubungkan trainer ke class.
func (h *ClassHandler) AssignTrainer(c *gin.Context) {
	classID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid class id", err.Error())
		return
	}

	var req dto.AssignUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := h.classService.AssignTrainer(c.Request.Context(), classID, req.UserID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to assign trainer", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Trainer assigned to class", nil)
}

// AssignTalent menghubungkan talent ke class.
func (h *ClassHandler) AssignTalent(c *gin.Context) {
	classID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid class id", err.Error())
		return
	}

	var req dto.AssignUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := h.classService.AssignTalent(c.Request.Context(), classID, req.UserID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to assign talent", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Talent assigned to class", nil)
}

// ListTrainersByClass menampilkan daftar trainer pada class tertentu.
func (h *ClassHandler) ListTrainersByClass(c *gin.Context) {
	classID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid class id", err.Error())
		return
	}

	trainers, err := h.classService.ListTrainersByClass(c.Request.Context(), classID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Failed to list class trainers", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Class trainers fetched", trainers)
}

// ListTalentsByClass menampilkan daftar talent pada class tertentu.
func (h *ClassHandler) ListTalentsByClass(c *gin.Context) {
	classID, err := parseIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid class id", err.Error())
		return
	}

	talents, err := h.classService.ListTalentsByClass(c.Request.Context(), classID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Failed to list class talents", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Class talents fetched", talents)
}

// parseIDParam mengubah parameter URL menjadi uint64.
func parseIDParam(c *gin.Context, name string) (uint64, error) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// parseQueryInt membaca query string integer dengan default value.
func parseQueryInt(c *gin.Context, name string, defaultValue int) int {
	value := c.Query(name)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
