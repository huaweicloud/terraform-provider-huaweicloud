package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateUrlAuthchainResponse Response Object
type CreateUrlAuthchainResponse struct {

	// 生成的鉴权串列表
	Keychain       *[]string `json:"keychain,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o CreateUrlAuthchainResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateUrlAuthchainResponse struct{}"
	}

	return strings.Join([]string{"CreateUrlAuthchainResponse", string(data)}, " ")
}
