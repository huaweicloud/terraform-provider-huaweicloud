package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteClouddcnSubnetsTagResponse Response Object
type DeleteClouddcnSubnetsTagResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteClouddcnSubnetsTagResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteClouddcnSubnetsTagResponse struct{}"
	}

	return strings.Join([]string{"DeleteClouddcnSubnetsTagResponse", string(data)}, " ")
}
