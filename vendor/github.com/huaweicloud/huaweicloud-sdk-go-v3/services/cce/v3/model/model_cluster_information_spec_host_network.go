package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ClusterInformationSpecHostNetwork 节点网络参数，包含了Node节点默认安群组设置
type ClusterInformationSpecHostNetwork struct {

	// 集群默认Node节点安全组需要放通部分端口来保证正常通信，[详细设置请参考[集群安全组规则配置](https://support.huaweicloud.com/cce_faq/cce_faq_00265.html)。](tag:hws)[详细设置请参考[集群安全组规则配置](https://support.huaweicloud.com/intl/zh-cn/cce_faq/cce_faq_00265.html)。](tag:hws_hk) 修改后的安全组只作用于新创建的节点和新纳管的节点，存量节点的安全组需手动修改。
	SecurityGroup *string `json:"SecurityGroup,omitempty"`
}

func (o ClusterInformationSpecHostNetwork) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterInformationSpecHostNetwork struct{}"
	}

	return strings.Join([]string{"ClusterInformationSpecHostNetwork", string(data)}, " ")
}
