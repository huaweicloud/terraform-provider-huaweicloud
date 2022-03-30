package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowGaussMySqlProxyFlavorsResponse struct {
	// 规格组信息。

	ProxyFlavorGroups *[]MysqlProxyFlavorGroups `json:"proxy_flavor_groups,omitempty"`
	HttpStatusCode    int                       `json:"-"`
}

func (o ShowGaussMySqlProxyFlavorsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGaussMySqlProxyFlavorsResponse struct{}"
	}

	return strings.Join([]string{"ShowGaussMySqlProxyFlavorsResponse", string(data)}, " ")
}
