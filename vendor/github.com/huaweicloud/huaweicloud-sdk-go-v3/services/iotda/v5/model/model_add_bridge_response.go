package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddBridgeResponse Response Object
type AddBridgeResponse struct {

	// 网桥ID，用于唯一标识一个网桥。在注册网桥时直接指定，或者由物联网平台分配获得。
	BridgeId *string `json:"bridge_id,omitempty"`

	// 网桥名称。
	BridgeName *string `json:"bridge_name,omitempty"`

	AuthInfo *BridgeAuthInfo `json:"auth_info,omitempty"`

	// 在物联网平台注册网桥的时间。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	CreateTime     *string `json:"create_time,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AddBridgeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddBridgeResponse struct{}"
	}

	return strings.Join([]string{"AddBridgeResponse", string(data)}, " ")
}
