package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateClusterNameResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateClusterNameResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateClusterNameResponse struct{}"
	}

	return strings.Join([]string{"UpdateClusterNameResponse", string(data)}, " ")
}
