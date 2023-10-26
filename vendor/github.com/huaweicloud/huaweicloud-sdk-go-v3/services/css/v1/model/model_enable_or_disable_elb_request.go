package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EnableOrDisableElbRequest Request Object
type EnableOrDisableElbRequest struct {

	// 指定待更改的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *UpdateEsElbRequestBody `json:"body,omitempty"`
}

func (o EnableOrDisableElbRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EnableOrDisableElbRequest struct{}"
	}

	return strings.Join([]string{"EnableOrDisableElbRequest", string(data)}, " ")
}
