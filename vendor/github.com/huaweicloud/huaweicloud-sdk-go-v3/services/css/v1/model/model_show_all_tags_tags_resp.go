package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAllTagsTagsResp 集群标签。
type ShowAllTagsTagsResp struct {

	// 集群标签的key值。
	Key *string `json:"key,omitempty"`

	// 集群标签的value值列表。
	Values *[]string `json:"values,omitempty"`
}

func (o ShowAllTagsTagsResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAllTagsTagsResp struct{}"
	}

	return strings.Join([]string{"ShowAllTagsTagsResp", string(data)}, " ")
}
