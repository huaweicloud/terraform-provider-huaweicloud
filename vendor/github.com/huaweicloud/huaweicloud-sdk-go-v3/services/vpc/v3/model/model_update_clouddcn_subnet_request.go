package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateClouddcnSubnetRequest Request Object
type UpdateClouddcnSubnetRequest struct {

	// clouddcn子网ID
	ClouddcnSubnetId string `json:"clouddcn_subnet_id"`

	Body *UpdateClouddcnSubnetRequestBody `json:"body,omitempty"`
}

func (o UpdateClouddcnSubnetRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateClouddcnSubnetRequest struct{}"
	}

	return strings.Join([]string{"UpdateClouddcnSubnetRequest", string(data)}, " ")
}
