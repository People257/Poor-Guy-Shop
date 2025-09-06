package cart

import (
	"time"

	"github.com/shopspring/decimal"
)

// ShoppingCart 购物车项实体（对应数据库模型）
type ShoppingCart struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	ProductID string          `json:"product_id"`
	SkuID     string          `json:"sku_id"`
	Quantity  int32           `json:"quantity"`
	Price     decimal.Decimal `json:"price"`
	Selected  bool            `json:"selected"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
	Version   int32           `json:"version"`
}

// CartItem 购物车商品项实体
type CartItem struct {
	ID           string          `json:"id"`
	UserID       string          `json:"user_id"`
	ProductID    string          `json:"product_id"`
	SkuID        string          `json:"sku_id"`
	ProductName  string          `json:"product_name"`
	ProductImage string          `json:"product_image"`
	SkuName      string          `json:"sku_name"`
	Price        decimal.Decimal `json:"price"`
	Quantity     int             `json:"quantity"`
	TotalAmount  decimal.Decimal `json:"total_amount"`
	Selected     bool            `json:"selected"`
	Available    bool            `json:"available"` // 商品是否可用（库存充足等）
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
	DeletedAt    *time.Time      `json:"deleted_at"`
	Version      int             `json:"version"`
}

// CartSummary 购物车汇总信息
type CartSummary struct {
	TotalItems     int             `json:"total_items"`
	SelectedItems  int             `json:"selected_items"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	SelectedAmount decimal.Decimal `json:"selected_amount"`
}

// UpdateQuantity 更新数量
func (c *CartItem) UpdateQuantity(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	c.Quantity = quantity
	c.TotalAmount = c.Price.Mul(decimal.NewFromInt(int64(quantity)))
	return nil
}

// Select 选择商品
func (c *CartItem) Select() {
	c.Selected = true
}

// Unselect 取消选择商品
func (c *CartItem) Unselect() {
	c.Selected = false
}

// SetAvailable 设置商品可用性
func (c *CartItem) SetAvailable(available bool) {
	c.Available = available
	if !available {
		c.Selected = false // 不可用商品自动取消选择
	}
}

// CalculateTotal 计算小计
func (c *CartItem) CalculateTotal() {
	c.TotalAmount = c.Price.Mul(decimal.NewFromInt(int64(c.Quantity)))
}

// Cart 购物车聚合根
type Cart struct {
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items"`
}

// AddItem 添加商品到购物车
func (c *Cart) AddItem(item CartItem) {
	item.UserID = c.UserID
	item.CalculateTotal()

	// 检查是否已存在相同商品SKU
	for i, existingItem := range c.Items {
		if existingItem.ProductID == item.ProductID && existingItem.SkuID == item.SkuID {
			// 如果已存在，更新数量
			c.Items[i].Quantity += item.Quantity
			c.Items[i].CalculateTotal()
			return
		}
	}

	// 不存在则添加新商品
	c.Items = append(c.Items, item)
}

// RemoveItem 移除商品
func (c *Cart) RemoveItem(itemID string) error {
	for i, item := range c.Items {
		if item.ID == itemID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return nil
		}
	}
	return ErrCartItemNotFound
}

// UpdateItemQuantity 更新商品数量
func (c *Cart) UpdateItemQuantity(itemID string, quantity int) error {
	for i, item := range c.Items {
		if item.ID == itemID {
			return c.Items[i].UpdateQuantity(quantity)
		}
	}
	return ErrCartItemNotFound
}

// SelectItems 选择商品
func (c *Cart) SelectItems(itemIDs []string, selected bool) {
	itemIDMap := make(map[string]bool)
	for _, id := range itemIDs {
		itemIDMap[id] = true
	}

	for i, item := range c.Items {
		if itemIDMap[item.ID] {
			c.Items[i].Selected = selected
		}
	}
}

// SelectAll 全选/全不选
func (c *Cart) SelectAll(selected bool) {
	for i := range c.Items {
		if c.Items[i].Available { // 只能选择可用商品
			c.Items[i].Selected = selected
		}
	}
}

// GetSelectedItems 获取已选择的商品
func (c *Cart) GetSelectedItems() []CartItem {
	var selectedItems []CartItem
	for _, item := range c.Items {
		if item.Selected {
			selectedItems = append(selectedItems, item)
		}
	}
	return selectedItems
}

// CalculateSummary 计算购物车汇总
func (c *Cart) CalculateSummary() CartSummary {
	summary := CartSummary{}

	for _, item := range c.Items {
		summary.TotalItems += item.Quantity
		summary.TotalAmount = summary.TotalAmount.Add(item.TotalAmount)

		if item.Selected {
			summary.SelectedItems += item.Quantity
			summary.SelectedAmount = summary.SelectedAmount.Add(item.TotalAmount)
		}
	}

	return summary
}

// Clear 清空购物车
func (c *Cart) Clear() {
	c.Items = []CartItem{}
}

// IsEmpty 检查购物车是否为空
func (c *Cart) IsEmpty() bool {
	return len(c.Items) == 0
}

// HasSelectedItems 检查是否有选中的商品
func (c *Cart) HasSelectedItems() bool {
	for _, item := range c.Items {
		if item.Selected {
			return true
		}
	}
	return false
}
