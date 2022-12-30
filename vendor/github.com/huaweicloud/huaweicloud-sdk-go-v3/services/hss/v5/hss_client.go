package v5

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5/model"
)

type HssClient struct {
	HcClient *http_client.HcHttpClient
}

func NewHssClient(hcClient *http_client.HcHttpClient) *HssClient {
	return &HssClient{HcClient: hcClient}
}

func HssClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// AddHostsGroup 创建服务器组
//
// 创建服务器组
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) AddHostsGroup(request *model.AddHostsGroupRequest) (*model.AddHostsGroupResponse, error) {
	requestDef := GenReqDefForAddHostsGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddHostsGroupResponse), nil
	}
}

// AddHostsGroupInvoker 创建服务器组
func (c *HssClient) AddHostsGroupInvoker(request *model.AddHostsGroupRequest) *AddHostsGroupInvoker {
	requestDef := GenReqDefForAddHostsGroup()
	return &AddHostsGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AssociatePolicyGroup 部署策略
//
// 部署策略
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) AssociatePolicyGroup(request *model.AssociatePolicyGroupRequest) (*model.AssociatePolicyGroupResponse, error) {
	requestDef := GenReqDefForAssociatePolicyGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AssociatePolicyGroupResponse), nil
	}
}

// AssociatePolicyGroupInvoker 部署策略
func (c *HssClient) AssociatePolicyGroupInvoker(request *model.AssociatePolicyGroupRequest) *AssociatePolicyGroupInvoker {
	requestDef := GenReqDefForAssociatePolicyGroup()
	return &AssociatePolicyGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchCreateTags 批量创建标签
//
// 批量创建标签
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) BatchCreateTags(request *model.BatchCreateTagsRequest) (*model.BatchCreateTagsResponse, error) {
	requestDef := GenReqDefForBatchCreateTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchCreateTagsResponse), nil
	}
}

// BatchCreateTagsInvoker 批量创建标签
func (c *HssClient) BatchCreateTagsInvoker(request *model.BatchCreateTagsRequest) *BatchCreateTagsInvoker {
	requestDef := GenReqDefForBatchCreateTags()
	return &BatchCreateTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ChangeEvent 处理告警事件
//
// 处理告警事件
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ChangeEvent(request *model.ChangeEventRequest) (*model.ChangeEventResponse, error) {
	requestDef := GenReqDefForChangeEvent()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeEventResponse), nil
	}
}

// ChangeEventInvoker 处理告警事件
func (c *HssClient) ChangeEventInvoker(request *model.ChangeEventRequest) *ChangeEventInvoker {
	requestDef := GenReqDefForChangeEvent()
	return &ChangeEventInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ChangeHostsGroup 编辑服务器组
//
// 编辑服务器组
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ChangeHostsGroup(request *model.ChangeHostsGroupRequest) (*model.ChangeHostsGroupResponse, error) {
	requestDef := GenReqDefForChangeHostsGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeHostsGroupResponse), nil
	}
}

// ChangeHostsGroupInvoker 编辑服务器组
func (c *HssClient) ChangeHostsGroupInvoker(request *model.ChangeHostsGroupRequest) *ChangeHostsGroupInvoker {
	requestDef := GenReqDefForChangeHostsGroup()
	return &ChangeHostsGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ChangeVulStatus 修改漏洞的状态
//
// 修改漏洞的状态
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ChangeVulStatus(request *model.ChangeVulStatusRequest) (*model.ChangeVulStatusResponse, error) {
	requestDef := GenReqDefForChangeVulStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeVulStatusResponse), nil
	}
}

// ChangeVulStatusInvoker 修改漏洞的状态
func (c *HssClient) ChangeVulStatusInvoker(request *model.ChangeVulStatusRequest) *ChangeVulStatusInvoker {
	requestDef := GenReqDefForChangeVulStatus()
	return &ChangeVulStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteHostsGroup 删除服务器组
//
// 删除服务器组
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) DeleteHostsGroup(request *model.DeleteHostsGroupRequest) (*model.DeleteHostsGroupResponse, error) {
	requestDef := GenReqDefForDeleteHostsGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteHostsGroupResponse), nil
	}
}

