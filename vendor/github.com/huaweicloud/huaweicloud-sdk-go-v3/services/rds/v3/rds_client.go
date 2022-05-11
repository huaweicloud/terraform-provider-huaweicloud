package v3

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"
)

type RdsClient struct {
	HcClient *http_client.HcHttpClient
}

func NewRdsClient(hcClient *http_client.HcHttpClient) *RdsClient {
	return &RdsClient{HcClient: hcClient}
}

func RdsClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// 应用参数模板
//
// 应用参数模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ApplyConfigurationAsync(request *model.ApplyConfigurationAsyncRequest) (*model.ApplyConfigurationAsyncResponse, error) {
	requestDef := GenReqDefForApplyConfigurationAsync()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ApplyConfigurationAsyncResponse), nil
	}
}

// 绑定和解绑弹性公网IP
//
// 绑定和解绑弹性公网IP。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) AttachEip(request *model.AttachEipRequest) (*model.AttachEipResponse, error) {
	requestDef := GenReqDefForAttachEip()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AttachEipResponse), nil
	}
}

// 批量添加标签
//
// 批量添加标签。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) BatchTagAddAction(request *model.BatchTagAddActionRequest) (*model.BatchTagAddActionResponse, error) {
	requestDef := GenReqDefForBatchTagAddAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchTagAddActionResponse), nil
	}
}

// 批量删除标签
//
// 批量删除标签。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) BatchTagDelAction(request *model.BatchTagDelActionRequest) (*model.BatchTagDelActionResponse, error) {
	requestDef := GenReqDefForBatchTagDelAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchTagDelActionResponse), nil
	}
}

// 更改主备实例的数据同步方式
//
// 更改主备实例的数据同步方式。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ChangeFailoverMode(request *model.ChangeFailoverModeRequest) (*model.ChangeFailoverModeResponse, error) {
	requestDef := GenReqDefForChangeFailoverMode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeFailoverModeResponse), nil
	}
}

// 切换主备实例的倒换策略
//
// 切换主备实例的倒换策略.
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ChangeFailoverStrategy(request *model.ChangeFailoverStrategyRequest) (*model.ChangeFailoverStrategyResponse, error) {
	requestDef := GenReqDefForChangeFailoverStrategy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeFailoverStrategyResponse), nil
	}
}

// 设置可维护时间段
//
// 设置可维护时间段
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ChangeOpsWindow(request *model.ChangeOpsWindowRequest) (*model.ChangeOpsWindowResponse, error) {
	requestDef := GenReqDefForChangeOpsWindow()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeOpsWindowResponse), nil
	}
}

// 创建参数模板
//
// 创建参数模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) CreateConfiguration(request *model.CreateConfigurationRequest) (*model.CreateConfigurationResponse, error) {
	requestDef := GenReqDefForCreateConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateConfigurationResponse), nil
	}
}

// 申请域名
//
// 申请域名
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) CreateDnsName(request *model.CreateDnsNameRequest) (*model.CreateDnsNameResponse, error) {
	requestDef := GenReqDefForCreateDnsName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDnsNameResponse), nil
	}
}

// 创建数据库实例
//
// 创建数据库实例。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) CreateInstance(request *model.CreateInstanceRequest) (*model.CreateInstanceResponse, error) {
	requestDef := GenReqDefForCreateInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateInstanceResponse), nil
	}
}

// 创建手动备份
//
// 创建手动备份。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) CreateManualBackup(request *model.CreateManualBackupRequest) (*model.CreateManualBackupResponse, error) {
	requestDef := GenReqDefForCreateManualBackup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateManualBackupResponse), nil
	}
}

// 恢复到新实例
//
// 恢复到新实例。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) CreateRestoreInstance(request *model.CreateRestoreInstanceRequest) (*model.CreateRestoreInstanceResponse, error) {
	requestDef := GenReqDefForCreateRestoreInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRestoreInstanceResponse), nil
	}
}

// 删除参数模板
//
// 删除参数模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) DeleteConfiguration(request *model.DeleteConfigurationRequest) (*model.DeleteConfigurationResponse, error) {
	requestDef := GenReqDefForDeleteConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteConfigurationResponse), nil
	}
}

// 删除数据库实例
//
// 删除数据库实例。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) DeleteInstance(request *model.DeleteInstanceRequest) (*model.DeleteInstanceResponse, error) {
	requestDef := GenReqDefForDeleteInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteInstanceResponse), nil
	}
}

// 删除手动备份
//
// 删除手动备份。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) DeleteManualBackup(request *model.DeleteManualBackupRequest) (*model.DeleteManualBackupResponse, error) {
	requestDef := GenReqDefForDeleteManualBackup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteManualBackupResponse), nil
	}
}

