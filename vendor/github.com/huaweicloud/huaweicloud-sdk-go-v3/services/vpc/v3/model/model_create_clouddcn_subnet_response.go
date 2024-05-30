package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClouddcnSubnetResponse Response Object
type CreateClouddcnSubnetResponse struct {
	ClouddcnSubnet *ClouddcnSubnet `json:"clouddcn_subnet,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o CreateClouddcnSubnetResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClouddcnSubnetResponse struct{}"
	}

	return strings.Join([]string{"CreateClouddcnSubnetResponse", string(data)}, " ")
}
