package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CacheUrlParameterFilter 缓存url参数配置。  > 此参数作为旧参数，将于近期下线，建议使用CacheRules设置URL参数。
type CacheUrlParameterFilter struct {

	// 缓存URL参数操作类型： - full_url：缓存所有参数； - ignore_url_params：忽略所有参数； - del_params：忽略指定URL参数； - reserve_params：保留指定URL参数。   > 本接口参数有调整，参数替换如下：   > - del_params替代del_args。   > - reserve_params替代reserve_args。
	Type *string `json:"type,omitempty"`

	// 参数值，多个参数使用分号分隔。
	Value *string `json:"value,omitempty"`
}

func (o CacheUrlParameterFilter) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CacheUrlParameterFilter struct{}"
	}

	return strings.Join([]string{"CacheUrlParameterFilter", string(data)}, " ")
}
