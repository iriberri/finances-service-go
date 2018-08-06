package api

type TokenInfo struct {
    Token         string `json:"token"`
    Status        string `json:"status"`
    Authenticated bool   `json:"authenticated"`
}

type TokenCreateIn struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