// DeleteHostsGroupInvoker 删除服务器组
func (c *HssClient) DeleteHostsGroupInvoker(request *model.DeleteHostsGroupRequest) *DeleteHostsGroupInvoker {
	requestDef := GenReqDefForDeleteHostsGroup()
	return &DeleteHostsGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteResourceInstanceTag 删除资源标签
//
// 删除单个资源下的标签
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) DeleteResourceInstanceTag(request *model.DeleteResourceInstanceTagRequest) (*model.DeleteResourceInstanceTagResponse, error) {
	requestDef := GenReqDefForDeleteResourceInstanceTag()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteResourceInstanceTagResponse), nil
	}
}

// DeleteResourceInstanceTagInvoker 删除资源标签
func (c *HssClient) DeleteResourceInstanceTagInvoker(request *model.DeleteResourceInstanceTagRequest) *DeleteResourceInstanceTagInvoker {
	requestDef := GenReqDefForDeleteResourceInstanceTag()
	return &DeleteResourceInstanceTagInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAlarmWhiteList 查询告警白名单列表
//
// 查询告警白名单列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListAlarmWhiteList(request *model.ListAlarmWhiteListRequest) (*model.ListAlarmWhiteListResponse, error) {
	requestDef := GenReqDefForListAlarmWhiteList()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAlarmWhiteListResponse), nil
	}
}

// ListAlarmWhiteListInvoker 查询告警白名单列表
func (c *HssClient) ListAlarmWhiteListInvoker(request *model.ListAlarmWhiteListRequest) *ListAlarmWhiteListInvoker {
	requestDef := GenReqDefForListAlarmWhiteList()
	return &ListAlarmWhiteListInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAppChangeHistories 资产指纹-软件信息-历史变动记录
//
// 资产指纹-软件信息-历史变动记录
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListAppChangeHistories(request *model.ListAppChangeHistoriesRequest) (*model.ListAppChangeHistoriesResponse, error) {
	requestDef := GenReqDefForListAppChangeHistories()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAppChangeHistoriesResponse), nil
	}
}

// ListAppChangeHistoriesInvoker 资产指纹-软件信息-历史变动记录
func (c *HssClient) ListAppChangeHistoriesInvoker(request *model.ListAppChangeHistoriesRequest) *ListAppChangeHistoriesInvoker {
	requestDef := GenReqDefForListAppChangeHistories()
	return &ListAppChangeHistoriesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAppStatistics 资产指纹-软件信息
//
// 资产指纹-软件信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListAppStatistics(request *model.ListAppStatisticsRequest) (*model.ListAppStatisticsResponse, error) {
	requestDef := GenReqDefForListAppStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAppStatisticsResponse), nil
	}
}

// ListAppStatisticsInvoker 资产指纹-软件信息
func (c *HssClient) ListAppStatisticsInvoker(request *model.ListAppStatisticsRequest) *ListAppStatisticsInvoker {
	requestDef := GenReqDefForListAppStatistics()
	return &ListAppStatisticsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListApps 单主机资产指纹-软件
//
// 单主机资产指纹-软件
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListApps(request *model.ListAppsRequest) (*model.ListAppsResponse, error) {
	requestDef := GenReqDefForListApps()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAppsResponse), nil
	}
}

// ListAppsInvoker 单主机资产指纹-软件
func (c *HssClient) ListAppsInvoker(request *model.ListAppsRequest) *ListAppsInvoker {
	requestDef := GenReqDefForListApps()
	return &ListAppsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutoLaunchChangeHistories 资产指纹-自启动项-历史变动记录
//
// 资产指纹-自启动项-历史变动记录
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListAutoLaunchChangeHistories(request *model.ListAutoLaunchChangeHistoriesRequest) (*model.ListAutoLaunchChangeHistoriesResponse, error) {
	requestDef := GenReqDefForListAutoLaunchChangeHistories()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutoLaunchChangeHistoriesResponse), nil
	}
}

