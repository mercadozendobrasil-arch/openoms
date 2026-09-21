package shopee

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ItemSummary struct {
	ItemID     int64  `json:"item_id"`
	ItemStatus string `json:"item_status,omitempty"`
	UpdateTime int64  `json:"update_time,omitempty"`
}

type ItemListResponse struct {
	Item      []ItemSummary `json:"item"`
	TotalCount int          `json:"total_count,omitempty"`
	HasNextPage bool        `json:"has_next_page,omitempty"`
	NextOffset int          `json:"next_offset,omitempty"`
}

type ItemBaseInfo struct {
	ItemID      int64   `json:"item_id"`
	CategoryID  int64   `json:"category_id,omitempty"`
	ItemName    string  `json:"item_name,omitempty"`
	ItemSKU     string  `json:"item_sku,omitempty"`
	Description string  `json:"description,omitempty"`
	ItemStatus  string  `json:"item_status,omitempty"`
	PriceInfo   []PriceInfo `json:"price_info,omitempty"`
	StockInfo   []StockInfo `json:"stock_info_v2,omitempty"`
}

type PriceInfo struct {
	Currency      string  `json:"currency,omitempty"`
	OriginalPrice float64 `json:"original_price,omitempty"`
	CurrentPrice  float64 `json:"current_price,omitempty"`
}

type StockInfo struct {
	SummaryInfo struct {
		TotalReservedStock int `json:"total_reserved_stock,omitempty"`
		TotalAvailableStock int `json:"total_available_stock,omitempty"`
	} `json:"summary_info,omitempty"`
}

type Model struct {
	ModelID   int64  `json:"model_id"`
	ModelSKU  string `json:"model_sku,omitempty"`
	ModelName string `json:"model_name,omitempty"`
	PriceInfo []PriceInfo `json:"price_info,omitempty"`
	StockInfo []StockInfo `json:"stock_info_v2,omitempty"`
	TierIndex []int `json:"tier_index,omitempty"`
}

func (c *Client) GetItemList(ctx context.Context, accessToken, shopID, itemStatus string, offset, pageSize int) (ItemListResponse, error) {
	if pageSize <= 0 { pageSize = 50 }
	q:=url.Values{}
	q.Set("offset",strconv.Itoa(offset)); q.Set("page_size",strconv.Itoa(pageSize))
	if itemStatus!="" { q.Set("item_status",itemStatus) }
	var out ItemListResponse
	err:=c.Do(ctx,Request{Method:http.MethodGet,Path:"/product/get_item_list",AccessToken:accessToken,ShopID:shopID,Query:q},&out)
	return out,err
}

func (c *Client) GetItemBaseInfo(ctx context.Context, accessToken, shopID string, itemIDs []int64) ([]ItemBaseInfo,error) {
	ids:=make([]string,len(itemIDs)); for i,id:=range itemIDs { ids[i]=strconv.FormatInt(id,10) }
	q:=url.Values{"item_id_list":[]string{strings.Join(ids,",")}}
	var out struct{ ItemList []ItemBaseInfo `json:"item_list"` }
	err:=c.Do(ctx,Request{Method:http.MethodGet,Path:"/product/get_item_base_info",AccessToken:accessToken,ShopID:shopID,Query:q},&out)
	return out.ItemList,err
}

func (c *Client) GetModelList(ctx context.Context, accessToken, shopID string, itemID int64) ([]Model,error) {
	q:=url.Values{"item_id":[]string{strconv.FormatInt(itemID,10)}}
	var out struct{ Model []Model `json:"model"` }
	err:=c.Do(ctx,Request{Method:http.MethodGet,Path:"/product/get_model_list",AccessToken:accessToken,ShopID:shopID,Query:q},&out)
	return out.Model,err
}
