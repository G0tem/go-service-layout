package http_router

import (
	ver1 "github.com/G0tem/go-service-layout/internal/product/controller/http_router/v1"
	"github.com/G0tem/go-service-layout/internal/product/usecase"
	"github.com/G0tem/go-service-layout/pkg/jwt"
	"github.com/G0tem/go-service-layout/pkg/ratelimit"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ProductRouter(r *gin.Engine, uc *usecase.UseCase, cfgRatelimit ratelimit.Config) {
	v1Group := r.Group("/api/v1")

	// Rate limiting middleware
	if cfgRatelimit.Enabled {
		v1Group.Use(ratelimit.RateLimiterMiddleware(cfgRatelimit))
		log.Info().
			Int("rps", cfgRatelimit.RequestsPerSecond).
			Int("burst", cfgRatelimit.Burst).
			Msg("Rate limiting enabled")
	}

	// Handlers
	v1Handler := ver1.New(uc)

	applesGroupJwt := v1Group.Group("/plant")

	// JWT Auth middleware
	applesGroupJwt.POST("/create_plant", v1Handler.CreatePlant) // POST /api/v1/plant/create_plant
	applesGroupJwt.GET("/get_plant/:id", v1Handler.GetPlant)    // GET /api/v1/plant/get_plant/{id}
	applesGroupJwt.Use(jwt.JWTAuth(uc.GetTokenManager()))

	// No JWT Auth middleware
	applesGroup := v1Group.Group("/apples")
	applesGroup.POST("/create_apple", v1Handler.CreateApple)       // POST /api/v1/apples/create_apple
	applesGroup.GET("/get_apple/:id", v1Handler.GetApple)          // GET /api/v1/apples/get_apple/{id}
	applesGroup.PUT("/update_apple/:id", v1Handler.UpdateApple)    // PUT /api/v1/apples/update_apple/{id}
	applesGroup.DELETE("/delete_apple/:id", v1Handler.DeleteApple) // DELETE /api/v1/apples/delete_apple/{id}

}
