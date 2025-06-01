package server

import (
	"github.com/Techeer-Hogwarts/slack-bot/cmd/auth"
	"github.com/Techeer-Hogwarts/slack-bot/cmd/handlers"
	"github.com/Techeer-Hogwarts/slack-bot/cmd/logger"
	"github.com/Techeer-Hogwarts/slack-bot/config"
	docs "github.com/Techeer-Hogwarts/slack-bot/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(handler *handlers.Handler) *gin.Engine {
	username := config.GetEnvVarAsString("SWAGGER_USERNAME", "admin")
	password := config.GetEnvVarAsString("SWAGGER_PASSWORD", "admin")
	authUsers := gin.Accounts{
		username: password,
	}

	router := gin.Default()
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Cookie", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Use(logger.JsonLoggerMiddleware())

	docs.SwaggerInfo.Title = "Techeerzip Slack Bot API"
	docs.SwaggerInfo.Description = "Slack Bot API for Slack and CICD"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	apiGroup := router.Group("/api/v1")
	{
		slackGroup := apiGroup.Group("/slack")
		slackGroup.Use(auth.VerifySlackRequest())
		{
			slackGroup.POST("/interactions", handler.SlackHandler.SlackInteractionHandler) // slack 봇 버튼 등 상호작용 처리
			slackGroup.POST("/commands", handler.SlackHandler.SlackCommandHandler)         // slack 봇 slash command 처리
		}

		profileGroup := apiGroup.Group("/profile") // legacy
		profileGroup.Use(auth.ValidateAPIKey())
		{
			profileGroup.POST("/picture", handler.ProfileHandler.ProfilePictureHandler) // 프로필 확인 및 사진 제공 legacy
		}

		alertGroup := apiGroup.Group("/alert")
		alertGroup.Use(auth.ValidateAPIKey())
		{
			alertGroup.POST("/messages", handler.AlertHandler.AlertMessageHandler)       // slack 봇으로 slack 메시지 전송
			alertGroup.POST("/find-member", handler.AlertHandler.AlertFindMemberHandler) // slack 봇으로 slack 메시지 전송
		}

		deployGroup := apiGroup.Group("/deploy")
		deployGroup.Use(auth.ValidateAPIKey())
		{
			deployGroup.POST("/image", handler.DeployHandler.DeployImageHandler)   // 배포 요청
			deployGroup.POST("/status", handler.DeployHandler.DeployStatusHandler) // 배포 요청
		}
	}
	swagger := router.Group("/swagger", gin.BasicAuth(authUsers))
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
