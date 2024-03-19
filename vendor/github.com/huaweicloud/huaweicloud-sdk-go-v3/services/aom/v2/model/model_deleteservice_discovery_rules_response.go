package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteserviceDiscoveryRulesResponse Response Object
type DeleteserviceDiscoveryRulesResponse struct {

	// 响应码。
	ErrorCode *string `json:"errorCode,omitempty"`

	// 响应信息描述。
	ErrorMessage *string `json:"errorMessage,omitempty"`

	// 响应状态码（不再使用）。
	ResponseStatus *int32 `json:"responseStatus,omitempty"`

	// 服务发现规则id列表，多AZ配置同步时使用。
	Id             *[]string `json:"id,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o DeleteserviceDiscoveryRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteserviceDiscoveryRulesResponse struct{}"
	}

	return strings.Join([]string{"DeleteserviceDiscoveryRulesResponse", string(data)}, " ")
}
