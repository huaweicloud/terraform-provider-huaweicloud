package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteserviceDiscoveryRulesRequest Request Object
type DeleteserviceDiscoveryRulesRequest struct {

	// 发现规则ID，传多个时以逗号分隔。不允许为空。
	AppRulesIds []string `json:"appRulesIds"`
}

func (o DeleteserviceDiscoveryRulesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteserviceDiscoveryRulesRequest struct{}"
	}

	return strings.Join([]string{"DeleteserviceDiscoveryRulesRequest", string(data)}, " ")
}
