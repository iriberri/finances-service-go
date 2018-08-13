package route

import (
    "net/http"

    "github.com/adeynack/finances-service-go/src/finances-service/controller"
    "github.com/gin-gonic/gin"
)

func Register(
    engine *gin.Engine,
    tokensController *controller.TokensController,
) {
    tokens := engine.Group("/tokens")
    tokens.POST("", tokensController.Create)
    tokens.GET("", tokensController.Validate)
    tokens.DELETE("", tokensController.Invalidate)

    books := engine.Group("/books", tokensController.AuthorizeMiddleware)
    books.GET("", func(c *gin.Context) {
        c.String(http.StatusOK, "List of books (protected by token validation)")
    })
}
