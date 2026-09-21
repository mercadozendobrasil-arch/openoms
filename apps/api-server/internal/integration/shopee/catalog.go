package shopee

import (
	"context"
	"fmt"
	"strconv"

	shopeesdk "github.com/openoms-org/openoms/packages/shopee-go-sdk"
)

type ListingSnapshot struct {
	ShopID       string
	ItemID       string
	ItemName     string
	ItemSKU      string
	Status       string
	Models       []ModelSnapshot
}

type ModelSnapshot struct {
	ModelID   string
	ModelSKU  string
	ModelName string
}

func (p *Provider) ReadListingPage(ctx context.Context, offset, pageSize int, status string) ([]ListingSnapshot, bool, int, error) {
	list,err:=p.client.GetItemList(ctx,p.creds.AccessToken,p.creds.ShopID,status,offset,pageSize)
	if err!=nil{return nil,false,offset,fmt.Errorf("shopee: list items: %w",err)}
	if len(list.Item)==0{return nil,list.HasNextPage,list.NextOffset,nil}
	ids:=make([]int64,len(list.Item)); statuses:=make(map[int64]string,len(list.Item))
	for i,item:=range list.Item{ids[i]=item.ItemID;statuses[item.ItemID]=item.ItemStatus}
	base,err:=p.client.GetItemBaseInfo(ctx,p.creds.AccessToken,p.creds.ShopID,ids)
	if err!=nil{return nil,false,offset,fmt.Errorf("shopee: item base info: %w",err)}
	out:=make([]ListingSnapshot,0,len(base))
	for _,item:=range base{
		models,err:=p.client.GetModelList(ctx,p.creds.AccessToken,p.creds.ShopID,item.ItemID)
		if err!=nil{return nil,false,offset,fmt.Errorf("shopee: models for item %d: %w",item.ItemID,err)}
		s:=ListingSnapshot{ShopID:p.creds.ShopID,ItemID:strconv.FormatInt(item.ItemID,10),ItemName:item.ItemName,ItemSKU:item.ItemSKU,Status:statuses[item.ItemID]}
		for _,m:=range models{s.Models=append(s.Models,modelSnapshot(m))}
		out=append(out,s)
	}
	next:=list.NextOffset
	if next==0 && list.HasNextPage{next=offset+pageSize}
	return out,list.HasNextPage,next,nil
}

func modelSnapshot(m shopeesdk.Model) ModelSnapshot {
	return ModelSnapshot{ModelID:strconv.FormatInt(m.ModelID,10),ModelSKU:m.ModelSKU,ModelName:m.ModelName}
}
