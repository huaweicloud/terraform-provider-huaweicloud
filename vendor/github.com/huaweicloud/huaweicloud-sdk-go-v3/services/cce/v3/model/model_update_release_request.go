package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateReleaseRequest Request Object
type UpdateReleaseRequest struct {

	// 模板实例名称
	Name string `json:"name"`

	// 模板实例所在的命名空间
	Namespace string `json:"namespace"`

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	Body *UpdateReleaseReqBody `json:"body,omitempty"`
}

func (o UpdateReleaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateReleaseRequest struct{}"
	}

	return strings.Join([]string{"UpdateReleaseRequest", string(data)}, " ")
}
