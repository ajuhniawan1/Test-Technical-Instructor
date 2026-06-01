package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RequireClassAccess(db *sql.DB, classIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, role, ok := getAuthUser(c)
		if !ok {
			abortUnauthorized(c)
			return
		}

		classID, err := strconv.Atoi(c.Param(classIDParam))
		if err != nil || classID <= 0 {
			abortBadRequest(c, "Class ID tidak valid")
			return
		}

		allowed, err := canAccessClass(c, db, userID, role, classID)
		if err != nil {
			abortServerError(c)
			return
		}

		if !allowed {
			abortForbidden(c, "Anda tidak memiliki akses ke class ini")
			return
		}

		c.Next()
	}
}

func RequireAssignmentAccess(db *sql.DB, assignmentIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, role, ok := getAuthUser(c)
		if !ok {
			abortUnauthorized(c)
			return
		}

		assignmentID, err := strconv.Atoi(c.Param(assignmentIDParam))
		if err != nil || assignmentID <= 0 {
			abortBadRequest(c, "Assignment ID tidak valid")
			return
		}

		classID, err := getClassIDByAssignmentID(c, db, assignmentID)
		if err == sql.ErrNoRows {
			abortNotFound(c, "Assignment tidak ditemukan")
			return
		}
		if err != nil {
			abortServerError(c)
			return
		}

		allowed, err := canAccessClass(c, db, userID, role, classID)
		if err != nil {
			abortServerError(c)
			return
		}

		if !allowed {
			abortForbidden(c, "Anda tidak memiliki akses ke assignment ini")
			return
		}

		c.Next()
	}
}

func RequireSubmissionAccess(db *sql.DB, submissionIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, role, ok := getAuthUser(c)
		if !ok {
			abortUnauthorized(c)
			return
		}

		submissionID, err := strconv.Atoi(c.Param(submissionIDParam))
		if err != nil || submissionID <= 0 {
			abortBadRequest(c, "Submission ID tidak valid")
			return
		}

		submissionOwnerID, classID, err := getSubmissionOwnerAndClassID(c, db, submissionID)
		if err == sql.ErrNoRows {
			abortNotFound(c, "Submission tidak ditemukan")
			return
		}
		if err != nil {
			abortServerError(c)
			return
		}

		if role == "admin" {
			c.Next()
			return
		}

		if role == "talent" {
			if userID != submissionOwnerID {
				abortForbidden(c, "Anda hanya boleh mengakses submission milik orang lain")
				return
			}

			c.Next()
			return
		}

		if role == "trainer" {
			allowed, err := isTrainerAssignedToClass(c, db, userID, classID)
			if err != nil {
				abortServerError(c)
				return
			}

			if !allowed {
				abortForbidden(c, "Trainer tidak memiliki akses ke submission class ini")
				return
			}

			c.Next()
			return
		}

		abortForbidden(c, "Role tidak memiliki akses")
	}
}

func RequireTalentProgressAccess(db *sql.DB, talentIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, role, ok := getAuthUser(c)
		if !ok {
			abortUnauthorized(c)
			return
		}

		talentID, err := strconv.Atoi(c.Param(talentIDParam))
		if err != nil || talentID <= 0 {
			abortBadRequest(c, "Talent ID tidak valid")
			return
		}

		if role == "admin" {
			c.Next()
			return
		}

		if role == "talent" {
			if userID != talentID {
				abortForbidden(c, "Talent hanya boleh melihat progress miliknya sendiri")
				return
			}

			c.Next()
			return
		}

		if role == "trainer" {
			allowed, err := isTrainerRelatedToTalent(c, db, userID, talentID)
			if err != nil {
				abortServerError(c)
				return
			}

			if !allowed {
				abortForbidden(c, "Trainer tidak memiliki akses ke progress talent ini")
				return
			}

			c.Next()
			return
		}

		abortForbidden(c, "Role tidak memiliki akses")
	}
}

