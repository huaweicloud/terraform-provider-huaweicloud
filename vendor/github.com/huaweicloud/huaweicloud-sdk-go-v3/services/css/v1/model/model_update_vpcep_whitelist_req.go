package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateVpcepWhitelistReq struct {

	// 白名单(用户的账号ID)。
	VpcPermissions []string `json:"vpcPermissions"`
}

func (o UpdateVpcepWhitelistReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateVpcepWhitelistReq struct{}"
	}

	return strings.Join([]string{"UpdateVpcepWhitelistReq", string(data)}, " ")
}