// 获取慢日志下载链接
//
// 获取慢日志下载链接。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) DownloadSlowlog(request *model.DownloadSlowlogRequest) (*model.DownloadSlowlogResponse, error) {
	requestDef := GenReqDefForDownloadSlowlog()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DownloadSlowlogResponse), nil
	}
}

// 应用参数模板
//
// 应用参数模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) EnableConfiguration(request *model.EnableConfigurationRequest) (*model.EnableConfigurationResponse, error) {
	requestDef := GenReqDefForEnableConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.EnableConfigurationResponse), nil
	}
}

// 获取审计日志列表
//
// 获取审计日志列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListAuditlogs(request *model.ListAuditlogsRequest) (*model.ListAuditlogsResponse, error) {
	requestDef := GenReqDefForListAuditlogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAuditlogsResponse), nil
	}
}

// 获取备份列表
//
// 获取备份列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListBackups(request *model.ListBackupsRequest) (*model.ListBackupsResponse, error) {
	requestDef := GenReqDefForListBackups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListBackupsResponse), nil
	}
}

// 查询SQLServer可用字符集
//
// 查询SQLServer可用字符集
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListCollations(request *model.ListCollationsRequest) (*model.ListCollationsResponse, error) {
	requestDef := GenReqDefForListCollations()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListCollationsResponse), nil
	}
}

// 获取参数模板列表
//
// 获取参数模板列表，包括所有数据库的默认参数模板和用户创建的参数模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListConfigurations(request *model.ListConfigurationsRequest) (*model.ListConfigurationsResponse, error) {
	requestDef := GenReqDefForListConfigurations()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListConfigurationsResponse), nil
	}
}

// 查询数据库引擎的版本
//
// 查询数据库引擎的版本。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListDatastores(request *model.ListDatastoresRequest) (*model.ListDatastoresResponse, error) {
	requestDef := GenReqDefForListDatastores()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDatastoresResponse), nil
	}
}

// 查询数据库错误日志
//
// 查询数据库错误日志。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListErrorLogs(request *model.ListErrorLogsRequest) (*model.ListErrorLogsResponse, error) {
	requestDef := GenReqDefForListErrorLogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListErrorLogsResponse), nil
	}
}

// 查询数据库错误日志
//
// 查询数据库错误日志。(与原v3接口相比修改offset,符合华为云服务开放 API遵从性规范3.0)
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListErrorLogsNew(request *model.ListErrorLogsNewRequest) (*model.ListErrorLogsNewResponse, error) {
	requestDef := GenReqDefForListErrorLogsNew()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListErrorLogsNewResponse), nil
	}
}

// 查询数据库规格
//
// 查询数据库规格。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListFlavors(request *model.ListFlavorsRequest) (*model.ListFlavorsResponse, error) {
	requestDef := GenReqDefForListFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListFlavorsResponse), nil
	}
}

// 查询数据库实例列表
//
// 查询数据库实例列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListInstances(request *model.ListInstancesRequest) (*model.ListInstancesResponse, error) {
	requestDef := GenReqDefForListInstances()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstancesResponse), nil
	}
}

// 获取指定ID的任务信息
//
// 获取指定ID的任务信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListJobInfo(request *model.ListJobInfoRequest) (*model.ListJobInfoResponse, error) {
	requestDef := GenReqDefForListJobInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListJobInfoResponse), nil
	}
}

// 获取指定实例和时间范围的任务信息（SQL Server）
//
// 获取指定实例和时间范围的任务信息（SQL Server）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListJobInfoDetail(request *model.ListJobInfoDetailRequest) (*model.ListJobInfoDetailResponse, error) {
	requestDef := GenReqDefForListJobInfoDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListJobInfoDetailResponse), nil
	}
}

// 查询跨区域备份列表
//
// 查询跨区域备份列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListOffSiteBackups(request *model.ListOffSiteBackupsRequest) (*model.ListOffSiteBackupsResponse, error) {
	requestDef := GenReqDefForListOffSiteBackups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListOffSiteBackupsResponse), nil
	}
}

// 查询跨区域备份实例列表
//
// 查询跨区域备份实例列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListOffSiteInstances(request *model.ListOffSiteInstancesRequest) (*model.ListOffSiteInstancesResponse, error) {
	requestDef := GenReqDefForListOffSiteInstances()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListOffSiteInstancesResponse), nil
	}
}

// 查询跨区域备份可恢复时间段
//
// 查询跨区域备份可恢复时间段。
// 如果您备份策略中的保存天数设置较长，建议您传入查询日期“date”。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListOffSiteRestoreTimes(request *model.ListOffSiteRestoreTimesRequest) (*model.ListOffSiteRestoreTimesResponse, error) {
	requestDef := GenReqDefForListOffSiteRestoreTimes()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListOffSiteRestoreTimesResponse), nil
	}
}