// ListAutoLaunchChangeHistoriesInvoker 资产指纹-自启动项-历史变动记录
func (c *HssClient) ListAutoLaunchChangeHistoriesInvoker(request *model.ListAutoLaunchChangeHistoriesRequest) *ListAutoLaunchChangeHistoriesInvoker {
	requestDef := GenReqDefForListAutoLaunchChangeHistories()
	return &ListAutoLaunchChangeHistoriesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutoLaunchStatistics 资产指纹-自启动项信息
//
// 资产指纹-自启动项信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListAutoLaunchStatistics(request *model.ListAutoLaunchStatisticsRequest) (*model.ListAutoLaunchStatisticsResponse, error) {
	requestDef := GenReqDefForListAutoLaunchStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutoLaunchStatisticsResponse), nil
	}
}

// ListAutoLaunchStatisticsInvoker 资产指纹-自启动项信息
func (c *HssClient) ListAutoLaunchStatisticsInvoker(request *model.ListAutoLaunchStatisticsRequest) *ListAutoLaunchStatisticsInvoker {
	requestDef := GenReqDefForListAutoLaunchStatistics()
	return &ListAutoLaunchStatisticsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutoLaunchs 单主机资产指纹-自启动项
//
// 单主机资产指纹-自启动项
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListAutoLaunchs(request *model.ListAutoLaunchsRequest) (*model.ListAutoLaunchsResponse, error) {
	requestDef := GenReqDefForListAutoLaunchs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutoLaunchsResponse), nil
	}
}

// ListAutoLaunchsInvoker 单主机资产指纹-自启动项
func (c *HssClient) ListAutoLaunchsInvoker(request *model.ListAutoLaunchsRequest) *ListAutoLaunchsInvoker {
	requestDef := GenReqDefForListAutoLaunchs()
	return &ListAutoLaunchsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListHostGroups 查询服务器组列表
//
// 查询服务器组列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListHostGroups(request *model.ListHostGroupsRequest) (*model.ListHostGroupsResponse, error) {
	requestDef := GenReqDefForListHostGroups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListHostGroupsResponse), nil
	}
}

// ListHostGroupsInvoker 查询服务器组列表
func (c *HssClient) ListHostGroupsInvoker(request *model.ListHostGroupsRequest) *ListHostGroupsInvoker {
	requestDef := GenReqDefForListHostGroups()
	return &ListHostGroupsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListHostProtectHistoryInfo 查询主机静态网页防篡改防护动态
//
// 查询主机静态网页防篡改防护动态
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListHostProtectHistoryInfo(request *model.ListHostProtectHistoryInfoRequest) (*model.ListHostProtectHistoryInfoResponse, error) {
	requestDef := GenReqDefForListHostProtectHistoryInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListHostProtectHistoryInfoResponse), nil
	}
}

// ListHostProtectHistoryInfoInvoker 查询主机静态网页防篡改防护动态
func (c *HssClient) ListHostProtectHistoryInfoInvoker(request *model.ListHostProtectHistoryInfoRequest) *ListHostProtectHistoryInfoInvoker {
	requestDef := GenReqDefForListHostProtectHistoryInfo()
	return &ListHostProtectHistoryInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListHostRaspProtectHistoryInfo 查询主机动态网页防篡改防护动态
//
// 查询主机动态网页防篡改防护动态
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListHostRaspProtectHistoryInfo(request *model.ListHostRaspProtectHistoryInfoRequest) (*model.ListHostRaspProtectHistoryInfoResponse, error) {
	requestDef := GenReqDefForListHostRaspProtectHistoryInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListHostRaspProtectHistoryInfoResponse), nil
	}
}

// ListHostRaspProtectHistoryInfoInvoker 查询主机动态网页防篡改防护动态
func (c *HssClient) ListHostRaspProtectHistoryInfoInvoker(request *model.ListHostRaspProtectHistoryInfoRequest) *ListHostRaspProtectHistoryInfoInvoker {
	requestDef := GenReqDefForListHostRaspProtectHistoryInfo()
	return &ListHostRaspProtectHistoryInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListHostStatus 查询云服务器列表
//
// 查询云服务器列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListHostStatus(request *model.ListHostStatusRequest) (*model.ListHostStatusResponse, error) {
	requestDef := GenReqDefForListHostStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListHostStatusResponse), nil
	}
}

