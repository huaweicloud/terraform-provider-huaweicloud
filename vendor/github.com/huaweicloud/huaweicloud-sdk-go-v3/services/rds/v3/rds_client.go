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

//绑定和解绑弹性公网IP。
func (c *RdsClient) AttachEip(request *model.AttachEipRequest) (*model.AttachEipResponse, error) {
	requestDef := GenReqDefForAttachEip()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AttachEipResponse), nil
	}
}

//批量添加标签。
func (c *RdsClient) BatchTagAddAction(request *model.BatchTagAddActionRequest) (*model.BatchTagAddActionResponse, error) {
	requestDef := GenReqDefForBatchTagAddAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchTagAddActionResponse), nil
	}
}

//批量删除标签。
func (c *RdsClient) BatchTagDelAction(request *model.BatchTagDelActionRequest) (*model.BatchTagDelActionResponse, error) {
	requestDef := GenReqDefForBatchTagDelAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchTagDelActionResponse), nil
	}
}

//更改主备实例的数据同步方式。
func (c *RdsClient) ChangeFailoverMode(request *model.ChangeFailoverModeRequest) (*model.ChangeFailoverModeResponse, error) {
	requestDef := GenReqDefForChangeFailoverMode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeFailoverModeResponse), nil
	}
}

//切换主备实例的倒换策略.
func (c *RdsClient) ChangeFailoverStrategy(request *model.ChangeFailoverStrategyRequest) (*model.ChangeFailoverStrategyResponse, error) {
	requestDef := GenReqDefForChangeFailoverStrategy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeFailoverStrategyResponse), nil
	}
}

//设置可维护时间段
func (c *RdsClient) ChangeOpsWindow(request *model.ChangeOpsWindowRequest) (*model.ChangeOpsWindowResponse, error) {
	requestDef := GenReqDefForChangeOpsWindow()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeOpsWindowResponse), nil
	}
}

//创建参数模板。
func (c *RdsClient) CreateConfiguration(request *model.CreateConfigurationRequest) (*model.CreateConfigurationResponse, error) {
	requestDef := GenReqDefForCreateConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateConfigurationResponse), nil
	}
}

//申请域名
func (c *RdsClient) CreateDnsName(request *model.CreateDnsNameRequest) (*model.CreateDnsNameResponse, error) {
	requestDef := GenReqDefForCreateDnsName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDnsNameResponse), nil
	}
}

//创建数据库实例。
func (c *RdsClient) CreateInstance(request *model.CreateInstanceRequest) (*model.CreateInstanceResponse, error) {
	requestDef := GenReqDefForCreateInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateInstanceResponse), nil
	}
}

//创建手动备份。
func (c *RdsClient) CreateManualBackup(request *model.CreateManualBackupRequest) (*model.CreateManualBackupResponse, error) {
	requestDef := GenReqDefForCreateManualBackup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateManualBackupResponse), nil
	}
}

//恢复到新实例。
func (c *RdsClient) CreateRestoreInstance(request *model.CreateRestoreInstanceRequest) (*model.CreateRestoreInstanceResponse, error) {
	requestDef := GenReqDefForCreateRestoreInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRestoreInstanceResponse), nil
	}
}

//删除参数模板。
func (c *RdsClient) DeleteConfiguration(request *model.DeleteConfigurationRequest) (*model.DeleteConfigurationResponse, error) {
	requestDef := GenReqDefForDeleteConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteConfigurationResponse), nil
	}
}

//删除数据库实例。
func (c *RdsClient) DeleteInstance(request *model.DeleteInstanceRequest) (*model.DeleteInstanceResponse, error) {
	requestDef := GenReqDefForDeleteInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteInstanceResponse), nil
	}
}

//删除手动备份。
func (c *RdsClient) DeleteManualBackup(request *model.DeleteManualBackupRequest) (*model.DeleteManualBackupResponse, error) {
	requestDef := GenReqDefForDeleteManualBackup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteManualBackupResponse), nil
	}
}

//获取慢日志下载链接。
func (c *RdsClient) DownloadSlowlog(request *model.DownloadSlowlogRequest) (*model.DownloadSlowlogResponse, error) {
	requestDef := GenReqDefForDownloadSlowlog()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DownloadSlowlogResponse), nil
	}
}