// 查询项目标签
//
// 查询项目标签。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListProjectTags(request *model.ListProjectTagsRequest) (*model.ListProjectTagsResponse, error) {
	requestDef := GenReqDefForListProjectTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProjectTagsResponse), nil
	}
}

// 查询可恢复时间段
//
// 查询可恢复时间段。
// 如果您备份策略中的保存天数设置较长，建议您传入查询日期“date”。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListRestoreTimes(request *model.ListRestoreTimesRequest) (*model.ListRestoreTimesResponse, error) {
	requestDef := GenReqDefForListRestoreTimes()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRestoreTimesResponse), nil
	}
}

// 查询慢日志文件列表
//
// 查询慢日志文件列表。
// 调用该接口取到慢日志文件名后，可以调用接口/v3/{project_id}/instances/{instance_id}/slowlog-download 获取慢日志文件下载链接
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListSlowLogFile(request *model.ListSlowLogFileRequest) (*model.ListSlowLogFileResponse, error) {
	requestDef := GenReqDefForListSlowLogFile()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowLogFileResponse), nil
	}
}

// 查询数据库慢日志
//
// 查询数据库慢日志。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListSlowLogs(request *model.ListSlowLogsRequest) (*model.ListSlowLogsResponse, error) {
	requestDef := GenReqDefForListSlowLogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowLogsResponse), nil
	}
}

// 查询数据库慢日志
//
// 查询数据库慢日志。(与原v3接口相比修改offset,符合华为云服务开放 API遵从性规范3.0)
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListSlowLogsNew(request *model.ListSlowLogsNewRequest) (*model.ListSlowLogsNewResponse, error) {
	requestDef := GenReqDefForListSlowLogsNew()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowLogsNewResponse), nil
	}
}

// 获取慢日志统计信息
//
// 获取慢日志统计信息
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListSlowlogStatistics(request *model.ListSlowlogStatisticsRequest) (*model.ListSlowlogStatisticsResponse, error) {
	requestDef := GenReqDefForListSlowlogStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowlogStatisticsResponse), nil
	}
}

// 查询数据库磁盘类型
//
// 查询数据库磁盘类型。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListStorageTypes(request *model.ListStorageTypesRequest) (*model.ListStorageTypesResponse, error) {
	requestDef := GenReqDefForListStorageTypes()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListStorageTypesResponse), nil
	}
}

// 迁移主备实例的备机
//
// 迁移主备实例的备机
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) MigrateFollower(request *model.MigrateFollowerRequest) (*model.MigrateFollowerResponse, error) {
	requestDef := GenReqDefForMigrateFollower()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.MigrateFollowerResponse), nil
	}
}

// 恢复到已有实例
//
// 恢复到已有实例。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) RestoreExistInstance(request *model.RestoreExistInstanceRequest) (*model.RestoreExistInstanceResponse, error) {
	requestDef := GenReqDefForRestoreExistInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RestoreExistInstanceResponse), nil
	}
}

// 表级时间点恢复(MySQL)
//
// 表级时间点恢复(MySQL)。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) RestoreTables(request *model.RestoreTablesRequest) (*model.RestoreTablesResponse, error) {
	requestDef := GenReqDefForRestoreTables()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RestoreTablesResponse), nil
	}
}

// 恢复到已有实例
//
// 恢复到已有实例。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) RestoreToExistingInstance(request *model.RestoreToExistingInstanceRequest) (*model.RestoreToExistingInstanceResponse, error) {
	requestDef := GenReqDefForRestoreToExistingInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RestoreToExistingInstanceResponse), nil
	}
}

// 设置审计日志策略
//
// 设置审计日志策略。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) SetAuditlogPolicy(request *model.SetAuditlogPolicyRequest) (*model.SetAuditlogPolicyResponse, error) {
	requestDef := GenReqDefForSetAuditlogPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetAuditlogPolicyResponse), nil
	}
}

// 设置自动备份策略
//
// 设置自动备份策略。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) SetBackupPolicy(request *model.SetBackupPolicyRequest) (*model.SetBackupPolicyResponse, error) {
	requestDef := GenReqDefForSetBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetBackupPolicyResponse), nil
	}
}

// 设置binlog本地保留时长
//
// 修改指定实例的binlog本地保留时长。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) SetBinlogClearPolicy(request *model.SetBinlogClearPolicyRequest) (*model.SetBinlogClearPolicyResponse, error) {
	requestDef := GenReqDefForSetBinlogClearPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetBinlogClearPolicyResponse), nil
	}
}

