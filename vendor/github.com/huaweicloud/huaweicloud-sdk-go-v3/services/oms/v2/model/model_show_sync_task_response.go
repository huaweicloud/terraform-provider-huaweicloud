package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowSyncTaskResponse Response Object
type ShowSyncTaskResponse struct {

	// 同步任务ID
	SyncTaskId *string `json:"sync_task_id,omitempty"`

	// 源端云服务提供商。  可选值有AWS、Azure、Aliyun、Tencent、HuaweiCloud、QingCloud、KingsoftCloud、Baidu、Qiniu、UCloud。默认值为Aliyun。
	SrcCloudType *ShowSyncTaskResponseSrcCloudType `json:"src_cloud_type,omitempty"`

	// 源端桶所处的区域
	SrcRegion *string `json:"src_region,omitempty"`

	// 源端桶
	SrcBucket *string `json:"src_bucket,omitempty"`

	// 同步任务创建时间（Unix时间戳，毫秒）
	CreateTime *int64 `json:"create_time,omitempty"`

	// 最近启动同步任务时间（Unix时间戳，毫秒）
	LastStartTime *int64 `json:"last_start_time,omitempty"`

	// 目的端桶。
	DstBucket *string `json:"dst_bucket,omitempty"`

	// 目的端region
	DstRegion *string `json:"dst_region,omitempty"`

	// 任务描述，不能超过255个字符，且不能包含<>()\"'&等特殊字符。
	Description *string `json:"description,omitempty"`

	// 同步任务状态 SYNCHRONIZING：同步中 STOPPED：已停止
	Status *ShowSyncTaskResponseStatus `json:"status,omitempty"`

	// 是否开启KMS加密，默认不开启。
	EnableKms *bool `json:"enable_kms,omitempty"`

	// 是否启用元数据迁移，默认否。不启用时，为保证迁移任务正常运行，仍将为您迁移ContentType元数据。
	EnableMetadataMigration *bool `json:"enable_metadata_migration,omitempty"`

	// 是否自动解冻归档数据，默认否。 开启后，如果遇到归档类型数据，会自动解冻再进行迁移。
	EnableRestore *bool `json:"enable_restore,omitempty"`

	// 目的端存储类型设置，当且仅当目的端为华为云OBS时需要，默认为标准存储 STANDARD：华为云OBS标准存储 IA：华为云OBS低频存储 ARCHIVE：华为云OBS归档存储 DEEP_ARCHIVE：华为云OBS深度归档存储 SRC_STORAGE_MAPPING：保留源端存储类型，将源端存储类型映射为华为云OBS存储类型
	DstStoragePolicy *ShowSyncTaskResponseDstStoragePolicy `json:"dst_storage_policy,omitempty"`

	// 当源端为腾讯云时，需要填写此参数。
	AppId *string `json:"app_id,omitempty"`

	// 当月接收同步请求对象数
	MonthlyAcceptanceRequest *int64 `json:"monthly_acceptance_request,omitempty"`

	// 当月同步成功对象数
	MonthlySuccessObject *int64 `json:"monthly_success_object,omitempty"`

	// 当月同步失败对象数
	MonthlyFailureObject *int64 `json:"monthly_failure_object,omitempty"`

	// 当月同步忽略对象数
	MonthlySkipObject *int64 `json:"monthly_skip_object,omitempty"`

	// 当月同步对象容量大小（Byte）。
	MonthlySize *int64 `json:"monthly_size,omitempty"`

	// 迁移前同名对象覆盖方式，用于迁移前判断源端与目的端有同名对象时，覆盖目的端或跳过迁移。默认SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE。 NO_OVERWRITE：不覆盖。迁移前源端对象与目的端对象同名时，不做对比直接跳过迁移。 SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE：大小/最后修改时间对比覆盖。默认配置。迁移前源端对象与目的端对象同名时，通过对比源端和目的端对象大小和最后修改时间，判断是否覆盖目的端，需满足源端/目的端对象的加密状态一致。源端与目的端同名对象大小不相同，或目的端对象的最后修改时间晚于源端对象的最后修改时间(源端较新)，覆盖目的端。 CRC64_COMPARISON_OVERWRITE：CRC64对比覆盖。目前仅支持华为/阿里/腾讯。迁移前源端对象与目的端对象同名时，通过对比源端和目的端对象元数据中CRC64值是否相同，判断是否覆盖目的端，需满足源端/目的端对象的加密状态一致。如果源端与目的端对象元数据中不存在CRC64值，则系统会默认使用SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE(大小/最后修改时间对比覆盖)来对比进行覆盖判断。 FULL_OVERWRITE：全覆盖。迁移前源端对象与目的端对象同名时，不做对比覆盖目的端。
	ObjectOverwriteMode *ShowSyncTaskResponseObjectOverwriteMode `json:"object_overwrite_mode,omitempty"`

	// 迁移后对象一致性校验方式，用于迁移后校验对象是否一致，所有校验方式需满足源端/目的端对象的加密状态一致，具体校验方式和校验结果可通过对象列表查看。默认size_last_modified。 size_last_modified：默认配置。迁移后，通过对比源端和目的端对象大小和最后修改时间，判断对象迁移后数据是否完整。源端与目的端同名对象大小相同，且目的端对象的最后修改时间不早于源端对象的最后修改时间，则代表该对象迁移成功。 crc64：目前仅支持华为/阿里/腾讯。迁移后，通过对比源端和目的端对象元数据中CRC64值是否相同，判断对象是否迁移完成。如果源端与目的端对象元数据中不存在CRC64值，则系统会默认使用大小/最后修改时间校验方式来校验。 no_check：目前仅支持HTTP/HTTPS数据源。当源端对象无法通过标准http协议中content-length字段获取数据大小时，默认数据下载成功即迁移成功，不对数据做额外校验。当源端对象能正常通过标准http协议中content-length字段获取数据大小时，则采用大小/最后修改时间校验方式来校验。
	ConsistencyCheck *ShowSyncTaskResponseConsistencyCheck `json:"consistency_check,omitempty"`
	HttpStatusCode   int                                   `json:"-"`
}

