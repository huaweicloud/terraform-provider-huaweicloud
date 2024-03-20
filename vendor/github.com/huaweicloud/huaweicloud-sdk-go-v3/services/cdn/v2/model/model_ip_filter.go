package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// IpFilter IP黑白名单。
type IpFilter struct {

	// IP黑白名单类型，off：关闭IP黑白名单，black：IP黑名单，white：IP白名单。
	Type string `json:"type"`

	// 配置IP黑白名单，当type=off时，非必传， 支持IPv6,支持配置IP地址和IP&掩码格式的网段, 多条规则用“,”分割,最多支持配置150个, 多个完全重复的IP/IP段将合并为一个,不支持带通配符的地址，如192.168.0.*。
	Value *string `json:"value,omitempty"`
}

func (o IpFilter) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IpFilter struct{}"
	}

	return strings.Join([]string{"IpFilter", string(data)}, " ")
}
