package v1

import (
	httpclient "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v1/model"
)

type CssClient struct {
	HcClient *httpclient.HcHttpClient
}

func NewCssClient(hcClient *httpclient.HcHttpClient) *CssClient {
	return &CssClient{HcClient: hcClient}
}

func CssClientBuilder() *httpclient.HcHttpClientBuilder {
	builder := httpclient.NewHcHttpClientBuilder()
	return builder
}

// AddIndependentNode 添加独立master、client
//
// 由于集群数据面业务的增长或者不确定性，很难在一开始就能够把集群的规模形态想明白，该接口能够在非独立master和client的集群上面独立master、client角色。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) AddIndependentNode(request *model.AddIndependentNodeRequest) (*model.AddIndependentNodeResponse, error) {
	requestDef := GenReqDefForAddIndependentNode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddIndependentNodeResponse), nil
	}
}

// AddIndependentNodeInvoker 添加独立master、client
func (c *CssClient) AddIndependentNodeInvoker(request *model.AddIndependentNodeRequest) *AddIndependentNodeInvoker {
	requestDef := GenReqDefForAddIndependentNode()
	return &AddIndependentNodeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ChangeMode 安全模式修改
//
// 该接口用于切换集群的安全模式。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ChangeMode(request *model.ChangeModeRequest) (*model.ChangeModeResponse, error) {
	requestDef := GenReqDefForChangeMode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeModeResponse), nil
	}
}

// ChangeModeInvoker 安全模式修改
func (c *CssClient) ChangeModeInvoker(request *model.ChangeModeRequest) *ChangeModeInvoker {
	requestDef := GenReqDefForChangeMode()
	return &ChangeModeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ChangeSecurityGroup 切换安全组
//
// 该接口可以在集群创建成功后，修改安全组。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ChangeSecurityGroup(request *model.ChangeSecurityGroupRequest) (*model.ChangeSecurityGroupResponse, error) {
	requestDef := GenReqDefForChangeSecurityGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeSecurityGroupResponse), nil
	}
}

// ChangeSecurityGroupInvoker 切换安全组
func (c *CssClient) ChangeSecurityGroupInvoker(request *model.ChangeSecurityGroupRequest) *ChangeSecurityGroupInvoker {
	requestDef := GenReqDefForChangeSecurityGroup()
	return &ChangeSecurityGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAiOps 创建一次集群检测任务
//
// 该接口用于创建一个集群检测任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) CreateAiOps(request *model.CreateAiOpsRequest) (*model.CreateAiOpsResponse, error) {
	requestDef := GenReqDefForCreateAiOps()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAiOpsResponse), nil
	}
}

// CreateAiOpsInvoker 创建一次集群检测任务
func (c *CssClient) CreateAiOpsInvoker(request *model.CreateAiOpsRequest) *CreateAiOpsInvoker {
	requestDef := GenReqDefForCreateAiOps()
	return &CreateAiOpsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAutoCreatePolicy 设置自动创建快照策略
//
// 该接口用于设置自动创建快照，默认一天创建一个快照。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) CreateAutoCreatePolicy(request *model.CreateAutoCreatePolicyRequest) (*model.CreateAutoCreatePolicyResponse, error) {
	requestDef := GenReqDefForCreateAutoCreatePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAutoCreatePolicyResponse), nil
	}
}

// CreateAutoCreatePolicyInvoker 设置自动创建快照策略
func (c *CssClient) CreateAutoCreatePolicyInvoker(request *model.CreateAutoCreatePolicyRequest) *CreateAutoCreatePolicyInvoker {
	requestDef := GenReqDefForCreateAutoCreatePolicy()
	return &CreateAutoCreatePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateBindPublic 开启公网访问
//
// 该接口用于开启公网访问。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) CreateBindPublic(request *model.CreateBindPublicRequest) (*model.CreateBindPublicResponse, error) {
	requestDef := GenReqDefForCreateBindPublic()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateBindPublicResponse), nil
	}
}

// CreateBindPublicInvoker 开启公网访问
func (c *CssClient) CreateBindPublicInvoker(request *model.CreateBindPublicRequest) *CreateBindPublicInvoker {
	requestDef := GenReqDefForCreateBindPublic()
	return &CreateBindPublicInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateCluster 创建集群
//
// 该接口用于创建集群。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) CreateCluster(request *model.CreateClusterRequest) (*model.CreateClusterResponse, error) {
	requestDef := GenReqDefForCreateCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateClusterResponse), nil
	}
}

// CreateClusterInvoker 创建集群
func (c *CssClient) CreateClusterInvoker(request *model.CreateClusterRequest) *CreateClusterInvoker {
	requestDef := GenReqDefForCreateCluster()
	return &CreateClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateClustersTags 添加指定集群标签
//
// 该接口用于给指定集群添加标签。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) CreateClustersTags(request *model.CreateClustersTagsRequest) (*model.CreateClustersTagsResponse, error) {
	requestDef := GenReqDefForCreateClustersTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateClustersTagsResponse), nil
	}
}

// CreateClustersTagsInvoker 添加指定集群标签
func (c *CssClient) CreateClustersTagsInvoker(request *model.CreateClustersTagsRequest) *CreateClustersTagsInvoker {
	requestDef := GenReqDefForCreateClustersTags()
	return &CreateClustersTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateElbListener es监听器配置。
//
// 该接口用于es监听器配置。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) CreateElbListener(request *model.CreateElbListenerRequest) (*model.CreateElbListenerResponse, error) {
	requestDef := GenReqDefForCreateElbListener()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateElbListenerResponse), nil
	}
}

// CreateElbListenerInvoker es监听器配置。
func (c *CssClient) CreateElbListenerInvoker(request *model.CreateElbListenerRequest) *CreateElbListenerInvoker {
	requestDef := GenReqDefForCreateElbListener()
	return &CreateElbListenerInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateLoadIkThesaurus 加载自定义词库
//
// 该接口用于加载存放于OBS的自定义词库。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) CreateLoadIkThesaurus(request *model.CreateLoadIkThesaurusRequest) (*model.CreateLoadIkThesaurusResponse, error) {
	requestDef := GenReqDefForCreateLoadIkThesaurus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateLoadIkThesaurusResponse), nil
	}
}

// CreateLoadIkThesaurusInvoker 加载自定义词库
func (c *CssClient) CreateLoadIkThesaurusInvoker(request *model.CreateLoadIkThesaurusRequest) *CreateLoadIkThesaurusInvoker {
	requestDef := GenReqDefForCreateLoadIkThesaurus()
	return &CreateLoadIkThesaurusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateLogBackup 备份日志
//
// 该接口用于备份日志。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) CreateLogBackup(request *model.CreateLogBackupRequest) (*model.CreateLogBackupResponse, error) {
	requestDef := GenReqDefForCreateLogBackup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateLogBackupResponse), nil
	}
}

// CreateLogBackupInvoker 备份日志
func (c *CssClient) CreateLogBackupInvoker(request *model.CreateLogBackupRequest) *CreateLogBackupInvoker {
	requestDef := GenReqDefForCreateLogBackup()
	return &CreateLogBackupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateSnapshot 手动创建快照
//
// 该接口用于手动创建一个快照。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) CreateSnapshot(request *model.CreateSnapshotRequest) (*model.CreateSnapshotResponse, error) {
	requestDef := GenReqDefForCreateSnapshot()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSnapshotResponse), nil
	}
}

