package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ConfigsGetBody 配置项。
type ConfigsGetBody struct {

	// 回源请求头配置。
	OriginRequestHeader *[]OriginRequestHeader `json:"origin_request_header,omitempty"`

	// http header配置
	HttpResponseHeader *[]HttpResponseHeader `json:"http_response_header,omitempty"`

	UrlAuth *UrlAuthGetBody `json:"url_auth,omitempty"`

	Https *HttpGetBody `json:"https,omitempty"`

	// 源站配置。
	Sources *[]SourcesConfig `json:"sources,omitempty"`

	// 回源跟随，on：开启，off：关闭。
	OriginFollow302Status *string `json:"origin_follow302_status,omitempty"`

	// 缓存规则。
	CacheRules *[]CacheRules `json:"cache_rules,omitempty"`

	IpFilter *IpFilter `json:"ip_filter,omitempty"`

	Referer *RefererConfig `json:"referer,omitempty"`

	// 回源协议。
	OriginProtocol *string `json:"origin_protocol,omitempty"`

	ForceRedirect *ForceRedirectConfig `json:"force_redirect,omitempty"`

	Compress *Compress `json:"compress,omitempty"`

	CacheUrlParameterFilter *CacheUrlParameterFilter `json:"cache_url_parameter_filter,omitempty"`

	// ipv6设置，1：打开；0：关闭。
	Ipv6Accelerate *int32 `json:"ipv6_accelerate,omitempty"`

	// 状态码缓存时间。
	ErrorCodeCache *[]ErrorCodeCache `json:"error_code_cache,omitempty"`

	// Range回源。
	OriginRangeStatus *string `json:"origin_range_status,omitempty"`

	UserAgentFilter *UserAgentFilter `json:"user_agent_filter,omitempty"`

	// 改写回源URL。
	OriginRequestUrlRewrite *[]OriginRequestUrlRewrite `json:"origin_request_url_rewrite,omitempty"`

	// 自定义错误页面。
	ErrorCodeRedirectRules *[]ErrorCodeRedirectRules `json:"error_code_redirect_rules,omitempty"`
}

func (o ConfigsGetBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfigsGetBody struct{}"
	}

	return strings.Join([]string{"ConfigsGetBody", string(data)}, " ")
}
