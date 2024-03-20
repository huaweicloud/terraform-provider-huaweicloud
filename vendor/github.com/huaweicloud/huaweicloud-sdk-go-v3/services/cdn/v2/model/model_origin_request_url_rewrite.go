package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// OriginRequestUrlRewrite 改写回源URL，最多配置20条。
type OriginRequestUrlRewrite struct {

	// 回源URL改写规则的优先级。 优先级设置具有唯一性，不支持多条回源URL改写规则设置同一优先级，且优先级不能输入为空。 多条规则下，不同规则中的相同资源内容，CDN按照优先级高的规则执行URL改写。 取值为1-100之间的整数，数值越大优先级越高。
	Priority int32 `json:"priority"`

	// 匹配类型， all：所有文件， file_path：URL路径， wildcard：通配符， full_path: 全路径。
	MatchType string `json:"match_type"`

	// 需要替换的URI。 改写后的URI以正斜线（/）开头的URI，不含http(s)://头及域名。 长度不超过512个字符。 支持通配符\\*匹配，如：/test/\\*_/\\*.mp4。 匹配方式为“所有文件”时，不支持配置参数。
	SourceUrl *string `json:"source_url,omitempty"`

	// 以正斜线（/）开头的URI，不含http(s)://头及域名。 长度不超过256个字符。 通配符 * 可通过$n捕获（n=1,2,3...，例如：/newtest/$1/$2.jpg）。
	TargetUrl string `json:"target_url"`
}

func (o OriginRequestUrlRewrite) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OriginRequestUrlRewrite struct{}"
	}

	return strings.Join([]string{"OriginRequestUrlRewrite", string(data)}, " ")
}