// CreateSnapshotInvoker 手动创建快照
func (c *CssClient) CreateSnapshotInvoker(request *model.CreateSnapshotRequest) *CreateSnapshotInvoker {
	requestDef := GenReqDefForCreateSnapshot()
	return &CreateSnapshotInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAiOps 删除一个检测任务记录
//
// 该接口用于删除一个检测任务记录。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) DeleteAiOps(request *model.DeleteAiOpsRequest) (*model.DeleteAiOpsResponse, error) {
	requestDef := GenReqDefForDeleteAiOps()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAiOpsResponse), nil
	}
}

// DeleteAiOpsInvoker 删除一个检测任务记录
func (c *CssClient) DeleteAiOpsInvoker(request *model.DeleteAiOpsRequest) *DeleteAiOpsInvoker {
	requestDef := GenReqDefForDeleteAiOps()
	return &DeleteAiOpsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteCluster 删除集群
//
// 此接口用于删除集群。集群删除将释放此集群的所有资源，包括客户数据。如果需要保留客户集群数据，建议在删除集群前先创建快照。
//
// &gt;此接口亦可用于包年/包月集群退订。公安冻结的集群不能删除。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) DeleteCluster(request *model.DeleteClusterRequest) (*model.DeleteClusterResponse, error) {
	requestDef := GenReqDefForDeleteCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteClusterResponse), nil
	}
}

// DeleteClusterInvoker 删除集群
func (c *CssClient) DeleteClusterInvoker(request *model.DeleteClusterRequest) *DeleteClusterInvoker {
	requestDef := GenReqDefForDeleteCluster()
	return &DeleteClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteClustersTags 删除集群标签
//
// 此接口用于删除集群标签。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) DeleteClustersTags(request *model.DeleteClustersTagsRequest) (*model.DeleteClustersTagsResponse, error) {
	requestDef := GenReqDefForDeleteClustersTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteClustersTagsResponse), nil
	}
}

// DeleteClustersTagsInvoker 删除集群标签
func (c *CssClient) DeleteClustersTagsInvoker(request *model.DeleteClustersTagsRequest) *DeleteClustersTagsInvoker {
	requestDef := GenReqDefForDeleteClustersTags()
	return &DeleteClustersTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteIkThesaurus 删除自定义词库
//
// 该接口用于删除自定义词库。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) DeleteIkThesaurus(request *model.DeleteIkThesaurusRequest) (*model.DeleteIkThesaurusResponse, error) {
	requestDef := GenReqDefForDeleteIkThesaurus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteIkThesaurusResponse), nil
	}
}

// DeleteIkThesaurusInvoker 删除自定义词库
func (c *CssClient) DeleteIkThesaurusInvoker(request *model.DeleteIkThesaurusRequest) *DeleteIkThesaurusInvoker {
	requestDef := GenReqDefForDeleteIkThesaurus()
	return &DeleteIkThesaurusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteSnapshot 删除快照
//
// 该接口用于删除快照。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) DeleteSnapshot(request *model.DeleteSnapshotRequest) (*model.DeleteSnapshotResponse, error) {
	requestDef := GenReqDefForDeleteSnapshot()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSnapshotResponse), nil
	}
}

// DeleteSnapshotInvoker 删除快照
func (c *CssClient) DeleteSnapshotInvoker(request *model.DeleteSnapshotRequest) *DeleteSnapshotInvoker {
	requestDef := GenReqDefForDeleteSnapshot()
	return &DeleteSnapshotInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DownloadCert 下载安全证书
//
// 该接口用于下载安全证书。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) DownloadCert(request *model.DownloadCertRequest) (*model.DownloadCertResponse, error) {
	requestDef := GenReqDefForDownloadCert()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DownloadCertResponse), nil
	}
}

// DownloadCertInvoker 下载安全证书
func (c *CssClient) DownloadCertInvoker(request *model.DownloadCertRequest) *DownloadCertInvoker {
	requestDef := GenReqDefForDownloadCert()
	return &DownloadCertInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// EnableOrDisableElb 打开或关闭es负载均衡器
//
// 该接口打开或关闭es负载均衡器。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) EnableOrDisableElb(request *model.EnableOrDisableElbRequest) (*model.EnableOrDisableElbResponse, error) {
	requestDef := GenReqDefForEnableOrDisableElb()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.EnableOrDisableElbResponse), nil
	}
}

// EnableOrDisableElbInvoker 打开或关闭es负载均衡器
func (c *CssClient) EnableOrDisableElbInvoker(request *model.EnableOrDisableElbRequest) *EnableOrDisableElbInvoker {
	requestDef := GenReqDefForEnableOrDisableElb()
	return &EnableOrDisableElbInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAiOps 获取智能运维任务列表及详情
//
// 该接口用于获取智能运维任务列表及详情。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListAiOps(request *model.ListAiOpsRequest) (*model.ListAiOpsResponse, error) {
	requestDef := GenReqDefForListAiOps()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAiOpsResponse), nil
	}
}

// ListAiOpsInvoker 获取智能运维任务列表及详情
func (c *CssClient) ListAiOpsInvoker(request *model.ListAiOpsRequest) *ListAiOpsInvoker {
	requestDef := GenReqDefForListAiOps()
	return &ListAiOpsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListClustersDetails 查询集群列表
//
// 该接口用于查询并显示集群列表以及集群的状态。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListClustersDetails(request *model.ListClustersDetailsRequest) (*model.ListClustersDetailsResponse, error) {
	requestDef := GenReqDefForListClustersDetails()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListClustersDetailsResponse), nil
	}
}

// ListClustersDetailsInvoker 查询集群列表
func (c *CssClient) ListClustersDetailsInvoker(request *model.ListClustersDetailsRequest) *ListClustersDetailsInvoker {
	requestDef := GenReqDefForListClustersDetails()
	return &ListClustersDetailsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListClustersTags 查询所有标签
//
// 该接口用于查询指定region下的所有标签集合。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListClustersTags(request *model.ListClustersTagsRequest) (*model.ListClustersTagsResponse, error) {
	requestDef := GenReqDefForListClustersTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListClustersTagsResponse), nil
	}
}

// ListClustersTagsInvoker 查询所有标签
func (c *CssClient) ListClustersTagsInvoker(request *model.ListClustersTagsRequest) *ListClustersTagsInvoker {
	requestDef := GenReqDefForListClustersTags()
	return &ListClustersTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListElbCerts 查询证书列表
//
// 该接口用于查询证书列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListElbCerts(request *model.ListElbCertsRequest) (*model.ListElbCertsResponse, error) {
	requestDef := GenReqDefForListElbCerts()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListElbCertsResponse), nil
	}
}

// ListElbCertsInvoker 查询证书列表
func (c *CssClient) ListElbCertsInvoker(request *model.ListElbCertsRequest) *ListElbCertsInvoker {
	requestDef := GenReqDefForListElbCerts()
	return &ListElbCertsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListElbs 查询集群支持的elbv3负载均衡器
//
// 展示查询集群支持的elbv3负载均衡器
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListElbs(request *model.ListElbsRequest) (*model.ListElbsResponse, error) {
	requestDef := GenReqDefForListElbs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListElbsResponse), nil
	}
}

// ListElbsInvoker 查询集群支持的elbv3负载均衡器
func (c *CssClient) ListElbsInvoker(request *model.ListElbsRequest) *ListElbsInvoker {
	requestDef := GenReqDefForListElbs()
	return &ListElbsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListFlavors 获取实例规格列表
//
// 该接口用于查询并显示支持的实例规格对应的ID。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListFlavors(request *model.ListFlavorsRequest) (*model.ListFlavorsResponse, error) {
	requestDef := GenReqDefForListFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListFlavorsResponse), nil
	}
}

