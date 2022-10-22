package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateClustersTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateClustersTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClustersTagsResponse struct{}"
	}

	return strings.Join([]string{"CreateClustersTagsResponse", string(data)}, " ")
}
