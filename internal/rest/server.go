package rest

import (
	"database/sql"
	"mistar-be-go/internal/config"
	infrastructurehandler "mistar-be-go/internal/rest/handler/infrastructure"
	infrastructuretypehandler "mistar-be-go/internal/rest/handler/infrastructure_type"
	userhandler "mistar-be-go/internal/rest/handler/user"
	"mistar-be-go/internal/rest/middleware"
	storepgsql "mistar-be-go/internal/store/pgsql"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func New(
	cfg *config.Config,
	zlogger zerolog.Logger,
	sqlDB *sql.DB,
) http.Handler {
	r := chi.NewRouter()

	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
		middleware.HTTPTracer,
		middleware.RequestID(zlogger),
		middleware.HTTPLogger,
	)

	infrastructureTypeStore := storepgsql.NewInfrastructureType(sqlDB)
	infrastructureStore := storepgsql.NewInfrastructure(sqlDB)
	userStore := storepgsql.NewUser(sqlDB)

	infrastructureTypeHandler := infrastructuretypehandler.NewInfrastructureTypeHandler(sqlDB, infrastructureTypeStore, cfg.API)
	infrastructureHandler := infrastructurehandler.NewInfrastructureHandler(sqlDB, infrastructureStore, infrastructureTypeStore, cfg.API)
	userHandler := userhandler.NewUserHandler(sqlDB, userStore, cfg.API)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mistar Backend APIs"))
	})

	r.Route("/infrastructure", func(r chi.Router) {
		r.Get("/", infrastructureHandler.GetAllInfrastructureList)
		r.Post("/", infrastructureHandler.CreateInfrastructure)
		r.Get("/type", infrastructureTypeHandler.GetInfrastructureTypeList)
		r.Get("/filter", infrastructureTypeHandler.GetInfrastructureSubTypeFilter)
		r.Get("/filter/data", infrastructureTypeHandler.GetInfrastructureSubTypeFilterData)

		r.Get("/{infrastructure_type_id}/sub", infrastructureTypeHandler.GetInfrastructureSubTypeList)
	})

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", userHandler.RegisterUser)
		r.Post("/login", userHandler.LoginUser)
	})

	// STATIC FILE SERVE (FOR DEVELOPMENT PURPOSE ONLY)
	fs := http.FileServer(http.Dir("static/icons"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	return r
}
