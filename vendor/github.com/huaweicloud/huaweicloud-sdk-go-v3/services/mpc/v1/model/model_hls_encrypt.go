package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type HlsEncrypt struct {

	// 内容加密秘钥
	Key string `json:"key"`

	// 秘钥获取服务的地址
	Url string `json:"url"`

	// 初始向量，base64binary，随机数
	Iv *string `json:"iv,omitempty"`

	// 加密算法。 - AES-128-CTR - AES-128-CBC - SM4CBC  默认值：AES-128-CTR
	Algorithm *string `json:"algorithm,omitempty"`
}

func (o HlsEncrypt) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HlsEncrypt struct{}"
	}

	return strings.Join([]string{"HlsEncrypt", string(data)}, " ")
}