//应用参数模板。
func (c *RdsClient) EnableConfiguration(request *model.EnableConfigurationRequest) (*model.EnableConfigurationResponse, error) {
	requestDef := GenReqDefForEnableConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.EnableConfigurationResponse), nil
	}
}

//获取审计日志列表。
func (c *RdsClient) ListAuditlogs(request *model.ListAuditlogsRequest) (*model.ListAuditlogsResponse, error) {
	requestDef := GenReqDefForListAuditlogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAuditlogsResponse), nil
	}
}

//获取备份列表。
func (c *RdsClient) ListBackups(request *model.ListBackupsRequest) (*model.ListBackupsResponse, error) {
	requestDef := GenReqDefForListBackups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListBackupsResponse), nil
	}
}

//查询SQLServer可用字符集
func (c *RdsClient) ListCollations(request *model.ListCollationsRequest) (*model.ListCollationsResponse, error) {
	requestDef := GenReqDefForListCollations()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListCollationsResponse), nil
	}
}

//获取参数模板列表，包括所有数据库的默认参数模板和用户创建的参数模板。
func (c *RdsClient) ListConfigurations(request *model.ListConfigurationsRequest) (*model.ListConfigurationsResponse, error) {
	requestDef := GenReqDefForListConfigurations()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListConfigurationsResponse), nil
	}
}

//查询数据库引擎的版本。
func (c *RdsClient) ListDatastores(request *model.ListDatastoresRequest) (*model.ListDatastoresResponse, error) {
	requestDef := GenReqDefForListDatastores()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDatastoresResponse), nil
	}
}

//查询数据库错误日志。
func (c *RdsClient) ListErrorLogs(request *model.ListErrorLogsRequest) (*model.ListErrorLogsResponse, error) {
	requestDef := GenReqDefForListErrorLogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListErrorLogsResponse), nil
	}
}

//查询数据库错误日志。(与原v3接口相比修改offset,符合华为云服务开放 API遵从性规范3.0)
func (c *RdsClient) ListErrorLogsNew(request *model.ListErrorLogsNewRequest) (*model.ListErrorLogsNewResponse, error) {
	requestDef := GenReqDefForListErrorLogsNew()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListErrorLogsNewResponse), nil
	}
}

//查询数据库规格。
func (c *RdsClient) ListFlavors(request *model.ListFlavorsRequest) (*model.ListFlavorsResponse, error) {
	requestDef := GenReqDefForListFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListFlavorsResponse), nil
	}
}

//查询数据库实例列表。
func (c *RdsClient) ListInstances(request *model.ListInstancesRequest) (*model.ListInstancesResponse, error) {
	requestDef := GenReqDefForListInstances()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstancesResponse), nil
	}
}

//获取指定ID的任务信息。
func (c *RdsClient) ListJobInfo(request *model.ListJobInfoRequest) (*model.ListJobInfoResponse, error) {
	requestDef := GenReqDefForListJobInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListJobInfoResponse), nil
	}
}

//获取指定实例和时间范围的任务信息（SQL Server）。
func (c *RdsClient) ListJobInfoDetail(request *model.ListJobInfoDetailRequest) (*model.ListJobInfoDetailResponse, error) {
	requestDef := GenReqDefForListJobInfoDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListJobInfoDetailResponse), nil
	}
}

//查询跨区域备份列表。
func (c *RdsClient) ListOffSiteBackups(request *model.ListOffSiteBackupsRequest) (*model.ListOffSiteBackupsResponse, error) {
	requestDef := GenReqDefForListOffSiteBackups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListOffSiteBackupsResponse), nil
	}
}

//查询跨区域备份实例列表。
func (c *RdsClient) ListOffSiteInstances(request *model.ListOffSiteInstancesRequest) (*model.ListOffSiteInstancesResponse, error) {
	requestDef := GenReqDefForListOffSiteInstances()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListOffSiteInstancesResponse), nil
	}
}

//查询跨区域备份可恢复时间段。 如果您备份策略中的保存天数设置较长，建议您传入查询日期“date”。
func (c *RdsClient) ListOffSiteRestoreTimes(request *model.ListOffSiteRestoreTimesRequest) (*model.ListOffSiteRestoreTimesResponse, error) {
	requestDef := GenReqDefForListOffSiteRestoreTimes()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListOffSiteRestoreTimesResponse), nil
	}
}

