package internal

import (
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/user-service/gen/gen/query"
)

func NewDB(cfg *db.DatabaseConfig) *db.DB[*query.Query] {
	return db.NewDB(cfg, query.Use)
}
