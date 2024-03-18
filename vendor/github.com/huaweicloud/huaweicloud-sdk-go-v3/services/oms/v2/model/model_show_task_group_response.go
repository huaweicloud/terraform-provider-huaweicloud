package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowTaskGroupResponse Response Object
type ShowTaskGroupResponse struct {

	// 任务组id
	GroupId *string `json:"group_id,omitempty"`

	// 迁移组任务状态。 0 – 等待中 1 – 执行中/创建中 2 – 监控任务执行 3 – 暂停 4 – 创建任务失败 5 – 迁移失败 6 – 迁移完成 7 – 暂停中 8 – 等待删除中 9 – 删除
	Status *int32 `json:"status,omitempty"`

	ErrorReason *ErrorReasonResp `json:"error_reason,omitempty"`

	SrcNode *TaskGroupSrcNodeResp `json:"src_node,omitempty"`

	// 任务描述，不能超过255个字符，且不能包含^<>&\"'等特殊字符。
	Description *string `json:"description,omitempty"`

	DstNode *TaskGroupDstNodeResp `json:"dst_node,omitempty"`

	// 是否启用元数据迁移，默认否。不启用时，为保证迁移任务正常运行，仍将为您迁移ContentType元数据。
	EnableMetadataMigration *bool `json:"enable_metadata_migration,omitempty"`

	// 是否开启记录失败对象
	EnableFailedObjectRecording *bool `json:"enable_failed_object_recording,omitempty"`

	// 是否自动解冻归档数据，（由于对象存储解冻需要源端存储等待一定时间，开启自动解冻会对迁移速度有较大影响，建议先完成归档存储数据解冻后再启动迁移）。 开启后，如果遇到归档类型数据，会自动解冻再进行迁移；如果遇到归档类型的对象直接跳过相应对象，系统默认对象迁移失败并记录相关信息到失败对象列表中。
	EnableRestore *bool `json:"enable_restore,omitempty"`

	// 存储入OBS时是否使用KMS加密。
	EnableKms *bool `json:"enable_kms,omitempty"`

	// 任务类型，默认为PREFIX。 LIST：对象列表迁移 URL_LIST：URL列表迁移， PREFIX：对象前缀迁移
	TaskType *ShowTaskGroupResponseTaskType `json:"task_type,omitempty"`

	// 配置流量控制策略。数组中一个元素对应一个时段的最大带宽，最多允许5个时段，且时段不能重叠。
	BandwidthPolicy *[]BandwidthPolicyDto `json:"bandwidth_policy,omitempty"`

	SmnConfig *SmnInfo `json:"smn_config,omitempty"`

	SourceCdn *SourceCdnResp `json:"source_cdn,omitempty"`

	// 迁移指定时间（时间戳，毫秒），表示仅迁移在指定时间之后修改的源端待迁移对象。默认为0，表示不设置迁移指定时间。
	MigrateSince *int64 `json:"migrate_since,omitempty"`

	// 任务组迁移速度（Byte/s）
	MigrateSpeed *int64 `json:"migrate_speed,omitempty"`

	// 迁移任务组总耗时(毫秒)
	TotalTime *int64 `json:"total_time,omitempty"`

	// 迁移任务组的启动时间(Unix时间戳，毫秒)
	StartTime *int64 `json:"start_time,omitempty"`

	// 任务组包含的迁移任务总数
	TotalTaskNum *int64 `json:"total_task_num,omitempty"`

	// 已创建的迁移任务数
	CreateTaskNum *int64 `json:"create_task_num,omitempty"`

	// 失败的迁移任务数
	FailedTaskNum *int64 `json:"failed_task_num,omitempty"`

	// 已完成的迁移任务数
	CompleteTaskNum *int64 `json:"complete_task_num,omitempty"`

	// 暂停的迁移任务数
	PausedTaskNum *int64 `json:"paused_task_num,omitempty"`

	// 正在运行的迁移任务数
	ExecutingTaskNum *int64 `json:"executing_task_num,omitempty"`

	// 等待中的迁移任务数
	WaitingTaskNum *int64 `json:"waiting_task_num,omitempty"`

	// 迁移任务组包含的对象总数量
	TotalNum *int64 `json:"total_num,omitempty"`

	// 已完成任务创建的对象总数量
	CreateCompleteNum *int64 `json:"create_complete_num,omitempty"`

	// 成功的对象数量
	SuccessNum *int64 `json:"success_num,omitempty"`

	// 失败的对象数量
	FailNum *int64 `json:"fail_num,omitempty"`

	// 忽略的对象数量
	SkipNum *int64 `json:"skip_num,omitempty"`

	// 任务迁移总大小(Byte)
	TotalSize *int64 `json:"total_size,omitempty"`

	// 已创建迁移任务包含的对象总大小(Byte)
	CreateCompleteSize *int64 `json:"create_complete_size,omitempty"`

	// 已迁移成功的对象总大小(Byte)
	CompleteSize *int64 `json:"complete_size,omitempty"`

	FailedObjectRecord *FailedObjectRecordDto `json:"failed_object_record,omitempty"`

	// 迁移前同名对象覆盖方式，用于迁移前判断源端与目的端有同名对象时，覆盖目的端或跳过迁移。默认SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE。 NO_OVERWRITE：不覆盖。迁移前源端对象与目的端对象同名时，不做对比直接跳过迁移。 SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE：大小/最后修改时间对比覆盖。默认配置。迁移前源端对象与目的端对象同名时，通过对比源端和目的端对象大小和最后修改时间，判断是否覆盖目的端，需满足源端/目的端对象的加密状态一致。源端与目的端同名对象大小不相同，或目的端对象的最后修改时间晚于源端对象的最后修改时间(源端较新)，覆盖目的端。 CRC64_COMPARISON_OVERWRITE：CRC64对比覆盖。目前仅支持华为/阿里/腾讯。迁移前源端对象与目的端对象同名时，通过对比源端和目的端对象元数据中CRC64值是否相同，判断是否覆盖目的端，需满足源端/目的端对象的加密状态一致。如果源端与目的端对象元数据中不存在CRC64值，则系统会默认使用SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE(大小/最后修改时间对比覆盖)来对比进行覆盖判断。 FULL_OVERWRITE：全覆盖。迁移前源端对象与目的端对象同名时，不做对比覆盖目的端。
	ObjectOverwriteMode *ShowTaskGroupResponseObjectOverwriteMode `json:"object_overwrite_mode,omitempty"`

	// 目的端存储类型设置，当且仅当目的端为华为云OBS时需要，默认为标准存储 STANDARD：华为云OBS标准存储 IA：华为云OBS低频存储 ARCHIVE：华为云OBS归档存储 DEEP_ARCHIVE：华为云OBS深度归档存储 SRC_STORAGE_MAPPING：保留源端存储类型，将源端存储类型映射为华为云OBS存储类型
	DstStoragePolicy *ShowTaskGroupResponseDstStoragePolicy `json:"dst_storage_policy,omitempty"`

	// 一致性校验方式，用于迁移前/后校验对象是否一致，所有校验方式需满足源端/目的端对象的加密状态一致，具体校验方式和校验结果可通过对象列表查看。默认size_last_modified。 size_last_modified：默认配置。迁移前后，通过对比源端和目的端对象大小+最后修改时间，判断对象是否已存在或迁移后数据是否完整。源端与目的端同名对象大小相同，且目的端对象的最后修改时间不早于源端对象的最后修改时间，则代表该对象已存在/迁移成功。 crc64：目前仅支持华为/阿里/腾讯。迁移前后，通过对比源端和目的端对象元数据中CRC64值是否相同，判断对象是否已存在/迁移完成。如果源端与目的端对象元数据中不存在CRC64值，则系统会默认使用大小/最后修改时间校验方式来校验。 no_check：目前仅支持HTTP/HTTPS数据源。当源端对象无法通过标准http协议中content-length字段获取数据大小时，默认数据下载成功即迁移成功，不对数据做额外校验，且迁移时源端对象默认覆盖目的端同名对象。当源端对象能正常通过标准http协议中content-length字段获取数据大小时，则采用大小/最后修改时间校验方式来校验。
	ConsistencyCheck *ShowTaskGroupResponseConsistencyCheck `json:"consistency_check,omitempty"`

	// 是否开启请求者付款，在启用后，请求者支付请求和数据传输费用。
	EnableRequesterPays *bool `json:"enable_requester_pays,omitempty"`
	HttpStatusCode      int   `json:"-"`
}

