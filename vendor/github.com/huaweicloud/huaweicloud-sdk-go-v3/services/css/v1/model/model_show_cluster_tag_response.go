package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowClusterTagResponse struct {

	// 集群标签列表。
	Tags           *[]ShowTagsTagsResp `json:"tags,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o ShowClusterTagResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterTagResponse struct{}"
	}

	return strings.Join([]string{"ShowClusterTagResponse", string(data)}, " ")
}
