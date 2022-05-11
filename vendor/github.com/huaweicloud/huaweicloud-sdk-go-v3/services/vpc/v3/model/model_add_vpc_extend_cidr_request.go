package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type AddVpcExtendCidrRequest struct {

	// VPC资源ID
	VpcId string `json:"vpc_id"`

	Body *AddVpcExtendCidrRequestBody `json:"body,omitempty"`
}

func (o AddVpcExtendCidrRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddVpcExtendCidrRequest struct{}"
	}

	return strings.Join([]string{"AddVpcExtendCidrRequest", string(data)}, " ")
}
