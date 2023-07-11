package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateShrinkClusterResponse Response Object
type UpdateShrinkClusterResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateShrinkClusterResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateShrinkClusterResponse struct{}"
	}

	return strings.Join([]string{"UpdateShrinkClusterResponse", string(data)}, " ")
}
