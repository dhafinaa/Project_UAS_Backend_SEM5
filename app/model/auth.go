package model

// Login Request
type LoginRequest struct {
    Login    string `json:"login"`    // username atau email
    Password string `json:"password"` // password plaintext
}

// Login Response
type LoginResponse struct {
    Token     string   `json:"token"`
    Refresh   string   `json:"refresh"`
    User      User     `json:"user"`
    Permissions     []string `json:"permissions"`
}