//查询项目标签。
func (c *RdsClient) ListProjectTags(request *model.ListProjectTagsRequest) (*model.ListProjectTagsResponse, error) {
	requestDef := GenReqDefForListProjectTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProjectTagsResponse), nil
	}
}

//查询可恢复时间段。 如果您备份策略中的保存天数设置较长，建议您传入查询日期“date”。
func (c *RdsClient) ListRestoreTimes(request *model.ListRestoreTimesRequest) (*model.ListRestoreTimesResponse, error) {
	requestDef := GenReqDefForListRestoreTimes()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRestoreTimesResponse), nil
	}
}

//查询慢日志文件列表。 调用该接口取到慢日志文件名后，可以调用接口/v3/{project_id}/instances/{instance_id}/slowlog-download 获取慢日志文件下载链接
func (c *RdsClient) ListSlowLogFile(request *model.ListSlowLogFileRequest) (*model.ListSlowLogFileResponse, error) {
	requestDef := GenReqDefForListSlowLogFile()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowLogFileResponse), nil
	}
}

//查询数据库慢日志。
func (c *RdsClient) ListSlowLogs(request *model.ListSlowLogsRequest) (*model.ListSlowLogsResponse, error) {
	requestDef := GenReqDefForListSlowLogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowLogsResponse), nil
	}
}

//查询数据库慢日志。(与原v3接口相比修改offset,符合华为云服务开放 API遵从性规范3.0)
func (c *RdsClient) ListSlowLogsNew(request *model.ListSlowLogsNewRequest) (*model.ListSlowLogsNewResponse, error) {
	requestDef := GenReqDefForListSlowLogsNew()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowLogsNewResponse), nil
	}
}

//获取慢日志统计信息
func (c *RdsClient) ListSlowlogStatistics(request *model.ListSlowlogStatisticsRequest) (*model.ListSlowlogStatisticsResponse, error) {
	requestDef := GenReqDefForListSlowlogStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowlogStatisticsResponse), nil
	}
}

//查询数据库磁盘类型。
func (c *RdsClient) ListStorageTypes(request *model.ListStorageTypesRequest) (*model.ListStorageTypesResponse, error) {
	requestDef := GenReqDefForListStorageTypes()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListStorageTypesResponse), nil
	}
}

//迁移主备实例的备机
func (c *RdsClient) MigrateFollower(request *model.MigrateFollowerRequest) (*model.MigrateFollowerResponse, error) {
	requestDef := GenReqDefForMigrateFollower()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.MigrateFollowerResponse), nil
	}
}

//恢复到已有实例。
func (c *RdsClient) RestoreExistInstance(request *model.RestoreExistInstanceRequest) (*model.RestoreExistInstanceResponse, error) {
	requestDef := GenReqDefForRestoreExistInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RestoreExistInstanceResponse), nil
	}
}

//表级时间点恢复(MySQL)。
func (c *RdsClient) RestoreTables(request *model.RestoreTablesRequest) (*model.RestoreTablesResponse, error) {
	requestDef := GenReqDefForRestoreTables()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RestoreTablesResponse), nil
	}
}

//恢复到已有实例。
func (c *RdsClient) RestoreToExistingInstance(request *model.RestoreToExistingInstanceRequest) (*model.RestoreToExistingInstanceResponse, error) {
	requestDef := GenReqDefForRestoreToExistingInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RestoreToExistingInstanceResponse), nil
	}
}

//设置审计日志策略。
func (c *RdsClient) SetAuditlogPolicy(request *model.SetAuditlogPolicyRequest) (*model.SetAuditlogPolicyResponse, error) {
	requestDef := GenReqDefForSetAuditlogPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetAuditlogPolicyResponse), nil
	}
}

//设置自动备份策略。
func (c *RdsClient) SetBackupPolicy(request *model.SetBackupPolicyRequest) (*model.SetBackupPolicyResponse, error) {
	requestDef := GenReqDefForSetBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetBackupPolicyResponse), nil
	}
}

//修改指定实例的binlog本地保留时长。
func (c *RdsClient) SetBinlogClearPolicy(request *model.SetBinlogClearPolicyRequest) (*model.SetBinlogClearPolicyResponse, error) {
	requestDef := GenReqDefForSetBinlogClearPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetBinlogClearPolicyResponse), nil
	}
}

