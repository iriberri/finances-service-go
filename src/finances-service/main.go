package main

import (
    "log"
    "os"

    "github.com/adeynack/finances-service-go/src/finances-service/controller"
    "github.com/adeynack/finances-service-go/src/finances-service/route"
    "github.com/adeynack/finances-service-go/src/finances-service/service"
    "github.com/gin-gonic/gin"
    "github.com/olebedev/config"
)

func main() {
    // Configuration
    conf := readConfiguration()

    // Dependencies resolution
    databaseService := service.NewDatabaseService(conf)
    userService := service.NewUserService(databaseService)
    tokenService := service.NewTokenService(userService)
    tokenController := controller.NewTokensController(tokenService)

    // Create route and start listening to requests.
    gin.SetMode(conf.UString("gin.mode"))
    routes := route.Register(tokenController)
    routes.Run(":3000")
}

func readConfiguration() *config.Config {
    envVarConfigFile := "FINANCES_SERVICE_CONFIG"
    envConfigFile, found := os.LookupEnv(envVarConfigFile)
    if found && len(envConfigFile) > 0 {
        log.Printf("Environment variable '%s' set. Using configuration file '%s'.", envVarConfigFile, envConfigFile)
    } else {
        defaultConfigFile := "config/dev.yaml"
        log.Printf("Environment variable '%s' not set. Using default configuration file '%s'.", envVarConfigFile, defaultConfigFile)
        envConfigFile = defaultConfigFile
    }
    conf, err := config.ParseYamlFile(envConfigFile)
    if err != nil {
        panic(err)
    }
    conf.EnvPrefix("FINANCES_SERVICE").Flag()
    return conf
}
