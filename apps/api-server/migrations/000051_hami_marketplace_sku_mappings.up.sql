CREATE TABLE marketplace_sku_mappings (
 id uuid PRIMARY KEY, tenant_id uuid NOT NULL, integration_id uuid NOT NULL REFERENCES integrations(id) ON DELETE CASCADE,
 product_id uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE, variant_id uuid REFERENCES product_variants(id) ON DELETE CASCADE,
 provider text NOT NULL, external_item_id text NOT NULL, external_model_id text NOT NULL DEFAULT '', external_sku text,
 mapping_source text NOT NULL DEFAULT 'manual', active boolean NOT NULL DEFAULT true,
 created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now(),
 UNIQUE (integration_id, external_item_id, external_model_id)
);
CREATE INDEX idx_marketplace_sku_mappings_tenant_product ON marketplace_sku_mappings (tenant_id, product_id);
CREATE INDEX idx_marketplace_sku_mappings_external_sku ON marketplace_sku_mappings (tenant_id, provider, external_sku) WHERE external_sku IS NOT NULL;
ALTER TABLE marketplace_sku_mappings ENABLE ROW LEVEL SECURITY;
ALTER TABLE marketplace_sku_mappings FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_sku_mappings_tenant_isolation ON marketplace_sku_mappings
USING (tenant_id = current_setting('app.current_tenant_id', true)::uuid)
WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::uuid);
