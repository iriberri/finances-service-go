package main

import (
    "github.com/adeynack/finances-service-go/src/finances-service/controller"
    "github.com/adeynack/finances-service-go/src/finances-service/route"
    "github.com/adeynack/finances-service-go/src/finances-service/service"
)

func main() {
    // Dependencies resolution
    userService := service.NewUserService()
    tokenService := service.NewTokenService(userService)
    tokenController := controller.NewTokensController(tokenService)

    // Create route and start listening to requests.
    routes := route.Register(tokenController)
    routes.Run(":3000")
}
