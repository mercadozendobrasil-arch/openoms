package service

import (
 "context"
 "encoding/json"
 "fmt"
 "strings"
 "time"

 "github.com/google/uuid"
 "github.com/jackc/pgx/v5"
 "github.com/jackc/pgx/v5/pgxpool"
 "github.com/openoms-org/openoms/apps/api-server/internal/database"
 shopeeint "github.com/openoms-org/openoms/apps/api-server/internal/integration/shopee"
 "github.com/openoms-org/openoms/apps/api-server/internal/model"
 "github.com/openoms-org/openoms/apps/api-server/internal/repository"
)

type ShopeeOrderImportResult struct{Created,Updated,UnresolvedItems,Errors int}

type ShopeeOrderImportService struct{
 pool *pgxpool.Pool
 orders *repository.OrderRepository
 mappings *repository.MarketplaceSKUMappingRepository
}

func NewShopeeOrderImportService(pool *pgxpool.Pool)*ShopeeOrderImportService{
 return &ShopeeOrderImportService{pool:pool,orders:repository.NewOrderRepository(),mappings:repository.NewMarketplaceSKUMappingRepository()}
}

func(s *ShopeeOrderImportService)ImportWindow(ctx context.Context,tenantID,integrationID uuid.UUID,p *shopeeint.Provider,from,to time.Time,cursor,status string)(*ShopeeOrderImportResult,string,bool,error){
 snapshots,next,more,err:=p.ReadOrders(ctx,from,to,cursor,status);if err!=nil{return nil,"",false,err}
 result:=&ShopeeOrderImportResult{}
 err=database.WithTenant(ctx,s.pool,tenantID,func(tx pgx.Tx)error{
  for _,snap:=range snapshots{if err:=s.upsert(ctx,tx,tenantID,integrationID,snap,result);err!=nil{result.Errors++}}
  return nil
 })
 if err!=nil{return nil,"",false,fmt.Errorf("shopee order import: %w",err)}
 return result,next,more,nil
}

func(s *ShopeeOrderImportService)upsert(ctx context.Context,tx pgx.Tx,tenantID,integrationID uuid.UUID,snap shopeeint.OrderSnapshot,result *ShopeeOrderImportResult)error{
 items:=make([]map[string]any,0,len(snap.Items))
 for _,it:=range snap.Items{
  resolved,err:=s.mappings.FindExternal(ctx,tx,integrationID,it.ItemID,it.ModelID);if err!=nil{return err}
  row:=map[string]any{"item_id":it.ItemID,"model_id":it.ModelID,"item_sku":it.ItemSKU,"model_sku":it.ModelSKU,"name":it.Name,"quantity":it.Quantity,"unit_price":it.UnitPrice}
  if resolved!=nil{row["product_id"]=resolved.ProductID.String();if resolved.VariantID!=nil{row["variant_id"]=resolved.VariantID.String()}}else{row["sku_resolution"]="unresolved";result.UnresolvedItems++}
  items=append(items,row)
 }
 itemJSON,_:=json.Marshal(items)
 meta,_:=json.Marshal(map[string]any{"shopee_order_status":snap.Status,"integration_id":integrationID.String()})
 address:=snap.ShippingAddress;if len(address)==0{address=json.RawMessage("{}")}
 status:=mapShopeeOrderStatus(snap.Status)
 payment:=snap.PaymentMethod
 existing,err:=s.orders.FindByExternalID(ctx,tx,"shopee",snap.OrderSN);if err!=nil{return err}
 if existing!=nil{
  currency:=snap.Currency;amount:=snap.TotalAmount
  req:=model.UpdateOrderRequest{Items:itemJSON,ShippingAddress:address,TotalAmount:&amount,Currency:&currency,PaymentMethod:&payment,Metadata:meta}
  if err:=s.orders.Update(ctx,tx,existing.ID,req);err!=nil{return err}
  if existing.Status!=status{if err:=s.orders.UpdateStatus(ctx,tx,existing.ID,status,nil,nil);err!=nil{return err}}
  result.Updated++;return nil
 }
 ext:=snap.OrderSN;integration:=integrationID;ordered:=snap.CreatedAt
 customer:=strings.TrimSpace(snap.BuyerUsername);if customer==""{customer="Shopee buyer"}
 o:=model.Order{ID:uuid.New(),TenantID:tenantID,ExternalID:&ext,Source:"shopee",IntegrationID:&integration,Status:status,CustomerName:customer,ShippingAddress:address,BillingAddress:json.RawMessage("{}"),Items:itemJSON,TotalAmount:snap.TotalAmount,Currency:snap.Currency,Metadata:meta,Tags:[]string{"shopee"},OrderedAt:&ordered,PaymentStatus:mapShopeePaymentStatus(snap.Status),PaymentMethod:&payment,Priority:"normal"}
 created,err:=s.orders.CreateIfExternalIDNotExists(ctx,tx,&o);if err!=nil{return err};if created{result.Created++}else{result.Updated++};return nil
}

func mapShopeeOrderStatus(v string)string{switch strings.ToUpper(v){case "UNPAID":return model.OrderStatusNew;case "READY_TO_SHIP","PROCESSED":return model.OrderStatusReadyToShip;case "SHIPPED":return model.OrderStatusShipped;case "COMPLETED":return model.OrderStatusCompleted;case "CANCELLED","IN_CANCEL":return model.OrderStatusCancelled;default:return model.OrderStatusProcessing}}
func mapShopeePaymentStatus(v string)string{if strings.EqualFold(v,"UNPAID"){return "unpaid"};if strings.EqualFold(v,"CANCELLED"){return "cancelled"};return "paid"}
