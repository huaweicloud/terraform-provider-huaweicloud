package v2

import (
	httpclient "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/aom/v2/model"
)

type AomClient struct {
	HcClient *httpclient.HcHttpClient
}

func NewAomClient(hcClient *httpclient.HcHttpClient) *AomClient {
	return &AomClient{HcClient: hcClient}
}

func AomClientBuilder() *httpclient.HcHttpClientBuilder {
	builder := httpclient.NewHcHttpClientBuilder()
	return builder
}

// AddActionRule 新增告警行动规则
//
// 新增告警行动规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) AddActionRule(request *model.AddActionRuleRequest) (*model.AddActionRuleResponse, error) {
	requestDef := GenReqDefForAddActionRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddActionRuleResponse), nil
	}
}

// AddActionRuleInvoker 新增告警行动规则
func (c *AomClient) AddActionRuleInvoker(request *model.AddActionRuleRequest) *AddActionRuleInvoker {
	requestDef := GenReqDefForAddActionRule()
	return &AddActionRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddAlarmRule 添加阈值规则
//
// 该接口用于添加一条阈值规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) AddAlarmRule(request *model.AddAlarmRuleRequest) (*model.AddAlarmRuleResponse, error) {
	requestDef := GenReqDefForAddAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddAlarmRuleResponse), nil
	}
}

// AddAlarmRuleInvoker 添加阈值规则
func (c *AomClient) AddAlarmRuleInvoker(request *model.AddAlarmRuleRequest) *AddAlarmRuleInvoker {
	requestDef := GenReqDefForAddAlarmRule()
	return &AddAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddEvent2alarmRule 新增一条事件类告警规则
//
// 新增一条事件类告警规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) AddEvent2alarmRule(request *model.AddEvent2alarmRuleRequest) (*model.AddEvent2alarmRuleResponse, error) {
	requestDef := GenReqDefForAddEvent2alarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddEvent2alarmRuleResponse), nil
	}
}

// AddEvent2alarmRuleInvoker 新增一条事件类告警规则
func (c *AomClient) AddEvent2alarmRuleInvoker(request *model.AddEvent2alarmRuleRequest) *AddEvent2alarmRuleInvoker {
	requestDef := GenReqDefForAddEvent2alarmRule()
	return &AddEvent2alarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddMetricData 添加监控数据
//
// 该接口用于向服务端添加一条或多条监控数据。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) AddMetricData(request *model.AddMetricDataRequest) (*model.AddMetricDataResponse, error) {
	requestDef := GenReqDefForAddMetricData()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddMetricDataResponse), nil
	}
}

// AddMetricDataInvoker 添加监控数据
func (c *AomClient) AddMetricDataInvoker(request *model.AddMetricDataRequest) *AddMetricDataInvoker {
	requestDef := GenReqDefForAddMetricData()
	return &AddMetricDataInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddMuteRules 新增静默规则
//
// 新增静默规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) AddMuteRules(request *model.AddMuteRulesRequest) (*model.AddMuteRulesResponse, error) {
	requestDef := GenReqDefForAddMuteRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddMuteRulesResponse), nil
	}
}

// AddMuteRulesInvoker 新增静默规则
func (c *AomClient) AddMuteRulesInvoker(request *model.AddMuteRulesRequest) *AddMuteRulesInvoker {
	requestDef := GenReqDefForAddMuteRules()
	return &AddMuteRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddOrUpdateMetricOrEventAlarmRule 添加或修改指标类或事件类告警规则
//
// 添加或修改AOM2.0指标类或事件类告警规则。(注：接口目前开放的region为：华东-上海一)
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) AddOrUpdateMetricOrEventAlarmRule(request *model.AddOrUpdateMetricOrEventAlarmRuleRequest) (*model.AddOrUpdateMetricOrEventAlarmRuleResponse, error) {
	requestDef := GenReqDefForAddOrUpdateMetricOrEventAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddOrUpdateMetricOrEventAlarmRuleResponse), nil
	}
}

