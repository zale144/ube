package model

import "fmt"

// Base is the base part of the product model
type Base struct {
	BaseKey
	Product          string `json:"product,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	ShortDescription string `json:"short_description,omitempty"`
}

var _ Entity = Base{}

type BaseKey struct {
	ProductID int32 `json:"product_id,omitempty"`
}

func (b BaseKey) PK() string {
	return fmt.Sprint(b.ProductID)
}

// GetKey returns the key to the entity
func (p Base) GetKey() Key {
	return p.BaseKey
}

// PriceBase holds the pricing data relating to the product
type PriceBase struct {
	ID        string  `json:"id,omitempty"`
	Price     float64 `json:"price"`
	PriceList string  `json:"price_list,omitempty"`
	Currency  string  `json:"currency,omitempty"`
}

// Campaign holds details relating to the price
type Campaign struct {
	ID       string  `json:"id,omitempty"`
	Campaign string  `json:"campaign,omitempty"`
	Price    float64 `json:"price"`
}

type EntityKey struct {
	ProductID int32  `json:"product_id,omitempty"`
	SKU       string `json:"sku"`
}

func (e EntityKey) PK() string {
	return fmt.Sprintf("%d_%s", e.ProductID, e.SKU)
}

// TODO: put below struct in a separate file?

// ApparelBase is a base part of a common apparel product model
type ApparelBase struct {
	EntityKey
	Base
	VariantSKU        string  `json:"variant_sku,omitempty"`
	VariantID         int32   `json:"variant_id,omitempty"`
	SizeSKU           string  `json:"size_sku,omitempty"`
	Brand             string  `json:"brand,omitempty"`
	Collection        string  `json:"collection,omitempty"`
	Variant           string  `json:"variant,omitempty"`
	Size              string  `json:"size,omitempty"`
	SizeComment       string  `json:"size_comment,omitempty"`
	StockItemID       int32   `json:"stock_item_id,omitempty"`
	Weight            float64 `json:"weight,omitempty"`
	WeightUnit        string  `json:"weight_unit,omitempty"`
	CountryOfOrigin   string  `json:"country_of_origin,omitempty"`
	Active            int32   `json:"active,omitempty"`
	Store             string  `json:"store,omitempty"`
	MetaTitle         string  `json:"meta_title,omitempty"`
	MetaDescription   string  `json:"meta_description,omitempty"`
	MetaKeywords      string  `json:"meta_keywords,omitempty"`
	CostPrice         float64 `json:"cost_price,omitempty"`
	CostPriceCurrency string  `json:"cost_price_currency,omitempty"`
}

// GetKey returns the custom composed key to the entity
func (p ApparelBase) GetKey() Key {
	return p.EntityKey
}

var _ Entity = ApparelBase{}
