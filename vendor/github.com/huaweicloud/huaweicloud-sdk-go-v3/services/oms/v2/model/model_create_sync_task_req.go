package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// CreateSyncTaskReq 创建同步请求body体
type CreateSyncTaskReq struct {

	// 源端云服务提供商。  可选值有AWS、Azure、Aliyun、Tencent、HuaweiCloud、QingCloud、KingsoftCloud、Baidu、Qiniu、Cloud。默认值为Aliyun。
	SrcCloudType *string `json:"src_cloud_type,omitempty"`

	// 源端桶所处的区域
	SrcRegion string `json:"src_region"`

	// 源端桶名
	SrcBucket string `json:"src_bucket"`

	// 源端桶的AK（最大长度100个字符）。
	SrcAk string `json:"src_ak"`

	// 源端桶的SK（最大长度100个字符）。
	SrcSk string `json:"src_sk"`

	// 目的端桶的AK（最大长度100个字符）。
	DstAk string `json:"dst_ak"`

	// 目的端桶的SK（最大长度100个字符）。
	DstSk string `json:"dst_sk"`

	// 目的端region
	DstRegion string `json:"dst_region"`

	// 目的端桶名
	DstBucket string `json:"dst_bucket"`

	// 任务描述，不能超过255个字符，且不能包含<>()\"'&等特殊字符。
	Description *string `json:"description,omitempty"`

	// 是否启用元数据迁移，默认否。不启用时，为保证迁移任务正常运行，仍将为您迁移ContentType元数据。
	EnableMetadataMigration *bool `json:"enable_metadata_migration,omitempty"`

	// 是否开启KMS加密，默认不开启。
	EnableKms *bool `json:"enable_kms,omitempty"`

	// 是否自动解冻归档数据，默认否。  开启后，如果遇到归档类型数据，会自动解冻再进行迁移。
	EnableRestore *bool `json:"enable_restore,omitempty"`

	// 目的端存储类型设置，当且仅当目的端为华为云OBS时需要，默认为标准存储 STANDARD：华为云OBS标准存储 IA：华为云OBS低频存储 ARCHIVE：华为云OBS归档存储 DEEP_ARCHIVE：华为云OBS深度归档存储 SRC_STORAGE_MAPPING：保留源端存储类型，将源端存储类型映射为华为云OBS存储类型
	DstStoragePolicy *CreateSyncTaskReqDstStoragePolicy `json:"dst_storage_policy,omitempty"`

	// 当源端为腾讯云时，需要填写此参数。
	AppId *string `json:"app_id,omitempty"`

	SourceCdn *SourceCdnReq `json:"source_cdn,omitempty"`

	// 一致性校验方式，用于迁移前/后校验对象是否一致，所有校验方式需满足源端/目的端对象的加密状态一致，具体校验方式和校验结果可通过对象列表查看。默认size_last_modified。 size_last_modified：默认配置。迁移前后，通过对比源端和目的端对象大小+最后修改时间，判断对象是否已存在或迁移后数据是否完整。源端与目的端同名对象大小相同，且目的端对象最后修改时间晚于源端对象最后修改时间，则代表该对象已存在/迁移成功。 crc64：目前仅支持华为/阿里/腾讯。迁移前后，通过对比源端和目的端对象元数据中CRC64值是否相同，判断对象是否已存在/迁移完成。如果源端与目的端对象元数据中不存在CRC64值，则系统会默认使用大小/最后修改时间校验方式来校验。 transmission：目前仅支持HTTP/HTTPS数据源。当源端对象无法通过标准http协议中content-length字段获取数据大小时，默认数据下载成功即迁移成功，不对数据做额外校验，且迁移时源端对象默认覆盖目的端同名对象。当源端对象能正常通过标准http协议中content-length字段获取数据大小时，则采用大小/最后修改时间校验方式来校验。
	ConsistencyCheck *CreateSyncTaskReqConsistencyCheck `json:"consistency_check,omitempty"`
}

func (o CreateSyncTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSyncTaskReq struct{}"
	}

	return strings.Join([]string{"CreateSyncTaskReq", string(data)}, " ")
}

type CreateSyncTaskReqDstStoragePolicy struct {
	value string
}

type CreateSyncTaskReqDstStoragePolicyEnum struct {
	STANDARD            CreateSyncTaskReqDstStoragePolicy
	IA                  CreateSyncTaskReqDstStoragePolicy
	ARCHIVE             CreateSyncTaskReqDstStoragePolicy
	DEEP_ARCHIVE        CreateSyncTaskReqDstStoragePolicy
	SRC_STORAGE_MAPPING CreateSyncTaskReqDstStoragePolicy
}

func GetCreateSyncTaskReqDstStoragePolicyEnum() CreateSyncTaskReqDstStoragePolicyEnum {
	return CreateSyncTaskReqDstStoragePolicyEnum{
		STANDARD: CreateSyncTaskReqDstStoragePolicy{
			value: "STANDARD",
		},
		IA: CreateSyncTaskReqDstStoragePolicy{
			value: "IA",
		},
		ARCHIVE: CreateSyncTaskReqDstStoragePolicy{
			value: "ARCHIVE",
		},
		DEEP_ARCHIVE: CreateSyncTaskReqDstStoragePolicy{
			value: "DEEP_ARCHIVE",
		},
		SRC_STORAGE_MAPPING: CreateSyncTaskReqDstStoragePolicy{
			value: "SRC_STORAGE_MAPPING",
		},
	}
}

func (c CreateSyncTaskReqDstStoragePolicy) Value() string {
	return c.value
}

func (c CreateSyncTaskReqDstStoragePolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateSyncTaskReqDstStoragePolicy) UnmarshalJSON(b []byte) error {
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

type CreateSyncTaskReqConsistencyCheck struct {
	value string
}

type CreateSyncTaskReqConsistencyCheckEnum struct {
	SIZE_LAST_MODIFIED CreateSyncTaskReqConsistencyCheck
	CRC64              CreateSyncTaskReqConsistencyCheck
	TRANSMISSION       CreateSyncTaskReqConsistencyCheck
}

func GetCreateSyncTaskReqConsistencyCheckEnum() CreateSyncTaskReqConsistencyCheckEnum {
	return CreateSyncTaskReqConsistencyCheckEnum{
		SIZE_LAST_MODIFIED: CreateSyncTaskReqConsistencyCheck{
			value: "size_last_modified",
		},
		CRC64: CreateSyncTaskReqConsistencyCheck{
			value: "crc64",
		},
		TRANSMISSION: CreateSyncTaskReqConsistencyCheck{
			value: "transmission",
		},
	}
}

func (c CreateSyncTaskReqConsistencyCheck) Value() string {
	return c.value
}

func (c CreateSyncTaskReqConsistencyCheck) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateSyncTaskReqConsistencyCheck) UnmarshalJSON(b []byte) error {
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
