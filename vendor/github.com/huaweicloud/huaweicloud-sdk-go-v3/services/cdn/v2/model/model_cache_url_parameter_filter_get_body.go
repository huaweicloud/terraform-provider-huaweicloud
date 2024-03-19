package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CacheUrlParameterFilterGetBody 缓存url参数配置查询响应体， >  此参数作为旧参数，将于近期下线。
type CacheUrlParameterFilterGetBody struct {

	// 缓存URL参数操作类型： - full_url：缓存所有参数； - ignore_url_params：忽略所有参数； - del_params：忽略指定URL参数； - reserve_params：保留指定URL参数。   > 本接口参数有调整，参数替换如下：   > - del_params替代del_args。   > - reserve_params替代reserve_args。
	Type *string `json:"type,omitempty"`

	// 参数值。
	Value *string `json:"value,omitempty"`
}

func (o CacheUrlParameterFilterGetBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CacheUrlParameterFilterGetBody struct{}"
	}

	return strings.Join([]string{"CacheUrlParameterFilterGetBody", string(data)}, " ")
}
