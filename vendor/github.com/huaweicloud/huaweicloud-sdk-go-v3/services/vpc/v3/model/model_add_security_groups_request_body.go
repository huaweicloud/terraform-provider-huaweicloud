package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddSecurityGroupsRequestBody This is a auto create Body Object
type AddSecurityGroupsRequestBody struct {
	Port *InsertSecurityGroupOption `json:"port"`
}

func (o AddSecurityGroupsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddSecurityGroupsRequestBody struct{}"
	}

	return strings.Join([]string{"AddSecurityGroupsRequestBody", string(data)}, " ")
}
