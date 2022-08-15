package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 修改私有桶开启关闭状态
type UpdatePrivateBucketAccessBody struct {

	// 桶开启关闭状态（true：开启；false：关闭）
	Status *bool `json:"status,omitempty"`
}

func (o UpdatePrivateBucketAccessBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePrivateBucketAccessBody struct{}"
	}

	return strings.Join([]string{"UpdatePrivateBucketAccessBody", string(data)}, " ")
}
