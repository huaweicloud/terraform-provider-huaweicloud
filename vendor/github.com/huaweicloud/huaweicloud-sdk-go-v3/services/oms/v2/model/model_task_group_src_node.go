package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TaskGroupSrcNode 迁移任务组的源端节点
type TaskGroupSrcNode struct {

	// 源端桶的AK（最大长度100个字符），task_type为非url_list时，本参数为必选。
	Ak *string `json:"ak,omitempty"`

	// 源端桶的SK（最大长度100个字符），task_type为非url_list时，本参数为必选。
	Sk *string `json:"sk,omitempty"`

	// 用于谷歌云Cloud Storage鉴权
	JsonAuthFile *string `json:"json_auth_file,omitempty"`

	// 当源端为腾讯云时，需要填写此参数。
	AppId *string `json:"app_id,omitempty"`

	// 源端桶所处的区域，task_type为非URL_LIST时，本参数为必选。
	Region *string `json:"region,omitempty"`

	// 任务类型为前缀迁移任务时，表示待迁移前缀。 整桶迁移时，此参数设置为[\"\"]。
	ObjectKey *[]string `json:"object_key,omitempty"`

	// 源端所在桶
	Bucket *string `json:"bucket,omitempty"`

	// 源端云服务提供商，当task_type为URL_LIST时，本参数为URLSource且必选。可选值有AWS、Azure、Aliyun、Tencent、HuaweiCloud、QingCloud、KingsoftCloud、Baidu、Qiniu、URLSource、Google或者UCloud。默认值为Aliyun。
	CloudType *string `json:"cloud_type,omitempty"`

	ListFile *ListFile `json:"list_file,omitempty"`
}

func (o TaskGroupSrcNode) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskGroupSrcNode struct{}"
	}

	return strings.Join([]string{"TaskGroupSrcNode", string(data)}, " ")
}
