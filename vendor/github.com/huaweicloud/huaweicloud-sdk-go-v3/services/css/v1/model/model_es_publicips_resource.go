package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type EsPublicipsResource struct {

	// 弹性公网ip配置id。
	PublicipId *string `json:"publicip_id,omitempty"`

	// IP地址。
	PublicipAddress *string `json:"publicip_address,omitempty"`

	// IP版本信息。 - 4：表示IPv4。 - 6：表示IPv6。
	IpVersion *string `json:"ip_version,omitempty"`
}

func (o EsPublicipsResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EsPublicipsResource struct{}"
	}

	return strings.Join([]string{"EsPublicipsResource", string(data)}, " ")
}
