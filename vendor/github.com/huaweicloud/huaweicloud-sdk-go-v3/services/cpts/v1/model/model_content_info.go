package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ContentInfo struct {
	// body_type

	BodyType *int32 `json:"body_type,omitempty"`
	// bodys

	Bodys *string `json:"bodys,omitempty"`
	// check_end_length

	CheckEndLength *string `json:"check_end_length,omitempty"`
	// check_end_str

	CheckEndStr *string `json:"check_end_str,omitempty"`
	// check_end_type

	CheckEndType *string `json:"check_end_type,omitempty"`
	// connect_timeout

	ConnectTimeout *int32 `json:"connect_timeout,omitempty"`
	// connect_type

	ConnectType *int32 `json:"connect_type,omitempty"`
	// headers

	Headers *[]ContentHeader `json:"headers,omitempty"`
	// http_version

	HttpVersion *string `json:"http_version,omitempty"`
	// method

	Method *string `json:"method,omitempty"`
	// name

	Name *string `json:"name,omitempty"`
	// protocol_type

	ProtocolType *int32 `json:"protocol_type,omitempty"`
	// return_timeout

	ReturnTimeout *int32 `json:"return_timeout,omitempty"`
	// return_timeout_param

	ReturnTimeoutParam *string `json:"return_timeout_param,omitempty"`
	// url

	Url *string `json:"url,omitempty"`
}

func (o ContentInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ContentInfo struct{}"
	}

	return strings.Join([]string{"ContentInfo", string(data)}, " ")
}
