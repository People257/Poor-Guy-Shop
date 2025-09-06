package user

import (
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "正常密码",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "复杂密码",
			password: "MySecureP@ssw0rd!",
			wantErr:  false,
		},
		{
			name:     "简单密码",
			password: "simple123",
			wantErr:  false,
		},
		{
			name:     "空密码",
			password: "",
			wantErr:  false, // HashPassword本身不验证密码格式
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 验证生成的哈希是否有效
				if len(hash) == 0 {
					t.Error("HashPassword() 返回空哈希")
				}

				// 验证哈希是否能正确验证原始密码
				err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password))
				if err != nil {
					t.Errorf("生成的哈希无法验证原始密码: %v", err)
				}
			}
		})
	}
}

func TestUser_VerifyPassword(t *testing.T) {
	// 预先生成的测试密码和哈希
	testPassword := "testPassword123"
	testHash, err := HashPassword(testPassword)
	if err != nil {
		t.Fatalf("无法生成测试哈希: %v", err)
	}

	tests := []struct {
		name         string
		user         *User
		password     string
		wantErr      bool
		errorMessage string
	}{
		{
			name: "正确密码",
			user: &User{
				PasswordHash: &testHash,
			},
			password: testPassword,
			wantErr:  false,
		},
		{
			name: "错误密码",
			user: &User{
				PasswordHash: &testHash,
			},
			password: "wrongPassword",
			wantErr:  true,
		},
		{
			name: "用户未设置密码",
			user: &User{
				PasswordHash: nil,
			},
			password:     "anyPassword",
			wantErr:      true,
			errorMessage: "用户未设置密码",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.VerifyPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.errorMessage != "" && err != nil {
				if err.Error() != tt.errorMessage {
					t.Errorf("User.VerifyPassword() error message = %v, want %v", err.Error(), tt.errorMessage)
				}
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		email    string
		password string
		phone    *string
		wantErr  bool
	}{
		{
			name:     "有效用户",
			username: "testuser",
			email:    "test@example.com",
			password: "password123",
			phone:    nil,
			wantErr:  false,
		},
		{
			name:     "有效用户带手机号",
			username: "testuser2",
			email:    "test2@example.com",
			password: "password123",
			phone:    stringPtr("13800138000"),
			wantErr:  false,
		},
		{
			name:     "无效用户名",
			username: "te",
			email:    "test@example.com",
			password: "password123",
			phone:    nil,
			wantErr:  true,
		},
		{
			name:     "无效邮箱",
			username: "testuser",
			email:    "invalid-email",
			password: "password123",
			phone:    nil,
			wantErr:  true,
		},
		{
			name:     "无效密码",
			username: "testuser",
			email:    "test@example.com",
			password: "weak",
			phone:    nil,
			wantErr:  true,
		},
		{
			name:     "无效手机号",
			username: "testuser",
			email:    "test@example.com",
			password: "password123",
			phone:    stringPtr("invalid-phone"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := CreateUser(tt.username, tt.email, tt.password, tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 验证创建的用户是否正确
				if user.Username != tt.username {
					t.Errorf("CreateUser() username = %v, want %v", user.Username, tt.username)
				}
				if user.Email == nil || *user.Email != tt.email {
					t.Errorf("CreateUser() email = %v, want %v", user.Email, tt.email)
				}
				if user.PasswordHash == nil {
					t.Error("CreateUser() 密码哈希为空")
				}
				if user.Status != UserStatusNormal {
					t.Errorf("CreateUser() status = %v, want %v", user.Status, UserStatusNormal)
				}

				// 验证密码是否正确加密
				err = user.VerifyPassword(tt.password)
				if err != nil {
					t.Errorf("CreateUser() 创建的用户无法验证密码: %v", err)
				}
			}
		})
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	// 创建一个测试用户
	user, err := CreateUser("testuser", "test@example.com", "oldPassword123", nil)
	if err != nil {
		t.Fatalf("无法创建测试用户: %v", err)
	}

	oldPasswordHash := *user.PasswordHash
	oldUpdateTime := user.UpdatedAt

	// 等待一小段时间确保更新时间不同
	time.Sleep(time.Millisecond * 10)

	// 更新密码
	newPassword := "newPassword456"
	err = user.UpdatePassword(newPassword)
	if err != nil {
		t.Errorf("UpdatePassword() error = %v", err)
	}

	// 验证密码哈希已更改
	if *user.PasswordHash == oldPasswordHash {
		t.Error("UpdatePassword() 密码哈希未更改")
	}

	// 验证新密码能够验证通过
	err = user.VerifyPassword(newPassword)
	if err != nil {
		t.Errorf("UpdatePassword() 新密码验证失败: %v", err)
	}

	// 验证旧密码不能验证通过
	err = user.VerifyPassword("oldPassword123")
	if err == nil {
		t.Error("UpdatePassword() 旧密码仍然能够验证通过")
	}

	// 验证更新时间已改变
	if !user.UpdatedAt.After(oldUpdateTime) {
		t.Error("UpdatePassword() 更新时间未改变")
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "有效密码",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "复杂有效密码",
			password: "MyComplexP@ssw0rd",
			wantErr:  false,
		},
		{
			name:     "密码太短",
			password: "pass1",
			wantErr:  true,
		},
		{
			name:     "密码太长",
			password: "a1" + string(make([]byte, 127)),
			wantErr:  true,
		},
		{
			name:     "只有字母",
			password: "onlyletters",
			wantErr:  true,
		},
		{
			name:     "只有数字",
			password: "12345678",
			wantErr:  true,
		},
		{
			name:     "最小有效密码",
			password: "abcdefg1",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// 辅助函数：创建字符串指针
func stringPtr(s string) *string {
	return &s
}

// 基准测试
func BenchmarkHashPassword(b *testing.B) {
	password := "testPassword123"
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "testPassword123"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatal(err)
	}

	user := &User{PasswordHash: &hash}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := user.VerifyPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}
