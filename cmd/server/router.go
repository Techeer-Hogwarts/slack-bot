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

func setupRouter(handler *handlers.Hanlder) *gin.Engine {
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
		userGroup := apiGroup.Group("/auth")
		userGroup.Use(auth.ValidateJWT())
		{
			userGroup.POST("/login", handler.UserHandler.LoginHandler)  // ui 로그인
			userGroup.PATCH("/reset", handler.UserHandler.LoginHandler) // ui 비밀번호 초기화
			userGroup.POST("/logout", handler.UserHandler.LoginHandler) // ui 로그아웃
			userGroup.GET("/tokens", handler.UserHandler.LoginHandler)  // 토큰 발급
		}

		slackGroup := apiGroup.Group("/slack")
		slackGroup.Use(auth.VerifySlackRequest())
		{
			slackGroup.POST("/interactions", handler.UserHandler.LoginHandler) // slack 봇 버튼 등 상호작용 처리
			slackGroup.POST("/commands", handler.UserHandler.LoginHandler)     // slack 봇 slash command 처리
		}

		alertGroup := apiGroup.Group("/alert")
		alertGroup.Use(auth.ValidateAPIKey())
		{
			slackGroup.GET("/profiles", handler.UserHandler.LoginHandler) // 프로필 확인 및 사진 제공

			alertGroup.POST("/events", handler.UserHandler.LoginHandler)      // slack 봇으로 이벤트 알림 전송
			alertGroup.POST("/blogs", handler.UserHandler.LoginHandler)       // slack 봇으로 블로그 알림 전송
			alertGroup.GET("/channels", handler.UserHandler.LoginHandler)     // 채널 알림
			alertGroup.GET("/users", handler.UserHandler.LoginHandler)        // 유저 알림
			alertGroup.GET("/dev/channels", handler.UserHandler.LoginHandler) // 개발용 채널 알림
			alertGroup.GET("/dev/users", handler.UserHandler.LoginHandler)    // 개발용 유저 알림
		}

		workflowGroup := apiGroup.Group("/workflows")
		workflowGroup.Use(auth.ValidateAPIKey())
		workflowGroup.Use(auth.ValidateJWT())
		{
			workflowGroup.POST("", handler.UserHandler.LoginHandler)               // 워크플로우 생성
			workflowGroup.GET("/:id", handler.UserHandler.LoginHandler)            // 워크플로우 조회
			workflowGroup.PATCH("/:id", handler.UserHandler.LoginHandler)          // 워크플로우 수정
			workflowGroup.GET("/events", handler.UserHandler.LoginHandler)         // 워크플로우 이벤트들 확인
			workflowGroup.POST("/events", handler.UserHandler.LoginHandler)        // 워크플로우 이벤트들 생성 (깃헙 액션에서 보내면 됨)
			workflowGroup.GET("/events/:id", handler.UserHandler.LoginHandler)     // 워크플로우 이벤트들 확인
			workflowGroup.POST("/events/deploy", handler.UserHandler.LoginHandler) // 워크플로우 배포
		}

		// githubGroup := apiGroup.Group("/github")
		// {
		// 	githubGroup.POST("/events", handler.UserHandler.LoginHandler)
		// }
	}
	swagger := router.Group("/swagger", gin.BasicAuth(authUsers))
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
