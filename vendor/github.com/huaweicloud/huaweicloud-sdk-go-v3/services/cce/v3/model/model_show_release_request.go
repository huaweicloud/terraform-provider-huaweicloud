package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowReleaseRequest Request Object
type ShowReleaseRequest struct {

	// 模板实例名称
	Name string `json:"name"`

	// 模板实例所在的命名空间
	Namespace string `json:"namespace"`

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`
}

func (o ShowReleaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowReleaseRequest struct{}"
	}

	return strings.Join([]string{"ShowReleaseRequest", string(data)}, " ")
}