func canAccessClass(c *gin.Context, db *sql.DB, userID int, role string, classID int) (bool, error) {
	if role == "admin" {
		return true, nil
	}

	if role == "trainer" {
		return isTrainerAssignedToClass(c, db, userID, classID)
	}

	if role == "talent" {
		return isTalentAssignedToClass(c, db, userID, classID)
	}

	return false, nil
}

func isTrainerAssignedToClass(c *gin.Context, db *sql.DB, trainerID int, classID int) (bool, error) {
	var exists bool

	err := db.QueryRowContext(
		c.Request.Context(),
		`
		SELECT EXISTS(
			SELECT 1
			FROM class_trainers
			WHERE class_id = ? AND trainer_id = ?
		)
		`,
		classID,
		trainerID,
	).Scan(&exists)

	return exists, err
}

func isTalentAssignedToClass(c *gin.Context, db *sql.DB, talentID int, classID int) (bool, error) {
	var exists bool

	err := db.QueryRowContext(
		c.Request.Context(),
		`
		SELECT EXISTS(
			SELECT 1
			FROM class_talents
			WHERE class_id = ? AND talent_id = ?
		)
		`,
		classID,
		talentID,
	).Scan(&exists)

	return exists, err
}

func isTrainerRelatedToTalent(c *gin.Context, db *sql.DB, trainerID int, talentID int) (bool, error) {
	var exists bool

	err := db.QueryRowContext(
		c.Request.Context(),
		`
		SELECT EXISTS(
			SELECT 1
			FROM class_trainers ct
			INNER JOIN class_talents clt ON clt.class_id = ct.class_id
			WHERE ct.trainer_id = ? AND clt.talent_id = ?
		)
		`,
		trainerID,
		talentID,
	).Scan(&exists)

	return exists, err
}

func getClassIDByAssignmentID(c *gin.Context, db *sql.DB, assignmentID int) (int, error) {
	var classID int

	err := db.QueryRowContext(
		c.Request.Context(),
		`
		SELECT class_id
		FROM assignments
		WHERE id = ?
		`,
		assignmentID,
	).Scan(&classID)

	return classID, err
}

func getSubmissionOwnerAndClassID(c *gin.Context, db *sql.DB, submissionID int) (int, int, error) {
	var talentID int
	var classID int

	err := db.QueryRowContext(
		c.Request.Context(),
		`
		SELECT s.talent_id, a.class_id
		FROM submissions s
		INNER JOIN assignments a ON a.id = s.assignment_id
		WHERE s.id = ?
		`,
		submissionID,
	).Scan(&talentID, &classID)

	return talentID, classID, err
}

func getAuthUser(c *gin.Context) (int, string, bool) {
	userID, ok := getIntFromContext(c, "user_id", "userID", "UserID", "id")
	if !ok {
		return 0, "", false
	}

	role, ok := getStringFromContext(c, "role", "Role", "user_role")
	if !ok || role == "" {
		return 0, "", false
	}

	return userID, role, true
}

func getIntFromContext(c *gin.Context, keys ...string) (int, bool) {
	for _, key := range keys {
		value, exists := c.Get(key)
		if !exists {
			continue
		}

		switch v := value.(type) {
		case int:
			return v, true
		case int64:
			return int(v), true
		case float64:
			return int(v), true
		case uint:
			return int(v), true
		case uint64:
			return int(v), true
		case string:
			parsed, err := strconv.Atoi(v)
			if err == nil {
				return parsed, true
			}
		default:
			parsed, err := strconv.Atoi(fmt.Sprint(v))
			if err == nil {
				return parsed, true
			}
		}
	}

	return 0, false
}

func getStringFromContext(c *gin.Context, keys ...string) (string, bool) {
	for _, key := range keys {
		value, exists := c.Get(key)
		if !exists {
			continue
		}

		return fmt.Sprint(value), true
	}

	return "", false
}

func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": "User belum terautentikasi",
	})
}

func abortForbidden(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"success": false,
		"message": message,
	})
}

func abortBadRequest(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": message,
	})
}

func abortNotFound(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": message,
	})
}

func abortServerError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": "Terjadi kesalahan pada server",
	})
}