// ListFlavorsInvoker 获取实例规格列表
func (c *CssClient) ListFlavorsInvoker(request *model.ListFlavorsRequest) *ListFlavorsInvoker {
	requestDef := GenReqDefForListFlavors()
	return &ListFlavorsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListImages 获取目标镜像ID
//
// 该接口用于获取当前集群的可升级目标镜像ID。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListImages(request *model.ListImagesRequest) (*model.ListImagesResponse, error) {
	requestDef := GenReqDefForListImages()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListImagesResponse), nil
	}
}

// ListImagesInvoker 获取目标镜像ID
func (c *CssClient) ListImagesInvoker(request *model.ListImagesRequest) *ListImagesInvoker {
	requestDef := GenReqDefForListImages()
	return &ListImagesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListLogsJob 查询作业列表
//
// 该接口用于查询具体某个集群的日志任务记录列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListLogsJob(request *model.ListLogsJobRequest) (*model.ListLogsJobResponse, error) {
	requestDef := GenReqDefForListLogsJob()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListLogsJobResponse), nil
	}
}

// ListLogsJobInvoker 查询作业列表
func (c *CssClient) ListLogsJobInvoker(request *model.ListLogsJobRequest) *ListLogsJobInvoker {
	requestDef := GenReqDefForListLogsJob()
	return &ListLogsJobInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSmnTopics 获取智能运维告警可用的SMN主题
//
// 该接口用于获取智能运维告警可用的SMN主题。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListSmnTopics(request *model.ListSmnTopicsRequest) (*model.ListSmnTopicsResponse, error) {
	requestDef := GenReqDefForListSmnTopics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSmnTopicsResponse), nil
	}
}

// ListSmnTopicsInvoker 获取智能运维告警可用的SMN主题
func (c *CssClient) ListSmnTopicsInvoker(request *model.ListSmnTopicsRequest) *ListSmnTopicsInvoker {
	requestDef := GenReqDefForListSmnTopics()
	return &ListSmnTopicsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSnapshots 查询快照列表
//
// 该接口用于查询集群的所有快照。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListSnapshots(request *model.ListSnapshotsRequest) (*model.ListSnapshotsResponse, error) {
	requestDef := GenReqDefForListSnapshots()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSnapshotsResponse), nil
	}
}

// ListSnapshotsInvoker 查询快照列表
func (c *CssClient) ListSnapshotsInvoker(request *model.ListSnapshotsRequest) *ListSnapshotsInvoker {
	requestDef := GenReqDefForListSnapshots()
	return &ListSnapshotsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListYmls 获取参数配置列表
//
// 该接口用于获取当前集群现有的参数配置列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListYmls(request *model.ListYmlsRequest) (*model.ListYmlsResponse, error) {
	requestDef := GenReqDefForListYmls()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListYmlsResponse), nil
	}
}

// ListYmlsInvoker 获取参数配置列表
func (c *CssClient) ListYmlsInvoker(request *model.ListYmlsRequest) *ListYmlsInvoker {
	requestDef := GenReqDefForListYmls()
	return &ListYmlsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListYmlsJob 获取参数配置任务列表
//
// 该接口可获取参数配置的任务操作列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListYmlsJob(request *model.ListYmlsJobRequest) (*model.ListYmlsJobResponse, error) {
	requestDef := GenReqDefForListYmlsJob()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListYmlsJobResponse), nil
	}
}

// ListYmlsJobInvoker 获取参数配置任务列表
func (c *CssClient) ListYmlsJobInvoker(request *model.ListYmlsJobRequest) *ListYmlsJobInvoker {
	requestDef := GenReqDefForListYmlsJob()
	return &ListYmlsJobInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ResetPassword 修改密码
//
// 该接口用于修改集群密码。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ResetPassword(request *model.ResetPasswordRequest) (*model.ResetPasswordResponse, error) {
	requestDef := GenReqDefForResetPassword()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ResetPasswordResponse), nil
	}
}

// ResetPasswordInvoker 修改密码
func (c *CssClient) ResetPasswordInvoker(request *model.ResetPasswordRequest) *ResetPasswordInvoker {
	requestDef := GenReqDefForResetPassword()
	return &ResetPasswordInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RestartCluster 重启集群
//
// 此接口用于重启集群，重启集群将导致业务中断。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) RestartCluster(request *model.RestartClusterRequest) (*model.RestartClusterResponse, error) {
	requestDef := GenReqDefForRestartCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RestartClusterResponse), nil
	}
}

// RestartClusterInvoker 重启集群
func (c *CssClient) RestartClusterInvoker(request *model.RestartClusterRequest) *RestartClusterInvoker {
	requestDef := GenReqDefForRestartCluster()
	return &RestartClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RestoreSnapshot 恢复快照
//
// 该接口用于手动恢复一个快照。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) RestoreSnapshot(request *model.RestoreSnapshotRequest) (*model.RestoreSnapshotResponse, error) {
	requestDef := GenReqDefForRestoreSnapshot()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RestoreSnapshotResponse), nil
	}
}

// RestoreSnapshotInvoker 恢复快照
func (c *CssClient) RestoreSnapshotInvoker(request *model.RestoreSnapshotRequest) *RestoreSnapshotInvoker {
	requestDef := GenReqDefForRestoreSnapshot()
	return &RestoreSnapshotInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RetryUpgradeTask 重试升级失败任务
//
// 由于升级过程时间较长，可能由于网络等原因导致升级失败，可以通过该接口重试该任务或终止该任务的影响。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) RetryUpgradeTask(request *model.RetryUpgradeTaskRequest) (*model.RetryUpgradeTaskResponse, error) {
	requestDef := GenReqDefForRetryUpgradeTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RetryUpgradeTaskResponse), nil
	}
}

// RetryUpgradeTaskInvoker 重试升级失败任务
func (c *CssClient) RetryUpgradeTaskInvoker(request *model.RetryUpgradeTaskRequest) *RetryUpgradeTaskInvoker {
	requestDef := GenReqDefForRetryUpgradeTask()
	return &RetryUpgradeTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutoCreatePolicy 查询自动创建快照的策略
//
// 该接口用于查询自动创建快照策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ShowAutoCreatePolicy(request *model.ShowAutoCreatePolicyRequest) (*model.ShowAutoCreatePolicyResponse, error) {
	requestDef := GenReqDefForShowAutoCreatePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutoCreatePolicyResponse), nil
	}
}

// ShowAutoCreatePolicyInvoker 查询自动创建快照的策略
func (c *CssClient) ShowAutoCreatePolicyInvoker(request *model.ShowAutoCreatePolicyRequest) *ShowAutoCreatePolicyInvoker {
	requestDef := GenReqDefForShowAutoCreatePolicy()
	return &ShowAutoCreatePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowClusterDetail 查询集群详情
//
// 该接口用于查询并显示单个集群详情。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ShowClusterDetail(request *model.ShowClusterDetailRequest) (*model.ShowClusterDetailResponse, error) {
	requestDef := GenReqDefForShowClusterDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowClusterDetailResponse), nil
	}
}

// ShowClusterDetailInvoker 查询集群详情
func (c *CssClient) ShowClusterDetailInvoker(request *model.ShowClusterDetailRequest) *ShowClusterDetailInvoker {
	requestDef := GenReqDefForShowClusterDetail()
	return &ShowClusterDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowClusterTag 查询指定集群的标签
//
// 该接口用于查询指定集群的标签信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ShowClusterTag(request *model.ShowClusterTagRequest) (*model.ShowClusterTagResponse, error) {
	requestDef := GenReqDefForShowClusterTag()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowClusterTagResponse), nil
	}
}

