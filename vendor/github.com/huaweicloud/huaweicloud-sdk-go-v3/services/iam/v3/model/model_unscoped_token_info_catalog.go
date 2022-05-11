package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UnscopedTokenInfoCatalog struct {

	// 终端节点ID。
	Id *string `json:"id,omitempty"`

	// 接口类型，描述接口在该终端节点的可见性。值为“public”，表示该接口为公开接口。
	Interface *string `json:"interface,omitempty"`

	// 终端节点所属区域。
	Region *string `json:"region,omitempty"`

	// 终端节点所属区域ID。
	RegionId *string `json:"region_id,omitempty"`

	// 终端节点的URL。
	Url *string `json:"url,omitempty"`
}

func (o UnscopedTokenInfoCatalog) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnscopedTokenInfoCatalog struct{}"
	}

	return strings.Join([]string{"UnscopedTokenInfoCatalog", string(data)}, " ")
}
