package v3

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/gaussdb/v3/model"
)

type GaussDBClient struct {
	HcClient *http_client.HcHttpClient
}

func NewGaussDBClient(hcClient *http_client.HcHttpClient) *GaussDBClient {
	return &GaussDBClient{HcClient: hcClient}
}

func GaussDBClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

//批量添加或删除指定实例的标签。
func (c *GaussDBClient) BatchTagAction(request *model.BatchTagActionRequest) (*model.BatchTagActionResponse, error) {
	requestDef := GenReqDefForBatchTagAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchTagActionResponse), nil
	}
}

//变更数据库实例的规格。
func (c *GaussDBClient) ChangeGaussMySqlInstanceSpecification(request *model.ChangeGaussMySqlInstanceSpecificationRequest) (*model.ChangeGaussMySqlInstanceSpecificationResponse, error) {
	requestDef := GenReqDefForChangeGaussMySqlInstanceSpecification()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeGaussMySqlInstanceSpecificationResponse), nil
	}
}

//创建手动备份
func (c *GaussDBClient) CreateGaussMySqlBackup(request *model.CreateGaussMySqlBackupRequest) (*model.CreateGaussMySqlBackupResponse, error) {
	requestDef := GenReqDefForCreateGaussMySqlBackup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateGaussMySqlBackupResponse), nil
	}
}

//创建云数据库 GaussDB(for MySQL)实例。
func (c *GaussDBClient) CreateGaussMySqlInstance(request *model.CreateGaussMySqlInstanceRequest) (*model.CreateGaussMySqlInstanceResponse, error) {
	requestDef := GenReqDefForCreateGaussMySqlInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateGaussMySqlInstanceResponse), nil
	}
}

//开启数据库代理，只支持ELB模式。
func (c *GaussDBClient) CreateGaussMySqlProxy(request *model.CreateGaussMySqlProxyRequest) (*model.CreateGaussMySqlProxyResponse, error) {
	requestDef := GenReqDefForCreateGaussMySqlProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateGaussMySqlProxyResponse), nil
	}
}

//创建只读节点。
func (c *GaussDBClient) CreateGaussMySqlReadonlyNode(request *model.CreateGaussMySqlReadonlyNodeRequest) (*model.CreateGaussMySqlReadonlyNodeResponse, error) {
	requestDef := GenReqDefForCreateGaussMySqlReadonlyNode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateGaussMySqlReadonlyNodeResponse), nil
	}
}

//删除数据库实例，不支持删除包周期实例。
func (c *GaussDBClient) DeleteGaussMySqlInstance(request *model.DeleteGaussMySqlInstanceRequest) (*model.DeleteGaussMySqlInstanceResponse, error) {
	requestDef := GenReqDefForDeleteGaussMySqlInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteGaussMySqlInstanceResponse), nil
	}
}

//关闭数据库代理。
func (c *GaussDBClient) DeleteGaussMySqlProxy(request *model.DeleteGaussMySqlProxyRequest) (*model.DeleteGaussMySqlProxyResponse, error) {
	requestDef := GenReqDefForDeleteGaussMySqlProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteGaussMySqlProxyResponse), nil
	}
}

//删除实例的只读节点。多可用区模式删除只读节点时，要保证删除后，剩余的只读节点和主节点在不同的可用区中，否则无法删除该只读节点。
func (c *GaussDBClient) DeleteGaussMySqlReadonlyNode(request *model.DeleteGaussMySqlReadonlyNodeRequest) (*model.DeleteGaussMySqlReadonlyNodeResponse, error) {
	requestDef := GenReqDefForDeleteGaussMySqlReadonlyNode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteGaussMySqlReadonlyNodeResponse), nil
	}
}

//包周期存储扩容
func (c *GaussDBClient) ExpandGaussMySqlInstanceVolume(request *model.ExpandGaussMySqlInstanceVolumeRequest) (*model.ExpandGaussMySqlInstanceVolumeResponse, error) {
	requestDef := GenReqDefForExpandGaussMySqlInstanceVolume()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ExpandGaussMySqlInstanceVolumeResponse), nil
	}
}

//扩容数据库代理节点的数量。 DeC专属云账号暂不支持数据库代理。
func (c *GaussDBClient) ExpandGaussMySqlProxy(request *model.ExpandGaussMySqlProxyRequest) (*model.ExpandGaussMySqlProxyResponse, error) {
	requestDef := GenReqDefForExpandGaussMySqlProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ExpandGaussMySqlProxyResponse), nil
	}
}

