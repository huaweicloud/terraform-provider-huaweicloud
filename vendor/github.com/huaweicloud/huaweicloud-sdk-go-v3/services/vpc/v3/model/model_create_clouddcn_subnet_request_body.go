package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClouddcnSubnetRequestBody 创建clouddcn子网对象
type CreateClouddcnSubnetRequestBody struct {
	ClouddcnSubnet *CreateClouddcnSubnetOption `json:"clouddcn_subnet"`
}

func (o CreateClouddcnSubnetRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClouddcnSubnetRequestBody struct{}"
	}

	return strings.Join([]string{"CreateClouddcnSubnetRequestBody", string(data)}, " ")
}
