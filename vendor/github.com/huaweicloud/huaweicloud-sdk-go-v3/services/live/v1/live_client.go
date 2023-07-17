package v1

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1/model"
)

type LiveClient struct {
	HcClient *http_client.HcHttpClient
}

func NewLiveClient(hcClient *http_client.HcHttpClient) *LiveClient {
	return &LiveClient{HcClient: hcClient}
}

func LiveClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// BatchShowIpBelongs 查询IP归属信息
//
// 查询IP归属信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) BatchShowIpBelongs(request *model.BatchShowIpBelongsRequest) (*model.BatchShowIpBelongsResponse, error) {
	requestDef := GenReqDefForBatchShowIpBelongs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchShowIpBelongsResponse), nil
	}
}

// BatchShowIpBelongsInvoker 查询IP归属信息
func (c *LiveClient) BatchShowIpBelongsInvoker(request *model.BatchShowIpBelongsRequest) *BatchShowIpBelongsInvoker {
	requestDef := GenReqDefForBatchShowIpBelongs()
	return &BatchShowIpBelongsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateDomain 创建直播域名
//
// 可单独创建直播播放域名或推流域名，每个租户最多可配置64条域名记录。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) CreateDomain(request *model.CreateDomainRequest) (*model.CreateDomainResponse, error) {
	requestDef := GenReqDefForCreateDomain()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDomainResponse), nil
	}
}

// CreateDomainInvoker 创建直播域名
func (c *LiveClient) CreateDomainInvoker(request *model.CreateDomainRequest) *CreateDomainInvoker {
	requestDef := GenReqDefForCreateDomain()
	return &CreateDomainInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateDomainMapping 域名映射
//
// 将用户已创建的播放域名和推流域名建立域名映射关系
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) CreateDomainMapping(request *model.CreateDomainMappingRequest) (*model.CreateDomainMappingResponse, error) {
	requestDef := GenReqDefForCreateDomainMapping()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDomainMappingResponse), nil
	}
}

// CreateDomainMappingInvoker 域名映射
func (c *LiveClient) CreateDomainMappingInvoker(request *model.CreateDomainMappingRequest) *CreateDomainMappingInvoker {
	requestDef := GenReqDefForCreateDomainMapping()
	return &CreateDomainMappingInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRecordCallbackConfig 创建录制回调配置
//
// 创建录制回调配置接口
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) CreateRecordCallbackConfig(request *model.CreateRecordCallbackConfigRequest) (*model.CreateRecordCallbackConfigResponse, error) {
	requestDef := GenReqDefForCreateRecordCallbackConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRecordCallbackConfigResponse), nil
	}
}

// CreateRecordCallbackConfigInvoker 创建录制回调配置
func (c *LiveClient) CreateRecordCallbackConfigInvoker(request *model.CreateRecordCallbackConfigRequest) *CreateRecordCallbackConfigInvoker {
	requestDef := GenReqDefForCreateRecordCallbackConfig()
	return &CreateRecordCallbackConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRecordIndex 创建录制视频索引文件
//
// Create Record Index
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) CreateRecordIndex(request *model.CreateRecordIndexRequest) (*model.CreateRecordIndexResponse, error) {
	requestDef := GenReqDefForCreateRecordIndex()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRecordIndexResponse), nil
	}
}

// CreateRecordIndexInvoker 创建录制视频索引文件
func (c *LiveClient) CreateRecordIndexInvoker(request *model.CreateRecordIndexRequest) *CreateRecordIndexInvoker {
	requestDef := GenReqDefForCreateRecordIndex()
	return &CreateRecordIndexInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRecordRule 创建录制规则
//
// 创建录制规则接口，录制规则对新推送的流生效，对已经推送中的流不生效
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) CreateRecordRule(request *model.CreateRecordRuleRequest) (*model.CreateRecordRuleResponse, error) {
	requestDef := GenReqDefForCreateRecordRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRecordRuleResponse), nil
	}
}