// 设置跨区域备份策略
//
// 设置跨区域备份策略。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) SetOffSiteBackupPolicy(request *model.SetOffSiteBackupPolicyRequest) (*model.SetOffSiteBackupPolicyResponse, error) {
	requestDef := GenReqDefForSetOffSiteBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetOffSiteBackupPolicyResponse), nil
	}
}

// 修改安全组
//
// 修改安全组
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) SetSecurityGroup(request *model.SetSecurityGroupRequest) (*model.SetSecurityGroupResponse, error) {
	requestDef := GenReqDefForSetSecurityGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetSecurityGroupResponse), nil
	}
}

// 生成审计日志下载链接
//
// 生成审计日志下载链接。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ShowAuditlogDownloadLink(request *model.ShowAuditlogDownloadLinkRequest) (*model.ShowAuditlogDownloadLinkResponse, error) {
	requestDef := GenReqDefForShowAuditlogDownloadLink()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAuditlogDownloadLinkResponse), nil
	}
}

// 查询审计日志策略
//
// 查询审计日志策略。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ShowAuditlogPolicy(request *model.ShowAuditlogPolicyRequest) (*model.ShowAuditlogPolicyResponse, error) {
	requestDef := GenReqDefForShowAuditlogPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAuditlogPolicyResponse), nil
	}
}

// 获取备份下载链接
//
// 获取备份下载链接。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ShowBackupDownloadLink(request *model.ShowBackupDownloadLinkRequest) (*model.ShowBackupDownloadLinkResponse, error) {
	requestDef := GenReqDefForShowBackupDownloadLink()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBackupDownloadLinkResponse), nil
	}
}

// 查询自动备份策略
//
// 查询自动备份策略。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ShowBackupPolicy(request *model.ShowBackupPolicyRequest) (*model.ShowBackupPolicyResponse, error) {
	requestDef := GenReqDefForShowBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBackupPolicyResponse), nil
	}
}

// 获取binlog本地保留时长
//
// 查寻指定实例的binlog本地保留时长。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ShowBinlogClearPolicy(request *model.ShowBinlogClearPolicyRequest) (*model.ShowBinlogClearPolicyResponse, error) {
	requestDef := GenReqDefForShowBinlogClearPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBinlogClearPolicyResponse), nil
	}
}

// 获取指定参数模板的参数
//
// 获取指定参数模板的参数。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ShowConfiguration(request *model.ShowConfigurationRequest) (*model.ShowConfigurationResponse, error) {
	requestDef := GenReqDefForShowConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowConfigurationResponse), nil
	}
}

// 查询跨云容灾复制状态
//
// 建立跨云容灾关系后，查询主实例和灾备实例间的复制状态及延迟。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ShowDrReplicaStatus(request *model.ShowDrReplicaStatusRequest) (*model.ShowDrReplicaStatusResponse, error) {
	requestDef := GenReqDefForShowDrReplicaStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDrReplicaStatusResponse), nil
	}
}

// 获取指定实例的参数模板
//
// 获取指定实例的参数模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ShowInstanceConfiguration(request *model.ShowInstanceConfigurationRequest) (*model.ShowInstanceConfigurationResponse, error) {
	requestDef := GenReqDefForShowInstanceConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowInstanceConfigurationResponse), nil
	}
}

// 查询跨区域备份策略
//
// 查询跨区域备份策略。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ShowOffSiteBackupPolicy(request *model.ShowOffSiteBackupPolicyRequest) (*model.ShowOffSiteBackupPolicyResponse, error) {
	requestDef := GenReqDefForShowOffSiteBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowOffSiteBackupPolicyResponse), nil
	}
}

// 查询配额
//
// 查询当前项目下资源配额情况。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ShowQuotas(request *model.ShowQuotasRequest) (*model.ShowQuotasResponse, error) {
	requestDef := GenReqDefForShowQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowQuotasResponse), nil
	}
}

// 手动倒换主备
//
// 手动倒换主备.
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) StartFailover(request *model.StartFailoverRequest) (*model.StartFailoverResponse, error) {
	requestDef := GenReqDefForStartFailover()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartFailoverResponse), nil
	}
}

// 扩容数据库实例的磁盘空间
//
// 扩容数据库实例的磁盘空间。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) StartInstanceEnlargeVolumeAction(request *model.StartInstanceEnlargeVolumeActionRequest) (*model.StartInstanceEnlargeVolumeActionResponse, error) {
	requestDef := GenReqDefForStartInstanceEnlargeVolumeAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartInstanceEnlargeVolumeActionResponse), nil
	}
}

