package shopee
import("context";"net/http")
type StockUpdate struct{ModelID int64 `json:"model_id"`;NormalStock int `json:"normal_stock"`;SellerStock []SellerStock `json:"seller_stock,omitempty"`}
type SellerStock struct{LocationID string `json:"location_id,omitempty"`;Stock int `json:"stock"`}
func(c *Client)UpdateStock(ctx context.Context,token,shopID string,itemID int64,stocks []StockUpdate)error{body:=map[string]any{"item_id":itemID,"stock_list":stocks};var out any;return c.Do(ctx,Request{Method:http.MethodPost,Path:"/product/update_stock",AccessToken:token,ShopID:shopID,Body:body},&out)}
