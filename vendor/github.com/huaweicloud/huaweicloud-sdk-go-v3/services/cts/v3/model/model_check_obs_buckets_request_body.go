package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CheckObsBucketsRequestBody struct {

	// 请求检查的OBS桶列表。
	Buckets *[]CheckBucketRequest `json:"buckets,omitempty"`
}

func (o CheckObsBucketsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckObsBucketsRequestBody struct{}"
	}

	return strings.Join([]string{"CheckObsBucketsRequestBody", string(data)}, " ")
}
