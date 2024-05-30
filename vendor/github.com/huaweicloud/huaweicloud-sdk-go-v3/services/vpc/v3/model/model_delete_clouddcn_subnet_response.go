package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteClouddcnSubnetResponse Response Object
type DeleteClouddcnSubnetResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteClouddcnSubnetResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteClouddcnSubnetResponse struct{}"
	}

	return strings.Join([]string{"DeleteClouddcnSubnetResponse", string(data)}, " ")
}
