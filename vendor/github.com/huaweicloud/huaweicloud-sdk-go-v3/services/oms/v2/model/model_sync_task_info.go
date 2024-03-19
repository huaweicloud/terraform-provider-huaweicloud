package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// SyncTaskInfo 查询同步任务列表时返回的单个同步任务信息实体
type SyncTaskInfo struct {

	// 同步任务ID
	SyncTaskId *string `json:"sync_task_id,omitempty"`

	// 源端云服务提供商。  可选值有AWS、Azure、Aliyun、Tencent、HuaweiCloud、QingCloud、KingsoftCloud、Baidu、Qiniu、UCloud。默认值为Aliyun。
	SrcCloudType *SyncTaskInfoSrcCloudType `json:"src_cloud_type,omitempty"`

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
	Status *SyncTaskInfoStatus `json:"status,omitempty"`

	// 是否开启KMS加密，默认不开启。
	EnableKms *bool `json:"enable_kms,omitempty"`

	// 是否启用元数据迁移，默认否。不启用时，为保证迁移任务正常运行，仍将为您迁移ContentType元数据。
	EnableMetadataMigration *bool `json:"enable_metadata_migration,omitempty"`

	// 是否自动解冻归档数据，默认否。 开启后，如果遇到归档类型数据，会自动解冻再进行迁移。
	EnableRestore *bool `json:"enable_restore,omitempty"`

	// 迁移前同名对象覆盖方式，用于迁移前判断源端与目的端有同名对象时，覆盖目的端或跳过迁移。默认SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE。 NO_OVERWRITE：不覆盖。迁移前源端对象与目的端对象同名时，不做对比直接跳过迁移。 SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE：大小/最后修改时间对比覆盖。默认配置。迁移前源端对象与目的端对象同名时，通过对比源端和目的端对象大小和最后修改时间，判断是否覆盖目的端，需满足源端/目的端对象的加密状态一致。源端与目的端同名对象大小不相同，或目的端对象的最后修改时间晚于源端对象的最后修改时间(源端较新)，覆盖目的端。 CRC64_COMPARISON_OVERWRITE：CRC64对比覆盖。目前仅支持华为/阿里/腾讯。迁移前源端对象与目的端对象同名时，通过对比源端和目的端对象元数据中CRC64值是否相同，判断是否覆盖目的端，需满足源端/目的端对象的加密状态一致。如果源端与目的端对象元数据中不存在CRC64值，则系统会默认使用SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE(大小/最后修改时间对比覆盖)来对比进行覆盖判断。 FULL_OVERWRITE：全覆盖。迁移前源端对象与目的端对象同名时，不做对比覆盖目的端。
	ObjectOverwriteMode *SyncTaskInfoObjectOverwriteMode `json:"object_overwrite_mode,omitempty"`

	// 目的端存储类型设置，当且仅当目的端为华为云OBS时需要，默认为标准存储 STANDARD：华为云OBS标准存储 IA：华为云OBS低频存储 ARCHIVE：华为云OBS归档存储 DEEP_ARCHIVE：华为云OBS深度归档存储 SRC_STORAGE_MAPPING：保留源端存储类型，将源端存储类型映射为华为云OBS存储类型
	DstStoragePolicy *SyncTaskInfoDstStoragePolicy `json:"dst_storage_policy,omitempty"`

	// 当源端为腾讯云时，需要填写此参数。
	AppId *string `json:"app_id,omitempty"`

	SourceCdn *SourceCdnResp `json:"source_cdn,omitempty"`

	// 迁移后对象一致性校验方式，用于迁移后校验对象是否一致，所有校验方式需满足源端/目的端对象的加密状态一致，具体校验方式和校验结果可通过对象列表查看。默认size_last_modified。 size_last_modified：默认配置。迁移后，通过对比源端和目的端对象大小和最后修改时间，判断对象迁移后数据是否完整。源端与目的端同名对象大小相同，且目的端对象的最后修改时间不早于源端对象的最后修改时间，则代表该对象迁移成功。 crc64：目前仅支持华为/阿里/腾讯。迁移后，通过对比源端和目的端对象元数据中CRC64值是否相同，判断对象是否迁移完成。如果源端与目的端对象元数据中不存在CRC64值，则系统会默认使用大小/最后修改时间校验方式来校验。 no_check：目前仅支持HTTP/HTTPS数据源。当源端对象无法通过标准http协议中content-length字段获取数据大小时，默认数据下载成功即迁移成功，不对数据做额外校验。当源端对象能正常通过标准http协议中content-length字段获取数据大小时，则采用大小/最后修改时间校验方式来校验。
	ConsistencyCheck *SyncTaskInfoConsistencyCheck `json:"consistency_check,omitempty"`
}

