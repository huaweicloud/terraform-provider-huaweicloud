package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListRdSforMysqlProxyFlavorsResponse Response Object
type ListRdSforMysqlProxyFlavorsResponse struct {

	// 规格组信息。
	ComputeFlavorGroups *[]MysqlProxyFlavorsResponseComputeFlavorGroups `json:"compute_flavor_groups,omitempty"`
	HttpStatusCode      int                                             `json:"-"`
}

func (o ListRdSforMysqlProxyFlavorsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRdSforMysqlProxyFlavorsResponse struct{}"
	}

	return strings.Join([]string{"ListRdSforMysqlProxyFlavorsResponse", string(data)}, " ")
}
