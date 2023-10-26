package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateEsListenerRequestBody struct {
	Listener *EsListenerRequest `json:"listener"`
}

func (o UpdateEsListenerRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateEsListenerRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateEsListenerRequestBody", string(data)}, " ")
}
