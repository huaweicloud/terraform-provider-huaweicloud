package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteClustersTagsResponse Response Object
type DeleteClustersTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteClustersTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteClustersTagsResponse struct{}"
	}

	return strings.Join([]string{"DeleteClustersTagsResponse", string(data)}, " ")
}
