package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowClouddcnSubnetRequest Request Object
type ShowClouddcnSubnetRequest struct {

	// clouddcn子网ID
	ClouddcnSubnetId string `json:"clouddcn_subnet_id"`
}

func (o ShowClouddcnSubnetRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClouddcnSubnetRequest struct{}"
	}

	return strings.Join([]string{"ShowClouddcnSubnetRequest", string(data)}, " ")
}