// AddOrUpdateMetricOrEventAlarmRuleInvoker 添加或修改指标类或事件类告警规则
func (c *AomClient) AddOrUpdateMetricOrEventAlarmRuleInvoker(request *model.AddOrUpdateMetricOrEventAlarmRuleRequest) *AddOrUpdateMetricOrEventAlarmRuleInvoker {
	requestDef := GenReqDefForAddOrUpdateMetricOrEventAlarmRule()
	return &AddOrUpdateMetricOrEventAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddOrUpdateServiceDiscoveryRules 添加或修改服务发现规则
//
// 该接口用于添加或修改一条或多条服务发现规则。同一projectid下可添加的规则上限为100条。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) AddOrUpdateServiceDiscoveryRules(request *model.AddOrUpdateServiceDiscoveryRulesRequest) (*model.AddOrUpdateServiceDiscoveryRulesResponse, error) {
	requestDef := GenReqDefForAddOrUpdateServiceDiscoveryRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddOrUpdateServiceDiscoveryRulesResponse), nil
	}
}

// AddOrUpdateServiceDiscoveryRulesInvoker 添加或修改服务发现规则
func (c *AomClient) AddOrUpdateServiceDiscoveryRulesInvoker(request *model.AddOrUpdateServiceDiscoveryRulesRequest) *AddOrUpdateServiceDiscoveryRulesInvoker {
	requestDef := GenReqDefForAddOrUpdateServiceDiscoveryRules()
	return &AddOrUpdateServiceDiscoveryRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CountEvents 统计事件告警信息
//
// 该接口用于分段统计指定条件下的事件、告警。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) CountEvents(request *model.CountEventsRequest) (*model.CountEventsResponse, error) {
	requestDef := GenReqDefForCountEvents()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CountEventsResponse), nil
	}
}

// CountEventsInvoker 统计事件告警信息
func (c *AomClient) CountEventsInvoker(request *model.CountEventsRequest) *CountEventsInvoker {
	requestDef := GenReqDefForCountEvents()
	return &CountEventsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteActionRule 删除告警行动规则
//
// 删除告警行动规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) DeleteActionRule(request *model.DeleteActionRuleRequest) (*model.DeleteActionRuleResponse, error) {
	requestDef := GenReqDefForDeleteActionRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteActionRuleResponse), nil
	}
}

// DeleteActionRuleInvoker 删除告警行动规则
func (c *AomClient) DeleteActionRuleInvoker(request *model.DeleteActionRuleRequest) *DeleteActionRuleInvoker {
	requestDef := GenReqDefForDeleteActionRule()
	return &DeleteActionRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAlarmRule 删除阈值规则
//
// 该接口用于删除阈值规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) DeleteAlarmRule(request *model.DeleteAlarmRuleRequest) (*model.DeleteAlarmRuleResponse, error) {
	requestDef := GenReqDefForDeleteAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAlarmRuleResponse), nil
	}
}

// DeleteAlarmRuleInvoker 删除阈值规则
func (c *AomClient) DeleteAlarmRuleInvoker(request *model.DeleteAlarmRuleRequest) *DeleteAlarmRuleInvoker {
	requestDef := GenReqDefForDeleteAlarmRule()
	return &DeleteAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAlarmRules 批量删除阈值规则
//
// 该接口用于批量删除阈值规则
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) DeleteAlarmRules(request *model.DeleteAlarmRulesRequest) (*model.DeleteAlarmRulesResponse, error) {
	requestDef := GenReqDefForDeleteAlarmRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAlarmRulesResponse), nil
	}
}

// DeleteAlarmRulesInvoker 批量删除阈值规则
func (c *AomClient) DeleteAlarmRulesInvoker(request *model.DeleteAlarmRulesRequest) *DeleteAlarmRulesInvoker {
	requestDef := GenReqDefForDeleteAlarmRules()
	return &DeleteAlarmRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteEvent2alarmRule 删除事件类告警规则
//
// 删除一条事件类告警规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) DeleteEvent2alarmRule(request *model.DeleteEvent2alarmRuleRequest) (*model.DeleteEvent2alarmRuleResponse, error) {
	requestDef := GenReqDefForDeleteEvent2alarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteEvent2alarmRuleResponse), nil
	}
}

