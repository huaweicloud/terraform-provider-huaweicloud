package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type AddIndependentNodeResponse struct {

	// 集群ID。
	Id             *string `json:"id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AddIndependentNodeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddIndependentNodeResponse struct{}"
	}

	return strings.Join([]string{"AddIndependentNodeResponse", string(data)}, " ")
}
