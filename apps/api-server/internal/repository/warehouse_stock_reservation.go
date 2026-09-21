package repository
import("context";"fmt";"github.com/google/uuid";"github.com/jackc/pgx/v5")
func(r *WarehouseStockRepository)AdjustReserved(ctx context.Context,tx pgx.Tx,warehouseID,productID uuid.UUID,variantID *uuid.UUID,delta int)error{
 var tag pgconnCommandTag;var err error
 if variantID!=nil{tag,err=tx.Exec(ctx,"UPDATE warehouse_stock SET reserved=reserved+$1,updated_at=NOW() WHERE warehouse_id=$2 AND product_id=$3 AND variant_id=$4 AND reserved+$1>=0 AND quantity-reserved-$1>=0",delta,warehouseID,productID,*variantID)}else{tag,err=tx.Exec(ctx,"UPDATE warehouse_stock SET reserved=reserved+$1,updated_at=NOW() WHERE warehouse_id=$2 AND product_id=$3 AND variant_id IS NULL AND reserved+$1>=0 AND quantity-reserved-$1>=0",delta,warehouseID,productID)}
 if err!=nil{return fmt.Errorf("adjust reserved stock: %w",err)};if tag.RowsAffected()==0{return fmt.Errorf("insufficient available stock or stock row missing")};return nil
}
type pgconnCommandTag interface{RowsAffected() int64}