// CreateRecordRuleInvoker 创建录制规则
func (c *LiveClient) CreateRecordRuleInvoker(request *model.CreateRecordRuleRequest) *CreateRecordRuleInvoker {
	requestDef := GenReqDefForCreateRecordRule()
	return &CreateRecordRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateSnapshotConfig 创建直播截图配置
//
// 创建直播截图配置接口
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) CreateSnapshotConfig(request *model.CreateSnapshotConfigRequest) (*model.CreateSnapshotConfigResponse, error) {
	requestDef := GenReqDefForCreateSnapshotConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSnapshotConfigResponse), nil
	}
}

// CreateSnapshotConfigInvoker 创建直播截图配置
func (c *LiveClient) CreateSnapshotConfigInvoker(request *model.CreateSnapshotConfigRequest) *CreateSnapshotConfigInvoker {
	requestDef := GenReqDefForCreateSnapshotConfig()
	return &CreateSnapshotConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateStreamForbidden 禁止直播推流
//
// 禁止直播推流
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) CreateStreamForbidden(request *model.CreateStreamForbiddenRequest) (*model.CreateStreamForbiddenResponse, error) {
	requestDef := GenReqDefForCreateStreamForbidden()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateStreamForbiddenResponse), nil
	}
}

// CreateStreamForbiddenInvoker 禁止直播推流
func (c *LiveClient) CreateStreamForbiddenInvoker(request *model.CreateStreamForbiddenRequest) *CreateStreamForbiddenInvoker {
	requestDef := GenReqDefForCreateStreamForbidden()
	return &CreateStreamForbiddenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTranscodingsTemplate 创建直播转码模板
//
// 创建直播转码模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) CreateTranscodingsTemplate(request *model.CreateTranscodingsTemplateRequest) (*model.CreateTranscodingsTemplateResponse, error) {
	requestDef := GenReqDefForCreateTranscodingsTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTranscodingsTemplateResponse), nil
	}
}

// CreateTranscodingsTemplateInvoker 创建直播转码模板
func (c *LiveClient) CreateTranscodingsTemplateInvoker(request *model.CreateTranscodingsTemplateRequest) *CreateTranscodingsTemplateInvoker {
	requestDef := GenReqDefForCreateTranscodingsTemplate()
	return &CreateTranscodingsTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDomain 删除直播域名
//
// 删除域名。只有在域名停用（off）状态时才能删除。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) DeleteDomain(request *model.DeleteDomainRequest) (*model.DeleteDomainResponse, error) {
	requestDef := GenReqDefForDeleteDomain()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDomainResponse), nil
	}
}

// DeleteDomainInvoker 删除直播域名
func (c *LiveClient) DeleteDomainInvoker(request *model.DeleteDomainRequest) *DeleteDomainInvoker {
	requestDef := GenReqDefForDeleteDomain()
	return &DeleteDomainInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDomainKeyChain 删除指定域名的Key防盗链配置
//
// 删除指定域名的Key防盗链配置
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) DeleteDomainKeyChain(request *model.DeleteDomainKeyChainRequest) (*model.DeleteDomainKeyChainResponse, error) {
	requestDef := GenReqDefForDeleteDomainKeyChain()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDomainKeyChainResponse), nil
	}
}

// DeleteDomainKeyChainInvoker 删除指定域名的Key防盗链配置
func (c *LiveClient) DeleteDomainKeyChainInvoker(request *model.DeleteDomainKeyChainRequest) *DeleteDomainKeyChainInvoker {
	requestDef := GenReqDefForDeleteDomainKeyChain()
	return &DeleteDomainKeyChainInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDomainMapping 删除直播域名映射关系
//
// 将播放域名和推流域名的域名映射关系删除
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) DeleteDomainMapping(request *model.DeleteDomainMappingRequest) (*model.DeleteDomainMappingResponse, error) {
	requestDef := GenReqDefForDeleteDomainMapping()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDomainMappingResponse), nil
	}
}

