package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"jwt_token"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	SessionId string `json:"session_id"`
}

type VerifyRequest struct {
	SessionId string `json:"session_id"`
	Code      string `json:"code"`
}

type VerifyResponse struct {
	Token string `json:"jwt_token"`
}
