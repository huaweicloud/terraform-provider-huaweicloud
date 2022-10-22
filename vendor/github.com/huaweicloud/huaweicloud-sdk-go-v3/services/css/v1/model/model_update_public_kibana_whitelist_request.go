package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdatePublicKibanaWhitelistRequest struct {

	// 指定修改kibana的访问权限的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *UpdatePublicKibanaWhitelistReq `json:"body,omitempty"`
}

func (o UpdatePublicKibanaWhitelistRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePublicKibanaWhitelistRequest struct{}"
	}

	return strings.Join([]string{"UpdatePublicKibanaWhitelistRequest", string(data)}, " ")
}
