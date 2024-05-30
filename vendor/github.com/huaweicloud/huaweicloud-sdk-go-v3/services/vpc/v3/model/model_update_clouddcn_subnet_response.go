package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateClouddcnSubnetResponse Response Object
type UpdateClouddcnSubnetResponse struct {
	ClouddcnSubnet *ClouddcnSubnet `json:"clouddcn_subnet,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o UpdateClouddcnSubnetResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateClouddcnSubnetResponse struct{}"
	}

	return strings.Join([]string{"UpdateClouddcnSubnetResponse", string(data)}, " ")
}
