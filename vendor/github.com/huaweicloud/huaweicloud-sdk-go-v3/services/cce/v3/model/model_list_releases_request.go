package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListReleasesRequest Request Object
type ListReleasesRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 模板ID
	ChartId *string `json:"chart_id,omitempty"`

	// 模板对应的命名空间
	Namespace *string `json:"namespace,omitempty"`
}

func (o ListReleasesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListReleasesRequest struct{}"
	}

	return strings.Join([]string{"ListReleasesRequest", string(data)}, " ")
}
