package controller

import (
    "github.com/adeynack/finances-service-go/src/finances-service/problem"
    "github.com/gin-gonic/gin"
    "github.com/go-http-utils/headers"
)

func WriteProblem(c *gin.Context, problem *problem.Problem) {
    c.Header(headers.ContentType, "application/problem+json")
    c.JSON(problem.Status, problem)
}
