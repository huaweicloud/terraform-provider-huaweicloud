package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 失败原因。如果集群处于正常状态，则不返回该参数。
type ClusterListFailedReasons struct {

	// 错误码。  - CSS.6000：表示集群创建失败。 - CSS.6001：表示集群扩容失败。 - CSS.6002：表示集群重启失败。 - CSS.6004：表示集群节点创建失败。 - CSS.6005：表示服务初始化失败。
	ErrorCode *string `json:"errorCode,omitempty"`

	// 详细错误信息。
	ErrorMsg *string `json:"errorMsg,omitempty"`
}

func (o ClusterListFailedReasons) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterListFailedReasons struct{}"
	}

	return strings.Join([]string{"ClusterListFailedReasons", string(data)}, " ")
}