//设置跨区域备份策略。
func (c *RdsClient) SetOffSiteBackupPolicy(request *model.SetOffSiteBackupPolicyRequest) (*model.SetOffSiteBackupPolicyResponse, error) {
	requestDef := GenReqDefForSetOffSiteBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetOffSiteBackupPolicyResponse), nil
	}
}

//修改安全组
func (c *RdsClient) SetSecurityGroup(request *model.SetSecurityGroupRequest) (*model.SetSecurityGroupResponse, error) {
	requestDef := GenReqDefForSetSecurityGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetSecurityGroupResponse), nil
	}
}

//生成审计日志下载链接。
func (c *RdsClient) ShowAuditlogDownloadLink(request *model.ShowAuditlogDownloadLinkRequest) (*model.ShowAuditlogDownloadLinkResponse, error) {
	requestDef := GenReqDefForShowAuditlogDownloadLink()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAuditlogDownloadLinkResponse), nil
	}
}

//查询审计日志策略。
func (c *RdsClient) ShowAuditlogPolicy(request *model.ShowAuditlogPolicyRequest) (*model.ShowAuditlogPolicyResponse, error) {
	requestDef := GenReqDefForShowAuditlogPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAuditlogPolicyResponse), nil
	}
}

//获取备份下载链接。
func (c *RdsClient) ShowBackupDownloadLink(request *model.ShowBackupDownloadLinkRequest) (*model.ShowBackupDownloadLinkResponse, error) {
	requestDef := GenReqDefForShowBackupDownloadLink()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBackupDownloadLinkResponse), nil
	}
}

//查询自动备份策略。
func (c *RdsClient) ShowBackupPolicy(request *model.ShowBackupPolicyRequest) (*model.ShowBackupPolicyResponse, error) {
	requestDef := GenReqDefForShowBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBackupPolicyResponse), nil
	}
}

//查寻指定实例的binlog本地保留时长。
func (c *RdsClient) ShowBinlogClearPolicy(request *model.ShowBinlogClearPolicyRequest) (*model.ShowBinlogClearPolicyResponse, error) {
	requestDef := GenReqDefForShowBinlogClearPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBinlogClearPolicyResponse), nil
	}
}

//获取指定参数模板的参数。
func (c *RdsClient) ShowConfiguration(request *model.ShowConfigurationRequest) (*model.ShowConfigurationResponse, error) {
	requestDef := GenReqDefForShowConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowConfigurationResponse), nil
	}
}

//建立跨云容灾关系后，查询主实例和灾备实例间的复制状态及延迟。
func (c *RdsClient) ShowDrReplicaStatus(request *model.ShowDrReplicaStatusRequest) (*model.ShowDrReplicaStatusResponse, error) {
	requestDef := GenReqDefForShowDrReplicaStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDrReplicaStatusResponse), nil
	}
}

//获取指定实例的参数模板。
func (c *RdsClient) ShowInstanceConfiguration(request *model.ShowInstanceConfigurationRequest) (*model.ShowInstanceConfigurationResponse, error) {
	requestDef := GenReqDefForShowInstanceConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowInstanceConfigurationResponse), nil
	}
}

//查询跨区域备份策略。
func (c *RdsClient) ShowOffSiteBackupPolicy(request *model.ShowOffSiteBackupPolicyRequest) (*model.ShowOffSiteBackupPolicyResponse, error) {
	requestDef := GenReqDefForShowOffSiteBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowOffSiteBackupPolicyResponse), nil
	}
}

//查询当前项目下资源配额情况。
func (c *RdsClient) ShowQuotas(request *model.ShowQuotasRequest) (*model.ShowQuotasResponse, error) {
	requestDef := GenReqDefForShowQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowQuotasResponse), nil
	}
}

//手动倒换主备.
func (c *RdsClient) StartFailover(request *model.StartFailoverRequest) (*model.StartFailoverResponse, error) {
	requestDef := GenReqDefForStartFailover()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartFailoverResponse), nil
	}
}

//扩容数据库实例的磁盘空间。
func (c *RdsClient) StartInstanceEnlargeVolumeAction(request *model.StartInstanceEnlargeVolumeActionRequest) (*model.StartInstanceEnlargeVolumeActionResponse, error) {
	requestDef := GenReqDefForStartInstanceEnlargeVolumeAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartInstanceEnlargeVolumeActionResponse), nil
	}
}

