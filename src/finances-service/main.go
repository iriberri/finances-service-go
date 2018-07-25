package main

import (
    "github.com/adeynack/finances-service-go/src/finances-service/controller"
    "github.com/adeynack/finances-service-go/src/finances-service/route"
)

func main() {
    tokenController := controller.NewTokens()

    routes := route.Register(tokenController)
    routes.Run(":3000")
}
