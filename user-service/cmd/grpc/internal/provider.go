package internal

import (
	"context"
	"reflect"
	"unsafe"

	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/user-service/gen/gen/query"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/google/wire"
)

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

// ProvideRedisClient 从UniversalClient转换为*redis.Client
func ProvideRedisClient(client redis.UniversalClient) *redis.Client {
	if c, ok := client.(*redis.Client); ok {
		return c
	}
	// 如果不是*redis.Client类型，创建一个新的（这种情况通常不会发生）
	panic("expected *redis.Client")
}

var InternalProviderSet = wire.NewSet(
	NewDB,
	db.NewRedis,
	ProvideQuery,
	ProvideGORMDB,
	ProvideRedisClient,
)
