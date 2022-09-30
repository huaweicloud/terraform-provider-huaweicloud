package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// kibana白名单信息。
type CreateClusterPublicKibanaElbWhiteList struct {

	// 白名单。需要添加白名单的网段或ip，以逗号隔开，不可重复。
	WhiteList string `json:"whiteList"`

	// 是否开启kibana访问控制。
	EnableWhiteList bool `json:"enableWhiteList"`
}

func (o CreateClusterPublicKibanaElbWhiteList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterPublicKibanaElbWhiteList struct{}"
	}

	return strings.Join([]string{"CreateClusterPublicKibanaElbWhiteList", string(data)}, " ")
}
