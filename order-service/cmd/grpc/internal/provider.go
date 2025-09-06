package internal

import (
	"reflect"
	"unsafe"

	"gorm.io/gorm"

	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/order-service/gen/gen/query"
)

// NewDatabase 创建数据库实例
func NewDatabase(config *db.DatabaseConfig) *db.DB[*query.Query] {
	return db.NewDB(config, query.Use)
}

// NewGormDB 创建GORM数据库实例
func NewGormDB(database *db.DB[*query.Query]) *gorm.DB {
	// 使用反射获取私有字段
	v := reflect.ValueOf(database).Elem()
	dbField := v.FieldByName("db")
	if !dbField.IsValid() {
		panic("无法获取数据库连接")
	}

	// 使用unsafe获取私有字段的值
	return (*gorm.DB)(unsafe.Pointer(dbField.UnsafeAddr()))
}

// NewQuery 创建查询实例
func NewQuery(database *db.DB[*query.Query]) *query.Query {
	// 使用反射获取私有字段
	v := reflect.ValueOf(database).Elem()
	queryField := v.FieldByName("query")
	if !queryField.IsValid() {
		panic("无法获取查询实例")
	}

	// 使用unsafe获取私有字段的值
	return (*query.Query)(unsafe.Pointer(queryField.UnsafeAddr()))
}
