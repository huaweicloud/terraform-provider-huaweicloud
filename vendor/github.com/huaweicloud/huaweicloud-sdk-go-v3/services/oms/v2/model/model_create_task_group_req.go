package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// CreateTaskGroupReq 迁移任务组创建请求体
type CreateTaskGroupReq struct {
	SrcNode *TaskGroupSrcNode `json:"src_node"`

	// 任务组描述，不能超过255个字符，且不能包含^<>&\"'等特殊字符。
	Description *string `json:"description,omitempty"`

	DstNode *TaskGroupDstNode `json:"dst_node"`

	// 是否启用元数据迁移，默认否。不启用时，为保证迁移任务正常运行，仍将为您迁移ContentType元数据。
	EnableMetadataMigration *bool `json:"enable_metadata_migration,omitempty"`

	// 是否开启记录失败对象
	EnableFailedObjectRecording *bool `json:"enable_failed_object_recording,omitempty"`

	// 是否自动解冻归档数据，默认否。 开启后，如果遇到归档类型数据，会自动解冻再进行迁移。
	EnableRestore *bool `json:"enable_restore,omitempty"`

	// 是否开启KMS加密，默认不开启。
	EnableKms bool `json:"enable_kms"`

	// 任务类型，默认为PREFIX。 LIST：对象列表迁移 URL_LIST：URL列表迁移， PREFIX：对象前缀迁移
	TaskType *CreateTaskGroupReqTaskType `json:"task_type,omitempty"`

	// 配置流量控制策略。数组中一个元素对应一个时段的最大带宽，最多允许5个时段，且时段不能重叠。
	BandwidthPolicy *[]BandwidthPolicyDto `json:"bandwidth_policy,omitempty"`

	SmnConfig *SmnConfig `json:"smn_config,omitempty"`

	SourceCdn *SourceCdnReq `json:"source_cdn,omitempty"`

	// 以时间戳方式表示的迁移指定时间（单位：秒），表示仅迁移在指定时间之后修改的源端待迁移对象。默认不设置迁移指定时间。
	MigrateSince *int64 `json:"migrate_since,omitempty"`

	// 迁移前同名对象覆盖方式，用于迁移前判断源端与目的端有同名对象时，覆盖目的端或跳过迁移。默认SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE。 NO_OVERWRITE：不覆盖。迁移前源端对象与目的端对象同名时，不做对比直接跳过迁移。 SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE：大小/最后修改时间对比覆盖。默认配置。迁移前源端对象与目的端对象同名时，通过对比源端和目的端对象大小和最后修改时间，判断是否覆盖目的端，需满足源端/目的端对象的加密状态一致。源端与目的端同名对象大小不相同，或目的端对象的最后修改时间晚于源端对象的最后修改时间(源端较新)，覆盖目的端。 CRC64_COMPARISON_OVERWRITE：CRC64对比覆盖。目前仅支持华为/阿里/腾讯。迁移前源端对象与目的端对象同名时，通过对比源端和目的端对象元数据中CRC64值是否相同，判断是否覆盖目的端，需满足源端/目的端对象的加密状态一致。如果源端与目的端对象元数据中不存在CRC64值，则系统会默认使用SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE(大小/最后修改时间对比覆盖)来对比进行覆盖判断。 FULL_OVERWRITE：全覆盖。迁移前源端对象与目的端对象同名时，不做对比覆盖目的端。
	ObjectOverwriteMode *CreateTaskGroupReqObjectOverwriteMode `json:"object_overwrite_mode,omitempty"`

	// 目的端存储类型设置，当且仅当目的端为华为云OBS时需要，默认为标准存储 STANDARD：华为云OBS标准存储 IA：华为云OBS低频存储 ARCHIVE：华为云OBS归档存储 DEEP_ARCHIVE：华为云OBS深度归档存储 SRC_STORAGE_MAPPING：保留源端存储类型，将源端存储类型映射为华为云OBS存储类型
	DstStoragePolicy *CreateTaskGroupReqDstStoragePolicy `json:"dst_storage_policy,omitempty"`

	// 一致性校验方式，用于迁移前/后校验对象是否一致，所有校验方式需满足源端/目的端对象的加密状态一致，具体校验方式和校验结果可通过对象列表查看。默认size_last_modified。 size_last_modified：默认配置。迁移前后，通过对比源端和目的端对象大小+最后修改时间，判断对象是否已存在或迁移后数据是否完整。源端与目的端同名对象大小相同，且目的端对象的最后修改时间不早于源端对象的最后修改时间，则代表该对象已存在/迁移成功。 crc64：目前仅支持华为/阿里/腾讯。迁移前后，通过对比源端和目的端对象元数据中CRC64值是否相同，判断对象是否已存在/迁移完成。如果源端与目的端对象元数据中不存在CRC64值，则系统会默认使用大小/最后修改时间校验方式来校验。 no_check：目前仅支持HTTP/HTTPS数据源。当源端对象无法通过标准http协议中content-length字段获取数据大小时，默认数据下载成功即迁移成功，不对数据做额外校验，且迁移时源端对象默认覆盖目的端同名对象。当源端对象能正常通过标准http协议中content-length字段获取数据大小时，则采用大小/最后修改时间校验方式来校验。
	ConsistencyCheck *CreateTaskGroupReqConsistencyCheck `json:"consistency_check,omitempty"`

	// 是否开启请求者付款，在启用后，请求者支付请求和数据传输费用。
	EnableRequesterPays *bool `json:"enable_requester_pays,omitempty"`
}

