package server

import (
	"GraphQL/api"
	"GraphQL/internal/user/delivery/graphql"
	"GraphQL/internal/user/repository/postgres"
	"GraphQL/internal/user/usecase"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func (s *Server) MapHandlers() error {
	userPGRepo := postgres.NewClientPGRepository(s.cfg, s.pgDB)

	userUC := usecase.NewUserUC(
		s.cfg,
		userPGRepo,
	)
	resolver := graphql.NewResolver(userUC, s.cfg)
	gqlHandler := handler.NewDefaultServer(api.NewExecutableSchema(api.Config{Resolvers: resolver}))

	// Create a new ServeMux and set up the GraphQL endpoint
	mux := http.NewServeMux()
	mux.Handle("/query", gqlHandler)
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// Assign the mux to the server
	s.server = &http.Server{
		Addr:    s.cfg.Server.Port,
		Handler: mux,
	}

	return nil
}
