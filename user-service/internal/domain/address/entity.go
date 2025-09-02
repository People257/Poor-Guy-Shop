package address

import (
	"errors"
	"regexp"
	"time"

	"github.com/shopspring/decimal"
)

// Address 地址领域实体
type Address struct {
	ID            string
	UserID        string
	ReceiverName  string
	ReceiverPhone string
	Province      string
	City          string
	District      string
	Street        string
	PostalCode    *string
	AddressLabel  string
	IsDefault     bool
	Longitude     *decimal.Decimal
	Latitude      *decimal.Decimal
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// AddressLabel 地址标签常量
const (
	LabelDefault = "默认"
	LabelHome    = "家"
	LabelCompany = "公司"
	LabelSchool  = "学校"
	LabelOther   = "其他"
)

// CreateAddress 创建新地址
func CreateAddress(userID, receiverName, receiverPhone, province, city, district, street string, opts *AddressOptions) (*Address, error) {
	// 验证必填字段
	if err := ValidateReceiverName(receiverName); err != nil {
		return nil, err
	}
	if err := ValidateReceiverPhone(receiverPhone); err != nil {
		return nil, err
	}
	if err := ValidateAddress(province, city, district, street); err != nil {
		return nil, err
	}

	now := time.Now()
	address := &Address{
		UserID:        userID,
		ReceiverName:  receiverName,
		ReceiverPhone: receiverPhone,
		Province:      province,
		City:          city,
		District:      district,
		Street:        street,
		AddressLabel:  LabelDefault,
		IsDefault:     false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// 应用可选参数
	if opts != nil {
		if opts.PostalCode != nil {
			if err := ValidatePostalCode(*opts.PostalCode); err != nil {
				return nil, err
			}
			address.PostalCode = opts.PostalCode
		}
		if opts.AddressLabel != "" {
			address.AddressLabel = opts.AddressLabel
		}
		if opts.Longitude != nil && opts.Latitude != nil {
			address.Longitude = opts.Longitude
			address.Latitude = opts.Latitude
		}
	}

	return address, nil
}

// AddressOptions 地址可选参数
type AddressOptions struct {
	PostalCode   *string
	AddressLabel string
	Longitude    *decimal.Decimal
	Latitude     *decimal.Decimal
}

// Update 更新地址信息
func (a *Address) Update(receiverName, receiverPhone, province, city, district, street string, opts *AddressOptions) error {
	// 验证必填字段
	if err := ValidateReceiverName(receiverName); err != nil {
		return err
	}
	if err := ValidateReceiverPhone(receiverPhone); err != nil {
		return err
	}
	if err := ValidateAddress(province, city, district, street); err != nil {
		return err
	}

	a.ReceiverName = receiverName
	a.ReceiverPhone = receiverPhone
	a.Province = province
	a.City = city
	a.District = district
	a.Street = street
	a.UpdatedAt = time.Now()

	// 应用可选参数
	if opts != nil {
		if opts.PostalCode != nil {
			if err := ValidatePostalCode(*opts.PostalCode); err != nil {
				return err
			}
			a.PostalCode = opts.PostalCode
		}
		if opts.AddressLabel != "" {
			a.AddressLabel = opts.AddressLabel
		}
		if opts.Longitude != nil && opts.Latitude != nil {
			a.Longitude = opts.Longitude
			a.Latitude = opts.Latitude
		}
	}

	return nil
}

// SetAsDefault 设置为默认地址
func (a *Address) SetAsDefault() {
	a.IsDefault = true
	a.UpdatedAt = time.Now()
}

// UnsetAsDefault 取消默认地址
func (a *Address) UnsetAsDefault() {
	a.IsDefault = false
	a.UpdatedAt = time.Now()
}

// GetFullAddress 获取完整地址字符串
func (a *Address) GetFullAddress() string {
	return a.Province + a.City + a.District + a.Street
}

// ValidateReceiverName 验证收货人姓名
func ValidateReceiverName(name string) error {
	if name == "" {
		return errors.New("收货人姓名不能为空")
	}
	if len(name) > 100 {
		return errors.New("收货人姓名长度不能超过100字符")
	}
	return nil
}

// ValidateReceiverPhone 验证收货人电话
func ValidateReceiverPhone(phone string) error {
	if phone == "" {
		return errors.New("收货人电话不能为空")
	}

	// 支持手机号和固定电话
	mobileRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	landlineRegex := regexp.MustCompile(`^0\d{2,3}-?\d{7,8}$`)

	if !mobileRegex.MatchString(phone) && !landlineRegex.MatchString(phone) {
		return errors.New("电话号码格式不正确")
	}

	return nil
}

// ValidateAddress 验证地址信息
func ValidateAddress(province, city, district, street string) error {
	if province == "" {
		return errors.New("省份不能为空")
	}
	if city == "" {
		return errors.New("城市不能为空")
	}
	if district == "" {
		return errors.New("区/县不能为空")
	}
	if street == "" {
		return errors.New("街道地址不能为空")
	}
	if len(street) > 200 {
		return errors.New("街道地址长度不能超过200字符")
	}
	return nil
}

// ValidatePostalCode 验证邮政编码
func ValidatePostalCode(code string) error {
	if code == "" {
		return nil // 邮政编码可以为空
	}

	postalRegex := regexp.MustCompile(`^\d{6}$`)
	if !postalRegex.MatchString(code) {
		return errors.New("邮政编码格式不正确")
	}

	return nil
}
