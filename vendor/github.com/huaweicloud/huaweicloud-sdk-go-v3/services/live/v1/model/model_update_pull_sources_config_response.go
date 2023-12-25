package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePullSourcesConfigResponse Response Object
type UpdatePullSourcesConfigResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdatePullSourcesConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePullSourcesConfigResponse struct{}"
	}

	return strings.Join([]string{"UpdatePullSourcesConfigResponse", string(data)}, " ")
}