// DeleteEvent2alarmRuleInvoker 删除事件类告警规则
func (c *AomClient) DeleteEvent2alarmRuleInvoker(request *model.DeleteEvent2alarmRuleRequest) *DeleteEvent2alarmRuleInvoker {
	requestDef := GenReqDefForDeleteEvent2alarmRule()
	return &DeleteEvent2alarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteMetricOrEventAlarmRule 删除指标类或事件类告警规则
//
// 删除AOM2.0指标类或事件类告警规则。(注：接口目前开放的region为：华东-上海一)
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) DeleteMetricOrEventAlarmRule(request *model.DeleteMetricOrEventAlarmRuleRequest) (*model.DeleteMetricOrEventAlarmRuleResponse, error) {
	requestDef := GenReqDefForDeleteMetricOrEventAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteMetricOrEventAlarmRuleResponse), nil
	}
}

// DeleteMetricOrEventAlarmRuleInvoker 删除指标类或事件类告警规则
func (c *AomClient) DeleteMetricOrEventAlarmRuleInvoker(request *model.DeleteMetricOrEventAlarmRuleRequest) *DeleteMetricOrEventAlarmRuleInvoker {
	requestDef := GenReqDefForDeleteMetricOrEventAlarmRule()
	return &DeleteMetricOrEventAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteMuteRules 删除静默规则
//
// 删除静默规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) DeleteMuteRules(request *model.DeleteMuteRulesRequest) (*model.DeleteMuteRulesResponse, error) {
	requestDef := GenReqDefForDeleteMuteRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteMuteRulesResponse), nil
	}
}

// DeleteMuteRulesInvoker 删除静默规则
func (c *AomClient) DeleteMuteRulesInvoker(request *model.DeleteMuteRulesRequest) *DeleteMuteRulesInvoker {
	requestDef := GenReqDefForDeleteMuteRules()
	return &DeleteMuteRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteserviceDiscoveryRules 删除服务发现规则
//
// 该接口用于删除服务发现规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) DeleteserviceDiscoveryRules(request *model.DeleteserviceDiscoveryRulesRequest) (*model.DeleteserviceDiscoveryRulesResponse, error) {
	requestDef := GenReqDefForDeleteserviceDiscoveryRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteserviceDiscoveryRulesResponse), nil
	}
}

// DeleteserviceDiscoveryRulesInvoker 删除服务发现规则
func (c *AomClient) DeleteserviceDiscoveryRulesInvoker(request *model.DeleteserviceDiscoveryRulesRequest) *DeleteserviceDiscoveryRulesInvoker {
	requestDef := GenReqDefForDeleteserviceDiscoveryRules()
	return &DeleteserviceDiscoveryRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListActionRule 获取告警行动规则列表
//
// 获取告警行动规则列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListActionRule(request *model.ListActionRuleRequest) (*model.ListActionRuleResponse, error) {
	requestDef := GenReqDefForListActionRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListActionRuleResponse), nil
	}
}

// ListActionRuleInvoker 获取告警行动规则列表
func (c *AomClient) ListActionRuleInvoker(request *model.ListActionRuleRequest) *ListActionRuleInvoker {
	requestDef := GenReqDefForListActionRule()
	return &ListActionRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAgents 查询主机安装的ICAgent信息
//
// 该接口用于查询集群主机或用户自定义主机安装的ICAgent信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListAgents(request *model.ListAgentsRequest) (*model.ListAgentsResponse, error) {
	requestDef := GenReqDefForListAgents()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAgentsResponse), nil
	}
}

// ListAgentsInvoker 查询主机安装的ICAgent信息
func (c *AomClient) ListAgentsInvoker(request *model.ListAgentsRequest) *ListAgentsInvoker {
	requestDef := GenReqDefForListAgents()
	return &ListAgentsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAlarmRule 查询阈值规则列表
//
// 该接口用于查询阈值规则列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListAlarmRule(request *model.ListAlarmRuleRequest) (*model.ListAlarmRuleResponse, error) {
	requestDef := GenReqDefForListAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAlarmRuleResponse), nil
	}
}

// ListAlarmRuleInvoker 查询阈值规则列表
func (c *AomClient) ListAlarmRuleInvoker(request *model.ListAlarmRuleRequest) *ListAlarmRuleInvoker {
	requestDef := GenReqDefForListAlarmRule()
	return &ListAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListEvent2alarmRule 查询事件类告警规则列表
//
// 查询事件类告警规则列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListEvent2alarmRule(request *model.ListEvent2alarmRuleRequest) (*model.ListEvent2alarmRuleResponse, error) {
	requestDef := GenReqDefForListEvent2alarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListEvent2alarmRuleResponse), nil
	}
}

