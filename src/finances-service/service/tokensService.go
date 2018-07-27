package service

import (
    "encoding/base64"
)

type TokenService interface {
    // Controls username and password and create a token for the user.
    // Returns a token when user is authenticated or an empty string if failed.
    CreateToken(username, password string) string

    // Validates if a given token is valid.
    ValidateToken(token string) bool

    // Invalidate given token.
    InvalidateToken(token string) bool
}

type tokenService struct {
    createChan     chan *createTokenRequest
    validateChan   chan *validateTokenRequest
    invalidateChan chan *invalidateTokenRequest
}

var _ TokenService = &tokenService{}

func NewTokenService() TokenService {
    s := &tokenService{
        createChan:     make(chan *createTokenRequest),
        validateChan:   make(chan *validateTokenRequest),
        invalidateChan: make(chan *invalidateTokenRequest),
    }
    go tokenServiceActor(s)
    return s
}

// This is of course a temporary situation, until this evolves to a
// centralised cache solution (database, Redis, ...).
var usersWithPassword = map[string]string{
    "max.mustermann": "maxisthebest",
    "laura.gärtner":  "mysafepassword",
}

func (s tokenService) CreateToken(username, password string) string {
    req := &createTokenRequest{
        Username:     username,
        Password:     password,
        ResponseChan: make(chan string),
    }
    s.createChan <- req
    token := <-req.ResponseChan
    return token
}

func (s tokenService) ValidateToken(token string) bool {
    req := &validateTokenRequest{
        Token:        token,
        ResponseChan: make(chan bool),
    }
    s.validateChan <- req
    validated := <-req.ResponseChan
    return validated
}

func (s tokenService) InvalidateToken(token string) bool {
    req := &invalidateTokenRequest{
        Token:        token,
        ResponseChan: make(chan bool),
    }
    s.invalidateChan <- req
    tokenWasValid := <-req.ResponseChan
    return tokenWasValid
}

//
// Internal Actor and its Messages
//

type tokenServiceInternal struct {
    tokenCache map[string]string
}

func tokenServiceActor(s *tokenService) {
    state := &tokenServiceInternal{
        tokenCache: map[string]string{},
    }
    for {
        select {
        case req := <-s.createChan:
            req.ResponseChan <- createToken(state, req)
        case req := <-s.validateChan:
            req.ResponseChan <- validateToken(state, req)
        case req := <-s.invalidateChan:
            req.ResponseChan <- invalidateToken(state, req)
        }
    }
}

func createToken(s *tokenServiceInternal, req *createTokenRequest) string {
    // Check username and password
    expectedPassword, found := usersWithPassword[req.Username]
    if !found || req.Password != expectedPassword {
        return ""
    }
    // Create token
    tokenPlain := req.Username + ":" + req.Password
    token := base64.StdEncoding.EncodeToString([]byte(tokenPlain))
    s.tokenCache[token] = req.Username
    return token
}

func validateToken(s *tokenServiceInternal, req *validateTokenRequest) bool {
    _, found := s.tokenCache[req.Token]
    if !found {
        return false
    }
    return true
}

func invalidateToken(s *tokenServiceInternal, req *invalidateTokenRequest) bool {
    _, found := s.tokenCache[req.Token]
    if !found {
        return false
    }
    delete(s.tokenCache, req.Token)
    return true
}

type createTokenRequest struct {
    Username     string
    Password     string
    ResponseChan chan string
}

type validateTokenRequest struct {
    Token        string
    ResponseChan chan bool
}

type invalidateTokenRequest struct {
    Token        string
    ResponseChan chan bool
}