// 重启数据库实例
//
// 重启数据库实例。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) StartInstanceRestartAction(request *model.StartInstanceRestartActionRequest) (*model.StartInstanceRestartActionResponse, error) {
	requestDef := GenReqDefForStartInstanceRestartAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartInstanceRestartActionResponse), nil
	}
}

// 单机转主备实例
//
// 单机转主备实例。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) StartInstanceSingleToHaAction(request *model.StartInstanceSingleToHaActionRequest) (*model.StartInstanceSingleToHaActionResponse, error) {
	requestDef := GenReqDefForStartInstanceSingleToHaAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartInstanceSingleToHaActionResponse), nil
	}
}

// 设置回收站策略
//
// 设置回收站策略。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) StartRecyclePolicy(request *model.StartRecyclePolicyRequest) (*model.StartRecyclePolicyResponse, error) {
	requestDef := GenReqDefForStartRecyclePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartRecyclePolicyResponse), nil
	}
}

// 变更数据库实例的规格
//
// 变更数据库实例的规格。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) StartResizeFlavorAction(request *model.StartResizeFlavorActionRequest) (*model.StartResizeFlavorActionResponse, error) {
	requestDef := GenReqDefForStartResizeFlavorAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartResizeFlavorActionResponse), nil
	}
}

// 开启实例
//
// 停止实例以节省费用，在停止数据库实例后，支持手动重新开启实例。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) StartupInstance(request *model.StartupInstanceRequest) (*model.StartupInstanceResponse, error) {
	requestDef := GenReqDefForStartupInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartupInstanceResponse), nil
	}
}

// 停止实例
//
// 实例进行关机，通过暂时停止按需实例以节省费用，实例默认停止七天。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) StopInstance(request *model.StopInstanceRequest) (*model.StopInstanceResponse, error) {
	requestDef := GenReqDefForStopInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopInstanceResponse), nil
	}
}

// 设置SSL数据加密
//
// 设置SSL数据加密。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) SwitchSsl(request *model.SwitchSslRequest) (*model.SwitchSslResponse, error) {
	requestDef := GenReqDefForSwitchSsl()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SwitchSslResponse), nil
	}
}

// 修改参数模板参数
//
// 修改参数模板参数。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) UpdateConfiguration(request *model.UpdateConfigurationRequest) (*model.UpdateConfigurationResponse, error) {
	requestDef := GenReqDefForUpdateConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateConfigurationResponse), nil
	}
}

// 修改内网地址
//
// 修改内网地址
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) UpdateDataIp(request *model.UpdateDataIpRequest) (*model.UpdateDataIpResponse, error) {
	requestDef := GenReqDefForUpdateDataIp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDataIpResponse), nil
	}
}

// 修改域名
//
// 修改域名
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) UpdateDnsName(request *model.UpdateDnsNameRequest) (*model.UpdateDnsNameResponse, error) {
	requestDef := GenReqDefForUpdateDnsName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDnsNameResponse), nil
	}
}

// 修改指定实例的参数
//
// 修改指定实例的参数。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) UpdateInstanceConfiguration(request *model.UpdateInstanceConfigurationRequest) (*model.UpdateInstanceConfigurationResponse, error) {
	requestDef := GenReqDefForUpdateInstanceConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateInstanceConfigurationResponse), nil
	}
}

// 修改指定实例的参数
//
// 修改指定实例的参数。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) UpdateInstanceConfigurationAsync(request *model.UpdateInstanceConfigurationAsyncRequest) (*model.UpdateInstanceConfigurationAsyncResponse, error) {
	requestDef := GenReqDefForUpdateInstanceConfigurationAsync()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateInstanceConfigurationAsyncResponse), nil
	}
}

// 修改实例名称
//
// 修改实例名称。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) UpdateInstanceName(request *model.UpdateInstanceNameRequest) (*model.UpdateInstanceNameResponse, error) {
	requestDef := GenReqDefForUpdateInstanceName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateInstanceNameResponse), nil
	}
}

// 修改数据库端口
//
// 修改数据库端口
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) UpdatePort(request *model.UpdatePortRequest) (*model.UpdatePortResponse, error) {
	requestDef := GenReqDefForUpdatePort()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePortResponse), nil
	}
}

// 修改实例备注信息
//
// 修改指定数据库实例的备注信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) UpdatePostgresqlInstanceAlias(request *model.UpdatePostgresqlInstanceAliasRequest) (*model.UpdatePostgresqlInstanceAliasResponse, error) {
	requestDef := GenReqDefForUpdatePostgresqlInstanceAlias()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePostgresqlInstanceAliasResponse), nil
	}
}

