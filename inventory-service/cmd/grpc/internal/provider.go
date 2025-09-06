package internal

import (
	"reflect"
	"unsafe"

	"gorm.io/gorm"

	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/inventory-service/gen/gen/query"
)

// NewDatabase 创建数据库连接
func NewDatabase(databaseConfig *db.DatabaseConfig) *db.DB[*query.Query] {
	database := db.NewDB(databaseConfig, query.Use)
	return database
}

// NewGormDB 从DB实例中提取*gorm.DB
func NewGormDB(database *db.DB[*query.Query]) *gorm.DB {
	// 使用反射获取私有字段
	dbValue := reflect.ValueOf(database).Elem()
	dbField := dbValue.FieldByName("db")

	// 使用unsafe访问私有字段
	return (*gorm.DB)(unsafe.Pointer(dbField.UnsafeAddr()))
}

// NewQuery 从DB实例中提取*query.Query
func NewQuery(database *db.DB[*query.Query]) *query.Query {
	// 使用反射获取私有字段
	dbValue := reflect.ValueOf(database).Elem()
	queryField := dbValue.FieldByName("query")

	// 使用unsafe访问私有字段
	return (*query.Query)(unsafe.Pointer(queryField.UnsafeAddr()))
}
