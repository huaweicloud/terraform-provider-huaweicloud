package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type DstNodeReq struct {

	// 目的端桶的AK（最大长度100个字符）。
	Ak string `json:"ak"`

	// 目的端桶的SK（最大长度100个字符）。
	Sk string `json:"sk"`

	// 目的端的临时Token（最大长度16384个字符）。
	SecurityToken *string `json:"security_token,omitempty"`

	// 目的端桶的名称。
	Bucket string `json:"bucket"`

	// 目的端桶内路径前缀（拼接在对象key前面,组成新的key,拼接后不能超过1024个字符）。
	SavePrefix *string `json:"save_prefix,omitempty"`

	// 目的端桶所处的区域。  请与Endpoint对应的区域保持一致。
	Region string `json:"region"`
}

func (o DstNodeReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DstNodeReq struct{}"
	}

	return strings.Join([]string{"DstNodeReq", string(data)}, " ")
}