// DeleteDomainMappingInvoker 删除直播域名映射关系
func (c *LiveClient) DeleteDomainMappingInvoker(request *model.DeleteDomainMappingRequest) *DeleteDomainMappingInvoker {
	requestDef := GenReqDefForDeleteDomainMapping()
	return &DeleteDomainMappingInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteRecordCallbackConfig 删除录制回调配置
//
// 删除录制回调配置接口
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) DeleteRecordCallbackConfig(request *model.DeleteRecordCallbackConfigRequest) (*model.DeleteRecordCallbackConfigResponse, error) {
	requestDef := GenReqDefForDeleteRecordCallbackConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteRecordCallbackConfigResponse), nil
	}
}

// DeleteRecordCallbackConfigInvoker 删除录制回调配置
func (c *LiveClient) DeleteRecordCallbackConfigInvoker(request *model.DeleteRecordCallbackConfigRequest) *DeleteRecordCallbackConfigInvoker {
	requestDef := GenReqDefForDeleteRecordCallbackConfig()
	return &DeleteRecordCallbackConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteRecordRule 删除录制规则
//
// 删除录制规则接口
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) DeleteRecordRule(request *model.DeleteRecordRuleRequest) (*model.DeleteRecordRuleResponse, error) {
	requestDef := GenReqDefForDeleteRecordRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteRecordRuleResponse), nil
	}
}

// DeleteRecordRuleInvoker 删除录制规则
func (c *LiveClient) DeleteRecordRuleInvoker(request *model.DeleteRecordRuleRequest) *DeleteRecordRuleInvoker {
	requestDef := GenReqDefForDeleteRecordRule()
	return &DeleteRecordRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteSnapshotConfig 删除直播截图配置
//
// 删除直播截图配置接口
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) DeleteSnapshotConfig(request *model.DeleteSnapshotConfigRequest) (*model.DeleteSnapshotConfigResponse, error) {
	requestDef := GenReqDefForDeleteSnapshotConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSnapshotConfigResponse), nil
	}
}

// DeleteSnapshotConfigInvoker 删除直播截图配置
func (c *LiveClient) DeleteSnapshotConfigInvoker(request *model.DeleteSnapshotConfigRequest) *DeleteSnapshotConfigInvoker {
	requestDef := GenReqDefForDeleteSnapshotConfig()
	return &DeleteSnapshotConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteStreamForbidden 禁推恢复
//
// 恢复直播推流接口
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) DeleteStreamForbidden(request *model.DeleteStreamForbiddenRequest) (*model.DeleteStreamForbiddenResponse, error) {
	requestDef := GenReqDefForDeleteStreamForbidden()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteStreamForbiddenResponse), nil
	}
}

// DeleteStreamForbiddenInvoker 禁推恢复
func (c *LiveClient) DeleteStreamForbiddenInvoker(request *model.DeleteStreamForbiddenRequest) *DeleteStreamForbiddenInvoker {
	requestDef := GenReqDefForDeleteStreamForbidden()
	return &DeleteStreamForbiddenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTranscodingsTemplate 删除直播转码模板
//
// 删除直播转码模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) DeleteTranscodingsTemplate(request *model.DeleteTranscodingsTemplateRequest) (*model.DeleteTranscodingsTemplateResponse, error) {
	requestDef := GenReqDefForDeleteTranscodingsTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTranscodingsTemplateResponse), nil
	}
}

// DeleteTranscodingsTemplateInvoker 删除直播转码模板
func (c *LiveClient) DeleteTranscodingsTemplateInvoker(request *model.DeleteTranscodingsTemplateRequest) *DeleteTranscodingsTemplateInvoker {
	requestDef := GenReqDefForDeleteTranscodingsTemplate()
	return &DeleteTranscodingsTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListLiveSampleLogs 获取直播播放日志
//
// 获取直播播放日志，基于域名以5分钟粒度进行打包，日志内容以 \&quot;|\&quot; 进行分隔。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ListLiveSampleLogs(request *model.ListLiveSampleLogsRequest) (*model.ListLiveSampleLogsResponse, error) {
	requestDef := GenReqDefForListLiveSampleLogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListLiveSampleLogsResponse), nil
	}
}

