package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListAllObsObjListResponse struct {

	// 返回OBS对象组
	Objects        *[]ObsObject `json:"objects,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ListAllObsObjListResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAllObsObjListResponse struct{}"
	}

	return strings.Join([]string{"ListAllObsObjListResponse", string(data)}, " ")
}
