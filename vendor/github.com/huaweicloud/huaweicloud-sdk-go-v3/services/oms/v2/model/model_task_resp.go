package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type TaskResp struct {

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

	// 是否启用元数据迁移，默认否。不启用时，为保证迁移任务正常运行，仍将为您迁移ContentType元数据。
	EnableMetadataMigration *bool `json:"enable_metadata_migration,omitempty"`

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

	// 任务状态。 1：等待调度 2：正在执行 3：停止 4：失败 5：成功 7：等待中
	Status *int32 `json:"status,omitempty"`

	// 迁移成功对象数量。
	SuccessfulNum *int64 `json:"successful_num,omitempty"`

	// 任务类型，为空默认设置为object。 list：对象列表迁移 object：文件/文件夹迁移 prefix：对象前缀迁移 url_list: url对象列表
	TaskType *TaskRespTaskType `json:"task_type,omitempty"`

	// 分组类型 NORMAL_TASK：一般迁移任务 SYNC_TASK：同步任务所属迁移任务 GROUP_TASK：任务组所属迁移任务
	GroupType *TaskRespGroupType `json:"group_type,omitempty"`

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

	// 迁移前同名对象覆盖方式，用于迁移前判断源端与目的端有同名对象时，覆盖目的端或跳过迁移。默认SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE。 NO_OVERWRITE：不覆盖。迁移前源端对象与目的端对象同名时，不做对比直接跳过迁移。 SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE：大小/最后修改时间对比覆盖。默认配置。迁移前源端对象与目的端对象同名时，通过对比源端和目的端对象大小和最后修改时间，判断是否覆盖目的端，需满足源端/目的端对象的加密状态一致。源端与目的端同名对象大小不相同，或目的端对象的最后修改时间晚于源端对象的最后修改时间(源端较新)，覆盖目的端。 CRC64_COMPARISON_OVERWRITE：CRC64对比覆盖。目前仅支持华为/阿里/腾讯。迁移前源端对象与目的端对象同名时，通过对比源端和目的端对象元数据中CRC64值是否相同，判断是否覆盖目的端，需满足源端/目的端对象的加密状态一致。如果源端与目的端对象元数据中不存在CRC64值，则系统会默认使用SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE(大小/最后修改时间对比覆盖)来对比进行覆盖判断。 FULL_OVERWRITE：全覆盖。迁移前源端对象与目的端对象同名时，不做对比覆盖目的端。
	ObjectOverwriteMode *TaskRespObjectOverwriteMode `json:"object_overwrite_mode,omitempty"`

	// 目的端存储类型设置，当且仅当目的端为华为云OBS时需要，默认为标准存储 STANDARD：华为云OBS标准存储 IA：华为云OBS低频存储 ARCHIVE：华为云OBS归档存储 DEEP_ARCHIVE：华为云OBS深度归档存储 SRC_STORAGE_MAPPING：保留源端存储类型，将源端存储类型映射为华为云OBS存储类型
	DstStoragePolicy *TaskRespDstStoragePolicy `json:"dst_storage_policy,omitempty"`

	// 一致性校验方式，用于迁移前/后校验对象是否一致，所有校验方式需满足源端/目的端对象的加密状态一致，具体校验方式和校验结果可通过对象列表查看。默认size_last_modified。 size_last_modified：默认配置。迁移前后，通过对比源端和目的端对象大小+最后修改时间，判断对象是否已存在或迁移后数据是否完整。源端与目的端同名对象大小相同，且目的端对象的最后修改时间不早于源端对象的最后修改时间，则代表该对象已存在/迁移成功。 crc64：目前仅支持华为/阿里/腾讯。迁移前后，通过对比源端和目的端对象元数据中CRC64值是否相同，判断对象是否已存在/迁移完成。如果源端与目的端对象元数据中不存在CRC64值，则系统会默认使用大小/最后修改时间校验方式来校验。 no_check：目前仅支持HTTP/HTTPS数据源。当源端对象无法通过标准http协议中content-length字段获取数据大小时，默认数据下载成功即迁移成功，不对数据做额外校验，且迁移时源端对象默认覆盖目的端同名对象。当源端对象能正常通过标准http协议中content-length字段获取数据大小时，则采用大小/最后修改时间校验方式来校验。
	ConsistencyCheck *TaskRespConsistencyCheck `json:"consistency_check,omitempty"`

	// 是否开启请求者付款，在启用后，请求者支付请求和数据传输费用。
	EnableRequesterPays *bool `json:"enable_requester_pays,omitempty"`

	// HIGH：高优先级 MEDIUM：中优先级 LOW：低优先级
	TaskPriority *TaskRespTaskPriority `json:"task_priority,omitempty"`
}

func (o TaskResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskResp struct{}"
	}

	return strings.Join([]string{"TaskResp", string(data)}, " ")
}

type TaskRespTaskType struct {
	value string
}

type TaskRespTaskTypeEnum struct {
	LIST     TaskRespTaskType
	OBJECT   TaskRespTaskType
	PREFIX   TaskRespTaskType
	URL_LIST TaskRespTaskType
}

func GetTaskRespTaskTypeEnum() TaskRespTaskTypeEnum {
	return TaskRespTaskTypeEnum{
		LIST: TaskRespTaskType{
			value: "list",
		},
		OBJECT: TaskRespTaskType{
			value: "object",
		},
		PREFIX: TaskRespTaskType{
			value: "prefix",
		},
		URL_LIST: TaskRespTaskType{
			value: "url_list",
		},
	}
}

