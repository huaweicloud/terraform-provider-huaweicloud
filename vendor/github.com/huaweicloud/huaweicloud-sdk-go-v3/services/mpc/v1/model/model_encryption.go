package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Encryption struct {
	HlsEncrypt *HlsEncrypt `json:"hls_encrypt,omitempty"`
}

func (o Encryption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Encryption struct{}"
	}

	return strings.Join([]string{"Encryption", string(data)}, " ")
}
