package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BrowserCacheRules 浏览器缓存过期时间。
type BrowserCacheRules struct {
	Condition *BrowserCacheRulesCondition `json:"condition"`

	// 缓存生效类型：   - follow_origin：遵循源站的缓存策略，即Cache-Control头部的设置，   - ttl：浏览器缓存遵循当前规则设置的过期时间，   - never：浏览器不缓存资源。
	CacheType string `json:"cache_type"`

	// 缓存过期时间，最大支持365天。   > 当缓存生效类型为ttl时必填。
	Ttl *int32 `json:"ttl,omitempty"`

	// 缓存过期时间单位，s：秒；m：分种；h：小时；d：天。   > 当缓存生效类型为ttl时必填。
	TtlUnit *string `json:"ttl_unit,omitempty"`
}

func (o BrowserCacheRules) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BrowserCacheRules struct{}"
	}

	return strings.Join([]string{"BrowserCacheRules", string(data)}, " ")
}
