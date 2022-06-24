package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// This is a auto create Body Object
type CreateTaskReq struct {

	// 任务类型，默认为object。  list：对象列表迁移 url_list：URL列表迁移 object：文件/文件夹迁移，默认 prefix：对象前缀迁移
	TaskType *CreateTaskReqTaskType `json:"task_type,omitempty"`

	SrcNode *SrcNodeReq `json:"src_node"`

	DstNode *DstNodeReq `json:"dst_node"`

	// 是否开启KMS加密，默认不开启。
	EnableKms *bool `json:"enable_kms,omitempty"`

	// 任务描述，不能超过255个字符，且不能包含^<>&\"'等特殊字符。
	Description *string `json:"description,omitempty"`

	// 以时间戳方式表示的迁移指定时间（单位：秒），表示仅迁移在指定时间之后修改的源端待迁移对象。默认不设置迁移指定时间。
	MigrateSince *int64 `json:"migrate_since,omitempty"`

	// 配置流量控制策略。数组中一个元素对应一个时段的最大带宽，最多允许5个时段，且时段不能重叠。
	BandwidthPolicy *[]BandwidthPolicyDto `json:"bandwidth_policy,omitempty"`

	SourceCdn *SourceCdnReq `json:"source_cdn,omitempty"`

	SmnConfig *SmnConfig `json:"smn_config,omitempty"`

	// 是否自动解冻归档数据，默认否。  开启后，如果遇到归档类型数据，会自动解冻再进行迁移。
	EnableRestore *bool `json:"enable_restore,omitempty"`

	// 是否记录失败对象，默认开启。  开启后，如果有迁移失败对象，会在目的端存储失败对象信息。
	EnableFailedObjectRecording *bool `json:"enable_failed_object_recording,omitempty"`
}

func (o CreateTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTaskReq struct{}"
	}

	return strings.Join([]string{"CreateTaskReq", string(data)}, " ")
}

type CreateTaskReqTaskType struct {
	value string
}

type CreateTaskReqTaskTypeEnum struct {
	LIST     CreateTaskReqTaskType
	URL_LIST CreateTaskReqTaskType
	OBJECT   CreateTaskReqTaskType
	PREFIX   CreateTaskReqTaskType
}

func GetCreateTaskReqTaskTypeEnum() CreateTaskReqTaskTypeEnum {
	return CreateTaskReqTaskTypeEnum{
		LIST: CreateTaskReqTaskType{
			value: "list",
		},
		URL_LIST: CreateTaskReqTaskType{
			value: "url_list",
		},
		OBJECT: CreateTaskReqTaskType{
			value: "object",
		},
		PREFIX: CreateTaskReqTaskType{
			value: "prefix",
		},
	}
}

func (c CreateTaskReqTaskType) Value() string {
	return c.value
}

func (c CreateTaskReqTaskType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateTaskReqTaskType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
