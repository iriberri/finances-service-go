package controller

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/adeynack/finances-service-go/src/finances-service/api"
    "github.com/adeynack/finances-service-go/src/finances-service/problem"
    "github.com/adeynack/finances-service-go/src/finances-service/service"
    "github.com/gin-gonic/gin"
    "github.com/go-http-utils/headers"
)

type TokensController struct {
    tokenService service.TokenService
}

func NewTokensController(tokenService service.TokenService) *TokensController {
    return &TokensController{
        tokenService: tokenService,
    }
}

func (ctrl *TokensController) Create(c *gin.Context) {
    req := api.TokenCreateIn{}
    if err := c.BindJSON(&req);
        err != nil || req.Email == "" || req.Password == "" {
        example := api.TokenCreateIn{
            Email:    "name@domain.com",
            Password: "something_very_secure",
        }
        exampleJson, err := json.Marshal(example)
        if err != nil {
            panic(err)
        }
        WriteProblem(c, &problem.Problem{
            Status: http.StatusBadRequest,
            Title:  "Unexpected body structure",
            Detail: fmt.Sprintf(`Expecting a body similar to: %s`, exampleJson),
        })
        return
    }

    token := ctrl.tokenService.CreateToken(req.Email, req.Password)
    if token == "" {
        WriteProblem(c, &problem.Problem{
            Status: http.StatusUnauthorized,
            Title:  "Invalid credentials",
            Detail: "The specified credentials do not represent a known user or the password was invalid.",
        })
        return
    }

    c.JSON(http.StatusCreated, &api.TokenInfo{
        Token:         token,
        Status:        "Valid",
        Authenticated: true,
    })
}

func (ctrl *TokensController) Validate(c *gin.Context) {
    token, authenticated := ctrl.withToken(c, ctrl.tokenService.ValidateToken)
    if authenticated {
        c.JSON(http.StatusOK, &api.TokenInfo{
            Token:         token,
            Status:        "Valid",
            Authenticated: true,
        })
    }
}

func (ctrl *TokensController) Invalidate(c *gin.Context) {
    token, found := ctrl.withToken(c, ctrl.tokenService.InvalidateToken)
    if found {
        c.JSON(http.StatusOK, &api.TokenInfo{
            Token:         token,
            Status:        "Invalidated",
            Authenticated: false,
        })
    }
}

// Perform an operation "op" with the content of the "Authorization" header.
// "op" needs to return `true` if the token was valid in the first place.
func (ctrl *TokensController) withToken(c *gin.Context, op func(string) bool) (string, bool) {
    authHeader := c.GetHeader(headers.Authorization)
    if authHeader == "" {
        WriteProblem(c, problem.Unauthorized(fmt.Sprintf(`Header "%s" not provided.`, headers.Authorization)))
        return "", false
    }
    tokenWasValid := op(authHeader)
    if !tokenWasValid {
        WriteProblem(c, problem.Unauthorized("Invalid token."))
        return "", false
    }
    return authHeader, true
}

func (ctrl *TokensController) AuthorizeMiddleware(c *gin.Context) {
    _, valid := ctrl.withToken(c, ctrl.tokenService.ValidateToken)
    if !valid {
        c.Abort()
    }
}
