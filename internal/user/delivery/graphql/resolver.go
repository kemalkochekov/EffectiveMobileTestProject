package graphql

import (
	"GraphQL/configs"
	"GraphQL/internal/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userUC user.UseCase
	cfg    *configs.Config
}

func NewResolver(userUC user.UseCase, cfg *configs.Config) *Resolver {
	return &Resolver{
		userUC: userUC,
		cfg:    cfg,
	}
}