func (o SyncTaskInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SyncTaskInfo struct{}"
	}

	return strings.Join([]string{"SyncTaskInfo", string(data)}, " ")
}

type SyncTaskInfoSrcCloudType struct {
	value string
}

type SyncTaskInfoSrcCloudTypeEnum struct {
	AWS            SyncTaskInfoSrcCloudType
	AZURE          SyncTaskInfoSrcCloudType
	ALIYUN         SyncTaskInfoSrcCloudType
	TENCENT        SyncTaskInfoSrcCloudType
	HUAWEI_CLOUD   SyncTaskInfoSrcCloudType
	QING_CLOUD     SyncTaskInfoSrcCloudType
	KINGSOFT_CLOUD SyncTaskInfoSrcCloudType
	BAIDU          SyncTaskInfoSrcCloudType
	QINIU          SyncTaskInfoSrcCloudType
	U_CLOUD        SyncTaskInfoSrcCloudType
}

func GetSyncTaskInfoSrcCloudTypeEnum() SyncTaskInfoSrcCloudTypeEnum {
	return SyncTaskInfoSrcCloudTypeEnum{
		AWS: SyncTaskInfoSrcCloudType{
			value: "AWS",
		},
		AZURE: SyncTaskInfoSrcCloudType{
			value: "Azure",
		},
		ALIYUN: SyncTaskInfoSrcCloudType{
			value: "Aliyun",
		},
		TENCENT: SyncTaskInfoSrcCloudType{
			value: "Tencent",
		},
		HUAWEI_CLOUD: SyncTaskInfoSrcCloudType{
			value: "HuaweiCloud",
		},
		QING_CLOUD: SyncTaskInfoSrcCloudType{
			value: "QingCloud",
		},
		KINGSOFT_CLOUD: SyncTaskInfoSrcCloudType{
			value: "KingsoftCloud",
		},
		BAIDU: SyncTaskInfoSrcCloudType{
			value: "Baidu",
		},
		QINIU: SyncTaskInfoSrcCloudType{
			value: "Qiniu",
		},
		U_CLOUD: SyncTaskInfoSrcCloudType{
			value: "UCloud",
		},
	}
}

func (c SyncTaskInfoSrcCloudType) Value() string {
	return c.value
}