//重启数据库实例。
func (c *RdsClient) StartInstanceRestartAction(request *model.StartInstanceRestartActionRequest) (*model.StartInstanceRestartActionResponse, error) {
	requestDef := GenReqDefForStartInstanceRestartAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartInstanceRestartActionResponse), nil
	}
}

//单机转主备实例。
func (c *RdsClient) StartInstanceSingleToHaAction(request *model.StartInstanceSingleToHaActionRequest) (*model.StartInstanceSingleToHaActionResponse, error) {
	requestDef := GenReqDefForStartInstanceSingleToHaAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartInstanceSingleToHaActionResponse), nil
	}
}

//设置回收站策略。
func (c *RdsClient) StartRecyclePolicy(request *model.StartRecyclePolicyRequest) (*model.StartRecyclePolicyResponse, error) {
	requestDef := GenReqDefForStartRecyclePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartRecyclePolicyResponse), nil
	}
}

//变更数据库实例的规格。
func (c *RdsClient) StartResizeFlavorAction(request *model.StartResizeFlavorActionRequest) (*model.StartResizeFlavorActionResponse, error) {
	requestDef := GenReqDefForStartResizeFlavorAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartResizeFlavorActionResponse), nil
	}
}

//停止实例以节省费用，在停止数据库实例后，支持手动重新开启实例。
func (c *RdsClient) StartupInstance(request *model.StartupInstanceRequest) (*model.StartupInstanceResponse, error) {
	requestDef := GenReqDefForStartupInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartupInstanceResponse), nil
	}
}

//实例进行关机，通过暂时停止按需实例以节省费用，实例默认停止七天。
func (c *RdsClient) StopInstance(request *model.StopInstanceRequest) (*model.StopInstanceResponse, error) {
	requestDef := GenReqDefForStopInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopInstanceResponse), nil
	}
}

//设置SSL数据加密。
func (c *RdsClient) SwitchSsl(request *model.SwitchSslRequest) (*model.SwitchSslResponse, error) {
	requestDef := GenReqDefForSwitchSsl()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SwitchSslResponse), nil
	}
}

//修改参数模板参数。
func (c *RdsClient) UpdateConfiguration(request *model.UpdateConfigurationRequest) (*model.UpdateConfigurationResponse, error) {
	requestDef := GenReqDefForUpdateConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateConfigurationResponse), nil
	}
}

//修改内网地址
func (c *RdsClient) UpdateDataIp(request *model.UpdateDataIpRequest) (*model.UpdateDataIpResponse, error) {
	requestDef := GenReqDefForUpdateDataIp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDataIpResponse), nil
	}
}

//修改域名
func (c *RdsClient) UpdateDnsName(request *model.UpdateDnsNameRequest) (*model.UpdateDnsNameResponse, error) {
	requestDef := GenReqDefForUpdateDnsName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDnsNameResponse), nil
	}
}

//修改指定实例的参数。
func (c *RdsClient) UpdateInstanceConfiguration(request *model.UpdateInstanceConfigurationRequest) (*model.UpdateInstanceConfigurationResponse, error) {
	requestDef := GenReqDefForUpdateInstanceConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateInstanceConfigurationResponse), nil
	}
}

//修改实例名称。
func (c *RdsClient) UpdateInstanceName(request *model.UpdateInstanceNameRequest) (*model.UpdateInstanceNameResponse, error) {
	requestDef := GenReqDefForUpdateInstanceName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateInstanceNameResponse), nil
	}
}

//修改数据库端口
func (c *RdsClient) UpdatePort(request *model.UpdatePortRequest) (*model.UpdatePortResponse, error) {
	requestDef := GenReqDefForUpdatePort()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePortResponse), nil
	}
}

//修改指定数据库实例的备注信息。
func (c *RdsClient) UpdatePostgresqlInstanceAlias(request *model.UpdatePostgresqlInstanceAliasRequest) (*model.UpdatePostgresqlInstanceAliasResponse, error) {
	requestDef := GenReqDefForUpdatePostgresqlInstanceAlias()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePostgresqlInstanceAliasResponse), nil
	}
}

//对实例进行小版本升级。
func (c *RdsClient) UpgradeDbVersion(request *model.UpgradeDbVersionRequest) (*model.UpgradeDbVersionResponse, error) {
	requestDef := GenReqDefForUpgradeDbVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeDbVersionResponse), nil
	}
}

