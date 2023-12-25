package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteAddonInstanceRequest Request Object
type DeleteAddonInstanceRequest struct {

	// 插件实例id
	Id string `json:"id"`

	// 集群 ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)
	ClusterId *string `json:"cluster_id,omitempty"`
}

func (o DeleteAddonInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAddonInstanceRequest struct{}"
	}

	return strings.Join([]string{"DeleteAddonInstanceRequest", string(data)}, " ")
}