// ListLiveSampleLogsInvoker 获取直播播放日志
func (c *LiveClient) ListLiveSampleLogsInvoker(request *model.ListLiveSampleLogsRequest) *ListLiveSampleLogsInvoker {
	requestDef := GenReqDefForListLiveSampleLogs()
	return &ListLiveSampleLogsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListLiveStreamsOnline 查询直播中的流信息
//
// 查询直播中的流信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ListLiveStreamsOnline(request *model.ListLiveStreamsOnlineRequest) (*model.ListLiveStreamsOnlineResponse, error) {
	requestDef := GenReqDefForListLiveStreamsOnline()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListLiveStreamsOnlineResponse), nil
	}
}

// ListLiveStreamsOnlineInvoker 查询直播中的流信息
func (c *LiveClient) ListLiveStreamsOnlineInvoker(request *model.ListLiveStreamsOnlineRequest) *ListLiveStreamsOnlineInvoker {
	requestDef := GenReqDefForListLiveStreamsOnline()
	return &ListLiveStreamsOnlineInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRecordCallbackConfigs 查询录制回调配置列表
//
// 查询录制回调配置列表接口。通过指定条件，查询满足条件的配置列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ListRecordCallbackConfigs(request *model.ListRecordCallbackConfigsRequest) (*model.ListRecordCallbackConfigsResponse, error) {
	requestDef := GenReqDefForListRecordCallbackConfigs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRecordCallbackConfigsResponse), nil
	}
}

// ListRecordCallbackConfigsInvoker 查询录制回调配置列表
func (c *LiveClient) ListRecordCallbackConfigsInvoker(request *model.ListRecordCallbackConfigsRequest) *ListRecordCallbackConfigsInvoker {
	requestDef := GenReqDefForListRecordCallbackConfigs()
	return &ListRecordCallbackConfigsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRecordContents 录制完成内容的查询
//
// 录制完成的内容查询
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ListRecordContents(request *model.ListRecordContentsRequest) (*model.ListRecordContentsResponse, error) {
	requestDef := GenReqDefForListRecordContents()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRecordContentsResponse), nil
	}
}

// ListRecordContentsInvoker 录制完成内容的查询
func (c *LiveClient) ListRecordContentsInvoker(request *model.ListRecordContentsRequest) *ListRecordContentsInvoker {
	requestDef := GenReqDefForListRecordContents()
	return &ListRecordContentsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRecordRules 查询录制规则列表
//
// 查询录制规则列表接口，通过指定条件，查询满足条件的录制规则列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ListRecordRules(request *model.ListRecordRulesRequest) (*model.ListRecordRulesResponse, error) {
	requestDef := GenReqDefForListRecordRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRecordRulesResponse), nil
	}
}

// ListRecordRulesInvoker 查询录制规则列表
func (c *LiveClient) ListRecordRulesInvoker(request *model.ListRecordRulesRequest) *ListRecordRulesInvoker {
	requestDef := GenReqDefForListRecordRules()
	return &ListRecordRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSnapshotConfigs 查询直播截图配置
//
// 查询直播截图配置接口
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ListSnapshotConfigs(request *model.ListSnapshotConfigsRequest) (*model.ListSnapshotConfigsResponse, error) {
	requestDef := GenReqDefForListSnapshotConfigs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSnapshotConfigsResponse), nil
	}
}

// ListSnapshotConfigsInvoker 查询直播截图配置
func (c *LiveClient) ListSnapshotConfigsInvoker(request *model.ListSnapshotConfigsRequest) *ListSnapshotConfigsInvoker {
	requestDef := GenReqDefForListSnapshotConfigs()
	return &ListSnapshotConfigsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListStreamForbidden 查询禁止直播推流列表
//
// 查询禁播黑名单列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ListStreamForbidden(request *model.ListStreamForbiddenRequest) (*model.ListStreamForbiddenResponse, error) {
	requestDef := GenReqDefForListStreamForbidden()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListStreamForbiddenResponse), nil
	}
}

// ListStreamForbiddenInvoker 查询禁止直播推流列表
func (c *LiveClient) ListStreamForbiddenInvoker(request *model.ListStreamForbiddenRequest) *ListStreamForbiddenInvoker {
	requestDef := GenReqDefForListStreamForbidden()
	return &ListStreamForbiddenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RunRecord 提交录制控制命令
//
// 对单条流的实时录制控制接口。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) RunRecord(request *model.RunRecordRequest) (*model.RunRecordResponse, error) {
	requestDef := GenReqDefForRunRecord()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RunRecordResponse), nil
	}
}