// 升级内核小版本
//
// 对实例进行小版本升级。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) UpgradeDbVersion(request *model.UpgradeDbVersionRequest) (*model.UpgradeDbVersionResponse, error) {
	requestDef := GenReqDefForUpgradeDbVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeDbVersionResponse), nil
	}
}

// 查询API版本列表
//
// 查询API版本列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListApiVersion(request *model.ListApiVersionRequest) (*model.ListApiVersionResponse, error) {
	requestDef := GenReqDefForListApiVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListApiVersionResponse), nil
	}
}

// 查询API版本列表
//
// 查询API版本列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListApiVersionNew(request *model.ListApiVersionNewRequest) (*model.ListApiVersionNewResponse, error) {
	requestDef := GenReqDefForListApiVersionNew()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListApiVersionNewResponse), nil
	}
}

// 查询指定的API版本信息
//
// 查询指定的API版本信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ShowApiVersion(request *model.ShowApiVersionRequest) (*model.ShowApiVersionResponse, error) {
	requestDef := GenReqDefForShowApiVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowApiVersionResponse), nil
	}
}

// 授权数据库帐号
//
// 授权数据库帐号。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) AllowDbUserPrivilege(request *model.AllowDbUserPrivilegeRequest) (*model.AllowDbUserPrivilegeResponse, error) {
	requestDef := GenReqDefForAllowDbUserPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AllowDbUserPrivilegeResponse), nil
	}
}

// 创建数据库
//
// 创建数据库。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) CreateDatabase(request *model.CreateDatabaseRequest) (*model.CreateDatabaseResponse, error) {
	requestDef := GenReqDefForCreateDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDatabaseResponse), nil
	}
}

// 创建数据库用户
//
// 创建数据库用户。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) CreateDbUser(request *model.CreateDbUserRequest) (*model.CreateDbUserResponse, error) {
	requestDef := GenReqDefForCreateDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDbUserResponse), nil
	}
}

// 删除数据库
//
// 删除数据库。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) DeleteDatabase(request *model.DeleteDatabaseRequest) (*model.DeleteDatabaseResponse, error) {
	requestDef := GenReqDefForDeleteDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDatabaseResponse), nil
	}
}

// 删除数据库用户
//
// 删除数据库用户。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) DeleteDbUser(request *model.DeleteDbUserRequest) (*model.DeleteDbUserResponse, error) {
	requestDef := GenReqDefForDeleteDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDbUserResponse), nil
	}
}

// 查询指定用户的已授权数据库
//
// 查询指定用户的已授权数据库。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListAuthorizedDatabases(request *model.ListAuthorizedDatabasesRequest) (*model.ListAuthorizedDatabasesResponse, error) {
	requestDef := GenReqDefForListAuthorizedDatabases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAuthorizedDatabasesResponse), nil
	}
}

// 查询指定数据库的已授权用户
//
// 查询指定数据库的已授权用户。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListAuthorizedDbUsers(request *model.ListAuthorizedDbUsersRequest) (*model.ListAuthorizedDbUsersResponse, error) {
	requestDef := GenReqDefForListAuthorizedDbUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAuthorizedDbUsersResponse), nil
	}
}

// 查询数据库列表
//
// 查询数据库列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListDatabases(request *model.ListDatabasesRequest) (*model.ListDatabasesResponse, error) {
	requestDef := GenReqDefForListDatabases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDatabasesResponse), nil
	}
}

// 查询数据库用户列表
//
// 查询数据库用户列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListDbUsers(request *model.ListDbUsersRequest) (*model.ListDbUsersResponse, error) {
	requestDef := GenReqDefForListDbUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDbUsersResponse), nil
	}
}

// 重置数据库密码
//
// 重置数据库密码.
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ResetPwd(request *model.ResetPwdRequest) (*model.ResetPwdResponse, error) {
	requestDef := GenReqDefForResetPwd()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ResetPwdResponse), nil
	}
}

// 解除数据库帐号权限
//
// 解除数据库帐号权限。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) Revoke(request *model.RevokeRequest) (*model.RevokeResponse, error) {
	requestDef := GenReqDefForRevoke()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RevokeResponse), nil
	}
}

// 设置数据库账号密码
//
// 设置数据库账号密码
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) SetDbUserPwd(request *model.SetDbUserPwdRequest) (*model.SetDbUserPwdResponse, error) {
	requestDef := GenReqDefForSetDbUserPwd()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetDbUserPwdResponse), nil
	}
}

// 修改指定实例的数据库备注
//
// 修改指定实例中的数据库备注。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) UpdateDatabase(request *model.UpdateDatabaseRequest) (*model.UpdateDatabaseResponse, error) {
	requestDef := GenReqDefForUpdateDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDatabaseResponse), nil
	}
}