// ShowClusterTagInvoker 查询指定集群的标签
func (c *CssClient) ShowClusterTagInvoker(request *model.ShowClusterTagRequest) *ShowClusterTagInvoker {
	requestDef := GenReqDefForShowClusterTag()
	return &ShowClusterTagInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowElbDetail 获取该esELB的信息，以及页面需要展示健康检查状态
//
// 该接口用于获取该esELB的信息，以及页面需要展示健康检查状态。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ShowElbDetail(request *model.ShowElbDetailRequest) (*model.ShowElbDetailResponse, error) {
	requestDef := GenReqDefForShowElbDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowElbDetailResponse), nil
	}
}

// ShowElbDetailInvoker 获取该esELB的信息，以及页面需要展示健康检查状态
func (c *CssClient) ShowElbDetailInvoker(request *model.ShowElbDetailRequest) *ShowElbDetailInvoker {
	requestDef := GenReqDefForShowElbDetail()
	return &ShowElbDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowGetLogSetting 查询日志基础配置
//
// 该接口用于日志基础配置查询。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ShowGetLogSetting(request *model.ShowGetLogSettingRequest) (*model.ShowGetLogSettingResponse, error) {
	requestDef := GenReqDefForShowGetLogSetting()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowGetLogSettingResponse), nil
	}
}

// ShowGetLogSettingInvoker 查询日志基础配置
func (c *CssClient) ShowGetLogSettingInvoker(request *model.ShowGetLogSettingRequest) *ShowGetLogSettingInvoker {
	requestDef := GenReqDefForShowGetLogSetting()
	return &ShowGetLogSettingInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowIkThesaurus 查询自定义词库状态
//
// 该接口用于查询自定义词库的加载状态。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ShowIkThesaurus(request *model.ShowIkThesaurusRequest) (*model.ShowIkThesaurusResponse, error) {
	requestDef := GenReqDefForShowIkThesaurus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowIkThesaurusResponse), nil
	}
}

// ShowIkThesaurusInvoker 查询自定义词库状态
func (c *CssClient) ShowIkThesaurusInvoker(request *model.ShowIkThesaurusRequest) *ShowIkThesaurusInvoker {
	requestDef := GenReqDefForShowIkThesaurus()
	return &ShowIkThesaurusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowLogBackup 查询日志
//
// 该接口用于查询日志信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ShowLogBackup(request *model.ShowLogBackupRequest) (*model.ShowLogBackupResponse, error) {
	requestDef := GenReqDefForShowLogBackup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowLogBackupResponse), nil
	}
}

// ShowLogBackupInvoker 查询日志
func (c *CssClient) ShowLogBackupInvoker(request *model.ShowLogBackupRequest) *ShowLogBackupInvoker {
	requestDef := GenReqDefForShowLogBackup()
	return &ShowLogBackupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowVpcepConnection 获取终端节点连接
//
// 该接口用于获取终端节点连接。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ShowVpcepConnection(request *model.ShowVpcepConnectionRequest) (*model.ShowVpcepConnectionResponse, error) {
	requestDef := GenReqDefForShowVpcepConnection()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowVpcepConnectionResponse), nil
	}
}

// ShowVpcepConnectionInvoker 获取终端节点连接
func (c *CssClient) ShowVpcepConnectionInvoker(request *model.ShowVpcepConnectionRequest) *ShowVpcepConnectionInvoker {
	requestDef := GenReqDefForShowVpcepConnection()
	return &ShowVpcepConnectionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartAutoSetting 自动设置集群快照的基础配置（不推荐使用）
//
// 该接口用于自动设置集群快照的基础配置，包括配置OBS桶和IAM委托。
//
// - “OBS桶”：快照存储的OBS桶位置。
//
// - “备份路径”：快照在OBS桶中的存放路径。
//
// - “IAM委托”：由于需要将快照保存在OBS中，所以需要在IAM中设置对应的委托获取对OBS服务的授权。
//
// &gt;自动设置集群快照接口将会自动创建快照OBS桶和委托。如果有多个集群，每个集群使用这个接口都会创建一个不一样的OBS桶，可能会导致OBS的配额不够，较多的OBS桶也难以维护。建议可以直接使用[修改集群快照的基础配置](UpdateSnapshotSetting.xml)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StartAutoSetting(request *model.StartAutoSettingRequest) (*model.StartAutoSettingResponse, error) {
	requestDef := GenReqDefForStartAutoSetting()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartAutoSettingResponse), nil
	}
}

// StartAutoSettingInvoker 自动设置集群快照的基础配置（不推荐使用）
func (c *CssClient) StartAutoSettingInvoker(request *model.StartAutoSettingRequest) *StartAutoSettingInvoker {
	requestDef := GenReqDefForStartAutoSetting()
	return &StartAutoSettingInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartLogAutoBackupPolicy 开启日志自动备份策略
//
// 该接口用于日志自动备份策略开启。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StartLogAutoBackupPolicy(request *model.StartLogAutoBackupPolicyRequest) (*model.StartLogAutoBackupPolicyResponse, error) {
	requestDef := GenReqDefForStartLogAutoBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartLogAutoBackupPolicyResponse), nil
	}
}

// StartLogAutoBackupPolicyInvoker 开启日志自动备份策略
func (c *CssClient) StartLogAutoBackupPolicyInvoker(request *model.StartLogAutoBackupPolicyRequest) *StartLogAutoBackupPolicyInvoker {
	requestDef := GenReqDefForStartLogAutoBackupPolicy()
	return &StartLogAutoBackupPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartLogs 开启日志功能
//
// 该接口用于开启日志功能。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StartLogs(request *model.StartLogsRequest) (*model.StartLogsResponse, error) {
	requestDef := GenReqDefForStartLogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartLogsResponse), nil
	}
}

// StartLogsInvoker 开启日志功能
func (c *CssClient) StartLogsInvoker(request *model.StartLogsRequest) *StartLogsInvoker {
	requestDef := GenReqDefForStartLogs()
	return &StartLogsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartPublicWhitelist 开启公网访问控制白名单
//
// 该接口用于开启公网访问控制白名单。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StartPublicWhitelist(request *model.StartPublicWhitelistRequest) (*model.StartPublicWhitelistResponse, error) {
	requestDef := GenReqDefForStartPublicWhitelist()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartPublicWhitelistResponse), nil
	}
}

// StartPublicWhitelistInvoker 开启公网访问控制白名单
func (c *CssClient) StartPublicWhitelistInvoker(request *model.StartPublicWhitelistRequest) *StartPublicWhitelistInvoker {
	requestDef := GenReqDefForStartPublicWhitelist()
	return &StartPublicWhitelistInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartTargetClusterConnectivityTest 连通性测试。
//
// 该接口用于连通性测试。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StartTargetClusterConnectivityTest(request *model.StartTargetClusterConnectivityTestRequest) (*model.StartTargetClusterConnectivityTestResponse, error) {
	requestDef := GenReqDefForStartTargetClusterConnectivityTest()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartTargetClusterConnectivityTestResponse), nil
	}
}

// StartTargetClusterConnectivityTestInvoker 连通性测试。
func (c *CssClient) StartTargetClusterConnectivityTestInvoker(request *model.StartTargetClusterConnectivityTestRequest) *StartTargetClusterConnectivityTestInvoker {
	requestDef := GenReqDefForStartTargetClusterConnectivityTest()
	return &StartTargetClusterConnectivityTestInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartVpecp 开启终端节点服务
//
// 该接口用于开启终端节点服务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StartVpecp(request *model.StartVpecpRequest) (*model.StartVpecpResponse, error) {
	requestDef := GenReqDefForStartVpecp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartVpecpResponse), nil
	}
}

