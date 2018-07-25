package problem

import (
    "net/http"
)

func Unauthorized(detail string) *Problem {
    return &Problem{
        Status: http.StatusUnauthorized,
        Title:  unauthorizedTitle,
        Detail: detail,
    }
}

const unauthorizedTitle = "Unauthorized"
