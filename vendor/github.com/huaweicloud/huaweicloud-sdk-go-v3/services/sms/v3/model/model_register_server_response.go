package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type RegisterServerResponse struct {
	// 源端id

	Id             *string `json:"id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o RegisterServerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RegisterServerResponse struct{}"
	}

	return strings.Join([]string{"RegisterServerResponse", string(data)}, " ")
}
