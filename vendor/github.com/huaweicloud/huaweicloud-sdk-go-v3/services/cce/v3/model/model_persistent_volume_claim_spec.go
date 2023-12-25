package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// PersistentVolumeClaimSpec
type PersistentVolumeClaimSpec struct {

	// 资源需为已经存在的存储资源 - 如果存储资源类型是SFS、EVS、SFS-Turbo，本参数需要填入对应资源的ID - 如果资源类型为OBS，本参数填入OBS名称
	VolumeID string `json:"volumeID"`

	// 云存储的类型，和volumeID搭配使用。即volumeID和storageType必须同时被配置。  - bs：EVS云存储 - nfs：SFS弹性文件存储 - obs：OBS对象存储 - efs：SFS Turbo极速文件存储
	StorageType string `json:"storageType"`

	// 指定volume应该具有的访问模式，列表中仅第一个配置参数有效。 - ReadWriteOnce：该卷可以被单个节点以读/写模式挂载   >集群版本为v1.13.10且storage-driver版本为1.0.19时，才支持此功能。 - ReadOnlyMany：该卷可以被多个节点以只读模式挂载（默认） - ReadWriteMany：该卷可以被多个节点以读/写模式挂载
	AccessModes []PersistentVolumeClaimSpecAccessModes `json:"accessModes"`

	// PVC的StorageClass名称
	StorageClassName *string `json:"storageClassName,omitempty"`

	// PVC绑定的PV名称
	VolumeName *string `json:"volumeName,omitempty"`

	Resources *ResourceRequirements `json:"resources,omitempty"`

	// PVC指定的PV类型
	VolumeMode *string `json:"volumeMode,omitempty"`
}

func (o PersistentVolumeClaimSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PersistentVolumeClaimSpec struct{}"
	}

	return strings.Join([]string{"PersistentVolumeClaimSpec", string(data)}, " ")
}

type PersistentVolumeClaimSpecAccessModes struct {
	value string
}

type PersistentVolumeClaimSpecAccessModesEnum struct {
	READ_ONLY_MANY  PersistentVolumeClaimSpecAccessModes
	READ_WRITE_MANY PersistentVolumeClaimSpecAccessModes
}

func GetPersistentVolumeClaimSpecAccessModesEnum() PersistentVolumeClaimSpecAccessModesEnum {
	return PersistentVolumeClaimSpecAccessModesEnum{
		READ_ONLY_MANY: PersistentVolumeClaimSpecAccessModes{
			value: "ReadOnlyMany",
		},
		READ_WRITE_MANY: PersistentVolumeClaimSpecAccessModes{
			value: "ReadWriteMany",
		},
	}
}

func (c PersistentVolumeClaimSpecAccessModes) Value() string {
	return c.value
}

func (c PersistentVolumeClaimSpecAccessModes) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PersistentVolumeClaimSpecAccessModes) UnmarshalJSON(b []byte) error {
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