func (o ShowTaskGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTaskGroupResponse struct{}"
	}

	return strings.Join([]string{"ShowTaskGroupResponse", string(data)}, " ")
}

type ShowTaskGroupResponseTaskType struct {
	value string
}

type ShowTaskGroupResponseTaskTypeEnum struct {
	LIST     ShowTaskGroupResponseTaskType
	URL_LIST ShowTaskGroupResponseTaskType
	PREFIX   ShowTaskGroupResponseTaskType
}

func GetShowTaskGroupResponseTaskTypeEnum() ShowTaskGroupResponseTaskTypeEnum {
	return ShowTaskGroupResponseTaskTypeEnum{
		LIST: ShowTaskGroupResponseTaskType{
			value: "LIST",
		},
		URL_LIST: ShowTaskGroupResponseTaskType{
			value: "URL_LIST",
		},
		PREFIX: ShowTaskGroupResponseTaskType{
			value: "PREFIX",
		},
	}
}

func (c ShowTaskGroupResponseTaskType) Value() string {
	return c.value
}

func (c ShowTaskGroupResponseTaskType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTaskGroupResponseTaskType) UnmarshalJSON(b []byte) error {
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

type ShowTaskGroupResponseObjectOverwriteMode struct {
	value string
}

type ShowTaskGroupResponseObjectOverwriteModeEnum struct {
	NO_OVERWRITE                            ShowTaskGroupResponseObjectOverwriteMode
	SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE ShowTaskGroupResponseObjectOverwriteMode
	CRC64_COMPARISON_OVERWRITE              ShowTaskGroupResponseObjectOverwriteMode
	FULL_OVERWRITE                          ShowTaskGroupResponseObjectOverwriteMode
}

func GetShowTaskGroupResponseObjectOverwriteModeEnum() ShowTaskGroupResponseObjectOverwriteModeEnum {
	return ShowTaskGroupResponseObjectOverwriteModeEnum{
		NO_OVERWRITE: ShowTaskGroupResponseObjectOverwriteMode{
			value: "NO_OVERWRITE",
		},
		SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE: ShowTaskGroupResponseObjectOverwriteMode{
			value: "SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE",
		},
		CRC64_COMPARISON_OVERWRITE: ShowTaskGroupResponseObjectOverwriteMode{
			value: "CRC64_COMPARISON_OVERWRITE",
		},
		FULL_OVERWRITE: ShowTaskGroupResponseObjectOverwriteMode{
			value: "FULL_OVERWRITE",
		},
	}
}

func (c ShowTaskGroupResponseObjectOverwriteMode) Value() string {
	return c.value
}

func (c ShowTaskGroupResponseObjectOverwriteMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTaskGroupResponseObjectOverwriteMode) UnmarshalJSON(b []byte) error {
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

type ShowTaskGroupResponseDstStoragePolicy struct {
	value string
}

type ShowTaskGroupResponseDstStoragePolicyEnum struct {
	STANDARD            ShowTaskGroupResponseDstStoragePolicy
	IA                  ShowTaskGroupResponseDstStoragePolicy
	ARCHIVE             ShowTaskGroupResponseDstStoragePolicy
	DEEP_ARCHIVE        ShowTaskGroupResponseDstStoragePolicy
	SRC_STORAGE_MAPPING ShowTaskGroupResponseDstStoragePolicy
}

func GetShowTaskGroupResponseDstStoragePolicyEnum() ShowTaskGroupResponseDstStoragePolicyEnum {
	return ShowTaskGroupResponseDstStoragePolicyEnum{
		STANDARD: ShowTaskGroupResponseDstStoragePolicy{
			value: "STANDARD",
		},
		IA: ShowTaskGroupResponseDstStoragePolicy{
			value: "IA",
		},
		ARCHIVE: ShowTaskGroupResponseDstStoragePolicy{
			value: "ARCHIVE",
		},
		DEEP_ARCHIVE: ShowTaskGroupResponseDstStoragePolicy{
			value: "DEEP_ARCHIVE",
		},
		SRC_STORAGE_MAPPING: ShowTaskGroupResponseDstStoragePolicy{
			value: "SRC_STORAGE_MAPPING",
		},
	}
}

func (c ShowTaskGroupResponseDstStoragePolicy) Value() string {
	return c.value
}

func (c ShowTaskGroupResponseDstStoragePolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTaskGroupResponseDstStoragePolicy) UnmarshalJSON(b []byte) error {
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

type ShowTaskGroupResponseConsistencyCheck struct {
	value string
}

type ShowTaskGroupResponseConsistencyCheckEnum struct {
	SIZE_LAST_MODIFIED ShowTaskGroupResponseConsistencyCheck
	CRC64              ShowTaskGroupResponseConsistencyCheck
	NO_CHECK           ShowTaskGroupResponseConsistencyCheck
}

func GetShowTaskGroupResponseConsistencyCheckEnum() ShowTaskGroupResponseConsistencyCheckEnum {
	return ShowTaskGroupResponseConsistencyCheckEnum{
		SIZE_LAST_MODIFIED: ShowTaskGroupResponseConsistencyCheck{
			value: "size_last_modified",
		},
		CRC64: ShowTaskGroupResponseConsistencyCheck{
			value: "crc64",
		},
		NO_CHECK: ShowTaskGroupResponseConsistencyCheck{
			value: "no_check",
		},
	}
}

func (c ShowTaskGroupResponseConsistencyCheck) Value() string {
	return c.value
}

func (c ShowTaskGroupResponseConsistencyCheck) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTaskGroupResponseConsistencyCheck) UnmarshalJSON(b []byte) error {
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
