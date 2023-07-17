package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ChangeSecurityGroupReq struct {

	// 期望安全组的ID。
	SecurityGroupIds string `json:"security_group_ids"`
}

func (o ChangeSecurityGroupReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeSecurityGroupReq struct{}"
	}

	return strings.Join([]string{"ChangeSecurityGroupReq", string(data)}, " ")
}