//查询API版本列表。
func (c *RdsClient) ListApiVersion(request *model.ListApiVersionRequest) (*model.ListApiVersionResponse, error) {
	requestDef := GenReqDefForListApiVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListApiVersionResponse), nil
	}
}

//查询API版本列表。
func (c *RdsClient) ListApiVersionNew(request *model.ListApiVersionNewRequest) (*model.ListApiVersionNewResponse, error) {
	requestDef := GenReqDefForListApiVersionNew()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListApiVersionNewResponse), nil
	}
}

//查询指定的API版本信息。
func (c *RdsClient) ShowApiVersion(request *model.ShowApiVersionRequest) (*model.ShowApiVersionResponse, error) {
	requestDef := GenReqDefForShowApiVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowApiVersionResponse), nil
	}
}

//授权数据库帐号。
func (c *RdsClient) AllowDbUserPrivilege(request *model.AllowDbUserPrivilegeRequest) (*model.AllowDbUserPrivilegeResponse, error) {
	requestDef := GenReqDefForAllowDbUserPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AllowDbUserPrivilegeResponse), nil
	}
}

//创建数据库。
func (c *RdsClient) CreateDatabase(request *model.CreateDatabaseRequest) (*model.CreateDatabaseResponse, error) {
	requestDef := GenReqDefForCreateDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDatabaseResponse), nil
	}
}

//创建数据库用户。
func (c *RdsClient) CreateDbUser(request *model.CreateDbUserRequest) (*model.CreateDbUserResponse, error) {
	requestDef := GenReqDefForCreateDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDbUserResponse), nil
	}
}

//删除数据库。
func (c *RdsClient) DeleteDatabase(request *model.DeleteDatabaseRequest) (*model.DeleteDatabaseResponse, error) {
	requestDef := GenReqDefForDeleteDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDatabaseResponse), nil
	}
}

//删除数据库用户。
func (c *RdsClient) DeleteDbUser(request *model.DeleteDbUserRequest) (*model.DeleteDbUserResponse, error) {
	requestDef := GenReqDefForDeleteDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDbUserResponse), nil
	}
}

//查询指定用户的已授权数据库。
func (c *RdsClient) ListAuthorizedDatabases(request *model.ListAuthorizedDatabasesRequest) (*model.ListAuthorizedDatabasesResponse, error) {
	requestDef := GenReqDefForListAuthorizedDatabases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAuthorizedDatabasesResponse), nil
	}
}

//查询指定数据库的已授权用户。
func (c *RdsClient) ListAuthorizedDbUsers(request *model.ListAuthorizedDbUsersRequest) (*model.ListAuthorizedDbUsersResponse, error) {
	requestDef := GenReqDefForListAuthorizedDbUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAuthorizedDbUsersResponse), nil
	}
}

//查询数据库列表。
func (c *RdsClient) ListDatabases(request *model.ListDatabasesRequest) (*model.ListDatabasesResponse, error) {
	requestDef := GenReqDefForListDatabases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDatabasesResponse), nil
	}
}

//查询数据库用户列表。
func (c *RdsClient) ListDbUsers(request *model.ListDbUsersRequest) (*model.ListDbUsersResponse, error) {
	requestDef := GenReqDefForListDbUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDbUsersResponse), nil
	}
}

//重置数据库密码.
func (c *RdsClient) ResetPwd(request *model.ResetPwdRequest) (*model.ResetPwdResponse, error) {
	requestDef := GenReqDefForResetPwd()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ResetPwdResponse), nil
	}
}

//解除数据库帐号权限。
func (c *RdsClient) Revoke(request *model.RevokeRequest) (*model.RevokeResponse, error) {
	requestDef := GenReqDefForRevoke()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RevokeResponse), nil
	}
}

//设置数据库账号密码
func (c *RdsClient) SetDbUserPwd(request *model.SetDbUserPwdRequest) (*model.SetDbUserPwdResponse, error) {
	requestDef := GenReqDefForSetDbUserPwd()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetDbUserPwdResponse), nil
	}
}

//修改指定实例中的数据库备注。
func (c *RdsClient) UpdateDatabase(request *model.UpdateDatabaseRequest) (*model.UpdateDatabaseResponse, error) {
	requestDef := GenReqDefForUpdateDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDatabaseResponse), nil
	}
}

