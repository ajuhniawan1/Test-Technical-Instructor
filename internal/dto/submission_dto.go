package dto

// SubmitAssignmentRequest adalah payload saat talent submit tugas.
// Minimal github_url atau deployment_url harus diisi.
type SubmitAssignmentRequest struct {
	GithubURL     string `json:"github_url"`
	DeploymentURL string `json:"deployment_url"`
	Notes         string `json:"notes"`
}

// ReviewSubmissionRequest adalah payload saat trainer memberi review.
type ReviewSubmissionRequest struct {
	Score    int    `json:"score" binding:"required"`
	Feedback string `json:"feedback" binding:"required"`
	Status   string `json:"status" binding:"required"` // reviewed atau revision_required
}
