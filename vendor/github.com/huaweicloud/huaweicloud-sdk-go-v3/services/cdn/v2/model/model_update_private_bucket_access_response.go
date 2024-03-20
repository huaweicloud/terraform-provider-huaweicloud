package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePrivateBucketAccessResponse Response Object
type UpdatePrivateBucketAccessResponse struct {

	// 桶开启关闭状态（true：开启；false：关闭）
	Status *bool `json:"status,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdatePrivateBucketAccessResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePrivateBucketAccessResponse struct{}"
	}

	return strings.Join([]string{"UpdatePrivateBucketAccessResponse", string(data)}, " ")
}