// ListHostStatusInvoker 查询云服务器列表
func (c *HssClient) ListHostStatusInvoker(request *model.ListHostStatusRequest) *ListHostStatusInvoker {
	requestDef := GenReqDefForListHostStatus()
	return &ListHostStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPasswordComplexity 查询口令复杂度策略检测报告
//
// 查询口令复杂度策略检测报告
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListPasswordComplexity(request *model.ListPasswordComplexityRequest) (*model.ListPasswordComplexityResponse, error) {
	requestDef := GenReqDefForListPasswordComplexity()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPasswordComplexityResponse), nil
	}
}

// ListPasswordComplexityInvoker 查询口令复杂度策略检测报告
func (c *HssClient) ListPasswordComplexityInvoker(request *model.ListPasswordComplexityRequest) *ListPasswordComplexityInvoker {
	requestDef := GenReqDefForListPasswordComplexity()
	return &ListPasswordComplexityInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPolicyGroup 查询策略组列表
//
// 查询策略组列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListPolicyGroup(request *model.ListPolicyGroupRequest) (*model.ListPolicyGroupResponse, error) {
	requestDef := GenReqDefForListPolicyGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPolicyGroupResponse), nil
	}
}

// ListPolicyGroupInvoker 查询策略组列表
func (c *HssClient) ListPolicyGroupInvoker(request *model.ListPolicyGroupRequest) *ListPolicyGroupInvoker {
	requestDef := GenReqDefForListPolicyGroup()
	return &ListPolicyGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPortStatistics 资产指纹-开放端口信息
//
// 资产指纹-开放端口信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListPortStatistics(request *model.ListPortStatisticsRequest) (*model.ListPortStatisticsResponse, error) {
	requestDef := GenReqDefForListPortStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPortStatisticsResponse), nil
	}
}

// ListPortStatisticsInvoker 资产指纹-开放端口信息
func (c *HssClient) ListPortStatisticsInvoker(request *model.ListPortStatisticsRequest) *ListPortStatisticsInvoker {
	requestDef := GenReqDefForListPortStatistics()
	return &ListPortStatisticsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPorts 单主机资产指纹-开放端口信息
//
// 单主机资产指纹-开放端口信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListPorts(request *model.ListPortsRequest) (*model.ListPortsResponse, error) {
	requestDef := GenReqDefForListPorts()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPortsResponse), nil
	}
}

// ListPortsInvoker 单主机资产指纹-开放端口信息
func (c *HssClient) ListPortsInvoker(request *model.ListPortsRequest) *ListPortsInvoker {
	requestDef := GenReqDefForListPorts()
	return &ListPortsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListProcessStatistics 资产指纹-进程信息
//
// 资产指纹-进程信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListProcessStatistics(request *model.ListProcessStatisticsRequest) (*model.ListProcessStatisticsResponse, error) {
	requestDef := GenReqDefForListProcessStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProcessStatisticsResponse), nil
	}
}

// ListProcessStatisticsInvoker 资产指纹-进程信息
func (c *HssClient) ListProcessStatisticsInvoker(request *model.ListProcessStatisticsRequest) *ListProcessStatisticsInvoker {
	requestDef := GenReqDefForListProcessStatistics()
	return &ListProcessStatisticsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListProtectionPolicy 查询防护策略列表
//
// 查询防护策略列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListProtectionPolicy(request *model.ListProtectionPolicyRequest) (*model.ListProtectionPolicyResponse, error) {
	requestDef := GenReqDefForListProtectionPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProtectionPolicyResponse), nil
	}
}

// ListProtectionPolicyInvoker 查询防护策略列表
func (c *HssClient) ListProtectionPolicyInvoker(request *model.ListProtectionPolicyRequest) *ListProtectionPolicyInvoker {
	requestDef := GenReqDefForListProtectionPolicy()
	return &ListProtectionPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListProtectionServer 查询勒索防护服务器列表
//
// 查询勒索防护服务器列表，与云备份服务配合使用。因此使用勒索相关接口之前确保该局点有云备份服务
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListProtectionServer(request *model.ListProtectionServerRequest) (*model.ListProtectionServerResponse, error) {
	requestDef := GenReqDefForListProtectionServer()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProtectionServerResponse), nil
	}
}

