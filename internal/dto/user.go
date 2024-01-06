package dto

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