func (o ShowSyncTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowSyncTaskResponse struct{}"
	}

	return strings.Join([]string{"ShowSyncTaskResponse", string(data)}, " ")
}

type ShowSyncTaskResponseSrcCloudType struct {
	value string
}

type ShowSyncTaskResponseSrcCloudTypeEnum struct {
	AWS            ShowSyncTaskResponseSrcCloudType
	AZURE          ShowSyncTaskResponseSrcCloudType
	ALIYUN         ShowSyncTaskResponseSrcCloudType
	TENCENT        ShowSyncTaskResponseSrcCloudType
	HUAWEI_CLOUD   ShowSyncTaskResponseSrcCloudType
	QING_CLOUD     ShowSyncTaskResponseSrcCloudType
	KINGSOFT_CLOUD ShowSyncTaskResponseSrcCloudType
	BAIDU          ShowSyncTaskResponseSrcCloudType
	QINIU          ShowSyncTaskResponseSrcCloudType
	U_CLOUD        ShowSyncTaskResponseSrcCloudType
}

func GetShowSyncTaskResponseSrcCloudTypeEnum() ShowSyncTaskResponseSrcCloudTypeEnum {
	return ShowSyncTaskResponseSrcCloudTypeEnum{
		AWS: ShowSyncTaskResponseSrcCloudType{
			value: "AWS",
		},
		AZURE: ShowSyncTaskResponseSrcCloudType{
			value: "Azure",
		},
		ALIYUN: ShowSyncTaskResponseSrcCloudType{
			value: "Aliyun",
		},
		TENCENT: ShowSyncTaskResponseSrcCloudType{
			value: "Tencent",
		},
		HUAWEI_CLOUD: ShowSyncTaskResponseSrcCloudType{
			value: "HuaweiCloud",
		},
		QING_CLOUD: ShowSyncTaskResponseSrcCloudType{
			value: "QingCloud",
		},
		KINGSOFT_CLOUD: ShowSyncTaskResponseSrcCloudType{
			value: "KingsoftCloud",
		},
		BAIDU: ShowSyncTaskResponseSrcCloudType{
			value: "Baidu",
		},
		QINIU: ShowSyncTaskResponseSrcCloudType{
			value: "Qiniu",
		},
		U_CLOUD: ShowSyncTaskResponseSrcCloudType{
			value: "UCloud",
		},
	}
}

func (c ShowSyncTaskResponseSrcCloudType) Value() string {
	return c.value
}

