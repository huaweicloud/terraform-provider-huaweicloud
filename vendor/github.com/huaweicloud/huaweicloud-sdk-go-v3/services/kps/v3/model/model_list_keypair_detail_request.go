package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListKeypairDetailRequest struct {

	// 密钥对名称
	KeypairName string `json:"keypair_name"`
}

func (o ListKeypairDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListKeypairDetailRequest struct{}"
	}

	return strings.Join([]string{"ListKeypairDetailRequest", string(data)}, " ")
}
