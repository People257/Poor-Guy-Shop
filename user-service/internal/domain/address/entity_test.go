package address

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateAddress(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		receiverName  string
		receiverPhone string
		province      string
		city          string
		district      string
		street        string
		opts          *AddressOptions
		wantErr       bool
	}{
		{
			name:          "valid address",
			userID:        "user-123",
			receiverName:  "张三",
			receiverPhone: "13800138000",
			province:      "广东省",
			city:          "深圳市",
			district:      "南山区",
			street:        "科技园南区",
			opts:          nil,
			wantErr:       false,
		},
		{
			name:          "valid address with options",
			userID:        "user-123",
			receiverName:  "李四",
			receiverPhone: "13900139000",
			province:      "北京市",
			city:          "北京市",
			district:      "海淀区",
			street:        "中关村大街1号",
			opts: &AddressOptions{
				PostalCode:   stringPtr("100000"),
				AddressLabel: LabelCompany,
				Longitude:    decimalPtr("116.3099"),
				Latitude:     decimalPtr("39.9042"),
			},
			wantErr: false,
		},
		{
			name:          "invalid receiver name - empty",
			userID:        "user-123",
			receiverName:  "",
			receiverPhone: "13800138000",
			province:      "广东省",
			city:          "深圳市",
			district:      "南山区",
			street:        "科技园南区",
			opts:          nil,
			wantErr:       true,
		},
		{
			name:          "invalid phone number",
			userID:        "user-123",
			receiverName:  "张三",
			receiverPhone: "invalid-phone",
			province:      "广东省",
			city:          "深圳市",
			district:      "南山区",
			street:        "科技园南区",
			opts:          nil,
			wantErr:       true,
		},
		{
			name:          "invalid postal code",
			userID:        "user-123",
			receiverName:  "张三",
			receiverPhone: "13800138000",
			province:      "广东省",
			city:          "深圳市",
			district:      "南山区",
			street:        "科技园南区",
			opts: &AddressOptions{
				PostalCode: stringPtr("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := CreateAddress(
				tt.userID,
				tt.receiverName,
				tt.receiverPhone,
				tt.province,
				tt.city,
				tt.district,
				tt.street,
				tt.opts,
			)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, addr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, addr)
				assert.Equal(t, tt.userID, addr.UserID)
				assert.Equal(t, tt.receiverName, addr.ReceiverName)
				assert.Equal(t, tt.receiverPhone, addr.ReceiverPhone)
				assert.Equal(t, tt.province, addr.Province)
				assert.Equal(t, tt.city, addr.City)
				assert.Equal(t, tt.district, addr.District)
				assert.Equal(t, tt.street, addr.Street)
				assert.False(t, addr.IsDefault) // 新创建的地址默认不是默认地址

				// 验证可选参数
				if tt.opts != nil {
					if tt.opts.PostalCode != nil {
						assert.Equal(t, tt.opts.PostalCode, addr.PostalCode)
					}
					if tt.opts.AddressLabel != "" {
						assert.Equal(t, tt.opts.AddressLabel, addr.AddressLabel)
					}
					if tt.opts.Longitude != nil && tt.opts.Latitude != nil {
						assert.Equal(t, tt.opts.Longitude, addr.Longitude)
						assert.Equal(t, tt.opts.Latitude, addr.Latitude)
					}
				}
			}
		})
	}
}

func TestAddress_Update(t *testing.T) {
	// 创建一个地址
	addr, err := CreateAddress(
		"user-123",
		"张三",
		"13800138000",
		"广东省",
		"深圳市",
		"南山区",
		"科技园南区",
		nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, addr)

	// 更新地址
	opts := &AddressOptions{
		PostalCode:   stringPtr("518000"),
		AddressLabel: LabelHome,
	}

	err = addr.Update(
		"李四",
		"13900139000",
		"北京市",
		"北京市",
		"海淀区",
		"中关村大街1号",
		opts,
	)

	assert.NoError(t, err)
	assert.Equal(t, "李四", addr.ReceiverName)
	assert.Equal(t, "13900139000", addr.ReceiverPhone)
	assert.Equal(t, "北京市", addr.Province)
	assert.Equal(t, "北京市", addr.City)
	assert.Equal(t, "海淀区", addr.District)
	assert.Equal(t, "中关村大街1号", addr.Street)
	assert.Equal(t, stringPtr("518000"), addr.PostalCode)
	assert.Equal(t, LabelHome, addr.AddressLabel)
}

func TestAddress_SetAsDefault(t *testing.T) {
	addr, err := CreateAddress(
		"user-123",
		"张三",
		"13800138000",
		"广东省",
		"深圳市",
		"南山区",
		"科技园南区",
		nil,
	)
	assert.NoError(t, err)
	assert.False(t, addr.IsDefault)

	addr.SetAsDefault()
	assert.True(t, addr.IsDefault)
}

func TestAddress_GetFullAddress(t *testing.T) {
	addr, err := CreateAddress(
		"user-123",
		"张三",
		"13800138000",
		"广东省",
		"深圳市",
		"南山区",
		"科技园南区腾讯大厦",
		nil,
	)
	assert.NoError(t, err)

	fullAddress := addr.GetFullAddress()
	expected := "广东省深圳市南山区科技园南区腾讯大厦"
	assert.Equal(t, expected, fullAddress)
}

func TestValidateReceiverPhone(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"valid mobile", "13800138000", false},
		{"valid mobile 2", "18900189000", false},
		{"valid landline", "0755-12345678", false},
		{"valid landline without dash", "075512345678", false},
		{"empty phone", "", true},
		{"invalid mobile", "12345678901", true},
		{"invalid landline", "0755-123", true},
		{"invalid format", "abc123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReceiverPhone(tt.phone)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePostalCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{"valid code", "518000", false},
		{"empty code", "", false}, // 邮政编码可以为空
		{"invalid code - too short", "12345", true},
		{"invalid code - too long", "1234567", true},
		{"invalid code - letters", "abc123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePostalCode(tt.code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// 辅助函数
func stringPtr(s string) *string {
	return &s
}

func decimalPtr(s string) *decimal.Decimal {
	d, _ := decimal.NewFromString(s)
	return &d
}
