package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ErrorCodeCache CDN状态码缓存时间。
type ErrorCodeCache struct {

	// 允许配置的错误码: 400, 403, 404, 405, 414, 500, 501, 502, 503, 504
	Code *int32 `json:"code,omitempty"`

	// 错误码缓存时间，单位为秒，范围0-31,536,000(一年默认为365天)。
	Ttl *int32 `json:"ttl,omitempty"`
}

func (o ErrorCodeCache) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ErrorCodeCache struct{}"
	}

	return strings.Join([]string{"ErrorCodeCache", string(data)}, " ")
}