// 授权数据库帐号
//
// 在指定实例的数据库中, 设置帐号的权限。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) AllowDbPrivilege(request *model.AllowDbPrivilegeRequest) (*model.AllowDbPrivilegeResponse, error) {
	requestDef := GenReqDefForAllowDbPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AllowDbPrivilegeResponse), nil
	}
}

// 数据库代理规格变更
//
// 数据库代理实例进行规格变更。
//
// - 调用接口前，您需要了解API 认证鉴权。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ChangeProxyScale(request *model.ChangeProxyScaleRequest) (*model.ChangeProxyScaleResponse, error) {
	requestDef := GenReqDefForChangeProxyScale()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeProxyScaleResponse), nil
	}
}

// 修改读写分离阈值
//
// 修改指定实例的读写分离延时阈值。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ChangeTheDelayThreshold(request *model.ChangeTheDelayThresholdRequest) (*model.ChangeTheDelayThresholdResponse, error) {
	requestDef := GenReqDefForChangeTheDelayThreshold()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeTheDelayThresholdResponse), nil
	}
}

// 创建数据库
//
// 在指定实例中创建数据库。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) CreatePostgresqlDatabase(request *model.CreatePostgresqlDatabaseRequest) (*model.CreatePostgresqlDatabaseResponse, error) {
	requestDef := GenReqDefForCreatePostgresqlDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePostgresqlDatabaseResponse), nil
	}
}

// 创建数据库SCHEMA
//
// 在指定实例的数据库中, 创建数据库schema。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) CreatePostgresqlDatabaseSchema(request *model.CreatePostgresqlDatabaseSchemaRequest) (*model.CreatePostgresqlDatabaseSchemaResponse, error) {
	requestDef := GenReqDefForCreatePostgresqlDatabaseSchema()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePostgresqlDatabaseSchemaResponse), nil
	}
}

// 创建数据库用户
//
// 在指定实例中创建数据库用户。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) CreatePostgresqlDbUser(request *model.CreatePostgresqlDbUserRequest) (*model.CreatePostgresqlDbUserResponse, error) {
	requestDef := GenReqDefForCreatePostgresqlDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePostgresqlDbUserResponse), nil
	}
}

// 查询数据库SCHEMA列表
//
// 查询指定实例的数据库SCHEMA列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListPostgresqlDatabaseSchemas(request *model.ListPostgresqlDatabaseSchemasRequest) (*model.ListPostgresqlDatabaseSchemasResponse, error) {
	requestDef := GenReqDefForListPostgresqlDatabaseSchemas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlDatabaseSchemasResponse), nil
	}
}

// 查询数据库列表
//
// 查询指定实例中的数据库列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListPostgresqlDatabases(request *model.ListPostgresqlDatabasesRequest) (*model.ListPostgresqlDatabasesResponse, error) {
	requestDef := GenReqDefForListPostgresqlDatabases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlDatabasesResponse), nil
	}
}

// 查询数据库用户列表
//
// 在指定实例中查询数据库用户列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListPostgresqlDbUserPaginated(request *model.ListPostgresqlDbUserPaginatedRequest) (*model.ListPostgresqlDbUserPaginatedResponse, error) {
	requestDef := GenReqDefForListPostgresqlDbUserPaginated()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlDbUserPaginatedResponse), nil
	}
}

// 查询数据库代理可变更的规格
//
// 查询数据库代理可变更的规格信息。
//
// - 调用接口前，您需要了解API 认证鉴权。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) SearchQueryScaleComputeFlavors(request *model.SearchQueryScaleComputeFlavorsRequest) (*model.SearchQueryScaleComputeFlavorsResponse, error) {
	requestDef := GenReqDefForSearchQueryScaleComputeFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SearchQueryScaleComputeFlavorsResponse), nil
	}
}

// 查询数据库代理可变更的规格
//
// 查询数据库代理可变更的规格信息。
//
// - 调用接口前，您需要了解API 认证鉴权。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) SearchQueryScaleFlavors(request *model.SearchQueryScaleFlavorsRequest) (*model.SearchQueryScaleFlavorsResponse, error) {
	requestDef := GenReqDefForSearchQueryScaleFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SearchQueryScaleFlavorsResponse), nil
	}
}

// 重置数据库帐号密码
//
// 重置指定数据库帐号的密码。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) SetPostgresqlDbUserPwd(request *model.SetPostgresqlDbUserPwdRequest) (*model.SetPostgresqlDbUserPwdResponse, error) {
	requestDef := GenReqDefForSetPostgresqlDbUserPwd()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetPostgresqlDbUserPwdResponse), nil
	}
}

