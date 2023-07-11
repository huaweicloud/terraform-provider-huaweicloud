package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeSecurityGroupRequest Request Object
type ChangeSecurityGroupRequest struct {

	// 待修改安全组的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *ChangeSecurityGroupReq `json:"body,omitempty"`
}

func (o ChangeSecurityGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeSecurityGroupRequest struct{}"
	}

	return strings.Join([]string{"ChangeSecurityGroupRequest", string(data)}, " ")
}
