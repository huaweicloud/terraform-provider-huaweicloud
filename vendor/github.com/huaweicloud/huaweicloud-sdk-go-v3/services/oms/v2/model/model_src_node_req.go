package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 源端节点信息。
type SrcNodeReq struct {

	// 源端云服务提供商，task_type为非url_list时，本参数为URLSource。  可选值有AWS、Azure、Aliyun、Tencent、HuaweiCloud、QingCloud、KingsoftCloud、Baidu、Qiniu、URLSource或者UCloud。默认值为Aliyun。
	CloudType *string `json:"cloud_type,omitempty"`

	// 源端桶所处的区域，task_type为非url_list时，本参数为必选。
	Region *string `json:"region,omitempty"`

	// 源端桶的AK（最大长度100个字符），task_type为非url_list时，本参数为必选。
	Ak *string `json:"ak,omitempty"`

	// 源端桶的SK（最大长度100个字符），task_type为非url_list时，本参数为必选。
	Sk *string `json:"sk,omitempty"`

	// 源端桶的临时Token（最大长度16384个字符）
	SecurityToken *string `json:"security_token,omitempty"`

	// 当源端为腾讯云时，需要填写此参数。
	AppId *string `json:"app_id,omitempty"`

	// 源端桶的名称，task_type为非url_list时，本参数为必选。
	Bucket *string `json:"bucket,omitempty"`

	// 任务类型为对象迁移任务时，表示待迁移对象名称（以“/”结尾的字符串代表待迁移的文件夹，非“/”结尾的字符串代表待迁移的文件。）； 任务类型为前缀迁移任务时，表示待迁移前缀。 整桶迁移时，此参数设置为[\"\"]。
	ObjectKey *[]string `json:"object_key,omitempty"`

	ListFile *ListFile `json:"list_file,omitempty"`
}

func (o SrcNodeReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SrcNodeReq struct{}"
	}

	return strings.Join([]string{"SrcNodeReq", string(data)}, " ")
}
