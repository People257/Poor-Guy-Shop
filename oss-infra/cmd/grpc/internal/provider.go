package internal

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/oss-infra/gen/gen/query"
	"gorm.io/gorm"
)

// ExtractGormDB 从 db.DB 中提取 gorm.DB
func ExtractGormDB(dbWrapper *db.DB[*query.Query]) *gorm.DB {
	return dbWrapper.DB
}

var InternalProviderSet = wire.NewSet(
	NewDB,
	ExtractGormDB,
	NewRedisClient,
)
