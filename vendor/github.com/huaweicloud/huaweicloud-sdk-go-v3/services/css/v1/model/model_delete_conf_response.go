package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteConfResponse Response Object
type DeleteConfResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteConfResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteConfResponse struct{}"
	}

	return strings.Join([]string{"DeleteConfResponse", string(data)}, " ")
}