// ListEvent2alarmRuleInvoker 查询事件类告警规则列表
func (c *AomClient) ListEvent2alarmRuleInvoker(request *model.ListEvent2alarmRuleRequest) *ListEvent2alarmRuleInvoker {
	requestDef := GenReqDefForListEvent2alarmRule()
	return &ListEvent2alarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListEvents 查询事件告警信息
//
// 该接口用于查询对应用户的事件、告警。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListEvents(request *model.ListEventsRequest) (*model.ListEventsResponse, error) {
	requestDef := GenReqDefForListEvents()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListEventsResponse), nil
	}
}

// ListEventsInvoker 查询事件告警信息
func (c *AomClient) ListEventsInvoker(request *model.ListEventsRequest) *ListEventsInvoker {
	requestDef := GenReqDefForListEvents()
	return &ListEventsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListLogItems 查询日志
//
// 该接口用于查询不同维度(例如集群、IP、应用等)下的日志内容，支持分页查询。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListLogItems(request *model.ListLogItemsRequest) (*model.ListLogItemsResponse, error) {
	requestDef := GenReqDefForListLogItems()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListLogItemsResponse), nil
	}
}

// ListLogItemsInvoker 查询日志
func (c *AomClient) ListLogItemsInvoker(request *model.ListLogItemsRequest) *ListLogItemsInvoker {
	requestDef := GenReqDefForListLogItems()
	return &ListLogItemsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListMetricItems 查询指标
//
// 该接口用于查询系统当前可监控的指标列表，可以指定指标命名空间、指标名称、维度、所属资源的编号（格式为：resType_resId），分页查询的起始位置和返回的最大记录条数。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListMetricItems(request *model.ListMetricItemsRequest) (*model.ListMetricItemsResponse, error) {
	requestDef := GenReqDefForListMetricItems()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListMetricItemsResponse), nil
	}
}

// ListMetricItemsInvoker 查询指标
func (c *AomClient) ListMetricItemsInvoker(request *model.ListMetricItemsRequest) *ListMetricItemsInvoker {
	requestDef := GenReqDefForListMetricItems()
	return &ListMetricItemsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListMetricOrEventAlarmRule 查询指标类或者事件类告警规则列表
//
// 查询AOM2.0指标类或者事件类告警规则列表。(注：接口目前开放的region为：华东-上海一)
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListMetricOrEventAlarmRule(request *model.ListMetricOrEventAlarmRuleRequest) (*model.ListMetricOrEventAlarmRuleResponse, error) {
	requestDef := GenReqDefForListMetricOrEventAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListMetricOrEventAlarmRuleResponse), nil
	}
}

// ListMetricOrEventAlarmRuleInvoker 查询指标类或者事件类告警规则列表
func (c *AomClient) ListMetricOrEventAlarmRuleInvoker(request *model.ListMetricOrEventAlarmRuleRequest) *ListMetricOrEventAlarmRuleInvoker {
	requestDef := GenReqDefForListMetricOrEventAlarmRule()
	return &ListMetricOrEventAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListMuteRule 获取静默规则列表
//
// 获取静默规则列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListMuteRule(request *model.ListMuteRuleRequest) (*model.ListMuteRuleResponse, error) {
	requestDef := GenReqDefForListMuteRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListMuteRuleResponse), nil
	}
}

// ListMuteRuleInvoker 获取静默规则列表
func (c *AomClient) ListMuteRuleInvoker(request *model.ListMuteRuleRequest) *ListMuteRuleInvoker {
	requestDef := GenReqDefForListMuteRule()
	return &ListMuteRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListNotifiedHistories 获取告警发送结果
//
// 获取告警发送结果。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListNotifiedHistories(request *model.ListNotifiedHistoriesRequest) (*model.ListNotifiedHistoriesResponse, error) {
	requestDef := GenReqDefForListNotifiedHistories()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListNotifiedHistoriesResponse), nil
	}
}

