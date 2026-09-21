package model
import ("time"; "github.com/google/uuid")
type MarketplaceSKUMapping struct {
 ID uuid.UUID; TenantID uuid.UUID; IntegrationID uuid.UUID; ProductID uuid.UUID; VariantID *uuid.UUID
 Provider string; ExternalItemID string; ExternalModelID string; ExternalSKU *string; MappingSource string; Active bool
 CreatedAt time.Time; UpdatedAt time.Time
}
