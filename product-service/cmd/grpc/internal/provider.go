package internal

import (
	"context"
	"reflect"
	"unsafe"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/product-service/gen/gen/query"
)

// InternalProviderSet 内部依赖提供者集合
var InternalProviderSet = wire.NewSet(
	NewDatabase,
	ProvideGORMDB,
	ProvideQuery,
)

// NewDatabase 创建数据库连接
func NewDatabase(cfg *db.DatabaseConfig) *db.DB[*query.Query] {
	return db.NewDB(cfg, query.Use)
}

// ProvideQuery 从DB wrapper中提取*query.Query
func ProvideQuery(dbWrapper *db.DB[*query.Query]) *query.Query {
	return dbWrapper.Get(context.Background())
}

// ProvideGORMDB 从DB config重新创建gorm.DB
func ProvideGORMDB(dbWrapper *db.DB[*query.Query]) *gorm.DB {
	// 获取Query，然后通过反射获取内部的db字段
	q := dbWrapper.Get(context.Background())

	v := reflect.ValueOf(q).Elem()
	dbField := v.FieldByName("db")
	if !dbField.IsValid() {
		panic("cannot access db field from Query")
	}

	// 通过unsafe获取私有字段值
	ptr := unsafe.Pointer(dbField.UnsafeAddr())
	return *(**gorm.DB)(ptr)
}