//获取参数模板列表，包括所有数据库的默认参数模板和用户创建的参数模板。
func (c *GaussDBClient) ListGaussMySqlConfigurations(request *model.ListGaussMySqlConfigurationsRequest) (*model.ListGaussMySqlConfigurationsResponse, error) {
	requestDef := GenReqDefForListGaussMySqlConfigurations()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListGaussMySqlConfigurationsResponse), nil
	}
}

//获取专属资源池列表，包括用户开通的所有专属资源池信息。
func (c *GaussDBClient) ListGaussMySqlDedicatedResources(request *model.ListGaussMySqlDedicatedResourcesRequest) (*model.ListGaussMySqlDedicatedResourcesResponse, error) {
	requestDef := GenReqDefForListGaussMySqlDedicatedResources()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListGaussMySqlDedicatedResourcesResponse), nil
	}
}

//查询数据库错误日志。
func (c *GaussDBClient) ListGaussMySqlErrorLog(request *model.ListGaussMySqlErrorLogRequest) (*model.ListGaussMySqlErrorLogResponse, error) {
	requestDef := GenReqDefForListGaussMySqlErrorLog()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListGaussMySqlErrorLogResponse), nil
	}
}

//根据指定条件查询实例列表。
func (c *GaussDBClient) ListGaussMySqlInstances(request *model.ListGaussMySqlInstancesRequest) (*model.ListGaussMySqlInstancesResponse, error) {
	requestDef := GenReqDefForListGaussMySqlInstances()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListGaussMySqlInstancesResponse), nil
	}
}

//查询数据库慢日志
func (c *GaussDBClient) ListGaussMySqlSlowLog(request *model.ListGaussMySqlSlowLogRequest) (*model.ListGaussMySqlSlowLogResponse, error) {
	requestDef := GenReqDefForListGaussMySqlSlowLog()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListGaussMySqlSlowLogResponse), nil
	}
}

//查询指定实例的标签信息。
func (c *GaussDBClient) ListInstanceTags(request *model.ListInstanceTagsRequest) (*model.ListInstanceTagsResponse, error) {
	requestDef := GenReqDefForListInstanceTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstanceTagsResponse), nil
	}
}

//查询指定project ID下实例的所有标签集合。
func (c *GaussDBClient) ListProjectTags(request *model.ListProjectTagsRequest) (*model.ListProjectTagsResponse, error) {
	requestDef := GenReqDefForListProjectTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProjectTagsResponse), nil
	}
}

//重置数据库密码
func (c *GaussDBClient) ResetGaussMySqlPassword(request *model.ResetGaussMySqlPasswordRequest) (*model.ResetGaussMySqlPasswordResponse, error) {
	requestDef := GenReqDefForResetGaussMySqlPassword()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ResetGaussMySqlPasswordResponse), nil
	}
}

//设置指定企业项目的资源配额。
func (c *GaussDBClient) SetGaussMySqlQuotas(request *model.SetGaussMySqlQuotasRequest) (*model.SetGaussMySqlQuotasResponse, error) {
	requestDef := GenReqDefForSetGaussMySqlQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetGaussMySqlQuotasResponse), nil
	}
}

//查询审计日志开关状态
func (c *GaussDBClient) ShowAuditLog(request *model.ShowAuditLogRequest) (*model.ShowAuditLogResponse, error) {
	requestDef := GenReqDefForShowAuditLog()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAuditLogResponse), nil
	}
}

//查询备份列表
func (c *GaussDBClient) ShowGaussMySqlBackupList(request *model.ShowGaussMySqlBackupListRequest) (*model.ShowGaussMySqlBackupListResponse, error) {
	requestDef := GenReqDefForShowGaussMySqlBackupList()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowGaussMySqlBackupListResponse), nil
	}
}

//查询自动备份策略。
func (c *GaussDBClient) ShowGaussMySqlBackupPolicy(request *model.ShowGaussMySqlBackupPolicyRequest) (*model.ShowGaussMySqlBackupPolicyResponse, error) {
	requestDef := GenReqDefForShowGaussMySqlBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowGaussMySqlBackupPolicyResponse), nil
	}
}

//获取指定数据库引擎对应的数据库版本信息。
func (c *GaussDBClient) ShowGaussMySqlEngineVersion(request *model.ShowGaussMySqlEngineVersionRequest) (*model.ShowGaussMySqlEngineVersionResponse, error) {
	requestDef := GenReqDefForShowGaussMySqlEngineVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowGaussMySqlEngineVersionResponse), nil
	}
}

