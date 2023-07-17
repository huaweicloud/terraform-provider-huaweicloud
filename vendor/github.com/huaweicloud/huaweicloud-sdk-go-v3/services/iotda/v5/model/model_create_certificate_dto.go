package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateCertificateDto 创建CA证书结构体。
type CreateCertificateDto struct {

	// 证书内容信息。
	Content string `json:"content"`

	// 资源空间ID。此参数为非必选参数，存在多资源空间的用户需要使用该接口时，建议携带该参数指定创建的证书归属到哪个资源空间下，否则创建的证书将会归属到[[默认资源空间](https://support.huaweicloud.com/usermanual-iothub/iot_01_0006.html#section0)](tag:hws)[[默认资源空间](https://support.huaweicloud.com/intl/zh-cn/usermanual-iothub/iot_01_0006.html#section0)](tag:hws_hk)下。
	AppId *string `json:"app_id,omitempty"`
}

func (o CreateCertificateDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCertificateDto struct{}"
	}

	return strings.Join([]string{"CreateCertificateDto", string(data)}, " ")
}
