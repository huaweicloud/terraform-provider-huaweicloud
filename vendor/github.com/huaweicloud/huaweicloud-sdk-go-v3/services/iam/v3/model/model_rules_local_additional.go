package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type RulesLocalAdditional struct {

	// user：联邦用户在本系统中的用户名称。 ``` \"user\":{\"name\":\"{0}\"} ```  group：联邦用户在本系统中所属用户组。 ``` \"group\":{\"name\":\"0cd5e9\"} ```
	Name *string `json:"name,omitempty"`
}

func (o RulesLocalAdditional) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RulesLocalAdditional struct{}"
	}

	return strings.Join([]string{"RulesLocalAdditional", string(data)}, " ")
}
