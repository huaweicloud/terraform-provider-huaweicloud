package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowVpcRequest Request Object
type ShowVpcRequest struct {

	// VPC资源ID
	VpcId string `json:"vpc_id"`
}

func (o ShowVpcRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowVpcRequest struct{}"
	}

	return strings.Join([]string{"ShowVpcRequest", string(data)}, " ")
}