// ListProtectionServerInvoker 查询勒索防护服务器列表
func (c *HssClient) ListProtectionServerInvoker(request *model.ListProtectionServerRequest) *ListProtectionServerInvoker {
	requestDef := GenReqDefForListProtectionServer()
	return &ListProtectionServerInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListQuotasDetail 查询配额详情
//
// 查询配额详情
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListQuotasDetail(request *model.ListQuotasDetailRequest) (*model.ListQuotasDetailResponse, error) {
	requestDef := GenReqDefForListQuotasDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListQuotasDetailResponse), nil
	}
}

// ListQuotasDetailInvoker 查询配额详情
func (c *HssClient) ListQuotasDetailInvoker(request *model.ListQuotasDetailRequest) *ListQuotasDetailInvoker {
	requestDef := GenReqDefForListQuotasDetail()
	return &ListQuotasDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRiskConfigCheckRules 查询指定安全配置项的检查项列表
//
// 查询指定安全配置项的检查项列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListRiskConfigCheckRules(request *model.ListRiskConfigCheckRulesRequest) (*model.ListRiskConfigCheckRulesResponse, error) {
	requestDef := GenReqDefForListRiskConfigCheckRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRiskConfigCheckRulesResponse), nil
	}
}

// ListRiskConfigCheckRulesInvoker 查询指定安全配置项的检查项列表
func (c *HssClient) ListRiskConfigCheckRulesInvoker(request *model.ListRiskConfigCheckRulesRequest) *ListRiskConfigCheckRulesInvoker {
	requestDef := GenReqDefForListRiskConfigCheckRules()
	return &ListRiskConfigCheckRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRiskConfigHosts 查询指定安全配置项的受影响服务器列表
//
// 查询指定安全配置项的受影响服务器列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListRiskConfigHosts(request *model.ListRiskConfigHostsRequest) (*model.ListRiskConfigHostsResponse, error) {
	requestDef := GenReqDefForListRiskConfigHosts()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRiskConfigHostsResponse), nil
	}
}

// ListRiskConfigHostsInvoker 查询指定安全配置项的受影响服务器列表
func (c *HssClient) ListRiskConfigHostsInvoker(request *model.ListRiskConfigHostsRequest) *ListRiskConfigHostsInvoker {
	requestDef := GenReqDefForListRiskConfigHosts()
	return &ListRiskConfigHostsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRiskConfigs 查询租户的服务器安全配置检测结果列表
//
// 查询租户的服务器安全配置检测结果列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListRiskConfigs(request *model.ListRiskConfigsRequest) (*model.ListRiskConfigsResponse, error) {
	requestDef := GenReqDefForListRiskConfigs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRiskConfigsResponse), nil
	}
}

// ListRiskConfigsInvoker 查询租户的服务器安全配置检测结果列表
func (c *HssClient) ListRiskConfigsInvoker(request *model.ListRiskConfigsRequest) *ListRiskConfigsInvoker {
	requestDef := GenReqDefForListRiskConfigs()
	return &ListRiskConfigsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSecurityEvents 查入侵事件列表
//
// 查入侵事件列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListSecurityEvents(request *model.ListSecurityEventsRequest) (*model.ListSecurityEventsResponse, error) {
	requestDef := GenReqDefForListSecurityEvents()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSecurityEventsResponse), nil
	}
}

// ListSecurityEventsInvoker 查入侵事件列表
func (c *HssClient) ListSecurityEventsInvoker(request *model.ListSecurityEventsRequest) *ListSecurityEventsInvoker {
	requestDef := GenReqDefForListSecurityEvents()
	return &ListSecurityEventsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListUserChangeHistories 获取账户变动历史信息
//
// 获取账户变动历史记录信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListUserChangeHistories(request *model.ListUserChangeHistoriesRequest) (*model.ListUserChangeHistoriesResponse, error) {
	requestDef := GenReqDefForListUserChangeHistories()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListUserChangeHistoriesResponse), nil
	}
}

