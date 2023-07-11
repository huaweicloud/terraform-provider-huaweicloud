package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FileLocation 升级包的位置
type FileLocation struct {
	ObsLocation *ObsLocation `json:"obs_location,omitempty"`
}

func (o FileLocation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FileLocation struct{}"
	}

	return strings.Join([]string{"FileLocation", string(data)}, " ")
}
