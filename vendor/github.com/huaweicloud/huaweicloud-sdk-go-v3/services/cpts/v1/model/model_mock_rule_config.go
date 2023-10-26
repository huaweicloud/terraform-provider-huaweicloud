package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// MockRuleConfig 全链路压测探针Mock规则配置
type MockRuleConfig struct {

	// 服务类型（当前只支持http）
	ServiceType *string `json:"service_type,omitempty"`

	// 请求地址
	RequestUrl *string `json:"request_url,omitempty"`

	// 请求方式（GET/POST/PUT/DELET/PATCH）
	RequestMethod *string `json:"request_method,omitempty"`

	// 重定向地址
	RedirectUrl *string `json:"redirect_url,omitempty"`

	// mock策略（redirect/json）
	MockStrategy *string `json:"mock_strategy,omitempty"`

	// 是否启用
	Enable *bool `json:"enable,omitempty"`

	// 规则名称
	Name *string `json:"name,omitempty"`

	// 项目id
	ProjectId *string `json:"project_id,omitempty"`

	// 规则id
	Id *int32 `json:"id,omitempty"`

	// 自定义响应头
	ResponseHeader map[string]string `json:"response_header,omitempty"`

	// 自定义响应体
	ResponseBody *string `json:"response_body,omitempty"`

	// 自定义响应时长
	ResponseTime *int32 `json:"response_time,omitempty"`

	// 自定义响应码
	ResponseCode *int32 `json:"response_code,omitempty"`
}

func (o MockRuleConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MockRuleConfig struct{}"
	}

	return strings.Join([]string{"MockRuleConfig", string(data)}, " ")
}
