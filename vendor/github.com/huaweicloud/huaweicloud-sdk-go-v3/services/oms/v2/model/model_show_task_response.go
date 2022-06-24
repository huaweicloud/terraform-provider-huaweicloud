package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Response Object
type ShowTaskResponse struct {

	// 流量控制策略，每个任务最多可设置5条限速策略。
	BandwidthPolicy *[]BandwidthPolicyDto `json:"bandwidth_policy,omitempty"`

	// 任务迁移完成大小（Byte）。
	CompleteSize *int64 `json:"complete_size,omitempty"`

	// 任务描述，没有设置时为空字符串。
	Description *string `json:"description,omitempty"`

	DstNode *DstNodeResp `json:"dst_node,omitempty"`

	// 是否记录失败对象。开启后，如果有迁移失败对象，会在目的端存储失败对象信息。
	EnableFailedObjectRecording *bool `json:"enable_failed_object_recording,omitempty"`

	// 存储入OBS时是否使用KMS加密。
	EnableKms *bool `json:"enable_kms,omitempty"`

	// 是否自动解冻归档数据，（由于对象存储解冻需要源端存储等待一定时间，开启自动解冻会对迁移速度有较大影响，建议先完成归档存储数据解冻后再启动迁移）。 开启后，如果遇到归档类型数据，会自动解冻再进行迁移；如果遇到归档类型的对象直接跳过相应对象，系统默认对象迁移失败并记录相关信息到失败对象列表中。
	EnableRestore *bool `json:"enable_restore,omitempty"`

	ErrorReason *ErrorReasonResp `json:"error_reason,omitempty"`

	// 迁移失败对象数量。
	FailedNum *int64 `json:"failed_num,omitempty"`

	FailedObjectRecord *FailedObjectRecordDto `json:"failed_object_record,omitempty"`

	// 迁移任务组ID，当任务由迁移任务组创建时会包含迁移任务组的id信息。
	GroupId *string `json:"group_id,omitempty"`

	// 任务ID。
	Id *int64 `json:"id,omitempty"`

	// 迁移任务是否完成源端对象统计数据扫描。
	IsQueryOver *bool `json:"is_query_over,omitempty"`

	// 任务剩余时间（毫秒）。
	LeftTime *int64 `json:"left_time,omitempty"`

	// 迁移指定时间（时间戳，毫秒），表示仅迁移在指定时间之后修改的源端待迁移对象。默认为0，表示不设置迁移指定时间。
	MigrateSince *int64 `json:"migrate_since,omitempty"`

	// 任务迁移速度（Byte/s）。
	MigrateSpeed *int64 `json:"migrate_speed,omitempty"`

	// 任务名称。
	Name *string `json:"name,omitempty"`

	// 任务进度，例如：0.522代表任务进度为52.2%，1代表任务进度为100%。
	Progress *float64 `json:"progress,omitempty"`

	// 实际迁移对象总大小（Byte），忽略对象的大小不会统计在内。
	RealSize *int64 `json:"real_size,omitempty"`

	// 迁移忽略对象数（存在以下两种情况会自动跳过：1.源端对象最后修改时间在迁移指定时间前；2.目的端已有该对象。）
	SkippedNum *int64 `json:"skipped_num,omitempty"`

	SrcNode *SrcNodeResp `json:"src_node,omitempty"`

	// 任务启动时间（Unix时间戳，毫秒）。
	StartTime *int64 `json:"start_time,omitempty"`

	// 任务状态。 1：等待调度 2：正在执行 3：停止 4：失败 5：成功
	Status *int32 `json:"status,omitempty"`

	// 迁移成功对象数量。
	SuccessfulNum *int64 `json:"successful_num,omitempty"`

	// 任务类型，为空默认设置为object。 list：对象列表迁移 object：文件/文件夹迁移 prefix：对象前缀迁移 url_list: url对象列表
	TaskType *ShowTaskResponseTaskType `json:"task_type,omitempty"`

	// 分组类型 NORMAL_TASK：一般迁移任务 SYNC_TASK：同步任务所属迁移任务 GROUP_TASK：任务组所属迁移任务
	GroupType *ShowTaskResponseGroupType `json:"group_type,omitempty"`

	// 迁移任务对象总数量。
	TotalNum *int64 `json:"total_num,omitempty"`

	// 任务迁移总大小（Byte）。
	TotalSize *int64 `json:"total_size,omitempty"`

	// 任务总耗时（毫秒）。
	TotalTime *int64 `json:"total_time,omitempty"`

	SmnInfo *SmnInfo `json:"smn_info,omitempty"`

	SourceCdn *SourceCdnResp `json:"source_cdn,omitempty"`

	// 迁移成功对象列表记录失败错误码，记录成功时为空
	SuccessRecordErrorReason *string `json:"success_record_error_reason,omitempty"`

	// 迁移忽略对象列表记录失败错误码,记录记录成功时为空。
	SkipRecordErrorReason *string `json:"skip_record_error_reason,omitempty"`
	HttpStatusCode        int     `json:"-"`
}

func (o ShowTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTaskResponse struct{}"
	}

	return strings.Join([]string{"ShowTaskResponse", string(data)}, " ")
}

type ShowTaskResponseTaskType struct {
	value string
}

type ShowTaskResponseTaskTypeEnum struct {
	LIST     ShowTaskResponseTaskType
	OBJECT   ShowTaskResponseTaskType
	PREFIX   ShowTaskResponseTaskType
	URL_LIST ShowTaskResponseTaskType
}

func GetShowTaskResponseTaskTypeEnum() ShowTaskResponseTaskTypeEnum {
	return ShowTaskResponseTaskTypeEnum{
		LIST: ShowTaskResponseTaskType{
			value: "list",
		},
		OBJECT: ShowTaskResponseTaskType{
			value: "object",
		},
		PREFIX: ShowTaskResponseTaskType{
			value: "prefix",
		},
		URL_LIST: ShowTaskResponseTaskType{
			value: "url_list",
		},
	}
}

func (c ShowTaskResponseTaskType) Value() string {
	return c.value
}

func (c ShowTaskResponseTaskType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTaskResponseTaskType) UnmarshalJSON(b []byte) error {
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

type ShowTaskResponseGroupType struct {
	value string
}

type ShowTaskResponseGroupTypeEnum struct {
	NORMAL_TASK ShowTaskResponseGroupType
	SYNC_TASK   ShowTaskResponseGroupType
	GROUP_TASK  ShowTaskResponseGroupType
}

func GetShowTaskResponseGroupTypeEnum() ShowTaskResponseGroupTypeEnum {
	return ShowTaskResponseGroupTypeEnum{
		NORMAL_TASK: ShowTaskResponseGroupType{
			value: "NORMAL_TASK",
		},
		SYNC_TASK: ShowTaskResponseGroupType{
			value: "SYNC_TASK",
		},
		GROUP_TASK: ShowTaskResponseGroupType{
			value: "GROUP_TASK",
		},
	}
}

func (c ShowTaskResponseGroupType) Value() string {
	return c.value
}

func (c ShowTaskResponseGroupType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTaskResponseGroupType) UnmarshalJSON(b []byte) error {
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
