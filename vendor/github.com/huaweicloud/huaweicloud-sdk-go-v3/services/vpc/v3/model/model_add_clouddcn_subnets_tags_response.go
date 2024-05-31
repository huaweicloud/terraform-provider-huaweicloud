package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddClouddcnSubnetsTagsResponse Response Object
type AddClouddcnSubnetsTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o AddClouddcnSubnetsTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddClouddcnSubnetsTagsResponse struct{}"
	}

	return strings.Join([]string{"AddClouddcnSubnetsTagsResponse", string(data)}, " ")
}