func (c TaskRespTaskType) Value() string {
	return c.value
}

func (c TaskRespTaskType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TaskRespTaskType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type TaskRespGroupType struct {
	value string
}

type TaskRespGroupTypeEnum struct {
	NORMAL_TASK TaskRespGroupType
	SYNC_TASK   TaskRespGroupType
	GROUP_TASK  TaskRespGroupType
}

func GetTaskRespGroupTypeEnum() TaskRespGroupTypeEnum {
	return TaskRespGroupTypeEnum{
		NORMAL_TASK: TaskRespGroupType{
			value: "NORMAL_TASK",
		},
		SYNC_TASK: TaskRespGroupType{
			value: "SYNC_TASK",
		},
		GROUP_TASK: TaskRespGroupType{
			value: "GROUP_TASK",
		},
	}
}

func (c TaskRespGroupType) Value() string {
	return c.value
}

func (c TaskRespGroupType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TaskRespGroupType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type TaskRespObjectOverwriteMode struct {
	value string
}

type TaskRespObjectOverwriteModeEnum struct {
	NO_OVERWRITE                            TaskRespObjectOverwriteMode
	SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE TaskRespObjectOverwriteMode
	CRC64_COMPARISON_OVERWRITE              TaskRespObjectOverwriteMode
	FULL_OVERWRITE                          TaskRespObjectOverwriteMode
}

func GetTaskRespObjectOverwriteModeEnum() TaskRespObjectOverwriteModeEnum {
	return TaskRespObjectOverwriteModeEnum{
		NO_OVERWRITE: TaskRespObjectOverwriteMode{
			value: "NO_OVERWRITE",
		},
		SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE: TaskRespObjectOverwriteMode{
			value: "SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE",
		},
		CRC64_COMPARISON_OVERWRITE: TaskRespObjectOverwriteMode{
			value: "CRC64_COMPARISON_OVERWRITE",
		},
		FULL_OVERWRITE: TaskRespObjectOverwriteMode{
			value: "FULL_OVERWRITE",
		},
	}
}

func (c TaskRespObjectOverwriteMode) Value() string {
	return c.value
}

func (c TaskRespObjectOverwriteMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TaskRespObjectOverwriteMode) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type TaskRespDstStoragePolicy struct {
	value string
}

type TaskRespDstStoragePolicyEnum struct {
	STANDARD            TaskRespDstStoragePolicy
	IA                  TaskRespDstStoragePolicy
	ARCHIVE             TaskRespDstStoragePolicy
	DEEP_ARCHIVE        TaskRespDstStoragePolicy
	SRC_STORAGE_MAPPING TaskRespDstStoragePolicy
}

func GetTaskRespDstStoragePolicyEnum() TaskRespDstStoragePolicyEnum {
	return TaskRespDstStoragePolicyEnum{
		STANDARD: TaskRespDstStoragePolicy{
			value: "STANDARD",
		},
		IA: TaskRespDstStoragePolicy{
			value: "IA",
		},
		ARCHIVE: TaskRespDstStoragePolicy{
			value: "ARCHIVE",
		},
		DEEP_ARCHIVE: TaskRespDstStoragePolicy{
			value: "DEEP_ARCHIVE",
		},
		SRC_STORAGE_MAPPING: TaskRespDstStoragePolicy{
			value: "SRC_STORAGE_MAPPING",
		},
	}
}

func (c TaskRespDstStoragePolicy) Value() string {
	return c.value
}

func (c TaskRespDstStoragePolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TaskRespDstStoragePolicy) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type TaskRespConsistencyCheck struct {
	value string
}

type TaskRespConsistencyCheckEnum struct {
	SIZE_LAST_MODIFIED TaskRespConsistencyCheck
	CRC64              TaskRespConsistencyCheck
	NO_CHECK           TaskRespConsistencyCheck
}

func GetTaskRespConsistencyCheckEnum() TaskRespConsistencyCheckEnum {
	return TaskRespConsistencyCheckEnum{
		SIZE_LAST_MODIFIED: TaskRespConsistencyCheck{
			value: "size_last_modified",
		},
		CRC64: TaskRespConsistencyCheck{
			value: "crc64",
		},
		NO_CHECK: TaskRespConsistencyCheck{
			value: "no_check",
		},
	}
}

func (c TaskRespConsistencyCheck) Value() string {
	return c.value
}

func (c TaskRespConsistencyCheck) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TaskRespConsistencyCheck) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type TaskRespTaskPriority struct {
	value string
}

type TaskRespTaskPriorityEnum struct {
	HIGH   TaskRespTaskPriority
	MEDIUM TaskRespTaskPriority
	LOW    TaskRespTaskPriority
}

func GetTaskRespTaskPriorityEnum() TaskRespTaskPriorityEnum {
	return TaskRespTaskPriorityEnum{
		HIGH: TaskRespTaskPriority{
			value: "HIGH",
		},
		MEDIUM: TaskRespTaskPriority{
			value: "MEDIUM",
		},
		LOW: TaskRespTaskPriority{
			value: "LOW",
		},
	}
}

func (c TaskRespTaskPriority) Value() string {
	return c.value
}

func (c TaskRespTaskPriority) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TaskRespTaskPriority) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
