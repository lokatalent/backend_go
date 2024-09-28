package response

type AuthResponse struct {
	TokensResponse `json:"tokens,omitempty"`
	UserResponse   `json:"user,omitempty"`
}