// 查询数据库代理信息
//
// 查询指定实例的数据库代理详细信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ShowInformationAboutDatabaseProxy(request *model.ShowInformationAboutDatabaseProxyRequest) (*model.ShowInformationAboutDatabaseProxyResponse, error) {
	requestDef := GenReqDefForShowInformationAboutDatabaseProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowInformationAboutDatabaseProxyResponse), nil
	}
}

// 开启数据库代理
//
// 为指定实例开启数据库代理。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) StartDatabaseProxy(request *model.StartDatabaseProxyRequest) (*model.StartDatabaseProxyResponse, error) {
	requestDef := GenReqDefForStartDatabaseProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartDatabaseProxyResponse), nil
	}
}

// 关闭数据库代理
//
// 为指定实例关闭数据库代理。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) StopDatabaseProxy(request *model.StopDatabaseProxyRequest) (*model.StopDatabaseProxyResponse, error) {
	requestDef := GenReqDefForStopDatabaseProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopDatabaseProxyResponse), nil
	}
}

// 修改读写分离权重
//
// 修改指定实例的读写分离权重。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) UpdateReadWeight(request *model.UpdateReadWeightRequest) (*model.UpdateReadWeightResponse, error) {
	requestDef := GenReqDefForUpdateReadWeight()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateReadWeightResponse), nil
	}
}

// 授权数据库帐号
//
// 授权数据库帐号。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) AllowSqlserverDbUserPrivilege(request *model.AllowSqlserverDbUserPrivilegeRequest) (*model.AllowSqlserverDbUserPrivilegeResponse, error) {
	requestDef := GenReqDefForAllowSqlserverDbUserPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AllowSqlserverDbUserPrivilegeResponse), nil
	}
}

// 创建数据库
//
// 创建数据库。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) CreateSqlserverDatabase(request *model.CreateSqlserverDatabaseRequest) (*model.CreateSqlserverDatabaseResponse, error) {
	requestDef := GenReqDefForCreateSqlserverDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSqlserverDatabaseResponse), nil
	}
}

// 创建数据库用户
//
// 创建数据库用户。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) CreateSqlserverDbUser(request *model.CreateSqlserverDbUserRequest) (*model.CreateSqlserverDbUserResponse, error) {
	requestDef := GenReqDefForCreateSqlserverDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSqlserverDbUserResponse), nil
	}
}

// 删除数据库
//
// 删除数据库。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) DeleteSqlserverDatabase(request *model.DeleteSqlserverDatabaseRequest) (*model.DeleteSqlserverDatabaseResponse, error) {
	requestDef := GenReqDefForDeleteSqlserverDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSqlserverDatabaseResponse), nil
	}
}

// 删除数据库
//
// 删除数据库。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) DeleteSqlserverDatabaseEx(request *model.DeleteSqlserverDatabaseExRequest) (*model.DeleteSqlserverDatabaseExResponse, error) {
	requestDef := GenReqDefForDeleteSqlserverDatabaseEx()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSqlserverDatabaseExResponse), nil
	}
}

// 删除数据库用户
//
// 删除数据库用户。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) DeleteSqlserverDbUser(request *model.DeleteSqlserverDbUserRequest) (*model.DeleteSqlserverDbUserResponse, error) {
	requestDef := GenReqDefForDeleteSqlserverDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSqlserverDbUserResponse), nil
	}
}

// 查询指定数据库的已授权用户
//
// 查询指定数据库的已授权用户。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListAuthorizedSqlserverDbUsers(request *model.ListAuthorizedSqlserverDbUsersRequest) (*model.ListAuthorizedSqlserverDbUsersResponse, error) {
	requestDef := GenReqDefForListAuthorizedSqlserverDbUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAuthorizedSqlserverDbUsersResponse), nil
	}
}

// 查询数据库列表
//
// 查询数据库列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListSqlserverDatabases(request *model.ListSqlserverDatabasesRequest) (*model.ListSqlserverDatabasesResponse, error) {
	requestDef := GenReqDefForListSqlserverDatabases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSqlserverDatabasesResponse), nil
	}
}

// 查询数据库用户列表
//
// 查询数据库用户列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) ListSqlserverDbUsers(request *model.ListSqlserverDbUsersRequest) (*model.ListSqlserverDbUsersResponse, error) {
	requestDef := GenReqDefForListSqlserverDbUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSqlserverDbUsersResponse), nil
	}
}

// 解除数据库帐号权限
//
// 解除数据库帐号权限。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *RdsClient) RevokeSqlserverDbUserPrivilege(request *model.RevokeSqlserverDbUserPrivilegeRequest) (*model.RevokeSqlserverDbUserPrivilegeResponse, error) {
	requestDef := GenReqDefForRevokeSqlserverDbUserPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RevokeSqlserverDbUserPrivilegeResponse), nil
	}
}
