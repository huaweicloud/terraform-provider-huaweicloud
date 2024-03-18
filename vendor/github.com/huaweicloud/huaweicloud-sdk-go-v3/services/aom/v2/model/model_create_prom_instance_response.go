package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePromInstanceResponse Response Object
type CreatePromInstanceResponse struct {

	// Prometheus实例名称列表。
	Prometheus     *[]PromInstanceEpsCreateModel `json:"prometheus,omitempty"`
	HttpStatusCode int                           `json:"-"`
}

func (o CreatePromInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePromInstanceResponse struct{}"
	}

	return strings.Join([]string{"CreatePromInstanceResponse", string(data)}, " ")
}
