package handlers

import (
	"context"

	"thousand.views_mine/internals/database/db_quaries"
)

type App struct {
	Quaries *db_quaries.Queries
	Ctx     context.Context
}
