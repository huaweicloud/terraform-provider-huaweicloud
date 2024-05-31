package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Configs 配置项。
type Configs struct {

	// 业务类型： - web：网站加速； - download：文件下载加速； - video：点播加速。  > 暂不支持“全站加速”变更为其它业务类型。
	BusinessType *string `json:"business_type,omitempty"`

	// 服务区域： - mainland_china：中国大陆； - global：全球； - outside_mainland_china：中国大陆境外。  > 暂不支持“中国大陆”与“中国大陆境外”互相直接切换。
	ServiceArea *string `json:"service_area,omitempty"`

	// 给域名添加备注，字符长度范围0-200。
	Remark *string `json:"remark,omitempty"`

	// 回源请求头改写 该功能将覆盖原有配置（清空之前的配置），在使用此接口时，请上传全量头部信息。
	OriginRequestHeader *[]OriginRequestHeader `json:"origin_request_header,omitempty"`

	// http header配置 该功能将覆盖原有配置（清空之前的配置），在使用此接口时，请上传全量头部信息。
	HttpResponseHeader *[]HttpResponseHeader `json:"http_response_header,omitempty"`

	UrlAuth *UrlAuth `json:"url_auth,omitempty"`

	Https *HttpPutBody `json:"https,omitempty"`

	// 源站配置。
	Sources *[]SourcesConfig `json:"sources,omitempty"`

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

	CacheUrlParameterFilter *CacheUrlParameterFilter `json:"cache_url_parameter_filter,omitempty"`

	// ipv6设置，1：打开；0：关闭。
	Ipv6Accelerate *int32 `json:"ipv6_accelerate,omitempty"`

	// 状态码缓存时间。
	ErrorCodeCache *[]ErrorCodeCache `json:"error_code_cache,omitempty"`

	// Range回源，即分片回源，开启: on，关闭: off。  > 开启Range回源的前提是您的源站支持Range请求，即HTTP请求头中包含Range字段，否则可能导致回源失败。
	OriginRangeStatus *string `json:"origin_range_status,omitempty"`

	UserAgentFilter *UserAgentFilter `json:"user_agent_filter,omitempty"`

	// 改写回源URL，最多配置20条。
	OriginRequestUrlRewrite *[]OriginRequestUrlRewrite `json:"origin_request_url_rewrite,omitempty"`

	// 高级回源，最多配置20条。
	FlexibleOrigin *[]FlexibleOrigins `json:"flexible_origin,omitempty"`

	// 回源是否校验ETag，on：开启，off：关闭。
	SliceEtagStatus *string `json:"slice_etag_status,omitempty"`

	// 回源超时时间，范围:5-60，单位：秒。
	OriginReceiveTimeout *int32 `json:"origin_receive_timeout,omitempty"`

	RemoteAuth *CommonRemoteAuth `json:"remote_auth,omitempty"`

	Websocket *WebSocketSeek `json:"websocket,omitempty"`

	VideoSeek *VideoSeek `json:"video_seek,omitempty"`

	// 请求限速配置。
	RequestLimitRules *[]RequestLimitRules `json:"request_limit_rules,omitempty"`

	IpFrequencyLimit *IpFrequencyLimit `json:"ip_frequency_limit,omitempty"`

	Hsts *Hsts `json:"hsts,omitempty"`

	Quic *Quic `json:"quic,omitempty"`

	// 自定义错误页面。
	ErrorCodeRedirectRules *[]ErrorCodeRedirectRules `json:"error_code_redirect_rules,omitempty"`

	Sni *Sni `json:"sni,omitempty"`

	// 访问URL重写。
	RequestUrlRewrite *[]RequestUrlRewrite `json:"request_url_rewrite,omitempty"`

	// 浏览器缓存过期时间。
	BrowserCacheRules *[]BrowserCacheRules `json:"browser_cache_rules,omitempty"`

	AccessAreaFilter *[]AccessAreaFilter `json:"access_area_filter,omitempty"`
}

func (o Configs) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Configs struct{}"
	}

	return strings.Join([]string{"Configs", string(data)}, " ")
}