// ListNotifiedHistoriesInvoker 获取告警发送结果
func (c *AomClient) ListNotifiedHistoriesInvoker(request *model.ListNotifiedHistoriesRequest) *ListNotifiedHistoriesInvoker {
	requestDef := GenReqDefForListNotifiedHistories()
	return &ListNotifiedHistoriesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPermissions 查询aom2.0相关云服务授权信息
//
// 该接口用于查询aom2.0相关云服务授权信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListPermissions(request *model.ListPermissionsRequest) (*model.ListPermissionsResponse, error) {
	requestDef := GenReqDefForListPermissions()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPermissionsResponse), nil
	}
}

// ListPermissionsInvoker 查询aom2.0相关云服务授权信息
func (c *AomClient) ListPermissionsInvoker(request *model.ListPermissionsRequest) *ListPermissionsInvoker {
	requestDef := GenReqDefForListPermissions()
	return &ListPermissionsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSample 查询时序数据
//
// 该接口用于查询指定时间范围内的监控时序数据，可以通过参数指定需要查询的数据维度，数据周期等。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListSample(request *model.ListSampleRequest) (*model.ListSampleResponse, error) {
	requestDef := GenReqDefForListSample()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSampleResponse), nil
	}
}

// ListSampleInvoker 查询时序数据
func (c *AomClient) ListSampleInvoker(request *model.ListSampleRequest) *ListSampleInvoker {
	requestDef := GenReqDefForListSample()
	return &ListSampleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSeries 查询时间序列
//
// 该接口用于查询系统当前可监控的时间序列列表，可以指定时间序列命名空间、名称、维度、所属资源的编号（格式为：resType_resId），分页查询的起始位置和返回的最大记录条数。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListSeries(request *model.ListSeriesRequest) (*model.ListSeriesResponse, error) {
	requestDef := GenReqDefForListSeries()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSeriesResponse), nil
	}
}

// ListSeriesInvoker 查询时间序列
func (c *AomClient) ListSeriesInvoker(request *model.ListSeriesRequest) *ListSeriesInvoker {
	requestDef := GenReqDefForListSeries()
	return &ListSeriesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListServiceDiscoveryRules 查询系统中已有服务发现规则
//
// 该接口用于查询系统当前已存在的服务发现规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListServiceDiscoveryRules(request *model.ListServiceDiscoveryRulesRequest) (*model.ListServiceDiscoveryRulesResponse, error) {
	requestDef := GenReqDefForListServiceDiscoveryRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListServiceDiscoveryRulesResponse), nil
	}
}

// ListServiceDiscoveryRulesInvoker 查询系统中已有服务发现规则
func (c *AomClient) ListServiceDiscoveryRulesInvoker(request *model.ListServiceDiscoveryRulesRequest) *ListServiceDiscoveryRulesInvoker {
	requestDef := GenReqDefForListServiceDiscoveryRules()
	return &ListServiceDiscoveryRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// PushEvents 上报事件告警信息
//
// 该接口用于上报对应用户的事件、告警。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) PushEvents(request *model.PushEventsRequest) (*model.PushEventsResponse, error) {
	requestDef := GenReqDefForPushEvents()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.PushEventsResponse), nil
	}
}

// PushEventsInvoker 上报事件告警信息
func (c *AomClient) PushEventsInvoker(request *model.PushEventsRequest) *PushEventsInvoker {
	requestDef := GenReqDefForPushEvents()
	return &PushEventsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowActionRule 通过规则名称获取告警行动规则
//
// 通过规则名称获取告警行动规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ShowActionRule(request *model.ShowActionRuleRequest) (*model.ShowActionRuleResponse, error) {
	requestDef := GenReqDefForShowActionRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowActionRuleResponse), nil
	}
}

// ShowActionRuleInvoker 通过规则名称获取告警行动规则
func (c *AomClient) ShowActionRuleInvoker(request *model.ShowActionRuleRequest) *ShowActionRuleInvoker {
	requestDef := GenReqDefForShowActionRule()
	return &ShowActionRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAlarmRule 查询单条阈值规则
//
// 该接口用于查询单条阈值规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ShowAlarmRule(request *model.ShowAlarmRuleRequest) (*model.ShowAlarmRuleResponse, error) {
	requestDef := GenReqDefForShowAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAlarmRuleResponse), nil
	}
}

