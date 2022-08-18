package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 缓存url参数配置
type CacheUrlParameterFilter struct {

	// 缓存URL参数操作类型（full_url：缓存所有参数，ignore_url_params：忽略所有参数，del_args：忽略指定URL参数，reserve_args：保留指定URL参数）
	Type *string `json:"type,omitempty"`

	// 参数值，多个参数使用分号分隔
	Value *string `json:"value,omitempty"`
}

func (o CacheUrlParameterFilter) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CacheUrlParameterFilter struct{}"
	}

	return strings.Join([]string{"CacheUrlParameterFilter", string(data)}, " ")
}
