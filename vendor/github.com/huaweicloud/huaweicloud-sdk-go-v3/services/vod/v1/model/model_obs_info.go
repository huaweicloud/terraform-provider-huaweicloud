package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ObsInfo struct {

	// OBS的bucket名称
	Bucket string `json:"bucket"`

	// OBS对象路径
	Object string `json:"object"`
}

func (o ObsInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ObsInfo struct{}"
	}

	return strings.Join([]string{"ObsInfo", string(data)}, " ")
}
