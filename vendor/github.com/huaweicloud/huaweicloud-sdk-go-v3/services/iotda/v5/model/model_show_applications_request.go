package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowApplicationsRequest struct {

	// **参数说明**：实例ID。物理多租下各实例的唯一标识，一般华为云租户无需携带该参数，仅在物理多租场景下从管理面访问API时需要携带该参数。
	InstanceId *string `json:"Instance-Id,omitempty"`

	// **参数说明**：默认资源空间标识，不携带则查询所有资源空间。 **取值范围**： - true：查询默认资源空间。 - false：查询非默认资源空间。
	DefaultApp *bool `json:"default_app,omitempty"`
}

func (o ShowApplicationsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowApplicationsRequest struct{}"
	}

	return strings.Join([]string{"ShowApplicationsRequest", string(data)}, " ")
}
