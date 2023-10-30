package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteTrafficMirrorFilterResponse Response Object
type DeleteTrafficMirrorFilterResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteTrafficMirrorFilterResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTrafficMirrorFilterResponse struct{}"
	}

	return strings.Join([]string{"DeleteTrafficMirrorFilterResponse", string(data)}, " ")
}
