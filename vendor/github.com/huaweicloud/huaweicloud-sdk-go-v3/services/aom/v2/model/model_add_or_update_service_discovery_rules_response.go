package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type AddOrUpdateServiceDiscoveryRulesResponse struct {

	// 响应码。
	ErrorCode *string `json:"errorCode,omitempty"`

	// 响应信息描述。
	ErrorMessage *string `json:"errorMessage,omitempty"`

	// 响应状态码。
	ResponseStatus *int32 `json:"responseStatus,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o AddOrUpdateServiceDiscoveryRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddOrUpdateServiceDiscoveryRulesResponse struct{}"
	}

	return strings.Join([]string{"AddOrUpdateServiceDiscoveryRulesResponse", string(data)}, " ")
}
