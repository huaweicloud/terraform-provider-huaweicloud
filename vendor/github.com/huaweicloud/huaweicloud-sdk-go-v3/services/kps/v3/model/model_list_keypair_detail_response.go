package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListKeypairDetailResponse struct {
	Keypair        *KeypairDetail `json:"keypair,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ListKeypairDetailResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListKeypairDetailResponse struct{}"
	}

	return strings.Join([]string{"ListKeypairDetailResponse", string(data)}, " ")
}
