package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type KeypairBean struct {

	// SSH密钥对名称。
	Name string `json:"name"`
}

func (o KeypairBean) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeypairBean struct{}"
	}

	return strings.Join([]string{"KeypairBean", string(data)}, " ")
}
