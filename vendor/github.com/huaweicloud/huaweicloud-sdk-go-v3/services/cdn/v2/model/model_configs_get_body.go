package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ConfigsGetBody 配置项。
type ConfigsGetBody struct {

	// 业务类型： - web：网站加速； - download：文件下载加速； - video：点播加速； - wholesite：全站加速。
	BusinessType *string `json:"business_type,omitempty"`

	// 服务区域： - mainland_china：中国大陆； - global：全球； - outside_mainland_china：中国大陆境外。
	ServiceArea *string `json:"service_area,omitempty"`

	// 域名备注。
	Remark *string `json:"remark,omitempty"`

	// 回源请求头配置
	OriginRequestHeader *[]OriginRequestHeader `json:"origin_request_header,omitempty"`

	// http header配置
	HttpResponseHeader *[]HttpResponseHeader `json:"http_response_header,omitempty"`

	UrlAuth *UrlAuthGetBody `json:"url_auth,omitempty"`

	Https *HttpGetBody `json:"https,omitempty"`

	// 源站配置。
	Sources *[]SourcesConfigResponseBody `json:"sources,omitempty"`

	// 回源协议，follow：协议跟随回源，http：HTTP回源(默认)，https：https回源。
	OriginProtocol *string `json:"origin_protocol,omitempty"`

	// 回源跟随，on：开启，off：关闭。
	OriginFollow302Status *string `json:"origin_follow302_status,omitempty"`

	// 缓存规则。
	CacheRules *[]CacheRules `json:"cache_rules,omitempty"`

	IpFilter *IpFilter `json:"ip_filter,omitempty"`

	Referer *RefererConfig `json:"referer,omitempty"`

	ForceRedirect *ForceRedirectConfig `json:"force_redirect,omitempty"`

	Compress *Compress `json:"compress,omitempty"`

	CacheUrlParameterFilter *CacheUrlParameterFilterGetBody `json:"cache_url_parameter_filter,omitempty"`

	// ipv6设置，1：打开；0：关闭。
	Ipv6Accelerate *int32 `json:"ipv6_accelerate,omitempty"`

	// 状态码缓存时间。
	ErrorCodeCache *[]ErrorCodeCache `json:"error_code_cache,omitempty"`

	// Range回源，开启: on，off:关闭。
	OriginRangeStatus *string `json:"origin_range_status,omitempty"`

	UserAgentFilter *UserAgentFilter `json:"user_agent_filter,omitempty"`

	// 改写回源URL。
	OriginRequestUrlRewrite *[]OriginRequestUrlRewrite `json:"origin_request_url_rewrite,omitempty"`

	// 高级回源。
	FlexibleOrigin *[]FlexibleOrigins `json:"flexible_origin,omitempty"`

	// 回源是否校验ETag，on：开启，off：关闭。
	SliceEtagStatus *string `json:"slice_etag_status,omitempty"`

	// 回源超时时间，单位：秒。
	OriginReceiveTimeout *int32 `json:"origin_receive_timeout,omitempty"`

	RemoteAuth *CommonRemoteAuth `json:"remote_auth,omitempty"`

	Websocket *WebSocketSeek `json:"websocket,omitempty"`

	VideoSeek *VideoSeek `json:"video_seek,omitempty"`

	// 请求限速。
	RequestLimitRules *[]RequestLimitRules `json:"request_limit_rules,omitempty"`

	IpFrequencyLimit *IpFrequencyLimitQuery `json:"ip_frequency_limit,omitempty"`

	Hsts *HstsQuery `json:"hsts,omitempty"`

	Quic *Quic `json:"quic,omitempty"`

	// 自定义错误页面
	ErrorCodeRedirectRules *[]ErrorCodeRedirectRules `json:"error_code_redirect_rules,omitempty"`

	Sni *Sni `json:"sni,omitempty"`

	// 访问URL重写。
	RequestUrlRewrite *[]RequestUrlRewrite `json:"request_url_rewrite,omitempty"`

	// 浏览器缓存过期时间。
	BrowserCacheRules *[]BrowserCacheRules `json:"browser_cache_rules,omitempty"`

	AccessAreaFilter *[]AccessAreaFilter `json:"access_area_filter,omitempty"`
}

func (o ConfigsGetBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfigsGetBody struct{}"
	}

	return strings.Join([]string{"ConfigsGetBody", string(data)}, " ")
}