//在指定实例的数据库中, 设置帐号的权限。
func (c *RdsClient) AllowDbPrivilege(request *model.AllowDbPrivilegeRequest) (*model.AllowDbPrivilegeResponse, error) {
	requestDef := GenReqDefForAllowDbPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AllowDbPrivilegeResponse), nil
	}
}

//数据库代理实例进行规格变更。  - 调用接口前，您需要了解API 认证鉴权。
func (c *RdsClient) ChangeProxyScale(request *model.ChangeProxyScaleRequest) (*model.ChangeProxyScaleResponse, error) {
	requestDef := GenReqDefForChangeProxyScale()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeProxyScaleResponse), nil
	}
}

//修改指定实例的读写分离延时阈值。
func (c *RdsClient) ChangeTheDelayThreshold(request *model.ChangeTheDelayThresholdRequest) (*model.ChangeTheDelayThresholdResponse, error) {
	requestDef := GenReqDefForChangeTheDelayThreshold()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeTheDelayThresholdResponse), nil
	}
}

//在指定实例中创建数据库。
func (c *RdsClient) CreatePostgresqlDatabase(request *model.CreatePostgresqlDatabaseRequest) (*model.CreatePostgresqlDatabaseResponse, error) {
	requestDef := GenReqDefForCreatePostgresqlDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePostgresqlDatabaseResponse), nil
	}
}

//在指定实例的数据库中, 创建数据库schema。
func (c *RdsClient) CreatePostgresqlDatabaseSchema(request *model.CreatePostgresqlDatabaseSchemaRequest) (*model.CreatePostgresqlDatabaseSchemaResponse, error) {
	requestDef := GenReqDefForCreatePostgresqlDatabaseSchema()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePostgresqlDatabaseSchemaResponse), nil
	}
}

//在指定实例中创建数据库用户。
func (c *RdsClient) CreatePostgresqlDbUser(request *model.CreatePostgresqlDbUserRequest) (*model.CreatePostgresqlDbUserResponse, error) {
	requestDef := GenReqDefForCreatePostgresqlDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePostgresqlDbUserResponse), nil
	}
}

//查询指定实例的数据库SCHEMA列表。
func (c *RdsClient) ListPostgresqlDatabaseSchemas(request *model.ListPostgresqlDatabaseSchemasRequest) (*model.ListPostgresqlDatabaseSchemasResponse, error) {
	requestDef := GenReqDefForListPostgresqlDatabaseSchemas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlDatabaseSchemasResponse), nil
	}
}

//查询指定实例中的数据库列表。
func (c *RdsClient) ListPostgresqlDatabases(request *model.ListPostgresqlDatabasesRequest) (*model.ListPostgresqlDatabasesResponse, error) {
	requestDef := GenReqDefForListPostgresqlDatabases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlDatabasesResponse), nil
	}
}

//在指定实例中查询数据库用户列表。
func (c *RdsClient) ListPostgresqlDbUserPaginated(request *model.ListPostgresqlDbUserPaginatedRequest) (*model.ListPostgresqlDbUserPaginatedResponse, error) {
	requestDef := GenReqDefForListPostgresqlDbUserPaginated()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlDbUserPaginatedResponse), nil
	}
}

//查询数据库代理可变更的规格信息。  - 调用接口前，您需要了解API 认证鉴权。
func (c *RdsClient) SearchQueryScaleComputeFlavors(request *model.SearchQueryScaleComputeFlavorsRequest) (*model.SearchQueryScaleComputeFlavorsResponse, error) {
	requestDef := GenReqDefForSearchQueryScaleComputeFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SearchQueryScaleComputeFlavorsResponse), nil
	}
}

//查询数据库代理可变更的规格信息。  - 调用接口前，您需要了解API 认证鉴权。
func (c *RdsClient) SearchQueryScaleFlavors(request *model.SearchQueryScaleFlavorsRequest) (*model.SearchQueryScaleFlavorsResponse, error) {
	requestDef := GenReqDefForSearchQueryScaleFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SearchQueryScaleFlavorsResponse), nil
	}
}

//重置指定数据库帐号的密码。
func (c *RdsClient) SetPostgresqlDbUserPwd(request *model.SetPostgresqlDbUserPwdRequest) (*model.SetPostgresqlDbUserPwdResponse, error) {
	requestDef := GenReqDefForSetPostgresqlDbUserPwd()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetPostgresqlDbUserPwdResponse), nil
	}
}

