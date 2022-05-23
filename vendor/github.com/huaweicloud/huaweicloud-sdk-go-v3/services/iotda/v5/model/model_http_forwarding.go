package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// httpserver配置信息
type HttpForwarding struct {

	// **参数说明**：用于接收满足规则条件数据的http服务器地址。
	Url string `json:"url"`

	// **参数说明**：证书id，请参见[获取证书ID](https://support.huaweicloud.com/usermanual-iothub/iot_01_0001.html#section3)
	CertId *string `json:"cert_id,omitempty"`

	// **参数说明**：当sni_enable为true时，此字段需要填写，内容为将要请求的服务端证书的域名,举例:domain:8443;当sni_enbale为false时，此字段默认不填写。
	CnName *string `json:"cn_name,omitempty"`

	// **参数说明**：需要https服务端和客户端都支持此功能，默认为false，设成true表明Https的客户端在发起请求时，需要携带cn_name；https服务端根据cn_name返回对应的证书；设为false可关闭此功能。
	SniEnable *bool `json:"sni_enable,omitempty"`
}

func (o HttpForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HttpForwarding struct{}"
	}

	return strings.Join([]string{"HttpForwarding", string(data)}, " ")
}
