package http_router

import (
	"net/http"

	ver1 "github.com/G0tem/go-service-layout/internal/apple/controller/http_router/v1"
	"github.com/G0tem/go-service-layout/internal/apple/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
)

func AppleRouter(r *chi.Mux, uc *usecase.UseCase) {
	r.Route("/api/v1", func(v1 chi.Router) {
		// Health check
		v1.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

		// // Swagger UI (доступен на /swagger/index.html)
		// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		// 	ginSwagger.URL("/swagger/doc.json"),     // URL для загрузки swagger.json
		// 	ginSwagger.DefaultModelsExpandDepth(-1), // скрыть модели по умолчанию
		// ))

		v1.Route("/apples", func(apples chi.Router) {
			apples.Use(otelchi.Middleware("Apple", otelchi.WithChiRoutes(apples)))
			// apples.Use(middleware.JWTAuth) // если нужно

			v1Handler := ver1.New(uc)
			apples.Post("/create_apple", v1Handler.CreateApple)        // POST /api/v1/apples/create_apple
			apples.Get("/get_apple/{id}", v1Handler.GetApple)          // GET /api/v1/apples/get_apple/{id}
			apples.Put("/update_apple/{id}", v1Handler.UpdateApple)    // PUT /api/v1/apples/update_apple/{id}
			apples.Delete("/delete_apple/{id}", v1Handler.DeleteApple) // DELETE /api/v1/apples/delete_apple/{id}
		})
	})
}
