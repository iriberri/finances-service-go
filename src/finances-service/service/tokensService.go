package service

import (
    "encoding/base64"
    "sync"
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
    userService UserService
    lock        *sync.RWMutex
    tokenCache  map[string]string
}

var _ TokenService = &tokenService{}

func NewTokenService(userService UserService) TokenService {
    s := &tokenService{
        userService: userService,
        lock:        &sync.RWMutex{},
        tokenCache:  map[string]string{},
    }
    return s
}

func (s tokenService) CreateToken(username, password string) string {
    s.lock.Lock()
    defer s.lock.Unlock()

    // Check username and password
    if !s.userService.AuthenticateUser(username, password) {
        return ""
    }
    // Create token
    tokenPlain := username + ":" + password
    token := base64.StdEncoding.EncodeToString([]byte(tokenPlain))

    s.tokenCache[token] = username
    return token
}

func (s tokenService) ValidateToken(token string) bool {
    s.lock.RLock()
    defer s.lock.RUnlock()

    _, found := s.tokenCache[token]
    return found
}

func (s tokenService) InvalidateToken(token string) bool {
    s.lock.Lock()
    defer s.lock.Unlock()

    _, found := s.tokenCache[token]
    if !found {
        return false
    }
    delete(s.tokenCache, token)
    return true
}
