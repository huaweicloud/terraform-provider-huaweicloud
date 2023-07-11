package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RefererRsp struct {

	// Referer类型。取值：0代表不设置Referer过滤；1代表黑名单；2代表白名单。默认取值为0。
	RefererType *int32 `json:"referer_type,omitempty"`

	// ：配置的referer地址。
	RefererList *string `json:"referer_list,omitempty"`

	// 是否包含空Referer。如果是黑名单并开启该选项，则表示无referer不允许访问。如果是白名单并开启该选项，则表示无referer允许访问。默认不包含,true：包含，false：不包含。
	IncludeEmpty *bool `json:"include_empty,omitempty"`
}

func (o RefererRsp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RefererRsp struct{}"
	}

	return strings.Join([]string{"RefererRsp", string(data)}, " ")
}
