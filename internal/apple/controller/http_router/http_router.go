package http_router

import (
	"net/http"

	ver1 "github.com/G0tem/go-service-layout/internal/apple/controller/http_router/v1"
	"github.com/G0tem/go-service-layout/internal/apple/usecase"
	"github.com/G0tem/go-service-layout/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AppleRouter(r *gin.Engine, uc *usecase.UseCase) {
	v1Group := r.Group("/api/v1")
	{
		// Health check
		v1Group.GET("/healthz", func(c *gin.Context) {
			log.Debug().Msg("healthz")
			c.String(http.StatusOK, "ok")
		})

		// Swagger UI (доступен на /swagger/index.html)
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL("/swagger/doc.json"),     // URL для загрузки swagger.json
			ginSwagger.DefaultModelsExpandDepth(-1), // скрыть модели по умолчанию
		))

		applesGroup := v1Group.Group("/apples")
		{
			// TODO: Add OpenTelemetry middleware for gin
			// applesGroup.Use(otelgin.Middleware("Apple"))

			// JWT Auth middleware
			applesGroup.Use(jwt.JWTAuth(uc.GetTokenManager()))

			v1Handler := ver1.New(uc)
			applesGroup.POST("/create_apple", v1Handler.CreateApple)       // POST /api/v1/apples/create_apple
			applesGroup.GET("/get_apple/:id", v1Handler.GetApple)          // GET /api/v1/apples/get_apple/{id}
			applesGroup.PUT("/update_apple/:id", v1Handler.UpdateApple)    // PUT /api/v1/apples/update_apple/{id}
			applesGroup.DELETE("/delete_apple/:id", v1Handler.DeleteApple) // DELETE /api/v1/apples/delete_apple/{id}
		}
	}
}