//查询指定实例的数据库代理详细信息。
func (c *RdsClient) ShowInformationAboutDatabaseProxy(request *model.ShowInformationAboutDatabaseProxyRequest) (*model.ShowInformationAboutDatabaseProxyResponse, error) {
	requestDef := GenReqDefForShowInformationAboutDatabaseProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowInformationAboutDatabaseProxyResponse), nil
	}
}

//为指定实例开启数据库代理。
func (c *RdsClient) StartDatabaseProxy(request *model.StartDatabaseProxyRequest) (*model.StartDatabaseProxyResponse, error) {
	requestDef := GenReqDefForStartDatabaseProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartDatabaseProxyResponse), nil
	}
}

//为指定实例关闭数据库代理。
func (c *RdsClient) StopDatabaseProxy(request *model.StopDatabaseProxyRequest) (*model.StopDatabaseProxyResponse, error) {
	requestDef := GenReqDefForStopDatabaseProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopDatabaseProxyResponse), nil
	}
}

//修改指定实例的读写分离权重。
func (c *RdsClient) UpdateReadWeight(request *model.UpdateReadWeightRequest) (*model.UpdateReadWeightResponse, error) {
	requestDef := GenReqDefForUpdateReadWeight()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateReadWeightResponse), nil
	}
}

//授权数据库帐号。
func (c *RdsClient) AllowSqlserverDbUserPrivilege(request *model.AllowSqlserverDbUserPrivilegeRequest) (*model.AllowSqlserverDbUserPrivilegeResponse, error) {
	requestDef := GenReqDefForAllowSqlserverDbUserPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AllowSqlserverDbUserPrivilegeResponse), nil
	}
}

//创建数据库。
func (c *RdsClient) CreateSqlserverDatabase(request *model.CreateSqlserverDatabaseRequest) (*model.CreateSqlserverDatabaseResponse, error) {
	requestDef := GenReqDefForCreateSqlserverDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSqlserverDatabaseResponse), nil
	}
}

//创建数据库用户。
func (c *RdsClient) CreateSqlserverDbUser(request *model.CreateSqlserverDbUserRequest) (*model.CreateSqlserverDbUserResponse, error) {
	requestDef := GenReqDefForCreateSqlserverDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSqlserverDbUserResponse), nil
	}
}

//删除数据库。
func (c *RdsClient) DeleteSqlserverDatabase(request *model.DeleteSqlserverDatabaseRequest) (*model.DeleteSqlserverDatabaseResponse, error) {
	requestDef := GenReqDefForDeleteSqlserverDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSqlserverDatabaseResponse), nil
	}
}

//删除数据库用户。
func (c *RdsClient) DeleteSqlserverDbUser(request *model.DeleteSqlserverDbUserRequest) (*model.DeleteSqlserverDbUserResponse, error) {
	requestDef := GenReqDefForDeleteSqlserverDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSqlserverDbUserResponse), nil
	}
}

//查询指定数据库的已授权用户。
func (c *RdsClient) ListAuthorizedSqlserverDbUsers(request *model.ListAuthorizedSqlserverDbUsersRequest) (*model.ListAuthorizedSqlserverDbUsersResponse, error) {
	requestDef := GenReqDefForListAuthorizedSqlserverDbUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAuthorizedSqlserverDbUsersResponse), nil
	}
}

//查询数据库列表。
func (c *RdsClient) ListSqlserverDatabases(request *model.ListSqlserverDatabasesRequest) (*model.ListSqlserverDatabasesResponse, error) {
	requestDef := GenReqDefForListSqlserverDatabases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSqlserverDatabasesResponse), nil
	}
}

//查询数据库用户列表。
func (c *RdsClient) ListSqlserverDbUsers(request *model.ListSqlserverDbUsersRequest) (*model.ListSqlserverDbUsersResponse, error) {
	requestDef := GenReqDefForListSqlserverDbUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSqlserverDbUsersResponse), nil
	}
}

//解除数据库帐号权限。
func (c *RdsClient) RevokeSqlserverDbUserPrivilege(request *model.RevokeSqlserverDbUserPrivilegeRequest) (*model.RevokeSqlserverDbUserPrivilegeResponse, error) {
	requestDef := GenReqDefForRevokeSqlserverDbUserPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RevokeSqlserverDbUserPrivilegeResponse), nil
	}
}
