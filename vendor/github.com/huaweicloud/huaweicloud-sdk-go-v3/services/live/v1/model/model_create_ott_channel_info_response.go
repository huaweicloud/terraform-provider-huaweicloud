package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateOttChannelInfoResponse Response Object
type CreateOttChannelInfoResponse struct {

	// 错误码
	ResultCode *string `json:"result_code,omitempty"`

	// 错误描述
	ResultMsg *string `json:"result_msg,omitempty"`

	// 推流域名
	Domain *string `json:"domain,omitempty"`

	// 组名或应用名，为必填项
	AppName *string `json:"app_name,omitempty"`

	// 频道ID。频道唯一标识，为必填项
	Id *string `json:"id,omitempty"`

	// 推流URL列表。创建频道时，只有入流协议为RTMP_PUSH时，会返回推流URL列表
	Sources        *[]SourceRsp `json:"sources,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o CreateOttChannelInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateOttChannelInfoResponse struct{}"
	}

	return strings.Join([]string{"CreateOttChannelInfoResponse", string(data)}, " ")
}
