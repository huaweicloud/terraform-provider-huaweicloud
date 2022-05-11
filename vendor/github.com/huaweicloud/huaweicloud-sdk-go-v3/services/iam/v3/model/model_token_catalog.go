package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type TokenCatalog struct {

	// 该接口所属服务。
	Type string `json:"type"`

	// 服务ID。
	Id string `json:"id"`

	// 服务名称。
	Name string `json:"name"`

	// 终端节点。
	Endpoints []TokenCatalogEndpoint `json:"endpoints"`
}

func (o TokenCatalog) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TokenCatalog struct{}"
	}

	return strings.Join([]string{"TokenCatalog", string(data)}, " ")
}
