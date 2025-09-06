# 用户服务测试数据

这个目录包含了用户服务的测试数据，包括用户信息和加密密码的测试用例。

## 文件说明

- `users.json` - JSON格式的测试用户数据
- `users.sql` - SQL插入语句，可直接用于数据库初始化
- `README.md` - 本说明文档

## 测试用户列表

| 用户名 | 邮箱 | 原始密码 | 手机号 | 状态 |
|--------|------|---------|--------|------|
| testuser1 | testuser1@example.com | password123 | 13800138001 | 正常 |
| testuser2 | testuser2@example.com | mySecurePass456 | - | 正常 |
| admin_user | admin@example.com | adminPassword789 | 13800138002 | 正常 |
| demo_user | demo@example.com | demoPassword321 | 13800138003 | 正常 |
| test_no_phone | nophone@example.com | noPhoneUser456 | - | 正常 |

## 密码加密说明

所有密码都使用bcrypt算法加密，成本因子为默认值（10）。每个用户的密码哈希都是唯一的，即使原始密码相同，生成的哈希也会不同。

### bcrypt密码哈希示例

```
原始密码: password123
哈希值: $2a$10$onYFQ6FxYN.vw5VfFT4l1.8xNSw2XPYGishaX525OvswJrCoQTUIO
```

哈希格式说明：
- `$2a$` - bcrypt算法版本
- `10$` - 成本因子（2^10 = 1024轮）
- 后续22字符 - 盐值
- 剩余31字符 - 实际哈希值

## 使用方法

### 1. 在单元测试中使用

```go
package user_test

import (
    "encoding/json"
    "os"
    "testing"
    
    "github.com/people257/poor-guy-shop/user-service/internal/domain/user"
)

type TestUserData struct {
    Username     string  `json:"username"`
    Email        string  `json:"email"`
    Password     string  `json:"password"`
    PasswordHash string  `json:"password_hash"`
    PhoneNumber  *string `json:"phone_number"`
    Status       int16   `json:"status"`
}

func loadTestUsers(t *testing.T) []TestUserData {
    data, err := os.ReadFile("../../testdata/users.json")
    if err != nil {
        t.Fatalf("无法读取测试数据: %v", err)
    }
    
    var result struct {
        Users []TestUserData `json:"users"`
    }
    
    err = json.Unmarshal(data, &result)
    if err != nil {
        t.Fatalf("无法解析测试数据: %v", err)
    }
    
    return result.Users
}

func TestPasswordVerification(t *testing.T) {
    testUsers := loadTestUsers(t)
    
    for _, userData := range testUsers {
        t.Run(userData.Username, func(t *testing.T) {
            user := &user.User{
                Username:     userData.Username,
                PasswordHash: &userData.PasswordHash,
            }
            
            // 验证正确密码
            err := user.VerifyPassword(userData.Password)
            if err != nil {
                t.Errorf("正确密码验证失败: %v", err)
            }
            
            // 验证错误密码
            err = user.VerifyPassword("wrongpassword")
            if err == nil {
                t.Error("错误密码验证应该失败")
            }
        })
    }
}
```

### 2. 在数据库初始化中使用

```bash
# 使用SQL文件初始化数据库
psql -d user_service -f testdata/users.sql
```

### 3. 在集成测试中使用

```go
// 在测试setup中插入测试数据
func setupTestData(db *gorm.DB) error {
    data, err := os.ReadFile("testdata/users.json")
    if err != nil {
        return err
    }
    
    var testData struct {
        Users []TestUserData `json:"users"`
    }
    
    err = json.Unmarshal(data, &testData)
    if err != nil {
        return err
    }
    
    for _, userData := range testData.Users {
        user := &model.User{
            Username:     userData.Username,
            Email:        &userData.Email,
            PasswordHash: &userData.PasswordHash,
            PhoneNumber:  userData.PhoneNumber,
            Status:       userData.Status,
        }
        
        if err := db.Create(user).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

## 生成新的测试数据

如果需要生成新的测试数据，可以运行测试数据生成器：

```bash
cd user-service/cmd/grpc
go run testdata_generator.go
```

这将生成新的密码哈希值和时间戳。

## 安全注意事项

1. **这些是测试数据** - 不要在生产环境中使用
2. **密码强度** - 测试密码符合系统的密码策略要求
3. **哈希安全** - 使用bcrypt确保密码安全存储
4. **数据隔离** - 确保测试数据不会污染生产数据

## 密码验证测试

所有测试用户的密码都经过了验证：
- ✅ 正确密码能够通过验证
- ✅ 错误密码会被正确拒绝
- ✅ 哈希格式符合bcrypt标准
- ✅ 每次生成的哈希都是唯一的
