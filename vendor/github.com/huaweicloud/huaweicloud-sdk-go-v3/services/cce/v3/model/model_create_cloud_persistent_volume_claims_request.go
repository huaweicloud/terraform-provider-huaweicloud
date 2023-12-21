package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateCloudPersistentVolumeClaimsRequest Request Object
type CreateCloudPersistentVolumeClaimsRequest struct {

	// 指定PersistentVolumeClaim所在的命名空间。  使用namespace有如下约束：  - 用户自定义的namespace，使用前必须先在集群中创建namespace  - 系统自带的namespace：default  - 不能使用kube-system与kube-public
	Namespace string `json:"namespace"`

	// 集群ID，使用**https://Endpoint/uri**这种URL格式时必须指定此参数。获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	XClusterID *string `json:"X-Cluster-ID,omitempty"`

	Body *PersistentVolumeClaim `json:"body,omitempty"`
}

func (o CreateCloudPersistentVolumeClaimsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCloudPersistentVolumeClaimsRequest struct{}"
	}

	return strings.Join([]string{"CreateCloudPersistentVolumeClaimsRequest", string(data)}, " ")
}
