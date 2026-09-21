package repository
import ("context";"fmt";"github.com/google/uuid";"github.com/jackc/pgx/v5";"github.com/openoms-org/openoms/apps/api-server/internal/model")
type MarketplaceSKUMappingRepository struct{}
func NewMarketplaceSKUMappingRepository()*MarketplaceSKUMappingRepository{return &MarketplaceSKUMappingRepository{}}
func (r *MarketplaceSKUMappingRepository) Upsert(ctx context.Context,tx pgx.Tx,m *model.MarketplaceSKUMapping) error {
 q:="INSERT INTO marketplace_sku_mappings (id,tenant_id,integration_id,product_id,variant_id,provider,external_item_id,external_model_id,external_sku,mapping_source,active) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT(integration_id,external_item_id,external_model_id) DO UPDATE SET product_id=EXCLUDED.product_id,variant_id=EXCLUDED.variant_id,external_sku=EXCLUDED.external_sku,mapping_source=EXCLUDED.mapping_source,active=EXCLUDED.active,updated_at=now() RETURNING created_at,updated_at"
 return tx.QueryRow(ctx,q,m.ID,m.TenantID,m.IntegrationID,m.ProductID,m.VariantID,m.Provider,m.ExternalItemID,m.ExternalModelID,m.ExternalSKU,m.MappingSource,m.Active).Scan(&m.CreatedAt,&m.UpdatedAt)
}
func (r *MarketplaceSKUMappingRepository) FindExternal(ctx context.Context,tx pgx.Tx,integrationID uuid.UUID,itemID,modelID string)(*model.MarketplaceSKUMapping,error){
 var m model.MarketplaceSKUMapping
 q:="SELECT id,tenant_id,integration_id,product_id,variant_id,provider,external_item_id,external_model_id,external_sku,mapping_source,active,created_at,updated_at FROM marketplace_sku_mappings WHERE integration_id=$1 AND external_item_id=$2 AND external_model_id=$3"
 err:=tx.QueryRow(ctx,q,integrationID,itemID,modelID).Scan(&m.ID,&m.TenantID,&m.IntegrationID,&m.ProductID,&m.VariantID,&m.Provider,&m.ExternalItemID,&m.ExternalModelID,&m.ExternalSKU,&m.MappingSource,&m.Active,&m.CreatedAt,&m.UpdatedAt)
 if err==pgx.ErrNoRows{return nil,nil};if err!=nil{return nil,fmt.Errorf("find marketplace sku mapping: %w",err)};return &m,nil
}
