package main

import (
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/adeynack/finances-service-go/src/finances-service/controller"
    "github.com/adeynack/finances-service-go/src/finances-service/route"
    "github.com/adeynack/finances-service-go/src/finances-service/service"
    "github.com/adeynack/finances-service-go/src/finances-service/util"
    "github.com/gin-gonic/gin"
    "github.com/olebedev/config"
    log "github.com/sirupsen/logrus"
    "github.com/toorop/gin-logrus"
)

func main() {
    // Configuration
    conf := readConfiguration()
    setupLogging(conf)
    printBanner(conf)

    // Dependencies resolution
    databaseService := service.NewDatabaseService(conf)
    userService := service.NewUserService(databaseService)
    tokenService := service.NewTokenService(userService)
    tokensController := controller.NewTokensController(tokenService)

    // Create route and start listening to requests.
    startHttpService(conf, tokensController)
}

func startHttpService(
    conf *util.ConfigReader,
    tokensController *controller.TokensController,
) {
    gin.SetMode(conf.UString("gin.mode", "debug"))
    engine := gin.New()
    if conf.UBool("gin.log_requests", false) {
        engine.Use(ginlogrus.Logger(log.StandardLogger()))
    }
    engine.Use(gin.Recovery())
    route.Register(engine, tokensController)
    engine.Run(":3000")
}

func readConfiguration() *util.ConfigReader {
    envVarConfigFile := "FINANCES_SERVICE_CONFIG"
    envConfigFile, found := os.LookupEnv(envVarConfigFile)
    if found && len(envConfigFile) > 0 {
        log.Infof("Environment variable '%s' set. Using configuration file '%s'.", envVarConfigFile, envConfigFile)
    } else {
        defaultConfigFile := "config/dev.yaml"
        log.Infof("Environment variable '%s' not set. Using default configuration file '%s'.", envVarConfigFile, defaultConfigFile)
        envConfigFile = defaultConfigFile
    }
    conf, err := config.ParseYamlFile(envConfigFile)
    if err != nil {
        panic(err)
    }
    conf.EnvPrefix("FINANCES_SERVICE").Flag()
    return &util.ConfigReader{Config: conf}
}

func setupLogging(conf *util.ConfigReader) {
    log.SetOutput(os.Stdout)
    formatter := conf.UString("log.formatter", "text")
    switch formatter {
    case "text":
        log.SetFormatter(&log.TextFormatter{
            TimestampFormat: time.RFC3339,
            FullTimestamp:   true,
        })
    case "json":
        log.SetFormatter(&log.JSONFormatter{
            TimestampFormat: time.RFC3339,
        })
    default:
        panic(fmt.Errorf("unsupported formatter: %s", formatter))
    }
}

func printBanner(conf *util.ConfigReader) {
    banner := conf.MustString("banner")
    bannerLines := strings.Split(banner, "\n")
    for _, line := range bannerLines {
        log.Info(line)
    }
}
