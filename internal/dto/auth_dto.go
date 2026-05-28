package dto

// LoginRequest adalah payload untuk endpoint login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse adalah response ketika login berhasil.
type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