// StartVpecpInvoker 开启终端节点服务
func (c *CssClient) StartVpecpInvoker(request *model.StartVpecpRequest) *StartVpecpInvoker {
	requestDef := GenReqDefForStartVpecp()
	return &StartVpecpInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopLogAutoBackupPolicy 关闭日志自动备份策略
//
// 该接口用于日志自动备份策略关闭。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StopLogAutoBackupPolicy(request *model.StopLogAutoBackupPolicyRequest) (*model.StopLogAutoBackupPolicyResponse, error) {
	requestDef := GenReqDefForStopLogAutoBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopLogAutoBackupPolicyResponse), nil
	}
}

// StopLogAutoBackupPolicyInvoker 关闭日志自动备份策略
func (c *CssClient) StopLogAutoBackupPolicyInvoker(request *model.StopLogAutoBackupPolicyRequest) *StopLogAutoBackupPolicyInvoker {
	requestDef := GenReqDefForStopLogAutoBackupPolicy()
	return &StopLogAutoBackupPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopLogs 关闭日志功能
//
// 该接口用于关闭日志功能。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StopLogs(request *model.StopLogsRequest) (*model.StopLogsResponse, error) {
	requestDef := GenReqDefForStopLogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopLogsResponse), nil
	}
}

// StopLogsInvoker 关闭日志功能
func (c *CssClient) StopLogsInvoker(request *model.StopLogsRequest) *StopLogsInvoker {
	requestDef := GenReqDefForStopLogs()
	return &StopLogsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopPublicWhitelist 关闭公网访问控制白名单
//
// 该接口用于关闭公网访问控制白名单。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StopPublicWhitelist(request *model.StopPublicWhitelistRequest) (*model.StopPublicWhitelistResponse, error) {
	requestDef := GenReqDefForStopPublicWhitelist()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopPublicWhitelistResponse), nil
	}
}

// StopPublicWhitelistInvoker 关闭公网访问控制白名单
func (c *CssClient) StopPublicWhitelistInvoker(request *model.StopPublicWhitelistRequest) *StopPublicWhitelistInvoker {
	requestDef := GenReqDefForStopPublicWhitelist()
	return &StopPublicWhitelistInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopSnapshot 停用快照功能
//
// 该接口用于停用快照功能。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StopSnapshot(request *model.StopSnapshotRequest) (*model.StopSnapshotResponse, error) {
	requestDef := GenReqDefForStopSnapshot()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopSnapshotResponse), nil
	}
}

// StopSnapshotInvoker 停用快照功能
func (c *CssClient) StopSnapshotInvoker(request *model.StopSnapshotRequest) *StopSnapshotInvoker {
	requestDef := GenReqDefForStopSnapshot()
	return &StopSnapshotInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopVpecp 关闭终端节点服务
//
// 该接口用于关闭终端节点服务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StopVpecp(request *model.StopVpecpRequest) (*model.StopVpecpResponse, error) {
	requestDef := GenReqDefForStopVpecp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopVpecpResponse), nil
	}
}

// StopVpecpInvoker 关闭终端节点服务
func (c *CssClient) StopVpecpInvoker(request *model.StopVpecpRequest) *StopVpecpInvoker {
	requestDef := GenReqDefForStopVpecp()
	return &StopVpecpInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAzByInstanceType 切换集群实例AZ
//
// 该接口通过指定节点类型切换AZ。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateAzByInstanceType(request *model.UpdateAzByInstanceTypeRequest) (*model.UpdateAzByInstanceTypeResponse, error) {
	requestDef := GenReqDefForUpdateAzByInstanceType()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAzByInstanceTypeResponse), nil
	}
}

// UpdateAzByInstanceTypeInvoker 切换集群实例AZ
func (c *CssClient) UpdateAzByInstanceTypeInvoker(request *model.UpdateAzByInstanceTypeRequest) *UpdateAzByInstanceTypeInvoker {
	requestDef := GenReqDefForUpdateAzByInstanceType()
	return &UpdateAzByInstanceTypeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateBatchClustersTags 批量添加或删除集群标签
//
// 该接口用于对集群批量添加或删除标签。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateBatchClustersTags(request *model.UpdateBatchClustersTagsRequest) (*model.UpdateBatchClustersTagsResponse, error) {
	requestDef := GenReqDefForUpdateBatchClustersTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateBatchClustersTagsResponse), nil
	}
}

// UpdateBatchClustersTagsInvoker 批量添加或删除集群标签
func (c *CssClient) UpdateBatchClustersTagsInvoker(request *model.UpdateBatchClustersTagsRequest) *UpdateBatchClustersTagsInvoker {
	requestDef := GenReqDefForUpdateBatchClustersTags()
	return &UpdateBatchClustersTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateClusterName 修改集群名称
//
// 该接口用于修改集群名称。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateClusterName(request *model.UpdateClusterNameRequest) (*model.UpdateClusterNameResponse, error) {
	requestDef := GenReqDefForUpdateClusterName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateClusterNameResponse), nil
	}
}

// UpdateClusterNameInvoker 修改集群名称
func (c *CssClient) UpdateClusterNameInvoker(request *model.UpdateClusterNameRequest) *UpdateClusterNameInvoker {
	requestDef := GenReqDefForUpdateClusterName()
	return &UpdateClusterNameInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateEsListener 更新es监听器
//
// 该接口用于更新es监听器。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateEsListener(request *model.UpdateEsListenerRequest) (*model.UpdateEsListenerResponse, error) {
	requestDef := GenReqDefForUpdateEsListener()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateEsListenerResponse), nil
	}
}

// UpdateEsListenerInvoker 更新es监听器
func (c *CssClient) UpdateEsListenerInvoker(request *model.UpdateEsListenerRequest) *UpdateEsListenerInvoker {
	requestDef := GenReqDefForUpdateEsListener()
	return &UpdateEsListenerInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateExtendCluster 扩容集群
//
// 该接口用于集群扩容实例（仅支持扩容elasticsearch实例）。只扩容普通节点，且只针对要扩容的集群实例不存在特殊节点（Master、Client、冷数据节点）的情况。
//
// 集群扩容实例的数量和存储容量，请参考[扩容实例的数量和存储容量](UpdateExtendInstanceStorage.xml)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateExtendCluster(request *model.UpdateExtendClusterRequest) (*model.UpdateExtendClusterResponse, error) {
	requestDef := GenReqDefForUpdateExtendCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateExtendClusterResponse), nil
	}
}

// UpdateExtendClusterInvoker 扩容集群
func (c *CssClient) UpdateExtendClusterInvoker(request *model.UpdateExtendClusterRequest) *UpdateExtendClusterInvoker {
	requestDef := GenReqDefForUpdateExtendCluster()
	return &UpdateExtendClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateExtendInstanceStorage 扩容实例的数量和存储容量
//
// 该接口用于集群扩容不同类型实例的个数以及存储容量。已经存在独立Master、Client、冷数据节点的集群使用该接口扩容。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateExtendInstanceStorage(request *model.UpdateExtendInstanceStorageRequest) (*model.UpdateExtendInstanceStorageResponse, error) {
	requestDef := GenReqDefForUpdateExtendInstanceStorage()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateExtendInstanceStorageResponse), nil
	}
}