// ShowAlarmRuleInvoker 查询单条阈值规则
func (c *AomClient) ShowAlarmRuleInvoker(request *model.ShowAlarmRuleRequest) *ShowAlarmRuleInvoker {
	requestDef := GenReqDefForShowAlarmRule()
	return &ShowAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowMetricsData 查询监控数据
//
// 该接口用于查询指定时间范围内指标的监控数据，可以通过参数指定需要查询的数据维度，数据周期等。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ShowMetricsData(request *model.ShowMetricsDataRequest) (*model.ShowMetricsDataResponse, error) {
	requestDef := GenReqDefForShowMetricsData()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowMetricsDataResponse), nil
	}
}

// ShowMetricsDataInvoker 查询监控数据
func (c *AomClient) ShowMetricsDataInvoker(request *model.ShowMetricsDataRequest) *ShowMetricsDataInvoker {
	requestDef := GenReqDefForShowMetricsData()
	return &ShowMetricsDataInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateActionRule 修改告警行动规则
//
// 修改告警行动规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) UpdateActionRule(request *model.UpdateActionRuleRequest) (*model.UpdateActionRuleResponse, error) {
	requestDef := GenReqDefForUpdateActionRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateActionRuleResponse), nil
	}
}

// UpdateActionRuleInvoker 修改告警行动规则
func (c *AomClient) UpdateActionRuleInvoker(request *model.UpdateActionRuleRequest) *UpdateActionRuleInvoker {
	requestDef := GenReqDefForUpdateActionRule()
	return &UpdateActionRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAlarmRule 修改阈值规则
//
// 该接口用于修改一条阈值规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) UpdateAlarmRule(request *model.UpdateAlarmRuleRequest) (*model.UpdateAlarmRuleResponse, error) {
	requestDef := GenReqDefForUpdateAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAlarmRuleResponse), nil
	}
}

// UpdateAlarmRuleInvoker 修改阈值规则
func (c *AomClient) UpdateAlarmRuleInvoker(request *model.UpdateAlarmRuleRequest) *UpdateAlarmRuleInvoker {
	requestDef := GenReqDefForUpdateAlarmRule()
	return &UpdateAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateEventRule 更新事件类告警规则
//
// 更新事件类告警规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) UpdateEventRule(request *model.UpdateEventRuleRequest) (*model.UpdateEventRuleResponse, error) {
	requestDef := GenReqDefForUpdateEventRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateEventRuleResponse), nil
	}
}

// UpdateEventRuleInvoker 更新事件类告警规则
func (c *AomClient) UpdateEventRuleInvoker(request *model.UpdateEventRuleRequest) *UpdateEventRuleInvoker {
	requestDef := GenReqDefForUpdateEventRule()
	return &UpdateEventRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateMuteRule 修改静默规则
//
// 修改静默规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) UpdateMuteRule(request *model.UpdateMuteRuleRequest) (*model.UpdateMuteRuleResponse, error) {
	requestDef := GenReqDefForUpdateMuteRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateMuteRuleResponse), nil
	}
}

// UpdateMuteRuleInvoker 修改静默规则
func (c *AomClient) UpdateMuteRuleInvoker(request *model.UpdateMuteRuleRequest) *UpdateMuteRuleInvoker {
	requestDef := GenReqDefForUpdateMuteRule()
	return &UpdateMuteRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreatePromInstance 新增Prometheus实例
//
// 该接口用于新增Prometheus实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) CreatePromInstance(request *model.CreatePromInstanceRequest) (*model.CreatePromInstanceResponse, error) {
	requestDef := GenReqDefForCreatePromInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePromInstanceResponse), nil
	}
}

// CreatePromInstanceInvoker 新增Prometheus实例
func (c *AomClient) CreatePromInstanceInvoker(request *model.CreatePromInstanceRequest) *CreatePromInstanceInvoker {
	requestDef := GenReqDefForCreatePromInstance()
	return &CreatePromInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRecordingRule 创建Prometheus实例的预聚合规则
//
// 该接口用于给Prometheus实例创建预聚合规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) CreateRecordingRule(request *model.CreateRecordingRuleRequest) (*model.CreateRecordingRuleResponse, error) {
	requestDef := GenReqDefForCreateRecordingRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRecordingRuleResponse), nil
	}
}

