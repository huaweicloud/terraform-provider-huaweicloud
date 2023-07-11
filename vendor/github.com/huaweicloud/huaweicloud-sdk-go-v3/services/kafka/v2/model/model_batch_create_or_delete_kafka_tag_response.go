package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateOrDeleteKafkaTagResponse Response Object
type BatchCreateOrDeleteKafkaTagResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchCreateOrDeleteKafkaTagResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateOrDeleteKafkaTagResponse struct{}"
	}

	return strings.Join([]string{"BatchCreateOrDeleteKafkaTagResponse", string(data)}, " ")
}
