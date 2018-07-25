package controller

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "github.com/go-http-utils/headers"
    "fmt"
    "github.com/adeynack/finances-service-go/src/finances-service/problem"
    "encoding/base64"
)

type Tokens struct {
    tokenCache map[string]string
}

func NewTokens() *Tokens {
    return &Tokens{
        tokenCache: map[string]string{},
    }
}

func (ctrl *Tokens) Create(c *gin.Context) {
    req := TokenCreateIn{}
    if err := c.BindJSON(&req);
        err != nil || req.Username == "" || req.Password == "" {
        WriteProblem(c, &problem.Problem{
            Status: http.StatusBadRequest,
            Title:  "Unexpected body structure",
            Detail: `Expecting a body similar to: {"username":"foo","password":"bar"}`,
        })
        return
    }

    tokenPlain := req.Username + ":" + req.Password
    token := base64.StdEncoding.EncodeToString([]byte(tokenPlain))
    ctrl.tokenCache[token] = req.Username

    c.JSON(http.StatusCreated, &TokenInfo{
        Token:  token,
        Status: "Valid",
    })
}

func (ctrl *Tokens) Validate(c *gin.Context) {
    if token, found := ctrl.authorize(c); found {
        c.JSON(http.StatusOK, &TokenInfo{
            Token:  token,
            Status: "Valid",
        })
    }
}

func (ctrl *Tokens) Invalidate(c *gin.Context) {
    if token, found := ctrl.authorize(c); found {
        delete(ctrl.tokenCache, token)
        c.JSON(http.StatusOK, &TokenInfo{
            Token:  token,
            Status: "Invalidated",
        })
    }
}

func (ctrl *Tokens) authorize(c *gin.Context) (string, bool) {
    authHeader := c.GetHeader(headers.Authorization)
    if authHeader == "" {
        WriteProblem(c, problem.Unauthorized(fmt.Sprintf(`Header "%s" not provided.`, headers.Authorization)))
        return "", false
    }
    _, found := ctrl.tokenCache[authHeader]
    if !found {
        WriteProblem(c, problem.Unauthorized("Invalid token."))
        return "", false
    }
    return authHeader, true
}

func (ctrl *Tokens) AuthorizeMiddleware(c *gin.Context) {
    if _, found := ctrl.authorize(c); !found {
        c.Abort()
    }
}

type TokenInfo struct {
    Token  string `json:"token"`
    Status string `json:"status"`
}

type TokenCreateIn struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
