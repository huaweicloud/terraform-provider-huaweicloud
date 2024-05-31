package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteClouddcnSubnetRequest Request Object
type DeleteClouddcnSubnetRequest struct {

	// clouddcn子网ID
	ClouddcnSubnetId string `json:"clouddcn_subnet_id"`
}

func (o DeleteClouddcnSubnetRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteClouddcnSubnetRequest struct{}"
	}

	return strings.Join([]string{"DeleteClouddcnSubnetRequest", string(data)}, " ")
}
