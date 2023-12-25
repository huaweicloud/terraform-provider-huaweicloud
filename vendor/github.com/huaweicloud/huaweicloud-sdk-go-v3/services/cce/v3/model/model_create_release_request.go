package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateReleaseRequest Request Object
type CreateReleaseRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	Body *CreateReleaseReqBody `json:"body,omitempty"`
}

func (o CreateReleaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateReleaseRequest struct{}"
	}

	return strings.Join([]string{"CreateReleaseRequest", string(data)}, " ")
}