// RunRecordInvoker 提交录制控制命令
func (c *LiveClient) RunRecordInvoker(request *model.RunRecordRequest) *RunRecordInvoker {
	requestDef := GenReqDefForRunRecord()
	return &RunRecordInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDomain 查询直播域名
//
// 查询直播域名
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ShowDomain(request *model.ShowDomainRequest) (*model.ShowDomainResponse, error) {
	requestDef := GenReqDefForShowDomain()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainResponse), nil
	}
}

// ShowDomainInvoker 查询直播域名
func (c *LiveClient) ShowDomainInvoker(request *model.ShowDomainRequest) *ShowDomainInvoker {
	requestDef := GenReqDefForShowDomain()
	return &ShowDomainInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDomainKeyChain 查询指定域名的Key防盗链配置
//
// 查询指定域名的Key防盗链配置
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ShowDomainKeyChain(request *model.ShowDomainKeyChainRequest) (*model.ShowDomainKeyChainResponse, error) {
	requestDef := GenReqDefForShowDomainKeyChain()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainKeyChainResponse), nil
	}
}

// ShowDomainKeyChainInvoker 查询指定域名的Key防盗链配置
func (c *LiveClient) ShowDomainKeyChainInvoker(request *model.ShowDomainKeyChainRequest) *ShowDomainKeyChainInvoker {
	requestDef := GenReqDefForShowDomainKeyChain()
	return &ShowDomainKeyChainInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowRecordCallbackConfig 查询录制回调配置
//
// 查询录制回调配置接口
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ShowRecordCallbackConfig(request *model.ShowRecordCallbackConfigRequest) (*model.ShowRecordCallbackConfigResponse, error) {
	requestDef := GenReqDefForShowRecordCallbackConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowRecordCallbackConfigResponse), nil
	}
}

// ShowRecordCallbackConfigInvoker 查询录制回调配置
func (c *LiveClient) ShowRecordCallbackConfigInvoker(request *model.ShowRecordCallbackConfigRequest) *ShowRecordCallbackConfigInvoker {
	requestDef := GenReqDefForShowRecordCallbackConfig()
	return &ShowRecordCallbackConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowRecordRule 查询录制规则配置
//
// 查询录制规则接口
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ShowRecordRule(request *model.ShowRecordRuleRequest) (*model.ShowRecordRuleResponse, error) {
	requestDef := GenReqDefForShowRecordRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowRecordRuleResponse), nil
	}
}

// ShowRecordRuleInvoker 查询录制规则配置
func (c *LiveClient) ShowRecordRuleInvoker(request *model.ShowRecordRuleRequest) *ShowRecordRuleInvoker {
	requestDef := GenReqDefForShowRecordRule()
	return &ShowRecordRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTranscodingsTemplate 查询直播转码模板
//
// 查询直播转码模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ShowTranscodingsTemplate(request *model.ShowTranscodingsTemplateRequest) (*model.ShowTranscodingsTemplateResponse, error) {
	requestDef := GenReqDefForShowTranscodingsTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTranscodingsTemplateResponse), nil
	}
}

// ShowTranscodingsTemplateInvoker 查询直播转码模板
func (c *LiveClient) ShowTranscodingsTemplateInvoker(request *model.ShowTranscodingsTemplateRequest) *ShowTranscodingsTemplateInvoker {
	requestDef := GenReqDefForShowTranscodingsTemplate()
	return &ShowTranscodingsTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDomain 修改直播域名
//
// 修改直播播放、RTMP推流加速域名相关信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) UpdateDomain(request *model.UpdateDomainRequest) (*model.UpdateDomainResponse, error) {
	requestDef := GenReqDefForUpdateDomain()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainResponse), nil
	}
}

