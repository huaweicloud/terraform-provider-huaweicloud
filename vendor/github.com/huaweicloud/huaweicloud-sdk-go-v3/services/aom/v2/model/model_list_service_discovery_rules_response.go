package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListServiceDiscoveryRulesResponse struct {

	// 查询结果规则信息。
	AppRules *[]AppRules `json:"appRules,omitempty"`

	// 响应码,AOM_INVENTORY_2000000代表正常返回。
	ErrorCode *string `json:"errorCode,omitempty"`

	// 响应信息描述。
	ErrorMessage   *string `json:"errorMessage,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListServiceDiscoveryRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListServiceDiscoveryRulesResponse struct{}"
	}

	return strings.Join([]string{"ListServiceDiscoveryRulesResponse", string(data)}, " ")
}