func (c ShowSyncTaskResponseSrcCloudType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowSyncTaskResponseSrcCloudType) UnmarshalJSON(b []byte) error {
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

type ShowSyncTaskResponseStatus struct {
	value string
}

type ShowSyncTaskResponseStatusEnum struct {
	SYNCHRONIZING ShowSyncTaskResponseStatus
	STOPPED       ShowSyncTaskResponseStatus
}

func GetShowSyncTaskResponseStatusEnum() ShowSyncTaskResponseStatusEnum {
	return ShowSyncTaskResponseStatusEnum{
		SYNCHRONIZING: ShowSyncTaskResponseStatus{
			value: "SYNCHRONIZING",
		},
		STOPPED: ShowSyncTaskResponseStatus{
			value: "STOPPED",
		},
	}
}

func (c ShowSyncTaskResponseStatus) Value() string {
	return c.value
}

func (c ShowSyncTaskResponseStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowSyncTaskResponseStatus) UnmarshalJSON(b []byte) error {
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

type ShowSyncTaskResponseDstStoragePolicy struct {
	value string
}

type ShowSyncTaskResponseDstStoragePolicyEnum struct {
	STANDARD            ShowSyncTaskResponseDstStoragePolicy
	IA                  ShowSyncTaskResponseDstStoragePolicy
	ARCHIVE             ShowSyncTaskResponseDstStoragePolicy
	DEEP_ARCHIVE        ShowSyncTaskResponseDstStoragePolicy
	SRC_STORAGE_MAPPING ShowSyncTaskResponseDstStoragePolicy
}

func GetShowSyncTaskResponseDstStoragePolicyEnum() ShowSyncTaskResponseDstStoragePolicyEnum {
	return ShowSyncTaskResponseDstStoragePolicyEnum{
		STANDARD: ShowSyncTaskResponseDstStoragePolicy{
			value: "STANDARD",
		},
		IA: ShowSyncTaskResponseDstStoragePolicy{
			value: "IA",
		},
		ARCHIVE: ShowSyncTaskResponseDstStoragePolicy{
			value: "ARCHIVE",
		},
		DEEP_ARCHIVE: ShowSyncTaskResponseDstStoragePolicy{
			value: "DEEP_ARCHIVE",
		},
		SRC_STORAGE_MAPPING: ShowSyncTaskResponseDstStoragePolicy{
			value: "SRC_STORAGE_MAPPING",
		},
	}
}

func (c ShowSyncTaskResponseDstStoragePolicy) Value() string {
	return c.value
}

func (c ShowSyncTaskResponseDstStoragePolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowSyncTaskResponseDstStoragePolicy) UnmarshalJSON(b []byte) error {
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

type ShowSyncTaskResponseObjectOverwriteMode struct {
	value string
}

type ShowSyncTaskResponseObjectOverwriteModeEnum struct {
	NO_OVERWRITE                            ShowSyncTaskResponseObjectOverwriteMode
	SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE ShowSyncTaskResponseObjectOverwriteMode
	CRC64_COMPARISON_OVERWRITE              ShowSyncTaskResponseObjectOverwriteMode
	FULL_OVERWRITE                          ShowSyncTaskResponseObjectOverwriteMode
}

func GetShowSyncTaskResponseObjectOverwriteModeEnum() ShowSyncTaskResponseObjectOverwriteModeEnum {
	return ShowSyncTaskResponseObjectOverwriteModeEnum{
		NO_OVERWRITE: ShowSyncTaskResponseObjectOverwriteMode{
			value: "NO_OVERWRITE",
		},
		SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE: ShowSyncTaskResponseObjectOverwriteMode{
			value: "SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE",
		},
		CRC64_COMPARISON_OVERWRITE: ShowSyncTaskResponseObjectOverwriteMode{
			value: "CRC64_COMPARISON_OVERWRITE",
		},
		FULL_OVERWRITE: ShowSyncTaskResponseObjectOverwriteMode{
			value: "FULL_OVERWRITE",
		},
	}
}

func (c ShowSyncTaskResponseObjectOverwriteMode) Value() string {
	return c.value
}

func (c ShowSyncTaskResponseObjectOverwriteMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowSyncTaskResponseObjectOverwriteMode) UnmarshalJSON(b []byte) error {
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

type ShowSyncTaskResponseConsistencyCheck struct {
	value string
}

type ShowSyncTaskResponseConsistencyCheckEnum struct {
	SIZE_LAST_MODIFIED ShowSyncTaskResponseConsistencyCheck
	CRC64              ShowSyncTaskResponseConsistencyCheck
	NO_CHECK           ShowSyncTaskResponseConsistencyCheck
}

func GetShowSyncTaskResponseConsistencyCheckEnum() ShowSyncTaskResponseConsistencyCheckEnum {
	return ShowSyncTaskResponseConsistencyCheckEnum{
		SIZE_LAST_MODIFIED: ShowSyncTaskResponseConsistencyCheck{
			value: "size_last_modified",
		},
		CRC64: ShowSyncTaskResponseConsistencyCheck{
			value: "crc64",
		},
		NO_CHECK: ShowSyncTaskResponseConsistencyCheck{
			value: "no_check",
		},
	}
}

func (c ShowSyncTaskResponseConsistencyCheck) Value() string {
	return c.value
}

func (c ShowSyncTaskResponseConsistencyCheck) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowSyncTaskResponseConsistencyCheck) UnmarshalJSON(b []byte) error {
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
