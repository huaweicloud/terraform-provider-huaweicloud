package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateClouddcnSubnetRequestBody
type UpdateClouddcnSubnetRequestBody struct {
	ClouddcnSubnet *UpdateClouddcnSubnetOption `json:"clouddcn_subnet"`
}

func (o UpdateClouddcnSubnetRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateClouddcnSubnetRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateClouddcnSubnetRequestBody", string(data)}, " ")
}
