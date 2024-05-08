package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RemoteAuthRule 远程鉴权配置。
type RemoteAuthRule struct {

	// 可访问的鉴权服务器地址。 输入的URL必须有“http”或“https”。不能是localhost或127.0.0.1这类本地地址。 不能是CDN的加速域名。
	AuthServer string `json:"auth_server"`

	// 鉴权服务器支持的请求方法，支持GET、POST、HEAD。
	RequestMethod string `json:"request_method"`

	// all（所有文件类型）：所有文件均参与鉴权。 specific_file（指定文件类型）：指定类型的文件参与鉴权。示例：jpg|MP4。 文件类型不区分大小写，即：jpg和JPG代表同一种文件类型，多个文件类型用“|”分割。
	FileTypeSetting string `json:"file_type_setting"`

	// 字符总数不能超过512,当file_type_setting等于specific_file时为必选，其余情况为空， 由大小写字母和数字构成，文件类型用竖线分隔，例如jpg|mp4，只有在必选情况下才会对该字段做校验。
	SpecifiedFileType *string `json:"specified_file_type,omitempty"`

	// 设置用户请求中需要参与鉴权的参数，可选reserve_all_args（保留所有URL参数）、reserve_specific_args（保留指定URL参数）、ignore_all_args（忽略所有URL参数）。
	ReserveArgsSetting string `json:"reserve_args_setting"`

	// 当reserve_args_setting等于reserve_specific_args时为必选，其余情况为空，要保留的参数，多个参数用竖线分隔：key1|key2。
	ReserveArgs *string `json:"reserve_args,omitempty"`

	// URL鉴权参数
	AddCustomArgsRules *[]CustomArgs `json:"add_custom_args_rules,omitempty"`

	// 设置用户请求中参与鉴权请求头，可选reserve_all_headers（保留所有请求头参数）、reserve_specific_headers（保留指定请求头参数）、ignore_all_headers（忽略所有请求头参数）。
	ReserveHeadersSetting string `json:"reserve_headers_setting"`

	// 请求头鉴权参数
	AddCustomHeadersRules *[]CustomArgs `json:"add_custom_headers_rules,omitempty"`

	// 设置鉴权成功时远程鉴权服务器返回给CDN节点的状态码。取值范围：2xx/3xx。
	AuthSuccessStatus string `json:"auth_success_status"`

	// 设置鉴权失败时远程鉴权服务器返回给CDN节点的状态码。取值范围：4xx/5xx。
	AuthFailedStatus string `json:"auth_failed_status"`

	// 设置鉴权失败时CDN节点返回给用户的状态码。取值范围：2xx/3xx/4xx/5xx。
	ResponseStatus string `json:"response_status"`

	// 设置鉴权超时时间，即从CDN转发鉴权请求开始，到CDN节点收到远程鉴权服务器返回的结果的时间。单位为毫秒，值为0或50-3000。
	Timeout int32 `json:"timeout"`

	// 设置鉴权超时后，CDN节点如何处理用户请求。 pass(鉴权失败放过)：鉴权超时后允许用户请求，返回对应的资源。 forbid(鉴权失败拒绝)：鉴权超时后拒绝用户请求，返回配置的响应自定义状态码给用户。
	TimeoutAction string `json:"timeout_action"`

	// 当reserve_headers_setting等于reserve_specific_headers时为必选，其余情况为空，要保留的请求头，多个请求头用竖线分隔：key1|key2。
	ReserveHeaders *string `json:"reserve_headers,omitempty"`
}

func (o RemoteAuthRule) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoteAuthRule struct{}"
	}

	return strings.Join([]string{"RemoteAuthRule", string(data)}, " ")
}
