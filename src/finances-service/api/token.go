package api

type TokenInfo struct {
    Token         string `json:"token"`
    Status        string `json:"status"`
    Authenticated bool   `json:"authenticated"`
}

type TokenCreateIn struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
