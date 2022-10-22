package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListClustersTagsResponse struct {

	// 集群的标签列表。
	Tags           *[]ShowAllTagsTagsResp `json:"tags,omitempty"`
	HttpStatusCode int                    `json:"-"`
}

func (o ListClustersTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClustersTagsResponse struct{}"
	}

	return strings.Join([]string{"ListClustersTagsResponse", string(data)}, " ")
}