// UpdateDomainInvoker 修改直播域名
func (c *LiveClient) UpdateDomainInvoker(request *model.UpdateDomainRequest) *UpdateDomainInvoker {
	requestDef := GenReqDefForUpdateDomain()
	return &UpdateDomainInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDomainIp6Switch 配置域名IPV6开关
//
// 配置IPV6开关
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) UpdateDomainIp6Switch(request *model.UpdateDomainIp6SwitchRequest) (*model.UpdateDomainIp6SwitchResponse, error) {
	requestDef := GenReqDefForUpdateDomainIp6Switch()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainIp6SwitchResponse), nil
	}
}

// UpdateDomainIp6SwitchInvoker 配置域名IPV6开关
func (c *LiveClient) UpdateDomainIp6SwitchInvoker(request *model.UpdateDomainIp6SwitchRequest) *UpdateDomainIp6SwitchInvoker {
	requestDef := GenReqDefForUpdateDomainIp6Switch()
	return &UpdateDomainIp6SwitchInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDomainKeyChain 更新指定域名的Key防盗链配置
//
// 更新指定域名的Key防盗链配置
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) UpdateDomainKeyChain(request *model.UpdateDomainKeyChainRequest) (*model.UpdateDomainKeyChainResponse, error) {
	requestDef := GenReqDefForUpdateDomainKeyChain()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainKeyChainResponse), nil
	}
}

// UpdateDomainKeyChainInvoker 更新指定域名的Key防盗链配置
func (c *LiveClient) UpdateDomainKeyChainInvoker(request *model.UpdateDomainKeyChainRequest) *UpdateDomainKeyChainInvoker {
	requestDef := GenReqDefForUpdateDomainKeyChain()
	return &UpdateDomainKeyChainInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateRecordCallbackConfig 修改录制回调配置
//
// 修改录制回调配置接口
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) UpdateRecordCallbackConfig(request *model.UpdateRecordCallbackConfigRequest) (*model.UpdateRecordCallbackConfigResponse, error) {
	requestDef := GenReqDefForUpdateRecordCallbackConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateRecordCallbackConfigResponse), nil
	}
}

// UpdateRecordCallbackConfigInvoker 修改录制回调配置
func (c *LiveClient) UpdateRecordCallbackConfigInvoker(request *model.UpdateRecordCallbackConfigRequest) *UpdateRecordCallbackConfigInvoker {
	requestDef := GenReqDefForUpdateRecordCallbackConfig()
	return &UpdateRecordCallbackConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateRecordRule 修改录制规则
//
// 修改录制规则接口，如果规则修改后，修改后的规则对正在录制的流无效，对新的流有效。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) UpdateRecordRule(request *model.UpdateRecordRuleRequest) (*model.UpdateRecordRuleResponse, error) {
	requestDef := GenReqDefForUpdateRecordRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateRecordRuleResponse), nil
	}
}

// UpdateRecordRuleInvoker 修改录制规则
func (c *LiveClient) UpdateRecordRuleInvoker(request *model.UpdateRecordRuleRequest) *UpdateRecordRuleInvoker {
	requestDef := GenReqDefForUpdateRecordRule()
	return &UpdateRecordRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateSnapshotConfig 修改直播截图配置
//
// 修改直播截图配置接口
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) UpdateSnapshotConfig(request *model.UpdateSnapshotConfigRequest) (*model.UpdateSnapshotConfigResponse, error) {
	requestDef := GenReqDefForUpdateSnapshotConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateSnapshotConfigResponse), nil
	}
}

// UpdateSnapshotConfigInvoker 修改直播截图配置
func (c *LiveClient) UpdateSnapshotConfigInvoker(request *model.UpdateSnapshotConfigRequest) *UpdateSnapshotConfigInvoker {
	requestDef := GenReqDefForUpdateSnapshotConfig()
	return &UpdateSnapshotConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateStreamForbidden 修改禁推属性
//
// 修改禁推属性
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) UpdateStreamForbidden(request *model.UpdateStreamForbiddenRequest) (*model.UpdateStreamForbiddenResponse, error) {
	requestDef := GenReqDefForUpdateStreamForbidden()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateStreamForbiddenResponse), nil
	}
}

