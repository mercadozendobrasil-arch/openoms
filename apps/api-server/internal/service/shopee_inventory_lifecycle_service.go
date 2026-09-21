package service

import(
 "context"
 "encoding/json"
 "fmt"
 "strings"
 "github.com/google/uuid"
 "github.com/jackc/pgx/v5"
 "github.com/jackc/pgx/v5/pgxpool"
 "github.com/openoms-org/openoms/apps/api-server/internal/database"
 "github.com/openoms-org/openoms/apps/api-server/internal/repository"
)

type ShopeeInventoryLifecycleService struct{pool *pgxpool.Pool;warehouses *repository.WarehouseRepository;stock *repository.WarehouseStockRepository;bundles *repository.BundleRepository}
func NewShopeeInventoryLifecycleService(pool *pgxpool.Pool)*ShopeeInventoryLifecycleService{return &ShopeeInventoryLifecycleService{pool:pool,warehouses:repository.NewWarehouseRepository(),stock:repository.NewWarehouseStockRepository(),bundles:repository.NewBundleRepository()}}

func(s *ShopeeInventoryLifecycleService)Apply(ctx context.Context,tenantID,integrationID uuid.UUID,orderSN,status string,itemsJSON json.RawMessage)error{
 return database.WithTenant(ctx,s.pool,tenantID,func(tx pgx.Tx)error{
  wh,err:=s.warehouses.FindDefault(ctx,tx);if err!=nil{return err};if wh==nil{return fmt.Errorf("default warehouse not configured")}
  switch strings.ToUpper(status){
  case "READY_TO_SHIP","PROCESSED":
   return s.once(ctx,tx,tenantID,integrationID,orderSN,"reserve",func()error{return s.adjust(ctx,tx,wh.ID,itemsJSON,1,true)})
  case "SHIPPED","COMPLETED":
   return s.once(ctx,tx,tenantID,integrationID,orderSN,"commit",func()error{
    if err:=s.adjust(ctx,tx,wh.ID,itemsJSON,-1,true);err!=nil{return err}
    return s.adjust(ctx,tx,wh.ID,itemsJSON,-1,false)
   })
  case "CANCELLED","IN_CANCEL":
   return s.once(ctx,tx,tenantID,integrationID,orderSN,"release",func()error{return s.adjust(ctx,tx,wh.ID,itemsJSON,-1,true)})
  default:return nil
  }
 })
}
func(s *ShopeeInventoryLifecycleService)once(ctx context.Context,tx pgx.Tx,tenantID,integrationID uuid.UUID,orderSN,event string,fn func()error)error{
 tag,err:=tx.Exec(ctx,"INSERT INTO marketplace_order_inventory_events(id,tenant_id,integration_id,external_order_id,event_type) VALUES($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING",uuid.New(),tenantID,integrationID,orderSN,event);if err!=nil{return err};if tag.RowsAffected()==0{return nil};return fn()
}
func(s *ShopeeInventoryLifecycleService)adjust(ctx context.Context,tx pgx.Tx,warehouseID uuid.UUID,raw json.RawMessage,direction int,reserved bool)error{
 var items []shopeeResolvedOrderItem;if err:=json.Unmarshal(raw,&items);err!=nil{return err}
 for _,it:=range items{if it.ProductID==""||it.Quantity<=0{continue};pid,err:=uuid.Parse(it.ProductID);if err!=nil{continue};var vid *uuid.UUID;if it.VariantID!=""{v,e:=uuid.Parse(it.VariantID);if e==nil{vid=&v}}
  components,err:=s.bundles.ListByBundleProduct(ctx,tx,pid);if err!=nil{return err}
  if len(components)>0{for _,c:=range components{delta:=direction*it.Quantity*c.Quantity;if reserved{if err:=s.stock.AdjustReserved(ctx,tx,warehouseID,c.ComponentProductID,c.ComponentVariantID,delta);err!=nil{return err}}else{if err:=s.stock.AdjustQuantity(ctx,tx,warehouseID,c.ComponentProductID,c.ComponentVariantID,delta);err!=nil{return err}}};continue}
  delta:=direction*it.Quantity;if reserved{if err:=s.stock.AdjustReserved(ctx,tx,warehouseID,pid,vid,delta);err!=nil{return err}}else{if err:=s.stock.AdjustQuantity(ctx,tx,warehouseID,pid,vid,delta);err!=nil{return err}}
 };return nil
}