// UpdateExtendInstanceStorageInvoker 扩容实例的数量和存储容量
func (c *CssClient) UpdateExtendInstanceStorageInvoker(request *model.UpdateExtendInstanceStorageRequest) *UpdateExtendInstanceStorageInvoker {
	requestDef := GenReqDefForUpdateExtendInstanceStorage()
	return &UpdateExtendInstanceStorageInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateFlavor 变更规格
//
// 该接口用于变更集群规格。只支持变更ess节点类型。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateFlavor(request *model.UpdateFlavorRequest) (*model.UpdateFlavorResponse, error) {
	requestDef := GenReqDefForUpdateFlavor()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateFlavorResponse), nil
	}
}

// UpdateFlavorInvoker 变更规格
func (c *CssClient) UpdateFlavorInvoker(request *model.UpdateFlavorRequest) *UpdateFlavorInvoker {
	requestDef := GenReqDefForUpdateFlavor()
	return &UpdateFlavorInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateFlavorByType 指定节点类型规格变更
//
// 修改集群规格。支持修改:
// - ess： 数据节点。
// - ess-cold: 冷数据节点。
// - ess-client: Client节点。
// - ess-master: Master节点。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateFlavorByType(request *model.UpdateFlavorByTypeRequest) (*model.UpdateFlavorByTypeResponse, error) {
	requestDef := GenReqDefForUpdateFlavorByType()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateFlavorByTypeResponse), nil
	}
}

// UpdateFlavorByTypeInvoker 指定节点类型规格变更
func (c *CssClient) UpdateFlavorByTypeInvoker(request *model.UpdateFlavorByTypeRequest) *UpdateFlavorByTypeInvoker {
	requestDef := GenReqDefForUpdateFlavorByType()
	return &UpdateFlavorByTypeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateInstance 节点替换
//
// 该接口用于替换失败节点。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateInstance(request *model.UpdateInstanceRequest) (*model.UpdateInstanceResponse, error) {
	requestDef := GenReqDefForUpdateInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateInstanceResponse), nil
	}
}

// UpdateInstanceInvoker 节点替换
func (c *CssClient) UpdateInstanceInvoker(request *model.UpdateInstanceRequest) *UpdateInstanceInvoker {
	requestDef := GenReqDefForUpdateInstance()
	return &UpdateInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateLogSetting 修改日志基础配置
//
// 该接口用于修改日志基础配置。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateLogSetting(request *model.UpdateLogSettingRequest) (*model.UpdateLogSettingResponse, error) {
	requestDef := GenReqDefForUpdateLogSetting()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateLogSettingResponse), nil
	}
}

// UpdateLogSettingInvoker 修改日志基础配置
func (c *CssClient) UpdateLogSettingInvoker(request *model.UpdateLogSettingRequest) *UpdateLogSettingInvoker {
	requestDef := GenReqDefForUpdateLogSetting()
	return &UpdateLogSettingInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateOndemandClusterToPeriod 按需集群转包周期
//
// 该接口用于按需集群转包周期集群。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateOndemandClusterToPeriod(request *model.UpdateOndemandClusterToPeriodRequest) (*model.UpdateOndemandClusterToPeriodResponse, error) {
	requestDef := GenReqDefForUpdateOndemandClusterToPeriod()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateOndemandClusterToPeriodResponse), nil
	}
}

// UpdateOndemandClusterToPeriodInvoker 按需集群转包周期
func (c *CssClient) UpdateOndemandClusterToPeriodInvoker(request *model.UpdateOndemandClusterToPeriodRequest) *UpdateOndemandClusterToPeriodInvoker {
	requestDef := GenReqDefForUpdateOndemandClusterToPeriod()
	return &UpdateOndemandClusterToPeriodInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdatePublicBandWidth 修改公网访问带宽
//
// 该接口用于修改公网访问带宽。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdatePublicBandWidth(request *model.UpdatePublicBandWidthRequest) (*model.UpdatePublicBandWidthResponse, error) {
	requestDef := GenReqDefForUpdatePublicBandWidth()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePublicBandWidthResponse), nil
	}
}

// UpdatePublicBandWidthInvoker 修改公网访问带宽
func (c *CssClient) UpdatePublicBandWidthInvoker(request *model.UpdatePublicBandWidthRequest) *UpdatePublicBandWidthInvoker {
	requestDef := GenReqDefForUpdatePublicBandWidth()
	return &UpdatePublicBandWidthInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateShrinkCluster 指定节点类型缩容
//
// 该接口用于集群对不同类型实例的个数以及存储容量进行缩容。包周期类型的集群不支持通过api进行指定节点类型缩容操作。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateShrinkCluster(request *model.UpdateShrinkClusterRequest) (*model.UpdateShrinkClusterResponse, error) {
	requestDef := GenReqDefForUpdateShrinkCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateShrinkClusterResponse), nil
	}
}

// UpdateShrinkClusterInvoker 指定节点类型缩容
func (c *CssClient) UpdateShrinkClusterInvoker(request *model.UpdateShrinkClusterRequest) *UpdateShrinkClusterInvoker {
	requestDef := GenReqDefForUpdateShrinkCluster()
	return &UpdateShrinkClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateShrinkNodes 指定节点缩容
//
// 该接口可以对集群现有节点中指定节点进行缩容。包周期类型的集群不支持通过api进行指定节点缩容操作。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateShrinkNodes(request *model.UpdateShrinkNodesRequest) (*model.UpdateShrinkNodesResponse, error) {
	requestDef := GenReqDefForUpdateShrinkNodes()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateShrinkNodesResponse), nil
	}
}

// UpdateShrinkNodesInvoker 指定节点缩容
func (c *CssClient) UpdateShrinkNodesInvoker(request *model.UpdateShrinkNodesRequest) *UpdateShrinkNodesInvoker {
	requestDef := GenReqDefForUpdateShrinkNodes()
	return &UpdateShrinkNodesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateSnapshotSetting 修改集群快照的基础配置
//
// 该接口用于修改集群快照的基础配置，可修改OBS桶和IAM委托。
//
// 可以使用该接口开启快照功能。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateSnapshotSetting(request *model.UpdateSnapshotSettingRequest) (*model.UpdateSnapshotSettingResponse, error) {
	requestDef := GenReqDefForUpdateSnapshotSetting()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateSnapshotSettingResponse), nil
	}
}

// UpdateSnapshotSettingInvoker 修改集群快照的基础配置
func (c *CssClient) UpdateSnapshotSettingInvoker(request *model.UpdateSnapshotSettingRequest) *UpdateSnapshotSettingInvoker {
	requestDef := GenReqDefForUpdateSnapshotSetting()
	return &UpdateSnapshotSettingInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateUnbindPublic 关闭公网访问
//
// 该接口用于关闭公网访问。包周期类型的集群不支持通过api进行关闭公网访问。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateUnbindPublic(request *model.UpdateUnbindPublicRequest) (*model.UpdateUnbindPublicResponse, error) {
	requestDef := GenReqDefForUpdateUnbindPublic()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateUnbindPublicResponse), nil
	}
}

// UpdateUnbindPublicInvoker 关闭公网访问
func (c *CssClient) UpdateUnbindPublicInvoker(request *model.UpdateUnbindPublicRequest) *UpdateUnbindPublicInvoker {
	requestDef := GenReqDefForUpdateUnbindPublic()
	return &UpdateUnbindPublicInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateVpcepConnection 更新终端节点连接
//
// 该接口用于更新终端节点连接。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateVpcepConnection(request *model.UpdateVpcepConnectionRequest) (*model.UpdateVpcepConnectionResponse, error) {
	requestDef := GenReqDefForUpdateVpcepConnection()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateVpcepConnectionResponse), nil
	}
}

