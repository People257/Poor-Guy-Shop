package optimisticlock

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gen"
)

var ErrOptimisticLock = status.Error(codes.Aborted, "data has been modified by other user")

type Versioned interface {
	// Version 返回当前版本号
	Version() Version
	// SetVersion 用来把实体的版本号设置为 v
	SetVersion(v int64)
}

// CheckOptimisticLockAndIncrementVersion 检查是否存在并发冲突，并增加版本号
func CheckOptimisticLockAndIncrementVersion(v Versioned, info gen.ResultInfo) error {
	if info.RowsAffected == 0 {
		return ErrOptimisticLock
	}
	v.SetVersion(v.Version().Int64 + 1)
	return nil
}
