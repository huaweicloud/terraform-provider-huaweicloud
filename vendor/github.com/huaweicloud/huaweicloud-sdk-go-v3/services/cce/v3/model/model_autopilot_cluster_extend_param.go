package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AutopilotClusterExtendParam struct {

	// 集群所属的企业项目ID。 >   - 需要开通企业项目功能后才可配置企业项目。 >   - 集群所属的企业项目与集群下所关联的其他云服务资源所属的企业项目必须保持一致。
	EnterpriseProjectId *string `json:"enterpriseProjectId,omitempty"`

	// 记录集群通过何种升级方式升级到当前版本。
	Upgradefrom *string `json:"upgradefrom,omitempty"`
}

func (o AutopilotClusterExtendParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotClusterExtendParam struct{}"
	}

	return strings.Join([]string{"AutopilotClusterExtendParam", string(data)}, " ")
}