// UpdateVpcepConnectionInvoker 更新终端节点连接
func (c *CssClient) UpdateVpcepConnectionInvoker(request *model.UpdateVpcepConnectionRequest) *UpdateVpcepConnectionInvoker {
	requestDef := GenReqDefForUpdateVpcepConnection()
	return &UpdateVpcepConnectionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateVpcepWhitelist 修改终端节点服务白名单
//
// 该接口用于修改终端节点服务白名单。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateVpcepWhitelist(request *model.UpdateVpcepWhitelistRequest) (*model.UpdateVpcepWhitelistResponse, error) {
	requestDef := GenReqDefForUpdateVpcepWhitelist()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateVpcepWhitelistResponse), nil
	}
}

// UpdateVpcepWhitelistInvoker 修改终端节点服务白名单
func (c *CssClient) UpdateVpcepWhitelistInvoker(request *model.UpdateVpcepWhitelistRequest) *UpdateVpcepWhitelistInvoker {
	requestDef := GenReqDefForUpdateVpcepWhitelist()
	return &UpdateVpcepWhitelistInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateYmls 修改参数配置
//
// 该接口用于修改参数配置。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateYmls(request *model.UpdateYmlsRequest) (*model.UpdateYmlsResponse, error) {
	requestDef := GenReqDefForUpdateYmls()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateYmlsResponse), nil
	}
}

// UpdateYmlsInvoker 修改参数配置
func (c *CssClient) UpdateYmlsInvoker(request *model.UpdateYmlsRequest) *UpdateYmlsInvoker {
	requestDef := GenReqDefForUpdateYmls()
	return &UpdateYmlsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpgradeCore 集群内核升级
//
// 该接口用于将低版本的ES升级到高版本或同版本ES。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpgradeCore(request *model.UpgradeCoreRequest) (*model.UpgradeCoreResponse, error) {
	requestDef := GenReqDefForUpgradeCore()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeCoreResponse), nil
	}
}

// UpgradeCoreInvoker 集群内核升级
func (c *CssClient) UpgradeCoreInvoker(request *model.UpgradeCoreRequest) *UpgradeCoreInvoker {
	requestDef := GenReqDefForUpgradeCore()
	return &UpgradeCoreInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpgradeDetail 获取升级详情信息
//
// 由于升级过程时间较长，该接口可以展示当前升级（切换AZ）节点的各个阶段信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpgradeDetail(request *model.UpgradeDetailRequest) (*model.UpgradeDetailResponse, error) {
	requestDef := GenReqDefForUpgradeDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeDetailResponse), nil
	}
}

// UpgradeDetailInvoker 获取升级详情信息
func (c *CssClient) UpgradeDetailInvoker(request *model.UpgradeDetailRequest) *UpgradeDetailInvoker {
	requestDef := GenReqDefForUpgradeDetail()
	return &UpgradeDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartKibanaPublic 开启Kibana公网访问
//
// 该接口用于开启Kibana公网访问。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StartKibanaPublic(request *model.StartKibanaPublicRequest) (*model.StartKibanaPublicResponse, error) {
	requestDef := GenReqDefForStartKibanaPublic()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartKibanaPublicResponse), nil
	}
}

// StartKibanaPublicInvoker 开启Kibana公网访问
func (c *CssClient) StartKibanaPublicInvoker(request *model.StartKibanaPublicRequest) *StartKibanaPublicInvoker {
	requestDef := GenReqDefForStartKibanaPublic()
	return &StartKibanaPublicInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopPublicKibanaWhitelist 关闭Kibana公网访问控制
//
// 该接口用于关闭Kibana公网访问控制。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StopPublicKibanaWhitelist(request *model.StopPublicKibanaWhitelistRequest) (*model.StopPublicKibanaWhitelistResponse, error) {
	requestDef := GenReqDefForStopPublicKibanaWhitelist()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopPublicKibanaWhitelistResponse), nil
	}
}

// StopPublicKibanaWhitelistInvoker 关闭Kibana公网访问控制
func (c *CssClient) StopPublicKibanaWhitelistInvoker(request *model.StopPublicKibanaWhitelistRequest) *StopPublicKibanaWhitelistInvoker {
	requestDef := GenReqDefForStopPublicKibanaWhitelist()
	return &StopPublicKibanaWhitelistInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAlterKibana 修改Kibana公网带宽
//
// 该接口用于修改Kibana公网带宽。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateAlterKibana(request *model.UpdateAlterKibanaRequest) (*model.UpdateAlterKibanaResponse, error) {
	requestDef := GenReqDefForUpdateAlterKibana()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAlterKibanaResponse), nil
	}
}

// UpdateAlterKibanaInvoker 修改Kibana公网带宽
func (c *CssClient) UpdateAlterKibanaInvoker(request *model.UpdateAlterKibanaRequest) *UpdateAlterKibanaInvoker {
	requestDef := GenReqDefForUpdateAlterKibana()
	return &UpdateAlterKibanaInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateCloseKibana 关闭Kibana公网访问
//
// 该接口用于关闭Kibana公网访问。包周期类型集群不支持通过api进行关闭Kibana公网访问。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateCloseKibana(request *model.UpdateCloseKibanaRequest) (*model.UpdateCloseKibanaResponse, error) {
	requestDef := GenReqDefForUpdateCloseKibana()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateCloseKibanaResponse), nil
	}
}

// UpdateCloseKibanaInvoker 关闭Kibana公网访问
func (c *CssClient) UpdateCloseKibanaInvoker(request *model.UpdateCloseKibanaRequest) *UpdateCloseKibanaInvoker {
	requestDef := GenReqDefForUpdateCloseKibana()
	return &UpdateCloseKibanaInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdatePublicKibanaWhitelist 修改Kibana公网访问控制
//
// 该接口通过修改kibana白名单，修改kibana的访问权限。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdatePublicKibanaWhitelist(request *model.UpdatePublicKibanaWhitelistRequest) (*model.UpdatePublicKibanaWhitelistResponse, error) {
	requestDef := GenReqDefForUpdatePublicKibanaWhitelist()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePublicKibanaWhitelistResponse), nil
	}
}

// UpdatePublicKibanaWhitelistInvoker 修改Kibana公网访问控制
func (c *CssClient) UpdatePublicKibanaWhitelistInvoker(request *model.UpdatePublicKibanaWhitelistRequest) *UpdatePublicKibanaWhitelistInvoker {
	requestDef := GenReqDefForUpdatePublicKibanaWhitelist()
	return &UpdatePublicKibanaWhitelistInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddFavorite 添加到自定义模板
//
// 该接口用于添加到自定义模板。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) AddFavorite(request *model.AddFavoriteRequest) (*model.AddFavoriteResponse, error) {
	requestDef := GenReqDefForAddFavorite()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddFavoriteResponse), nil
	}
}

// AddFavoriteInvoker 添加到自定义模板
func (c *CssClient) AddFavoriteInvoker(request *model.AddFavoriteRequest) *AddFavoriteInvoker {
	requestDef := GenReqDefForAddFavorite()
	return &AddFavoriteInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateCnf 创建配置文件
//
// 该接口用于创建配置文件。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) CreateCnf(request *model.CreateCnfRequest) (*model.CreateCnfResponse, error) {
	requestDef := GenReqDefForCreateCnf()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateCnfResponse), nil
	}
}