// CreateRecordingRuleInvoker 创建Prometheus实例的预聚合规则
func (c *AomClient) CreateRecordingRuleInvoker(request *model.CreateRecordingRuleRequest) *CreateRecordingRuleInvoker {
	requestDef := GenReqDefForCreateRecordingRule()
	return &CreateRecordingRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeletePromInstance 卸载托管Prometheus实例
//
// 该接口用于卸载托管Prometheus实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) DeletePromInstance(request *model.DeletePromInstanceRequest) (*model.DeletePromInstanceResponse, error) {
	requestDef := GenReqDefForDeletePromInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeletePromInstanceResponse), nil
	}
}

// DeletePromInstanceInvoker 卸载托管Prometheus实例
func (c *AomClient) DeletePromInstanceInvoker(request *model.DeletePromInstanceRequest) *DeletePromInstanceInvoker {
	requestDef := GenReqDefForDeletePromInstance()
	return &DeletePromInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAccessCode 获取Prometheus实例调用凭证
//
// 该接口用于获取Prometheus实例调用凭证。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListAccessCode(request *model.ListAccessCodeRequest) (*model.ListAccessCodeResponse, error) {
	requestDef := GenReqDefForListAccessCode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAccessCodeResponse), nil
	}
}

// ListAccessCodeInvoker 获取Prometheus实例调用凭证
func (c *AomClient) ListAccessCodeInvoker(request *model.ListAccessCodeRequest) *ListAccessCodeInvoker {
	requestDef := GenReqDefForListAccessCode()
	return &ListAccessCodeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListInstantQueryAomPromGet GET方法查询瞬时数据
//
// 该接口使用GET方法查询PromQL(Prometheus Query Language)在特定时间点下的计算结果。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListInstantQueryAomPromGet(request *model.ListInstantQueryAomPromGetRequest) (*model.ListInstantQueryAomPromGetResponse, error) {
	requestDef := GenReqDefForListInstantQueryAomPromGet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstantQueryAomPromGetResponse), nil
	}
}

// ListInstantQueryAomPromGetInvoker GET方法查询瞬时数据
func (c *AomClient) ListInstantQueryAomPromGetInvoker(request *model.ListInstantQueryAomPromGetRequest) *ListInstantQueryAomPromGetInvoker {
	requestDef := GenReqDefForListInstantQueryAomPromGet()
	return &ListInstantQueryAomPromGetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListInstantQueryAomPromPost （推荐）POST方法查询瞬时数据
//
// 该接口使用POST方法查询PromQL(Prometheus Query Language) 在特定时间点下的计算结果。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListInstantQueryAomPromPost(request *model.ListInstantQueryAomPromPostRequest) (*model.ListInstantQueryAomPromPostResponse, error) {
	requestDef := GenReqDefForListInstantQueryAomPromPost()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstantQueryAomPromPostResponse), nil
	}
}

// ListInstantQueryAomPromPostInvoker （推荐）POST方法查询瞬时数据
func (c *AomClient) ListInstantQueryAomPromPostInvoker(request *model.ListInstantQueryAomPromPostRequest) *ListInstantQueryAomPromPostInvoker {
	requestDef := GenReqDefForListInstantQueryAomPromPost()
	return &ListInstantQueryAomPromPostInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListLabelValuesAomPromGet 查询标签值
//
// 该接口用于查询带有指定标签的时间序列列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListLabelValuesAomPromGet(request *model.ListLabelValuesAomPromGetRequest) (*model.ListLabelValuesAomPromGetResponse, error) {
	requestDef := GenReqDefForListLabelValuesAomPromGet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListLabelValuesAomPromGetResponse), nil
	}
}

// ListLabelValuesAomPromGetInvoker 查询标签值
func (c *AomClient) ListLabelValuesAomPromGetInvoker(request *model.ListLabelValuesAomPromGetRequest) *ListLabelValuesAomPromGetInvoker {
	requestDef := GenReqDefForListLabelValuesAomPromGet()
	return &ListLabelValuesAomPromGetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListLabelsAomPromGet GET方法获取标签名列表
//
// 该接口使用GET方法获取标签名列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListLabelsAomPromGet(request *model.ListLabelsAomPromGetRequest) (*model.ListLabelsAomPromGetResponse, error) {
	requestDef := GenReqDefForListLabelsAomPromGet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListLabelsAomPromGetResponse), nil
	}
}

