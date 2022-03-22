package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowGaussMySqlFlavorsResponse struct {
	// 实例规格信息列表

	Flavors        *[]MysqlFlavorsInfo `json:"flavors,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o ShowGaussMySqlFlavorsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGaussMySqlFlavorsResponse struct{}"
	}

	return strings.Join([]string{"ShowGaussMySqlFlavorsResponse", string(data)}, " ")
}
