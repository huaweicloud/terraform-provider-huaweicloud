package datastores

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type DataStoresResult struct {
	golangsdk.Result
}
type DataStores struct {
	DataStores []dataStores `json:"dataStores" `
}
type dataStores struct {
	Id   string `json:"id" `
	Name string `json:"name"`
}

type DataStoresPage struct {
	pagination.SinglePageBase
}

func (r DataStoresPage) IsEmpty() (bool, error) {
	data, err := ExtractDataStores(r)
	if err != nil {
		return false, err
	}
	return len(data.DataStores) == 0, err
}

func ExtractDataStores(r pagination.Page) (DataStores, error) {
	var s DataStores
	err := (r.(DataStoresPage)).ExtractInto(&s)
	return s, err
}
