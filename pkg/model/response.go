package model

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type RegisterResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	User    struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"user"`
}

type LoginResponse struct {
	Status      int    `json:"status"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expiresIn"`
	User        struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"user"`
}
