package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// MasterEipRequestSpecSpec 待绑定的弹性IP配置属性
type MasterEipRequestSpecSpec struct {

	// 弹性网卡ID，绑定时必选，解绑时该字段无效
	Id *string `json:"id,omitempty"`
}

func (o MasterEipRequestSpecSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MasterEipRequestSpecSpec struct{}"
	}

	return strings.Join([]string{"MasterEipRequestSpecSpec", string(data)}, " ")
}