// CreateCnfInvoker 创建配置文件
func (c *CssClient) CreateCnfInvoker(request *model.CreateCnfRequest) *CreateCnfInvoker {
	requestDef := GenReqDefForCreateCnf()
	return &CreateCnfInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteConf 删除配置文件
//
// 删除配置文件。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) DeleteConf(request *model.DeleteConfRequest) (*model.DeleteConfResponse, error) {
	requestDef := GenReqDefForDeleteConf()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteConfResponse), nil
	}
}

// DeleteConfInvoker 删除配置文件
func (c *CssClient) DeleteConfInvoker(request *model.DeleteConfRequest) *DeleteConfInvoker {
	requestDef := GenReqDefForDeleteConf()
	return &DeleteConfInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteConfig 删除配置文件V2
//
// 删除配置文件。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) DeleteConfig(request *model.DeleteConfigRequest) (*model.DeleteConfigResponse, error) {
	requestDef := GenReqDefForDeleteConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteConfigResponse), nil
	}
}

// DeleteConfigInvoker 删除配置文件V2
func (c *CssClient) DeleteConfigInvoker(request *model.DeleteConfigRequest) *DeleteConfigInvoker {
	requestDef := GenReqDefForDeleteConfig()
	return &DeleteConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTemplate 删除自定义模板
//
// 该接口用于删除自定义模板。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) DeleteTemplate(request *model.DeleteTemplateRequest) (*model.DeleteTemplateResponse, error) {
	requestDef := GenReqDefForDeleteTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTemplateResponse), nil
	}
}

// DeleteTemplateInvoker 删除自定义模板
func (c *CssClient) DeleteTemplateInvoker(request *model.DeleteTemplateRequest) *DeleteTemplateInvoker {
	requestDef := GenReqDefForDeleteTemplate()
	return &DeleteTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListActions 查询操作记录
//
// 该接口用于查询操作记录。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListActions(request *model.ListActionsRequest) (*model.ListActionsResponse, error) {
	requestDef := GenReqDefForListActions()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListActionsResponse), nil
	}
}

// ListActionsInvoker 查询操作记录
func (c *CssClient) ListActionsInvoker(request *model.ListActionsRequest) *ListActionsInvoker {
	requestDef := GenReqDefForListActions()
	return &ListActionsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListCerts 查询证书列表
//
// 该接口用于查询证书列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListCerts(request *model.ListCertsRequest) (*model.ListCertsResponse, error) {
	requestDef := GenReqDefForListCerts()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListCertsResponse), nil
	}
}

// ListCertsInvoker 查询证书列表
func (c *CssClient) ListCertsInvoker(request *model.ListCertsRequest) *ListCertsInvoker {
	requestDef := GenReqDefForListCerts()
	return &ListCertsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListConfs 查询配置文件列表
//
// 该接口用于查询配置文件列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListConfs(request *model.ListConfsRequest) (*model.ListConfsResponse, error) {
	requestDef := GenReqDefForListConfs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListConfsResponse), nil
	}
}

// ListConfsInvoker 查询配置文件列表
func (c *CssClient) ListConfsInvoker(request *model.ListConfsRequest) *ListConfsInvoker {
	requestDef := GenReqDefForListConfs()
	return &ListConfsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPipelines 查询pipeline列表
//
// 该接口用于查询pipeline列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListPipelines(request *model.ListPipelinesRequest) (*model.ListPipelinesResponse, error) {
	requestDef := GenReqDefForListPipelines()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPipelinesResponse), nil
	}
}

// ListPipelinesInvoker 查询pipeline列表
func (c *CssClient) ListPipelinesInvoker(request *model.ListPipelinesRequest) *ListPipelinesInvoker {
	requestDef := GenReqDefForListPipelines()
	return &ListPipelinesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTemplates 查询模板列表
//
// 该接口用于查询模板列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ListTemplates(request *model.ListTemplatesRequest) (*model.ListTemplatesResponse, error) {
	requestDef := GenReqDefForListTemplates()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTemplatesResponse), nil
	}
}

// ListTemplatesInvoker 查询模板列表
func (c *CssClient) ListTemplatesInvoker(request *model.ListTemplatesRequest) *ListTemplatesInvoker {
	requestDef := GenReqDefForListTemplates()
	return &ListTemplatesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowGetConfDetail 查询配置文件内容
//
// 该接口用于查询配置文件内容。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) ShowGetConfDetail(request *model.ShowGetConfDetailRequest) (*model.ShowGetConfDetailResponse, error) {
	requestDef := GenReqDefForShowGetConfDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowGetConfDetailResponse), nil
	}
}

// ShowGetConfDetailInvoker 查询配置文件内容
func (c *CssClient) ShowGetConfDetailInvoker(request *model.ShowGetConfDetailRequest) *ShowGetConfDetailInvoker {
	requestDef := GenReqDefForShowGetConfDetail()
	return &ShowGetConfDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartConnectivityTest 连通性测试
//
// 该接口用于连通性测试。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StartConnectivityTest(request *model.StartConnectivityTestRequest) (*model.StartConnectivityTestResponse, error) {
	requestDef := GenReqDefForStartConnectivityTest()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartConnectivityTestResponse), nil
	}
}

// StartConnectivityTestInvoker 连通性测试
func (c *CssClient) StartConnectivityTestInvoker(request *model.StartConnectivityTestRequest) *StartConnectivityTestInvoker {
	requestDef := GenReqDefForStartConnectivityTest()
	return &StartConnectivityTestInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartPipeline 启动pipeline迁移数据
//
// 该接口用于启动pipeline迁移数据。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StartPipeline(request *model.StartPipelineRequest) (*model.StartPipelineResponse, error) {
	requestDef := GenReqDefForStartPipeline()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartPipelineResponse), nil
	}
}

// StartPipelineInvoker 启动pipeline迁移数据
func (c *CssClient) StartPipelineInvoker(request *model.StartPipelineRequest) *StartPipelineInvoker {
	requestDef := GenReqDefForStartPipeline()
	return &StartPipelineInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopHotPipeline 热停止pipeline迁移数据。
//
// 该接口用于热停止pipeline迁移数据。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StopHotPipeline(request *model.StopHotPipelineRequest) (*model.StopHotPipelineResponse, error) {
	requestDef := GenReqDefForStopHotPipeline()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopHotPipelineResponse), nil
	}
}

// StopHotPipelineInvoker 热停止pipeline迁移数据。
func (c *CssClient) StopHotPipelineInvoker(request *model.StopHotPipelineRequest) *StopHotPipelineInvoker {
	requestDef := GenReqDefForStopHotPipeline()
	return &StopHotPipelineInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopPipeline 停止pipeline迁移数据
//
// 该接口用于停止pipeline迁移数据。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) StopPipeline(request *model.StopPipelineRequest) (*model.StopPipelineResponse, error) {
	requestDef := GenReqDefForStopPipeline()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopPipelineResponse), nil
	}
}

// StopPipelineInvoker 停止pipeline迁移数据
func (c *CssClient) StopPipelineInvoker(request *model.StopPipelineRequest) *StopPipelineInvoker {
	requestDef := GenReqDefForStopPipeline()
	return &StopPipelineInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateCnf 更新配置文件
//
// 该接口用于更新配置文件。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CssClient) UpdateCnf(request *model.UpdateCnfRequest) (*model.UpdateCnfResponse, error) {
	requestDef := GenReqDefForUpdateCnf()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateCnfResponse), nil
	}
}

// UpdateCnfInvoker 更新配置文件
func (c *CssClient) UpdateCnfInvoker(request *model.UpdateCnfRequest) *UpdateCnfInvoker {
	requestDef := GenReqDefForUpdateCnf()
	return &UpdateCnfInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
