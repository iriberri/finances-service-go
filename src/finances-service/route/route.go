package route

import (
    "github.com/gin-gonic/gin"
    "github.com/adeynack/finances-service-go/src/finances-service/controller"
    "net/http"
)

func Register(
    tokensController *controller.Tokens,
) *gin.Engine {

    r := gin.Default()

    tokens := r.Group("/tokens")
    tokens.POST("", tokensController.Create)
    tokens.GET("", tokensController.Validate)
    tokens.DELETE("", tokensController.Invalidate)

    books := r.Group("/books", tokensController.AuthorizeMiddleware)
    books.GET("", func(c *gin.Context) {
        c.String(http.StatusOK, "List of books (protected by token validation)")
    })

    return r
}
