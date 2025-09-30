package handlers

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"thousand.views_mine/internals/database/db_quaries"
)

type App struct {
	Quaries *db_quaries.Queries
	Ctx     context.Context
	Token   *jwtauth.JWTAuth
}