// ListUserChangeHistoriesInvoker 获取账户变动历史信息
func (c *HssClient) ListUserChangeHistoriesInvoker(request *model.ListUserChangeHistoriesRequest) *ListUserChangeHistoriesInvoker {
	requestDef := GenReqDefForListUserChangeHistories()
	return &ListUserChangeHistoriesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListUserStatistics 资产指纹-账号信息
//
// 资产指纹-账号信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListUserStatistics(request *model.ListUserStatisticsRequest) (*model.ListUserStatisticsResponse, error) {
	requestDef := GenReqDefForListUserStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListUserStatisticsResponse), nil
	}
}

// ListUserStatisticsInvoker 资产指纹-账号信息
func (c *HssClient) ListUserStatisticsInvoker(request *model.ListUserStatisticsRequest) *ListUserStatisticsInvoker {
	requestDef := GenReqDefForListUserStatistics()
	return &ListUserStatisticsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListUsers 获取资产的账号列表
//
// 获取资产的账号列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListUsers(request *model.ListUsersRequest) (*model.ListUsersResponse, error) {
	requestDef := GenReqDefForListUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListUsersResponse), nil
	}
}

// ListUsersInvoker 获取资产的账号列表
func (c *HssClient) ListUsersInvoker(request *model.ListUsersRequest) *ListUsersInvoker {
	requestDef := GenReqDefForListUsers()
	return &ListUsersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListVulHosts 查询单个漏洞影响的云服务器信息
//
// 查询单个漏洞影响的云服务器信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListVulHosts(request *model.ListVulHostsRequest) (*model.ListVulHostsResponse, error) {
	requestDef := GenReqDefForListVulHosts()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListVulHostsResponse), nil
	}
}

// ListVulHostsInvoker 查询单个漏洞影响的云服务器信息
func (c *HssClient) ListVulHostsInvoker(request *model.ListVulHostsRequest) *ListVulHostsInvoker {
	requestDef := GenReqDefForListVulHosts()
	return &ListVulHostsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListVulnerabilities 查询漏洞列表
//
// 查询漏洞列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListVulnerabilities(request *model.ListVulnerabilitiesRequest) (*model.ListVulnerabilitiesResponse, error) {
	requestDef := GenReqDefForListVulnerabilities()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListVulnerabilitiesResponse), nil
	}
}

// ListVulnerabilitiesInvoker 查询漏洞列表
func (c *HssClient) ListVulnerabilitiesInvoker(request *model.ListVulnerabilitiesRequest) *ListVulnerabilitiesInvoker {
	requestDef := GenReqDefForListVulnerabilities()
	return &ListVulnerabilitiesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListWeakPasswordUsers 查询弱口令检测结果列表
//
// 查询弱口令检测结果列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListWeakPasswordUsers(request *model.ListWeakPasswordUsersRequest) (*model.ListWeakPasswordUsersResponse, error) {
	requestDef := GenReqDefForListWeakPasswordUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListWeakPasswordUsersResponse), nil
	}
}

// ListWeakPasswordUsersInvoker 查询弱口令检测结果列表
func (c *HssClient) ListWeakPasswordUsersInvoker(request *model.ListWeakPasswordUsersRequest) *ListWeakPasswordUsersInvoker {
	requestDef := GenReqDefForListWeakPasswordUsers()
	return &ListWeakPasswordUsersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListWtpProtectHost 查询防护列表
//
// 查询防护列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ListWtpProtectHost(request *model.ListWtpProtectHostRequest) (*model.ListWtpProtectHostResponse, error) {
	requestDef := GenReqDefForListWtpProtectHost()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListWtpProtectHostResponse), nil
	}
}

// ListWtpProtectHostInvoker 查询防护列表
func (c *HssClient) ListWtpProtectHostInvoker(request *model.ListWtpProtectHostRequest) *ListWtpProtectHostInvoker {
	requestDef := GenReqDefForListWtpProtectHost()
	return &ListWtpProtectHostInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetRaspSwitch 开启/关闭动态网页防篡改防护
//
// 开启/关闭动态网页防篡改防护
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) SetRaspSwitch(request *model.SetRaspSwitchRequest) (*model.SetRaspSwitchResponse, error) {
	requestDef := GenReqDefForSetRaspSwitch()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetRaspSwitchResponse), nil
	}
}

