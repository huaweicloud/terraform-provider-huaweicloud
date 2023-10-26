package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CheckObsBucketsResponse Response Object
type CheckObsBucketsResponse struct {

	// 检查OBS桶状态响应体。
	Buckets        *[]Bucket `json:"buckets,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o CheckObsBucketsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckObsBucketsResponse struct{}"
	}

	return strings.Join([]string{"CheckObsBucketsResponse", string(data)}, " ")
}
