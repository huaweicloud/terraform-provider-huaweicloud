package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowClouddcnSubnetResponse Response Object
type ShowClouddcnSubnetResponse struct {
	ClouddcnSubnet *ClouddcnSubnet `json:"clouddcn_subnet,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ShowClouddcnSubnetResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClouddcnSubnetResponse struct{}"
	}

	return strings.Join([]string{"ShowClouddcnSubnetResponse", string(data)}, " ")
}
