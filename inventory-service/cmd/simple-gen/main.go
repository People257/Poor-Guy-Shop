package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=47.99.147.94 user=postgres password=wdh961100 dbname=inventory-service port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}

	// 创建代码生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./gen/query",
		ModelPkgPath: "./gen/model",
		Mode:         gen.WithDefaultQuery,
	})
	g.UseDB(db)

	// 生成所有表
	g.GenerateAllTable()

	// 执行生成
	g.Execute()

	fmt.Println("Code generation completed!")
}
