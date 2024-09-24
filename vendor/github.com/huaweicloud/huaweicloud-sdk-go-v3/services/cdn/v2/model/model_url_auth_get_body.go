package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UrlAuthGetBody URL鉴权查询响应体。
type UrlAuthGetBody struct {

	// 是否开启URL鉴权，on：开启,off：关闭。
	Status string `json:"status"`

	// 鉴权方式， type_a：鉴权方式A， type_b：鉴权方式B， type_c1：鉴权方式C1， type_c2：鉴权方式C2。
	Type *string `json:"type,omitempty"`

	// 过期时间，单位：秒。
	ExpireTime *int32 `json:"expire_time,omitempty"`

	// 加密算法。
	SignMethod *string `json:"sign_method,omitempty"`

	// 鉴权范围。
	MatchType *string `json:"match_type,omitempty"`

	InheritConfig *InheritConfigQuery `json:"inherit_config,omitempty"`

	// 鉴权KEY。
	Key *string `json:"key,omitempty"`

	// 鉴权KEY（备）。
	BackupKey *string `json:"backup_key,omitempty"`

	// 鉴权参数。
	SignArg *string `json:"sign_arg,omitempty"`

	// 时间格式， dec：十进制, hex：十六进制。
	TimeFormat *string `json:"time_format,omitempty"`

	// 时间参数。
	TimeArg *string `json:"time_arg,omitempty"`
}

func (o UrlAuthGetBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UrlAuthGetBody struct{}"
	}

	return strings.Join([]string{"UrlAuthGetBody", string(data)}, " ")
}