func (c SyncTaskInfoSrcCloudType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SyncTaskInfoSrcCloudType) UnmarshalJSON(b []byte) error {
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

type SyncTaskInfoStatus struct {
	value string
}

type SyncTaskInfoStatusEnum struct {
	SYNCHRONIZING SyncTaskInfoStatus
	STOPPED       SyncTaskInfoStatus
}

func GetSyncTaskInfoStatusEnum() SyncTaskInfoStatusEnum {
	return SyncTaskInfoStatusEnum{
		SYNCHRONIZING: SyncTaskInfoStatus{
			value: "SYNCHRONIZING",
		},
		STOPPED: SyncTaskInfoStatus{
			value: "STOPPED",
		},
	}
}

func (c SyncTaskInfoStatus) Value() string {
	return c.value
}

func (c SyncTaskInfoStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SyncTaskInfoStatus) UnmarshalJSON(b []byte) error {
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

type SyncTaskInfoObjectOverwriteMode struct {
	value string
}

type SyncTaskInfoObjectOverwriteModeEnum struct {
	NO_OVERWRITE                            SyncTaskInfoObjectOverwriteMode
	SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE SyncTaskInfoObjectOverwriteMode
	CRC64_COMPARISON_OVERWRITE              SyncTaskInfoObjectOverwriteMode
	FULL_OVERWRITE                          SyncTaskInfoObjectOverwriteMode
}

func GetSyncTaskInfoObjectOverwriteModeEnum() SyncTaskInfoObjectOverwriteModeEnum {
	return SyncTaskInfoObjectOverwriteModeEnum{
		NO_OVERWRITE: SyncTaskInfoObjectOverwriteMode{
			value: "NO_OVERWRITE",
		},
		SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE: SyncTaskInfoObjectOverwriteMode{
			value: "SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE",
		},
		CRC64_COMPARISON_OVERWRITE: SyncTaskInfoObjectOverwriteMode{
			value: "CRC64_COMPARISON_OVERWRITE",
		},
		FULL_OVERWRITE: SyncTaskInfoObjectOverwriteMode{
			value: "FULL_OVERWRITE",
		},
	}
}

func (c SyncTaskInfoObjectOverwriteMode) Value() string {
	return c.value
}

func (c SyncTaskInfoObjectOverwriteMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SyncTaskInfoObjectOverwriteMode) UnmarshalJSON(b []byte) error {
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

type SyncTaskInfoDstStoragePolicy struct {
	value string
}

type SyncTaskInfoDstStoragePolicyEnum struct {
	STANDARD            SyncTaskInfoDstStoragePolicy
	IA                  SyncTaskInfoDstStoragePolicy
	ARCHIVE             SyncTaskInfoDstStoragePolicy
	DEEP_ARCHIVE        SyncTaskInfoDstStoragePolicy
	SRC_STORAGE_MAPPING SyncTaskInfoDstStoragePolicy
}

func GetSyncTaskInfoDstStoragePolicyEnum() SyncTaskInfoDstStoragePolicyEnum {
	return SyncTaskInfoDstStoragePolicyEnum{
		STANDARD: SyncTaskInfoDstStoragePolicy{
			value: "STANDARD",
		},
		IA: SyncTaskInfoDstStoragePolicy{
			value: "IA",
		},
		ARCHIVE: SyncTaskInfoDstStoragePolicy{
			value: "ARCHIVE",
		},
		DEEP_ARCHIVE: SyncTaskInfoDstStoragePolicy{
			value: "DEEP_ARCHIVE",
		},
		SRC_STORAGE_MAPPING: SyncTaskInfoDstStoragePolicy{
			value: "SRC_STORAGE_MAPPING",
		},
	}
}

func (c SyncTaskInfoDstStoragePolicy) Value() string {
	return c.value
}

func (c SyncTaskInfoDstStoragePolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SyncTaskInfoDstStoragePolicy) UnmarshalJSON(b []byte) error {
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

type SyncTaskInfoConsistencyCheck struct {
	value string
}

type SyncTaskInfoConsistencyCheckEnum struct {
	SIZE_LAST_MODIFIED SyncTaskInfoConsistencyCheck
	CRC64              SyncTaskInfoConsistencyCheck
	NO_CHECK           SyncTaskInfoConsistencyCheck
}

func GetSyncTaskInfoConsistencyCheckEnum() SyncTaskInfoConsistencyCheckEnum {
	return SyncTaskInfoConsistencyCheckEnum{
		SIZE_LAST_MODIFIED: SyncTaskInfoConsistencyCheck{
			value: "size_last_modified",
		},
		CRC64: SyncTaskInfoConsistencyCheck{
			value: "crc64",
		},
		NO_CHECK: SyncTaskInfoConsistencyCheck{
			value: "no_check",
		},
	}
}

func (c SyncTaskInfoConsistencyCheck) Value() string {
	return c.value
}

func (c SyncTaskInfoConsistencyCheck) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SyncTaskInfoConsistencyCheck) UnmarshalJSON(b []byte) error {
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
