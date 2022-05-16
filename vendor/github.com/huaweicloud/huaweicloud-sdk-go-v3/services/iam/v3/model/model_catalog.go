package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type Catalog struct {

	// 终端节点信息。
	Endpoints []CatalogEndpoints `json:"endpoints"`

	// 服务ID。
	Id string `json:"id"`

	// 服务名。
	Name string `json:"name"`

	// 服务类型。
	Type string `json:"type"`
}

func (o Catalog) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Catalog struct{}"
	}

	return strings.Join([]string{"Catalog", string(data)}, " ")
}
