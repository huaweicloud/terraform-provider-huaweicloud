package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateFlavorByTypeResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateFlavorByTypeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateFlavorByTypeResponse struct{}"
	}

	return strings.Join([]string{"UpdateFlavorByTypeResponse", string(data)}, " ")
}
