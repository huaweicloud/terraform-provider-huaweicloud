package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UrlAuth URL鉴权。
type UrlAuth struct {

	// 是否开启URL鉴权，on：开启,off：关闭。
	Status string `json:"status"`

	// 鉴权方式 type_a：鉴权方式A type_b：鉴权方式B type_c1：鉴权方式C1 type_c2：鉴权方式C2
	Type *string `json:"type,omitempty"`

	// 过期时间：范围：0-31536000单位为秒。
	ExpireTime *int32 `json:"expire_time,omitempty"`

	// 加密的算法 可选择md5或sha256。
	SignMethod *string `json:"sign_method,omitempty"`

	// 鉴权范围，目前仅支持配置所有文件参与鉴权，all：所有文件。
	MatchType *string `json:"match_type,omitempty"`

	InheritConfig *InheritConfig `json:"inherit_config,omitempty"`

	// 鉴权KEY 由6-32位大小写字母、数字构成。
	Key *string `json:"key,omitempty"`

	// 鉴权KEY（备） 由6-32位大小写字母、数字构成。
	BackupKey *string `json:"backup_key,omitempty"`

	// 鉴权参数：1-100位可以由大小写字母、数字、下划线构成（不能以数字开头）。
	SignArg *string `json:"sign_arg,omitempty"`

	// 时间格式 dec：十进制 hex：十六进制 鉴权方式A：只支持十进制 鉴权方式B：只支持十进制 鉴权方式C1：只支持十六进制鉴权方式 鉴权方式C2：支持十进制/十六进制
	TimeFormat *string `json:"time_format,omitempty"`

	// 时间参数名称：当type为“type_c2”时，该字段可选。由1-100个字符组成。必须以字母开头，后面可以衔接任意数字、字母和下划线，不允许出现其他特殊符号。鉴权方式为C2时，不传或传空默认设置为“timestamp”。
	TimeArg *string `json:"time_arg,omitempty"`
}

func (o UrlAuth) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UrlAuth struct{}"
	}

	return strings.Join([]string{"UrlAuth", string(data)}, " ")
}
