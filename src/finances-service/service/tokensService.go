package service

import (
    "encoding/base64"

    "github.com/adeynack/finances-service-go/src/finances-service/util"
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
    syncChan   chan func()
    tokenCache map[string]string
}

var _ TokenService = &tokenService{}

func NewTokenService() TokenService {
    s := &tokenService{
        syncChan:   make(chan func()),
        tokenCache: map[string]string{},
    }
    util.StartSyncDispatcher(s.syncChan)
    return s
}

// This is of course a temporary situation, until this evolves to a
// centralised cache solution (database, Redis, ...).
var usersWithPassword = map[string]string{
    "max.mustermann": "maxisthebest",
    "laura.gärtner":  "mysafepassword",
}

func (s tokenService) CreateToken(username, password string) string {
    responseChan := make(chan string)
    defer close(responseChan)
    s.syncChan <- func() {
        // Check username and password
        expectedPassword, found := usersWithPassword[username]
        if !found || password != expectedPassword {
            responseChan <- ""
            return
        }
        // Create token
        tokenPlain := username + ":" + password
        token := base64.StdEncoding.EncodeToString([]byte(tokenPlain))
        s.tokenCache[token] = username
        responseChan <- token
    }
    return <-responseChan
}

func (s tokenService) ValidateToken(token string) bool {
    responseChan := make(chan bool)
    defer close(responseChan)
    s.syncChan <- func() {
        _, found := s.tokenCache[token]
        if !found {
            responseChan <- false
            return
        }
        responseChan <- true
    }
    return <-responseChan
}

func (s tokenService) InvalidateToken(token string) bool {
    responseChan := make(chan bool)
    defer close(responseChan)
    s.syncChan <- func() {
        _, found := s.tokenCache[token]
        if !found {
            responseChan <- false
            return
        }
        delete(s.tokenCache, token)
        responseChan <- true
    }
    return <-responseChan
}
