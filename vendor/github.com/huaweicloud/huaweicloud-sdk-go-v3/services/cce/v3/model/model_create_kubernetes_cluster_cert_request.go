package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateKubernetesClusterCertRequest Request Object
type CreateKubernetesClusterCertRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	Body *CertDuration `json:"body,omitempty"`
}

func (o CreateKubernetesClusterCertRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateKubernetesClusterCertRequest struct{}"
	}

	return strings.Join([]string{"CreateKubernetesClusterCertRequest", string(data)}, " ")
}
