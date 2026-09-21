CREATE TABLE marketplace_order_inventory_events (id uuid PRIMARY KEY, tenant_id uuid NOT NULL, integration_id uuid NOT NULL REFERENCES integrations(id) ON DELETE CASCADE, external_order_id text NOT NULL, event_type text NOT NULL, created_at timestamptz NOT NULL DEFAULT now(), UNIQUE(integration_id,external_order_id,event_type));
ALTER TABLE marketplace_order_inventory_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE marketplace_order_inventory_events FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_order_inventory_events_tenant_isolation ON marketplace_order_inventory_events USING (tenant_id=current_setting('app.current_tenant_id',true)::uuid) WITH CHECK (tenant_id=current_setting('app.current_tenant_id',true)::uuid);
