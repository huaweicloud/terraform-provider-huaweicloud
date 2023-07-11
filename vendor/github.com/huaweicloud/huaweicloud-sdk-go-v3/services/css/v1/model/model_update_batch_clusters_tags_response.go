package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateBatchClustersTagsResponse Response Object
type UpdateBatchClustersTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateBatchClustersTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBatchClustersTagsResponse struct{}"
	}

	return strings.Join([]string{"UpdateBatchClustersTagsResponse", string(data)}, " ")
}
