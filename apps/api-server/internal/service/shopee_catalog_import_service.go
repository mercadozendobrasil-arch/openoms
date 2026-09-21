package service

import (
 "context"
 "encoding/json"
 "fmt"
 "github.com/google/uuid"
 "github.com/jackc/pgx/v5"
 "github.com/jackc/pgx/v5/pgxpool"
 "github.com/openoms-org/openoms/apps/api-server/internal/database"
 shopeeint "github.com/openoms-org/openoms/apps/api-server/internal/integration/shopee"
 "github.com/openoms-org/openoms/apps/api-server/internal/model"
 "github.com/openoms-org/openoms/apps/api-server/internal/repository"
)

type ShopeeCatalogImportResult struct{ Created, Updated, Variants, Mappings, Errors int }

type ShopeeCatalogImportService struct {
 pool *pgxpool.Pool
 products *repository.ProductRepository
 variants *repository.VariantRepository
 listings *repository.ProductListingRepository
 mappings *repository.MarketplaceSKUMappingRepository
}

func NewShopeeCatalogImportService(pool *pgxpool.Pool)*ShopeeCatalogImportService{
 return &ShopeeCatalogImportService{pool:pool,products:repository.NewProductRepository(),variants:repository.NewVariantRepository(),listings:repository.NewProductListingRepository(),mappings:repository.NewMarketplaceSKUMappingRepository()}
}

func (s *ShopeeCatalogImportService) ImportPage(ctx context.Context,tenantID,integrationID uuid.UUID,p *shopeeint.Provider,offset,pageSize int,status string)(*ShopeeCatalogImportResult,bool,int,error){
 snapshots,hasNext,next,err:=p.ReadListingPage(ctx,offset,pageSize,status);if err!=nil{return nil,false,offset,err}
 result:=&ShopeeCatalogImportResult{}
 err=database.WithTenant(ctx,s.pool,tenantID,func(tx pgx.Tx) error{
  for _,snap:=range snapshots{
   if err:=s.upsertSnapshot(ctx,tx,tenantID,integrationID,snap,result);err!=nil{result.Errors++}
  }
  return nil
 })
 if err!=nil{return nil,false,offset,fmt.Errorf("shopee catalog import: %w",err)}
 return result,hasNext,next,nil
}

func (s *ShopeeCatalogImportService) upsertSnapshot(ctx context.Context,tx pgx.Tx,tenantID,integrationID uuid.UUID,snap shopeeint.ListingSnapshot,result *ShopeeCatalogImportResult) error{
 existing,err:=s.listings.FindByExternalIDAndIntegration(ctx,tx,snap.ItemID,integrationID);if err!=nil{return err}
 if existing!=nil{
  status:="inactive";if snap.Status=="NORMAL"{status="active"}
  metaBytes,_:=json.Marshal(map[string]any{"shop_id":snap.ShopID,"item_id":snap.ItemID});meta:=json.RawMessage(metaBytes)
  if err:=s.listings.Update(ctx,tx,existing.ID,&model.UpdateProductListingRequest{Status:&status,Metadata:&meta});err!=nil{return err}
  result.Updated++
  return s.upsertExistingMappings(ctx,tx,tenantID,integrationID,existing.ProductID,snap,result)
 }
 plan:=shopeeint.BuildCatalogMapping(tenantID,integrationID,snap)
 if err:=s.products.Create(ctx,tx,&plan.Product);err!=nil{return err}
 for i:=range plan.Variants{if err:=s.variants.Create(ctx,tx,&plan.Variants[i]);err!=nil{return err};result.Variants++}
 if err:=s.listings.Create(ctx,tx,&plan.Listing);err!=nil{return err}
 for _,m:=range shopeeint.BuildSKUMappings(tenantID,integrationID,snap,plan){mm:=m;if err:=s.mappings.Upsert(ctx,tx,&mm);err!=nil{return err};result.Mappings++}
 result.Created++;return nil
}

func (s *ShopeeCatalogImportService) upsertExistingMappings(ctx context.Context,tx pgx.Tx,tenantID,integrationID,productID uuid.UUID,snap shopeeint.ListingSnapshot,result *ShopeeCatalogImportResult) error{
 for _,m:=range snap.Models{
  found,err:=s.mappings.FindExternal(ctx,tx,integrationID,snap.ItemID,m.ModelID);if err!=nil{return err}
  if found!=nil{continue}
  // Existing listings are never auto-linked to an arbitrary variant. Record item/model identity first.
  var sku *string;if m.ModelSKU!=""{v:=m.ModelSKU;sku=&v}
  row:=model.MarketplaceSKUMapping{ID:uuid.New(),TenantID:tenantID,IntegrationID:integrationID,ProductID:productID,Provider:"shopee",ExternalItemID:snap.ItemID,ExternalModelID:m.ModelID,ExternalSKU:sku,MappingSource:"import_unresolved",Active:true}
  if err:=s.mappings.Upsert(ctx,tx,&row);err!=nil{return err};result.Mappings++
 }
 return nil
}
