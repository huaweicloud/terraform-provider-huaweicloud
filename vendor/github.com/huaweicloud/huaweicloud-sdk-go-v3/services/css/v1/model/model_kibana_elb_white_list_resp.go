package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// KibanaElbWhiteListResp Kibana公网访问信息。
type KibanaElbWhiteListResp struct {

	// 是否开启kibana访问控制。 - true: 开启访问控制。 - false: 关闭访问控制
	EnableWhiteList *bool `json:"enableWhiteList,omitempty"`

	// kibana公网访问白名单。
	WhiteList *string `json:"whiteList,omitempty"`
}

func (o KibanaElbWhiteListResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KibanaElbWhiteListResp struct{}"
	}

	return strings.Join([]string{"KibanaElbWhiteListResp", string(data)}, " ")
}