//获取指定数据库引擎版本对应的规格信息。
func (c *GaussDBClient) ShowGaussMySqlFlavors(request *model.ShowGaussMySqlFlavorsRequest) (*model.ShowGaussMySqlFlavorsResponse, error) {
	requestDef := GenReqDefForShowGaussMySqlFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowGaussMySqlFlavorsResponse), nil
	}
}

//查询实例详情信息
func (c *GaussDBClient) ShowGaussMySqlInstanceInfo(request *model.ShowGaussMySqlInstanceInfoRequest) (*model.ShowGaussMySqlInstanceInfoResponse, error) {
	requestDef := GenReqDefForShowGaussMySqlInstanceInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowGaussMySqlInstanceInfoResponse), nil
	}
}

//获取指定ID的任务信息。
func (c *GaussDBClient) ShowGaussMySqlJobInfo(request *model.ShowGaussMySqlJobInfoRequest) (*model.ShowGaussMySqlJobInfoResponse, error) {
	requestDef := GenReqDefForShowGaussMySqlJobInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowGaussMySqlJobInfoResponse), nil
	}
}

//获取指定租户的资源配额。
func (c *GaussDBClient) ShowGaussMySqlProjectQuotas(request *model.ShowGaussMySqlProjectQuotasRequest) (*model.ShowGaussMySqlProjectQuotasResponse, error) {
	requestDef := GenReqDefForShowGaussMySqlProjectQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowGaussMySqlProjectQuotasResponse), nil
	}
}

//查询数据库代理信息。
func (c *GaussDBClient) ShowGaussMySqlProxy(request *model.ShowGaussMySqlProxyRequest) (*model.ShowGaussMySqlProxyResponse, error) {
	requestDef := GenReqDefForShowGaussMySqlProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowGaussMySqlProxyResponse), nil
	}
}

//查询数据库代理规格信息。
func (c *GaussDBClient) ShowGaussMySqlProxyFlavors(request *model.ShowGaussMySqlProxyFlavorsRequest) (*model.ShowGaussMySqlProxyFlavorsResponse, error) {
	requestDef := GenReqDefForShowGaussMySqlProxyFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowGaussMySqlProxyFlavorsResponse), nil
	}
}

//获取指定企业项目的资源配额。
func (c *GaussDBClient) ShowGaussMySqlQuotas(request *model.ShowGaussMySqlQuotasRequest) (*model.ShowGaussMySqlQuotasResponse, error) {
	requestDef := GenReqDefForShowGaussMySqlQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowGaussMySqlQuotasResponse), nil
	}
}

//查询实例秒级监控频率。
func (c *GaussDBClient) ShowInstanceMonitorExtend(request *model.ShowInstanceMonitorExtendRequest) (*model.ShowInstanceMonitorExtendResponse, error) {
	requestDef := GenReqDefForShowInstanceMonitorExtend()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowInstanceMonitorExtendResponse), nil
	}
}

//开启或者关闭审计日志
func (c *GaussDBClient) UpdateAuditLog(request *model.UpdateAuditLogRequest) (*model.UpdateAuditLogResponse, error) {
	requestDef := GenReqDefForUpdateAuditLog()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAuditLogResponse), nil
	}
}

//修改备份策略
func (c *GaussDBClient) UpdateGaussMySqlBackupPolicy(request *model.UpdateGaussMySqlBackupPolicyRequest) (*model.UpdateGaussMySqlBackupPolicyResponse, error) {
	requestDef := GenReqDefForUpdateGaussMySqlBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateGaussMySqlBackupPolicyResponse), nil
	}
}

//修改实例名称
func (c *GaussDBClient) UpdateGaussMySqlInstanceName(request *model.UpdateGaussMySqlInstanceNameRequest) (*model.UpdateGaussMySqlInstanceNameResponse, error) {
	requestDef := GenReqDefForUpdateGaussMySqlInstanceName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateGaussMySqlInstanceNameResponse), nil
	}
}

//修改指定企业项目的资源配额。
func (c *GaussDBClient) UpdateGaussMySqlQuotas(request *model.UpdateGaussMySqlQuotasRequest) (*model.UpdateGaussMySqlQuotasResponse, error) {
	requestDef := GenReqDefForUpdateGaussMySqlQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateGaussMySqlQuotasResponse), nil
	}
}

//打开/关闭/修改实例秒级监控。
func (c *GaussDBClient) UpdateInstanceMonitor(request *model.UpdateInstanceMonitorRequest) (*model.UpdateInstanceMonitorResponse, error) {
	requestDef := GenReqDefForUpdateInstanceMonitor()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateInstanceMonitorResponse), nil
	}
}
