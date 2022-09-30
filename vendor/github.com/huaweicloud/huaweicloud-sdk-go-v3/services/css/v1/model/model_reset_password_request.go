package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ResetPasswordRequest struct {

	// 指定待修改集群密码的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *ResetPasswordReq `json:"body,omitempty"`
}

func (o ResetPasswordRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetPasswordRequest struct{}"
	}

	return strings.Join([]string{"ResetPasswordRequest", string(data)}, " ")
}