// UpdateStreamForbiddenInvoker 修改禁推属性
func (c *LiveClient) UpdateStreamForbiddenInvoker(request *model.UpdateStreamForbiddenRequest) *UpdateStreamForbiddenInvoker {
	requestDef := GenReqDefForUpdateStreamForbidden()
	return &UpdateStreamForbiddenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTranscodingsTemplate 配置直播转码模板
//
// 修改直播转码模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) UpdateTranscodingsTemplate(request *model.UpdateTranscodingsTemplateRequest) (*model.UpdateTranscodingsTemplateResponse, error) {
	requestDef := GenReqDefForUpdateTranscodingsTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTranscodingsTemplateResponse), nil
	}
}

// UpdateTranscodingsTemplateInvoker 配置直播转码模板
func (c *LiveClient) UpdateTranscodingsTemplateInvoker(request *model.UpdateTranscodingsTemplateRequest) *UpdateTranscodingsTemplateInvoker {
	requestDef := GenReqDefForUpdateTranscodingsTemplate()
	return &UpdateTranscodingsTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDomainHttpsCert 删除指定域名的https证书配置
//
// 删除指定域名的https证书配置
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) DeleteDomainHttpsCert(request *model.DeleteDomainHttpsCertRequest) (*model.DeleteDomainHttpsCertResponse, error) {
	requestDef := GenReqDefForDeleteDomainHttpsCert()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDomainHttpsCertResponse), nil
	}
}

// DeleteDomainHttpsCertInvoker 删除指定域名的https证书配置
func (c *LiveClient) DeleteDomainHttpsCertInvoker(request *model.DeleteDomainHttpsCertRequest) *DeleteDomainHttpsCertInvoker {
	requestDef := GenReqDefForDeleteDomainHttpsCert()
	return &DeleteDomainHttpsCertInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDomainHttpsCert 查询指定域名的https证书配置
//
// 查询指定域名的https证书配置
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) ShowDomainHttpsCert(request *model.ShowDomainHttpsCertRequest) (*model.ShowDomainHttpsCertResponse, error) {
	requestDef := GenReqDefForShowDomainHttpsCert()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainHttpsCertResponse), nil
	}
}

// ShowDomainHttpsCertInvoker 查询指定域名的https证书配置
func (c *LiveClient) ShowDomainHttpsCertInvoker(request *model.ShowDomainHttpsCertRequest) *ShowDomainHttpsCertInvoker {
	requestDef := GenReqDefForShowDomainHttpsCert()
	return &ShowDomainHttpsCertInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDomainHttpsCert 修改指定域名的https证书配置
//
// 修改指定域名的https证书配置
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) UpdateDomainHttpsCert(request *model.UpdateDomainHttpsCertRequest) (*model.UpdateDomainHttpsCertResponse, error) {
	requestDef := GenReqDefForUpdateDomainHttpsCert()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainHttpsCertResponse), nil
	}
}

// UpdateDomainHttpsCertInvoker 修改指定域名的https证书配置
func (c *LiveClient) UpdateDomainHttpsCertInvoker(request *model.UpdateDomainHttpsCertRequest) *UpdateDomainHttpsCertInvoker {
	requestDef := GenReqDefForUpdateDomainHttpsCert()
	return &UpdateDomainHttpsCertInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateObsBucketAuthorityPublic OBS桶授权及取消授权
//
// OBS桶授权及取消授权
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *LiveClient) UpdateObsBucketAuthorityPublic(request *model.UpdateObsBucketAuthorityPublicRequest) (*model.UpdateObsBucketAuthorityPublicResponse, error) {
	requestDef := GenReqDefForUpdateObsBucketAuthorityPublic()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateObsBucketAuthorityPublicResponse), nil
	}
}

// UpdateObsBucketAuthorityPublicInvoker OBS桶授权及取消授权
func (c *LiveClient) UpdateObsBucketAuthorityPublicInvoker(request *model.UpdateObsBucketAuthorityPublicRequest) *UpdateObsBucketAuthorityPublicInvoker {
	requestDef := GenReqDefForUpdateObsBucketAuthorityPublic()
	return &UpdateObsBucketAuthorityPublicInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
