package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListClouddcnSubnetsResponse Response Object
type ListClouddcnSubnetsResponse struct {

	// clouddcn subnet对象列表
	ClouddcnSubnets *[]ClouddcnSubnet `json:"clouddcn_subnets,omitempty"`
	HttpStatusCode  int               `json:"-"`
}

func (o ListClouddcnSubnetsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClouddcnSubnetsResponse struct{}"
	}

	return strings.Join([]string{"ListClouddcnSubnetsResponse", string(data)}, " ")
}
