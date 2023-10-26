package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteTrafficMirrorSessionResponse Response Object
type DeleteTrafficMirrorSessionResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteTrafficMirrorSessionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTrafficMirrorSessionResponse struct{}"
	}

	return strings.Join([]string{"DeleteTrafficMirrorSessionResponse", string(data)}, " ")
}