// SetRaspSwitchInvoker 开启/关闭动态网页防篡改防护
func (c *HssClient) SetRaspSwitchInvoker(request *model.SetRaspSwitchRequest) *SetRaspSwitchInvoker {
	requestDef := GenReqDefForSetRaspSwitch()
	return &SetRaspSwitchInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetWtpProtectionStatusInfo 开启关闭网页防篡改防护
//
// 开启关闭网页防篡改防护
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) SetWtpProtectionStatusInfo(request *model.SetWtpProtectionStatusInfoRequest) (*model.SetWtpProtectionStatusInfoResponse, error) {
	requestDef := GenReqDefForSetWtpProtectionStatusInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetWtpProtectionStatusInfoResponse), nil
	}
}

// SetWtpProtectionStatusInfoInvoker 开启关闭网页防篡改防护
func (c *HssClient) SetWtpProtectionStatusInfoInvoker(request *model.SetWtpProtectionStatusInfoRequest) *SetWtpProtectionStatusInfoInvoker {
	requestDef := GenReqDefForSetWtpProtectionStatusInfo()
	return &SetWtpProtectionStatusInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAssetStatistic 统计资产信息，账号、端口、进程等
//
// 资产统计信息，账号、端口、进程等
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ShowAssetStatistic(request *model.ShowAssetStatisticRequest) (*model.ShowAssetStatisticResponse, error) {
	requestDef := GenReqDefForShowAssetStatistic()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAssetStatisticResponse), nil
	}
}

// ShowAssetStatisticInvoker 统计资产信息，账号、端口、进程等
func (c *HssClient) ShowAssetStatisticInvoker(request *model.ShowAssetStatisticRequest) *ShowAssetStatisticInvoker {
	requestDef := GenReqDefForShowAssetStatistic()
	return &ShowAssetStatisticInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowBackupPolicyInfo 查询备份策略信息
//
// 查询备份策略信息,确保已经购买了勒索防护存储库，可以从cbr云备份服务进行验证，确保已经存在HSS_projectid命名的存储库已经购买
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ShowBackupPolicyInfo(request *model.ShowBackupPolicyInfoRequest) (*model.ShowBackupPolicyInfoResponse, error) {
	requestDef := GenReqDefForShowBackupPolicyInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBackupPolicyInfoResponse), nil
	}
}

// ShowBackupPolicyInfoInvoker 查询备份策略信息
func (c *HssClient) ShowBackupPolicyInfoInvoker(request *model.ShowBackupPolicyInfoRequest) *ShowBackupPolicyInfoInvoker {
	requestDef := GenReqDefForShowBackupPolicyInfo()
	return &ShowBackupPolicyInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowCheckRuleDetail 查询配置检查项检测报告
//
// 查询配置检查项检测报告
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ShowCheckRuleDetail(request *model.ShowCheckRuleDetailRequest) (*model.ShowCheckRuleDetailResponse, error) {
	requestDef := GenReqDefForShowCheckRuleDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowCheckRuleDetailResponse), nil
	}
}

// ShowCheckRuleDetailInvoker 查询配置检查项检测报告
func (c *HssClient) ShowCheckRuleDetailInvoker(request *model.ShowCheckRuleDetailRequest) *ShowCheckRuleDetailInvoker {
	requestDef := GenReqDefForShowCheckRuleDetail()
	return &ShowCheckRuleDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowResourceQuotas 查询配额信息
//
// 查询配额信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ShowResourceQuotas(request *model.ShowResourceQuotasRequest) (*model.ShowResourceQuotasResponse, error) {
	requestDef := GenReqDefForShowResourceQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowResourceQuotasResponse), nil
	}
}