func (o CreateTaskGroupReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTaskGroupReq struct{}"
	}

	return strings.Join([]string{"CreateTaskGroupReq", string(data)}, " ")
}

type CreateTaskGroupReqTaskType struct {
	value string
}

type CreateTaskGroupReqTaskTypeEnum struct {
	LIST     CreateTaskGroupReqTaskType
	URL_LIST CreateTaskGroupReqTaskType
	PREFIX   CreateTaskGroupReqTaskType
}

func GetCreateTaskGroupReqTaskTypeEnum() CreateTaskGroupReqTaskTypeEnum {
	return CreateTaskGroupReqTaskTypeEnum{
		LIST: CreateTaskGroupReqTaskType{
			value: "LIST",
		},
		URL_LIST: CreateTaskGroupReqTaskType{
			value: "URL_LIST",
		},
		PREFIX: CreateTaskGroupReqTaskType{
			value: "PREFIX",
		},
	}
}

func (c CreateTaskGroupReqTaskType) Value() string {
	return c.value
}

func (c CreateTaskGroupReqTaskType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateTaskGroupReqTaskType) UnmarshalJSON(b []byte) error {
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

type CreateTaskGroupReqObjectOverwriteMode struct {
	value string
}

type CreateTaskGroupReqObjectOverwriteModeEnum struct {
	NO_OVERWRITE                            CreateTaskGroupReqObjectOverwriteMode
	SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE CreateTaskGroupReqObjectOverwriteMode
	CRC64_COMPARISON_OVERWRITE              CreateTaskGroupReqObjectOverwriteMode
	FULL_OVERWRITE                          CreateTaskGroupReqObjectOverwriteMode
}

func GetCreateTaskGroupReqObjectOverwriteModeEnum() CreateTaskGroupReqObjectOverwriteModeEnum {
	return CreateTaskGroupReqObjectOverwriteModeEnum{
		NO_OVERWRITE: CreateTaskGroupReqObjectOverwriteMode{
			value: "NO_OVERWRITE",
		},
		SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE: CreateTaskGroupReqObjectOverwriteMode{
			value: "SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE",
		},
		CRC64_COMPARISON_OVERWRITE: CreateTaskGroupReqObjectOverwriteMode{
			value: "CRC64_COMPARISON_OVERWRITE",
		},
		FULL_OVERWRITE: CreateTaskGroupReqObjectOverwriteMode{
			value: "FULL_OVERWRITE",
		},
	}
}

func (c CreateTaskGroupReqObjectOverwriteMode) Value() string {
	return c.value
}

func (c CreateTaskGroupReqObjectOverwriteMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateTaskGroupReqObjectOverwriteMode) UnmarshalJSON(b []byte) error {
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

type CreateTaskGroupReqDstStoragePolicy struct {
	value string
}

type CreateTaskGroupReqDstStoragePolicyEnum struct {
	STANDARD            CreateTaskGroupReqDstStoragePolicy
	IA                  CreateTaskGroupReqDstStoragePolicy
	ARCHIVE             CreateTaskGroupReqDstStoragePolicy
	DEEP_ARCHIVE        CreateTaskGroupReqDstStoragePolicy
	SRC_STORAGE_MAPPING CreateTaskGroupReqDstStoragePolicy
}

func GetCreateTaskGroupReqDstStoragePolicyEnum() CreateTaskGroupReqDstStoragePolicyEnum {
	return CreateTaskGroupReqDstStoragePolicyEnum{
		STANDARD: CreateTaskGroupReqDstStoragePolicy{
			value: "STANDARD",
		},
		IA: CreateTaskGroupReqDstStoragePolicy{
			value: "IA",
		},
		ARCHIVE: CreateTaskGroupReqDstStoragePolicy{
			value: "ARCHIVE",
		},
		DEEP_ARCHIVE: CreateTaskGroupReqDstStoragePolicy{
			value: "DEEP_ARCHIVE",
		},
		SRC_STORAGE_MAPPING: CreateTaskGroupReqDstStoragePolicy{
			value: "SRC_STORAGE_MAPPING",
		},
	}
}

func (c CreateTaskGroupReqDstStoragePolicy) Value() string {
	return c.value
}

func (c CreateTaskGroupReqDstStoragePolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateTaskGroupReqDstStoragePolicy) UnmarshalJSON(b []byte) error {
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

type CreateTaskGroupReqConsistencyCheck struct {
	value string
}

type CreateTaskGroupReqConsistencyCheckEnum struct {
	SIZE_LAST_MODIFIED CreateTaskGroupReqConsistencyCheck
	CRC64              CreateTaskGroupReqConsistencyCheck
	NO_CHECK           CreateTaskGroupReqConsistencyCheck
}

func GetCreateTaskGroupReqConsistencyCheckEnum() CreateTaskGroupReqConsistencyCheckEnum {
	return CreateTaskGroupReqConsistencyCheckEnum{
		SIZE_LAST_MODIFIED: CreateTaskGroupReqConsistencyCheck{
			value: "size_last_modified",
		},
		CRC64: CreateTaskGroupReqConsistencyCheck{
			value: "crc64",
		},
		NO_CHECK: CreateTaskGroupReqConsistencyCheck{
			value: "no_check",
		},
	}
}

func (c CreateTaskGroupReqConsistencyCheck) Value() string {
	return c.value
}

func (c CreateTaskGroupReqConsistencyCheck) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateTaskGroupReqConsistencyCheck) UnmarshalJSON(b []byte) error {
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