// ListLabelsAomPromGetInvoker GET方法获取标签名列表
func (c *AomClient) ListLabelsAomPromGetInvoker(request *model.ListLabelsAomPromGetRequest) *ListLabelsAomPromGetInvoker {
	requestDef := GenReqDefForListLabelsAomPromGet()
	return &ListLabelsAomPromGetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListLabelsAomPromPost （推荐）POST方法获取标签名列表
//
// 该接口使用POST方法获取标签名列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListLabelsAomPromPost(request *model.ListLabelsAomPromPostRequest) (*model.ListLabelsAomPromPostResponse, error) {
	requestDef := GenReqDefForListLabelsAomPromPost()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListLabelsAomPromPostResponse), nil
	}
}

// ListLabelsAomPromPostInvoker （推荐）POST方法获取标签名列表
func (c *AomClient) ListLabelsAomPromPostInvoker(request *model.ListLabelsAomPromPostRequest) *ListLabelsAomPromPostInvoker {
	requestDef := GenReqDefForListLabelsAomPromPost()
	return &ListLabelsAomPromPostInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListMetadataAomPromGet 元数据查询
//
// 该接口用于查询序列及序列标签的元数据。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListMetadataAomPromGet(request *model.ListMetadataAomPromGetRequest) (*model.ListMetadataAomPromGetResponse, error) {
	requestDef := GenReqDefForListMetadataAomPromGet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListMetadataAomPromGetResponse), nil
	}
}

// ListMetadataAomPromGetInvoker 元数据查询
func (c *AomClient) ListMetadataAomPromGetInvoker(request *model.ListMetadataAomPromGetRequest) *ListMetadataAomPromGetInvoker {
	requestDef := GenReqDefForListMetadataAomPromGet()
	return &ListMetadataAomPromGetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPromInstance 查询Prometheus实例
//
// 该接口用于查询Prometheus实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListPromInstance(request *model.ListPromInstanceRequest) (*model.ListPromInstanceResponse, error) {
	requestDef := GenReqDefForListPromInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPromInstanceResponse), nil
	}
}

// ListPromInstanceInvoker 查询Prometheus实例
func (c *AomClient) ListPromInstanceInvoker(request *model.ListPromInstanceRequest) *ListPromInstanceInvoker {
	requestDef := GenReqDefForListPromInstance()
	return &ListPromInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRangeQueryAomPromGet GET方法查询区间数据
//
// 该接口使用GET方法查询PromQL(Prometheus Query Language)在一段时间返回内的计算结果。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListRangeQueryAomPromGet(request *model.ListRangeQueryAomPromGetRequest) (*model.ListRangeQueryAomPromGetResponse, error) {
	requestDef := GenReqDefForListRangeQueryAomPromGet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRangeQueryAomPromGetResponse), nil
	}
}

// ListRangeQueryAomPromGetInvoker GET方法查询区间数据
func (c *AomClient) ListRangeQueryAomPromGetInvoker(request *model.ListRangeQueryAomPromGetRequest) *ListRangeQueryAomPromGetInvoker {
	requestDef := GenReqDefForListRangeQueryAomPromGet()
	return &ListRangeQueryAomPromGetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRangeQueryAomPromPost （推荐）POST方法查询区间数据
//
// 该接口使用POST方法查询PromQL(Prometheus Query Language)在一段时间返回内的计算结果。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *AomClient) ListRangeQueryAomPromPost(request *model.ListRangeQueryAomPromPostRequest) (*model.ListRangeQueryAomPromPostResponse, error) {
	requestDef := GenReqDefForListRangeQueryAomPromPost()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRangeQueryAomPromPostResponse), nil
	}
}

// ListRangeQueryAomPromPostInvoker （推荐）POST方法查询区间数据
func (c *AomClient) ListRangeQueryAomPromPostInvoker(request *model.ListRangeQueryAomPromPostRequest) *ListRangeQueryAomPromPostInvoker {
	requestDef := GenReqDefForListRangeQueryAomPromPost()
	return &ListRangeQueryAomPromPostInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