// ShowResourceQuotasInvoker 查询配额信息
func (c *HssClient) ShowResourceQuotasInvoker(request *model.ShowResourceQuotasRequest) *ShowResourceQuotasInvoker {
	requestDef := GenReqDefForShowResourceQuotas()
	return &ShowResourceQuotasInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowRiskConfigDetail 查询指定安全配置项的检查结果
//
// 查询指定安全配置项的检查结果
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) ShowRiskConfigDetail(request *model.ShowRiskConfigDetailRequest) (*model.ShowRiskConfigDetailResponse, error) {
	requestDef := GenReqDefForShowRiskConfigDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowRiskConfigDetailResponse), nil
	}
}

// ShowRiskConfigDetailInvoker 查询指定安全配置项的检查结果
func (c *HssClient) ShowRiskConfigDetailInvoker(request *model.ShowRiskConfigDetailRequest) *ShowRiskConfigDetailInvoker {
	requestDef := GenReqDefForShowRiskConfigDetail()
	return &ShowRiskConfigDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartProtection 开启勒索病毒防护
//
// 开启勒索病毒防护,请保证该region有cbr云备份服务，勒索服务与云备份服务有关联关系
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) StartProtection(request *model.StartProtectionRequest) (*model.StartProtectionResponse, error) {
	requestDef := GenReqDefForStartProtection()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartProtectionResponse), nil
	}
}

// StartProtectionInvoker 开启勒索病毒防护
func (c *HssClient) StartProtectionInvoker(request *model.StartProtectionRequest) *StartProtectionInvoker {
	requestDef := GenReqDefForStartProtection()
	return &StartProtectionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopProtection 关闭勒索病毒防护
//
// 关闭勒索病毒防护
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) StopProtection(request *model.StopProtectionRequest) (*model.StopProtectionResponse, error) {
	requestDef := GenReqDefForStopProtection()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopProtectionResponse), nil
	}
}

// StopProtectionInvoker 关闭勒索病毒防护
func (c *HssClient) StopProtectionInvoker(request *model.StopProtectionRequest) *StopProtectionInvoker {
	requestDef := GenReqDefForStopProtection()
	return &StopProtectionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SwitchHostsProtectStatus 切换防护状态
//
// 切换防护状态
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) SwitchHostsProtectStatus(request *model.SwitchHostsProtectStatusRequest) (*model.SwitchHostsProtectStatusResponse, error) {
	requestDef := GenReqDefForSwitchHostsProtectStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SwitchHostsProtectStatusResponse), nil
	}
}

// SwitchHostsProtectStatusInvoker 切换防护状态
func (c *HssClient) SwitchHostsProtectStatusInvoker(request *model.SwitchHostsProtectStatusRequest) *SwitchHostsProtectStatusInvoker {
	requestDef := GenReqDefForSwitchHostsProtectStatus()
	return &SwitchHostsProtectStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateBackupPolicyInfo 修改备份策略
//
// 修改备份策略
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) UpdateBackupPolicyInfo(request *model.UpdateBackupPolicyInfoRequest) (*model.UpdateBackupPolicyInfoResponse, error) {
	requestDef := GenReqDefForUpdateBackupPolicyInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateBackupPolicyInfoResponse), nil
	}
}

// UpdateBackupPolicyInfoInvoker 修改备份策略
func (c *HssClient) UpdateBackupPolicyInfoInvoker(request *model.UpdateBackupPolicyInfoRequest) *UpdateBackupPolicyInfoInvoker {
	requestDef := GenReqDefForUpdateBackupPolicyInfo()
	return &UpdateBackupPolicyInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateProtectionPolicy 修改防护策略
//
// 修改防护策略
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *HssClient) UpdateProtectionPolicy(request *model.UpdateProtectionPolicyRequest) (*model.UpdateProtectionPolicyResponse, error) {
	requestDef := GenReqDefForUpdateProtectionPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateProtectionPolicyResponse), nil
	}
}

// UpdateProtectionPolicyInvoker 修改防护策略
func (c *HssClient) UpdateProtectionPolicyInvoker(request *model.UpdateProtectionPolicyRequest) *UpdateProtectionPolicyInvoker {
	requestDef := GenReqDefForUpdateProtectionPolicy()
	return &UpdateProtectionPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
