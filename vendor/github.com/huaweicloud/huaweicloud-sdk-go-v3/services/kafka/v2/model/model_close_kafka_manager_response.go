package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CloseKafkaManagerResponse Response Object
type CloseKafkaManagerResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CloseKafkaManagerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CloseKafkaManagerResponse struct{}"
	}

	return strings.Join([]string{"CloseKafkaManagerResponse", string(data)}, " ")
}
