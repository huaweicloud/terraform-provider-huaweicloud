package v3

import (
	httpclient "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cce/v3/model"
)

type CceClient struct {
	HcClient *httpclient.HcHttpClient
}

func NewCceClient(hcClient *httpclient.HcHttpClient) *CceClient {
	return &CceClient{HcClient: hcClient}
}

func CceClientBuilder() *httpclient.HcHttpClientBuilder {
	builder := httpclient.NewHcHttpClientBuilder()
	return builder
}

// AddNode 纳管节点
//
// 该API用于在指定集群下纳管节点。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) AddNode(request *model.AddNodeRequest) (*model.AddNodeResponse, error) {
	requestDef := GenReqDefForAddNode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddNodeResponse), nil
	}
}

// AddNodeInvoker 纳管节点
func (c *CceClient) AddNodeInvoker(request *model.AddNodeRequest) *AddNodeInvoker {
	requestDef := GenReqDefForAddNode()
	return &AddNodeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddNodesToNodePool 自定义节点池纳管节点
//
// 该API用于在指定集群自定义节点池下纳管节点。竞价实例不支持纳管。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) AddNodesToNodePool(request *model.AddNodesToNodePoolRequest) (*model.AddNodesToNodePoolResponse, error) {
	requestDef := GenReqDefForAddNodesToNodePool()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddNodesToNodePoolResponse), nil
	}
}

// AddNodesToNodePoolInvoker 自定义节点池纳管节点
func (c *CceClient) AddNodesToNodePoolInvoker(request *model.AddNodesToNodePoolRequest) *AddNodesToNodePoolInvoker {
	requestDef := GenReqDefForAddNodesToNodePool()
	return &AddNodesToNodePoolInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AwakeCluster 集群唤醒
//
// 集群唤醒用于唤醒已休眠的集群，唤醒后，将继续收取控制节点资源费用。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) AwakeCluster(request *model.AwakeClusterRequest) (*model.AwakeClusterResponse, error) {
	requestDef := GenReqDefForAwakeCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AwakeClusterResponse), nil
	}
}

// AwakeClusterInvoker 集群唤醒
func (c *CceClient) AwakeClusterInvoker(request *model.AwakeClusterRequest) *AwakeClusterInvoker {
	requestDef := GenReqDefForAwakeCluster()
	return &AwakeClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchCreateClusterTags 批量添加指定集群的资源标签
//
// 该API用于批量添加指定集群的资源标签。
// &gt; - 每个集群支持最多20个资源标签。
// &gt; - 此接口为幂等接口：创建时，如果创建的标签已经存在（key/value均相同视为重复），默认处理成功；key相同，value不同时会覆盖原有标签。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) BatchCreateClusterTags(request *model.BatchCreateClusterTagsRequest) (*model.BatchCreateClusterTagsResponse, error) {
	requestDef := GenReqDefForBatchCreateClusterTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchCreateClusterTagsResponse), nil
	}
}

// BatchCreateClusterTagsInvoker 批量添加指定集群的资源标签
func (c *CceClient) BatchCreateClusterTagsInvoker(request *model.BatchCreateClusterTagsRequest) *BatchCreateClusterTagsInvoker {
	requestDef := GenReqDefForBatchCreateClusterTags()
	return &BatchCreateClusterTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchDeleteClusterTags 批量删除指定集群的资源标签
//
// 该API用于批量删除指定集群的资源标签。
// &gt; - 此接口为幂等接口：删除时，如果删除的标签key不存在，默认处理成功。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) BatchDeleteClusterTags(request *model.BatchDeleteClusterTagsRequest) (*model.BatchDeleteClusterTagsResponse, error) {
	requestDef := GenReqDefForBatchDeleteClusterTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchDeleteClusterTagsResponse), nil
	}
}

// BatchDeleteClusterTagsInvoker 批量删除指定集群的资源标签
func (c *CceClient) BatchDeleteClusterTagsInvoker(request *model.BatchDeleteClusterTagsRequest) *BatchDeleteClusterTagsInvoker {
	requestDef := GenReqDefForBatchDeleteClusterTags()
	return &BatchDeleteClusterTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchSyncNodes 批量同步节点
//
// 该API用于批量同步节点。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) BatchSyncNodes(request *model.BatchSyncNodesRequest) (*model.BatchSyncNodesResponse, error) {
	requestDef := GenReqDefForBatchSyncNodes()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchSyncNodesResponse), nil
	}
}

// BatchSyncNodesInvoker 批量同步节点
func (c *CceClient) BatchSyncNodesInvoker(request *model.BatchSyncNodesRequest) *BatchSyncNodesInvoker {
	requestDef := GenReqDefForBatchSyncNodes()
	return &BatchSyncNodesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ContinueUpgradeClusterTask 继续执行集群升级任务
//
// 继续执行被暂停的集群升级任务。
// &gt; - 集群升级涉及多维度的组件升级操作，强烈建议统一通过CCE控制台执行交互式升级，降低集群升级过程的业务意外受损风险；
// &gt; - 当前集群升级相关接口受限开放。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ContinueUpgradeClusterTask(request *model.ContinueUpgradeClusterTaskRequest) (*model.ContinueUpgradeClusterTaskResponse, error) {
	requestDef := GenReqDefForContinueUpgradeClusterTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ContinueUpgradeClusterTaskResponse), nil
	}
}

// ContinueUpgradeClusterTaskInvoker 继续执行集群升级任务
func (c *CceClient) ContinueUpgradeClusterTaskInvoker(request *model.ContinueUpgradeClusterTaskRequest) *ContinueUpgradeClusterTaskInvoker {
	requestDef := GenReqDefForContinueUpgradeClusterTask()
	return &ContinueUpgradeClusterTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAddonInstance 创建AddonInstance
//
// 根据提供的插件模板，安装插件实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateAddonInstance(request *model.CreateAddonInstanceRequest) (*model.CreateAddonInstanceResponse, error) {
	requestDef := GenReqDefForCreateAddonInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAddonInstanceResponse), nil
	}
}

// CreateAddonInstanceInvoker 创建AddonInstance
func (c *CceClient) CreateAddonInstanceInvoker(request *model.CreateAddonInstanceRequest) *CreateAddonInstanceInvoker {
	requestDef := GenReqDefForCreateAddonInstance()
	return &CreateAddonInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateCloudPersistentVolumeClaims 创建PVC（待废弃）
//
// 该API用于在指定的Namespace下通过云存储服务中的云存储（EVS、SFS、OBS）去创建PVC（PersistentVolumeClaim）。该API待废弃，请使用Kubernetes PVC相关接口。
//
// &gt;存储管理的URL格式为：https://{clusterid}.Endpoint/uri。其中{clusterid}为集群ID，uri为资源路径，也即API访问的路径。如果使用https://Endpoint/uri，则必须指定请求header中的X-Cluster-ID参数。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateCloudPersistentVolumeClaims(request *model.CreateCloudPersistentVolumeClaimsRequest) (*model.CreateCloudPersistentVolumeClaimsResponse, error) {
	requestDef := GenReqDefForCreateCloudPersistentVolumeClaims()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateCloudPersistentVolumeClaimsResponse), nil
	}
}

// CreateCloudPersistentVolumeClaimsInvoker 创建PVC（待废弃）
func (c *CceClient) CreateCloudPersistentVolumeClaimsInvoker(request *model.CreateCloudPersistentVolumeClaimsRequest) *CreateCloudPersistentVolumeClaimsInvoker {
	requestDef := GenReqDefForCreateCloudPersistentVolumeClaims()
	return &CreateCloudPersistentVolumeClaimsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateCluster 创建集群
//
// 该API用于创建一个空集群（即只有控制节点Master，没有工作节点Node）。请在调用本接口完成集群创建之后，通过[创建节点](cce_02_0242.xml)添加节点。
//
// &gt;   - 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
// &gt;   - 调用该接口创建集群时，默认不安装ICAgent，若需安装ICAgent，可在请求Body参数的annotations中加入\&quot;cluster.install.addons.external/install\&quot;:\&quot;[{\&quot;addonTemplateName\&quot;:\&quot;icagent\&quot;}]\&quot;的集群注解，将在创建集群时自动安装ICAgent。ICAgent是应用性能管理APM的采集代理，运行在应用所在的服务器上，用于实时采集探针所获取的数据，安装ICAgent是使用应用性能管理APM的前提。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateCluster(request *model.CreateClusterRequest) (*model.CreateClusterResponse, error) {
	requestDef := GenReqDefForCreateCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateClusterResponse), nil
	}
}

// CreateClusterInvoker 创建集群
func (c *CceClient) CreateClusterInvoker(request *model.CreateClusterRequest) *CreateClusterInvoker {
	requestDef := GenReqDefForCreateCluster()
	return &CreateClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateClusterMasterSnapshot 集群备份
//
// 集群备份
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateClusterMasterSnapshot(request *model.CreateClusterMasterSnapshotRequest) (*model.CreateClusterMasterSnapshotResponse, error) {
	requestDef := GenReqDefForCreateClusterMasterSnapshot()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateClusterMasterSnapshotResponse), nil
	}
}

// CreateClusterMasterSnapshotInvoker 集群备份
func (c *CceClient) CreateClusterMasterSnapshotInvoker(request *model.CreateClusterMasterSnapshotRequest) *CreateClusterMasterSnapshotInvoker {
	requestDef := GenReqDefForCreateClusterMasterSnapshot()
	return &CreateClusterMasterSnapshotInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateKubernetesClusterCert 获取集群证书
//
// 该API用于获取指定集群的证书信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateKubernetesClusterCert(request *model.CreateKubernetesClusterCertRequest) (*model.CreateKubernetesClusterCertResponse, error) {
	requestDef := GenReqDefForCreateKubernetesClusterCert()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateKubernetesClusterCertResponse), nil
	}
}

// CreateKubernetesClusterCertInvoker 获取集群证书
func (c *CceClient) CreateKubernetesClusterCertInvoker(request *model.CreateKubernetesClusterCertRequest) *CreateKubernetesClusterCertInvoker {
	requestDef := GenReqDefForCreateKubernetesClusterCert()
	return &CreateKubernetesClusterCertInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateNode 创建节点
//
// 该API用于在指定集群下创建节点。
// &gt; - 若无集群，请先[创建集群](cce_02_0236.xml)。
// &gt; - 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateNode(request *model.CreateNodeRequest) (*model.CreateNodeResponse, error) {
	requestDef := GenReqDefForCreateNode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateNodeResponse), nil
	}
}

// CreateNodeInvoker 创建节点
func (c *CceClient) CreateNodeInvoker(request *model.CreateNodeRequest) *CreateNodeInvoker {
	requestDef := GenReqDefForCreateNode()
	return &CreateNodeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateNodePool 创建节点池
//
// 该API用于在指定集群下创建节点池。仅支持集群在处于可用、扩容、缩容状态时调用。
//
// 1.21版本的集群创建节点池时支持绑定安全组，每个节点池最多绑定五个安全组。
//
// 更新节点池的安全组后，只针对新创的pod生效，建议驱逐节点上原有的pod。
//
// &gt; 若无集群，请先[创建集群](cce_02_0236.xml)。
// &gt; 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateNodePool(request *model.CreateNodePoolRequest) (*model.CreateNodePoolResponse, error) {
	requestDef := GenReqDefForCreateNodePool()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateNodePoolResponse), nil
	}
}

// CreateNodePoolInvoker 创建节点池
func (c *CceClient) CreateNodePoolInvoker(request *model.CreateNodePoolRequest) *CreateNodePoolInvoker {
	requestDef := GenReqDefForCreateNodePool()
	return &CreateNodePoolInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreatePartition 创建分区
//
// 创建分区
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreatePartition(request *model.CreatePartitionRequest) (*model.CreatePartitionResponse, error) {
	requestDef := GenReqDefForCreatePartition()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePartitionResponse), nil
	}
}

// CreatePartitionInvoker 创建分区
func (c *CceClient) CreatePartitionInvoker(request *model.CreatePartitionRequest) *CreatePartitionInvoker {
	requestDef := GenReqDefForCreatePartition()
	return &CreatePartitionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreatePostCheck 集群升级后确认
//
// 集群升级后确认，该接口建议配合Console使用，主要用于升级步骤完成后，客户确认集群状态和业务正常后做反馈。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreatePostCheck(request *model.CreatePostCheckRequest) (*model.CreatePostCheckResponse, error) {
	requestDef := GenReqDefForCreatePostCheck()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePostCheckResponse), nil
	}
}

// CreatePostCheckInvoker 集群升级后确认
func (c *CceClient) CreatePostCheckInvoker(request *model.CreatePostCheckRequest) *CreatePostCheckInvoker {
	requestDef := GenReqDefForCreatePostCheck()
	return &CreatePostCheckInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreatePreCheck 集群升级前检查
//
// 集群升级前检查
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreatePreCheck(request *model.CreatePreCheckRequest) (*model.CreatePreCheckResponse, error) {
	requestDef := GenReqDefForCreatePreCheck()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePreCheckResponse), nil
	}
}

// CreatePreCheckInvoker 集群升级前检查
func (c *CceClient) CreatePreCheckInvoker(request *model.CreatePreCheckRequest) *CreatePreCheckInvoker {
	requestDef := GenReqDefForCreatePreCheck()
	return &CreatePreCheckInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRelease 创建模板实例
//
// 创建模板实例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateRelease(request *model.CreateReleaseRequest) (*model.CreateReleaseResponse, error) {
	requestDef := GenReqDefForCreateRelease()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateReleaseResponse), nil
	}
}

// CreateReleaseInvoker 创建模板实例
func (c *CceClient) CreateReleaseInvoker(request *model.CreateReleaseRequest) *CreateReleaseInvoker {
	requestDef := GenReqDefForCreateRelease()
	return &CreateReleaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateUpgradeWorkFlow 开启集群升级流程引导任务
//
// 该API用于创建一个集群升级流程引导任务。请在调用本接口完成引导任务创建之后，通过集群升级前检查开始检查任务。
// 升级流程任务用于控制集群升级任务的执行流程，执行流程为 升级前检查 &#x3D;&gt; 集群升级 &#x3D;&gt; 升级后检查。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateUpgradeWorkFlow(request *model.CreateUpgradeWorkFlowRequest) (*model.CreateUpgradeWorkFlowResponse, error) {
	requestDef := GenReqDefForCreateUpgradeWorkFlow()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateUpgradeWorkFlowResponse), nil
	}
}

// CreateUpgradeWorkFlowInvoker 开启集群升级流程引导任务
func (c *CceClient) CreateUpgradeWorkFlowInvoker(request *model.CreateUpgradeWorkFlowRequest) *CreateUpgradeWorkFlowInvoker {
	requestDef := GenReqDefForCreateUpgradeWorkFlow()
	return &CreateUpgradeWorkFlowInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAddonInstance 删除AddonInstance
//
// 删除插件实例的功能。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DeleteAddonInstance(request *model.DeleteAddonInstanceRequest) (*model.DeleteAddonInstanceResponse, error) {
	requestDef := GenReqDefForDeleteAddonInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAddonInstanceResponse), nil
	}
}

// DeleteAddonInstanceInvoker 删除AddonInstance
func (c *CceClient) DeleteAddonInstanceInvoker(request *model.DeleteAddonInstanceRequest) *DeleteAddonInstanceInvoker {
	requestDef := GenReqDefForDeleteAddonInstance()
	return &DeleteAddonInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteChart 删除模板
//
// 删除模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DeleteChart(request *model.DeleteChartRequest) (*model.DeleteChartResponse, error) {
	requestDef := GenReqDefForDeleteChart()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteChartResponse), nil
	}
}

// DeleteChartInvoker 删除模板
func (c *CceClient) DeleteChartInvoker(request *model.DeleteChartRequest) *DeleteChartInvoker {
	requestDef := GenReqDefForDeleteChart()
	return &DeleteChartInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteCloudPersistentVolumeClaims 删除PVC（待废弃）
//
// 该API用于删除指定Namespace下的PVC（PersistentVolumeClaim）对象，并可以选择保留后端的云存储。该API待废弃，请使用Kubernetes PVC相关接口。
// &gt;存储管理的URL格式为：https://{clusterid}.Endpoint/uri。其中{clusterid}为集群ID，uri为资源路径，也即API访问的路径。如果使用https://Endpoint/uri，则必须指定请求header中的X-Cluster-ID参数。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DeleteCloudPersistentVolumeClaims(request *model.DeleteCloudPersistentVolumeClaimsRequest) (*model.DeleteCloudPersistentVolumeClaimsResponse, error) {
	requestDef := GenReqDefForDeleteCloudPersistentVolumeClaims()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteCloudPersistentVolumeClaimsResponse), nil
	}
}

// DeleteCloudPersistentVolumeClaimsInvoker 删除PVC（待废弃）
func (c *CceClient) DeleteCloudPersistentVolumeClaimsInvoker(request *model.DeleteCloudPersistentVolumeClaimsRequest) *DeleteCloudPersistentVolumeClaimsInvoker {
	requestDef := GenReqDefForDeleteCloudPersistentVolumeClaims()
	return &DeleteCloudPersistentVolumeClaimsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteCluster 删除集群
//
// 该API用于删除一个指定的集群。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DeleteCluster(request *model.DeleteClusterRequest) (*model.DeleteClusterResponse, error) {
	requestDef := GenReqDefForDeleteCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteClusterResponse), nil
	}
}

// DeleteClusterInvoker 删除集群
func (c *CceClient) DeleteClusterInvoker(request *model.DeleteClusterRequest) *DeleteClusterInvoker {
	requestDef := GenReqDefForDeleteCluster()
	return &DeleteClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteNode 删除节点
//
// 该API用于删除指定的节点。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DeleteNode(request *model.DeleteNodeRequest) (*model.DeleteNodeResponse, error) {
	requestDef := GenReqDefForDeleteNode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteNodeResponse), nil
	}
}

// DeleteNodeInvoker 删除节点
func (c *CceClient) DeleteNodeInvoker(request *model.DeleteNodeRequest) *DeleteNodeInvoker {
	requestDef := GenReqDefForDeleteNode()
	return &DeleteNodeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteNodePool 删除节点池
//
// 该API用于删除指定的节点池。
// &gt; 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DeleteNodePool(request *model.DeleteNodePoolRequest) (*model.DeleteNodePoolResponse, error) {
	requestDef := GenReqDefForDeleteNodePool()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteNodePoolResponse), nil
	}
}

// DeleteNodePoolInvoker 删除节点池
func (c *CceClient) DeleteNodePoolInvoker(request *model.DeleteNodePoolRequest) *DeleteNodePoolInvoker {
	requestDef := GenReqDefForDeleteNodePool()
	return &DeleteNodePoolInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteRelease 删除指定模板实例
//
// 删除指定模板实例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DeleteRelease(request *model.DeleteReleaseRequest) (*model.DeleteReleaseResponse, error) {
	requestDef := GenReqDefForDeleteRelease()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteReleaseResponse), nil
	}
}

// DeleteReleaseInvoker 删除指定模板实例
func (c *CceClient) DeleteReleaseInvoker(request *model.DeleteReleaseRequest) *DeleteReleaseInvoker {
	requestDef := GenReqDefForDeleteRelease()
	return &DeleteReleaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DownloadChart 下载模板
//
// 下载模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DownloadChart(request *model.DownloadChartRequest) (*model.DownloadChartResponse, error) {
	requestDef := GenReqDefForDownloadChart()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DownloadChartResponse), nil
	}
}

// DownloadChartInvoker 下载模板
func (c *CceClient) DownloadChartInvoker(request *model.DownloadChartRequest) *DownloadChartInvoker {
	requestDef := GenReqDefForDownloadChart()
	return &DownloadChartInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// HibernateCluster 集群休眠
//
// 集群休眠用于将运行中的集群置于休眠状态，休眠后，将不再收取控制节点资源费用。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) HibernateCluster(request *model.HibernateClusterRequest) (*model.HibernateClusterResponse, error) {
	requestDef := GenReqDefForHibernateCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.HibernateClusterResponse), nil
	}
}

// HibernateClusterInvoker 集群休眠
func (c *CceClient) HibernateClusterInvoker(request *model.HibernateClusterRequest) *HibernateClusterInvoker {
	requestDef := GenReqDefForHibernateCluster()
	return &HibernateClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAddonInstances 获取AddonInstance列表
//
// 获取集群所有已安装插件实例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAddonInstances(request *model.ListAddonInstancesRequest) (*model.ListAddonInstancesResponse, error) {
	requestDef := GenReqDefForListAddonInstances()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAddonInstancesResponse), nil
	}
}

// ListAddonInstancesInvoker 获取AddonInstance列表
func (c *CceClient) ListAddonInstancesInvoker(request *model.ListAddonInstancesRequest) *ListAddonInstancesInvoker {
	requestDef := GenReqDefForListAddonInstances()
	return &ListAddonInstancesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAddonTemplates 查询AddonTemplates列表
//
// 插件模板查询接口，查询插件信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAddonTemplates(request *model.ListAddonTemplatesRequest) (*model.ListAddonTemplatesResponse, error) {
	requestDef := GenReqDefForListAddonTemplates()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAddonTemplatesResponse), nil
	}
}

// ListAddonTemplatesInvoker 查询AddonTemplates列表
func (c *CceClient) ListAddonTemplatesInvoker(request *model.ListAddonTemplatesRequest) *ListAddonTemplatesInvoker {
	requestDef := GenReqDefForListAddonTemplates()
	return &ListAddonTemplatesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListCharts 获取模板列表
//
// 获取模板列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListCharts(request *model.ListChartsRequest) (*model.ListChartsResponse, error) {
	requestDef := GenReqDefForListCharts()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListChartsResponse), nil
	}
}

// ListChartsInvoker 获取模板列表
func (c *CceClient) ListChartsInvoker(request *model.ListChartsRequest) *ListChartsInvoker {
	requestDef := GenReqDefForListCharts()
	return &ListChartsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListClusterMasterSnapshotTasks 获取集群备份任务详情列表
//
// 获取集群备份任务详情列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListClusterMasterSnapshotTasks(request *model.ListClusterMasterSnapshotTasksRequest) (*model.ListClusterMasterSnapshotTasksResponse, error) {
	requestDef := GenReqDefForListClusterMasterSnapshotTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListClusterMasterSnapshotTasksResponse), nil
	}
}

// ListClusterMasterSnapshotTasksInvoker 获取集群备份任务详情列表
func (c *CceClient) ListClusterMasterSnapshotTasksInvoker(request *model.ListClusterMasterSnapshotTasksRequest) *ListClusterMasterSnapshotTasksInvoker {
	requestDef := GenReqDefForListClusterMasterSnapshotTasks()
	return &ListClusterMasterSnapshotTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListClusterUpgradeFeatureGates 获取集群升级特性开关配置
//
// 获取集群升级特性开关配置
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListClusterUpgradeFeatureGates(request *model.ListClusterUpgradeFeatureGatesRequest) (*model.ListClusterUpgradeFeatureGatesResponse, error) {
	requestDef := GenReqDefForListClusterUpgradeFeatureGates()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListClusterUpgradeFeatureGatesResponse), nil
	}
}

// ListClusterUpgradeFeatureGatesInvoker 获取集群升级特性开关配置
func (c *CceClient) ListClusterUpgradeFeatureGatesInvoker(request *model.ListClusterUpgradeFeatureGatesRequest) *ListClusterUpgradeFeatureGatesInvoker {
	requestDef := GenReqDefForListClusterUpgradeFeatureGates()
	return &ListClusterUpgradeFeatureGatesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListClusterUpgradePaths 获取集群升级路径
//
// 获取集群升级路径
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListClusterUpgradePaths(request *model.ListClusterUpgradePathsRequest) (*model.ListClusterUpgradePathsResponse, error) {
	requestDef := GenReqDefForListClusterUpgradePaths()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListClusterUpgradePathsResponse), nil
	}
}

// ListClusterUpgradePathsInvoker 获取集群升级路径
func (c *CceClient) ListClusterUpgradePathsInvoker(request *model.ListClusterUpgradePathsRequest) *ListClusterUpgradePathsInvoker {
	requestDef := GenReqDefForListClusterUpgradePaths()
	return &ListClusterUpgradePathsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListClusters 获取指定项目下的集群
//
// 该API用于获取指定项目下所有集群的详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListClusters(request *model.ListClustersRequest) (*model.ListClustersResponse, error) {
	requestDef := GenReqDefForListClusters()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListClustersResponse), nil
	}
}

// ListClustersInvoker 获取指定项目下的集群
func (c *CceClient) ListClustersInvoker(request *model.ListClustersRequest) *ListClustersInvoker {
	requestDef := GenReqDefForListClusters()
	return &ListClustersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListNodePools 获取集群下所有节点池
//
// 该API用于获取集群下所有节点池。
// &gt; - 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径
// &gt; - nodepool是集群中具有相同配置的节点实例的子集。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListNodePools(request *model.ListNodePoolsRequest) (*model.ListNodePoolsResponse, error) {
	requestDef := GenReqDefForListNodePools()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListNodePoolsResponse), nil
	}
}

// ListNodePoolsInvoker 获取集群下所有节点池
func (c *CceClient) ListNodePoolsInvoker(request *model.ListNodePoolsRequest) *ListNodePoolsInvoker {
	requestDef := GenReqDefForListNodePools()
	return &ListNodePoolsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListNodes 获取集群下所有节点
//
// 该API用于通过集群ID获取指定集群下所有节点的详细信息。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListNodes(request *model.ListNodesRequest) (*model.ListNodesResponse, error) {
	requestDef := GenReqDefForListNodes()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListNodesResponse), nil
	}
}

// ListNodesInvoker 获取集群下所有节点
func (c *CceClient) ListNodesInvoker(request *model.ListNodesRequest) *ListNodesInvoker {
	requestDef := GenReqDefForListNodes()
	return &ListNodesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPartitions 获取分区列表
//
// 获取分区列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListPartitions(request *model.ListPartitionsRequest) (*model.ListPartitionsResponse, error) {
	requestDef := GenReqDefForListPartitions()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPartitionsResponse), nil
	}
}

// ListPartitionsInvoker 获取分区列表
func (c *CceClient) ListPartitionsInvoker(request *model.ListPartitionsRequest) *ListPartitionsInvoker {
	requestDef := GenReqDefForListPartitions()
	return &ListPartitionsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPreCheckTasks 获取集群升级前检查任务详情列表
//
// 获取集群升级前检查任务详情列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListPreCheckTasks(request *model.ListPreCheckTasksRequest) (*model.ListPreCheckTasksResponse, error) {
	requestDef := GenReqDefForListPreCheckTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPreCheckTasksResponse), nil
	}
}

// ListPreCheckTasksInvoker 获取集群升级前检查任务详情列表
func (c *CceClient) ListPreCheckTasksInvoker(request *model.ListPreCheckTasksRequest) *ListPreCheckTasksInvoker {
	requestDef := GenReqDefForListPreCheckTasks()
	return &ListPreCheckTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListReleases 获取模板实例列表
//
// 获取模板实例列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListReleases(request *model.ListReleasesRequest) (*model.ListReleasesResponse, error) {
	requestDef := GenReqDefForListReleases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListReleasesResponse), nil
	}
}

// ListReleasesInvoker 获取模板实例列表
func (c *CceClient) ListReleasesInvoker(request *model.ListReleasesRequest) *ListReleasesInvoker {
	requestDef := GenReqDefForListReleases()
	return &ListReleasesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListUpgradeClusterTasks 获取集群升级任务详情列表
//
// 获取集群升级任务详情列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListUpgradeClusterTasks(request *model.ListUpgradeClusterTasksRequest) (*model.ListUpgradeClusterTasksResponse, error) {
	requestDef := GenReqDefForListUpgradeClusterTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListUpgradeClusterTasksResponse), nil
	}
}

// ListUpgradeClusterTasksInvoker 获取集群升级任务详情列表
func (c *CceClient) ListUpgradeClusterTasksInvoker(request *model.ListUpgradeClusterTasksRequest) *ListUpgradeClusterTasksInvoker {
	requestDef := GenReqDefForListUpgradeClusterTasks()
	return &ListUpgradeClusterTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListUpgradeWorkFlows 获取UpgradeWorkFlows列表
//
// 获取历史集群升级引导任务列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListUpgradeWorkFlows(request *model.ListUpgradeWorkFlowsRequest) (*model.ListUpgradeWorkFlowsResponse, error) {
	requestDef := GenReqDefForListUpgradeWorkFlows()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListUpgradeWorkFlowsResponse), nil
	}
}

// ListUpgradeWorkFlowsInvoker 获取UpgradeWorkFlows列表
func (c *CceClient) ListUpgradeWorkFlowsInvoker(request *model.ListUpgradeWorkFlowsRequest) *ListUpgradeWorkFlowsInvoker {
	requestDef := GenReqDefForListUpgradeWorkFlows()
	return &ListUpgradeWorkFlowsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// MigrateNode 节点迁移
//
// 该API用于在指定集群下迁移节点到另一集群。
//
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) MigrateNode(request *model.MigrateNodeRequest) (*model.MigrateNodeResponse, error) {
	requestDef := GenReqDefForMigrateNode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.MigrateNodeResponse), nil
	}
}

// MigrateNodeInvoker 节点迁移
func (c *CceClient) MigrateNodeInvoker(request *model.MigrateNodeRequest) *MigrateNodeInvoker {
	requestDef := GenReqDefForMigrateNode()
	return &MigrateNodeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// PauseUpgradeClusterTask 暂停集群升级任务
//
// 暂停集群升级任务。
// &gt; - 集群升级涉及多维度的组件升级操作，强烈建议统一通过CCE控制台执行交互式升级，降低集群升级过程的业务意外受损风险；
// &gt; - 当前集群升级相关接口受限开放。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) PauseUpgradeClusterTask(request *model.PauseUpgradeClusterTaskRequest) (*model.PauseUpgradeClusterTaskResponse, error) {
	requestDef := GenReqDefForPauseUpgradeClusterTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.PauseUpgradeClusterTaskResponse), nil
	}
}

// PauseUpgradeClusterTaskInvoker 暂停集群升级任务
func (c *CceClient) PauseUpgradeClusterTaskInvoker(request *model.PauseUpgradeClusterTaskRequest) *PauseUpgradeClusterTaskInvoker {
	requestDef := GenReqDefForPauseUpgradeClusterTask()
	return &PauseUpgradeClusterTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RemoveNode 节点移除
//
// 该API用于在指定集群下移除节点。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) RemoveNode(request *model.RemoveNodeRequest) (*model.RemoveNodeResponse, error) {
	requestDef := GenReqDefForRemoveNode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RemoveNodeResponse), nil
	}
}

// RemoveNodeInvoker 节点移除
func (c *CceClient) RemoveNodeInvoker(request *model.RemoveNodeRequest) *RemoveNodeInvoker {
	requestDef := GenReqDefForRemoveNode()
	return &RemoveNodeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ResetNode 重置节点
//
// 该API用于在指定集群下重置节点。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ResetNode(request *model.ResetNodeRequest) (*model.ResetNodeResponse, error) {
	requestDef := GenReqDefForResetNode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ResetNodeResponse), nil
	}
}

// ResetNodeInvoker 重置节点
func (c *CceClient) ResetNodeInvoker(request *model.ResetNodeRequest) *ResetNodeInvoker {
	requestDef := GenReqDefForResetNode()
	return &ResetNodeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ResizeCluster 变更集群规格
//
// 该API用于变更一个指定集群的规格。
//
// &gt;   - 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
// &gt;   [- 使用限制请参考[变更集群规格](https://support.huaweicloud.com/usermanual-cce/cce_10_0403.html)。](tag:hws)
// &gt;   [- 使用限制请参考[变更集群规格](https://support.huaweicloud.com/intl/zh-cn/usermanual-cce/cce_10_0403.html)](tag:hws_hk)
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ResizeCluster(request *model.ResizeClusterRequest) (*model.ResizeClusterResponse, error) {
	requestDef := GenReqDefForResizeCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ResizeClusterResponse), nil
	}
}

// ResizeClusterInvoker 变更集群规格
func (c *CceClient) ResizeClusterInvoker(request *model.ResizeClusterRequest) *ResizeClusterInvoker {
	requestDef := GenReqDefForResizeCluster()
	return &ResizeClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RetryUpgradeClusterTask 重试集群升级任务
//
// 重新执行失败的集群升级任务。
// &gt; - 集群升级涉及多维度的组件升级操作，强烈建议统一通过CCE控制台执行交互式升级，降低集群升级过程的业务意外受损风险；
// &gt; - 当前集群升级相关接口受限开放。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) RetryUpgradeClusterTask(request *model.RetryUpgradeClusterTaskRequest) (*model.RetryUpgradeClusterTaskResponse, error) {
	requestDef := GenReqDefForRetryUpgradeClusterTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RetryUpgradeClusterTaskResponse), nil
	}
}

// RetryUpgradeClusterTaskInvoker 重试集群升级任务
func (c *CceClient) RetryUpgradeClusterTaskInvoker(request *model.RetryUpgradeClusterTaskRequest) *RetryUpgradeClusterTaskInvoker {
	requestDef := GenReqDefForRetryUpgradeClusterTask()
	return &RetryUpgradeClusterTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RevokeKubernetesClusterCert 吊销用户的集群证书
//
// 该API用于吊销指定集群的用户证书
//
// &gt; 吊销操作完成后，此证书申请人之前下载的证书和 kubectl 配置文件无法再用于连接集群。此证书申请人可以重新下载证书或 kubectl 配置文件，并使用新下载的文件连接集群
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) RevokeKubernetesClusterCert(request *model.RevokeKubernetesClusterCertRequest) (*model.RevokeKubernetesClusterCertResponse, error) {
	requestDef := GenReqDefForRevokeKubernetesClusterCert()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RevokeKubernetesClusterCertResponse), nil
	}
}

// RevokeKubernetesClusterCertInvoker 吊销用户的集群证书
func (c *CceClient) RevokeKubernetesClusterCertInvoker(request *model.RevokeKubernetesClusterCertRequest) *RevokeKubernetesClusterCertInvoker {
	requestDef := GenReqDefForRevokeKubernetesClusterCert()
	return &RevokeKubernetesClusterCertInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RollbackAddonInstance 回滚AddonInstance
//
// 将插件实例回滚到升级前的版本。只有在当前插件实例版本支持回滚到升级前的版本（status.isRollbackable为true），且插件实例状态为running（运行中）、available（可用）、abnormal（不可用）、upgradeFailed（升级失败）、rollbackFailed（回滚失败）时支持回滚。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) RollbackAddonInstance(request *model.RollbackAddonInstanceRequest) (*model.RollbackAddonInstanceResponse, error) {
	requestDef := GenReqDefForRollbackAddonInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RollbackAddonInstanceResponse), nil
	}
}

// RollbackAddonInstanceInvoker 回滚AddonInstance
func (c *CceClient) RollbackAddonInstanceInvoker(request *model.RollbackAddonInstanceRequest) *RollbackAddonInstanceInvoker {
	requestDef := GenReqDefForRollbackAddonInstance()
	return &RollbackAddonInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ScaleNodePool 伸缩节点池
//
// 该API用于伸缩指定的节点池
// &gt; 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ScaleNodePool(request *model.ScaleNodePoolRequest) (*model.ScaleNodePoolResponse, error) {
	requestDef := GenReqDefForScaleNodePool()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ScaleNodePoolResponse), nil
	}
}

// ScaleNodePoolInvoker 伸缩节点池
func (c *CceClient) ScaleNodePoolInvoker(request *model.ScaleNodePoolRequest) *ScaleNodePoolInvoker {
	requestDef := GenReqDefForScaleNodePool()
	return &ScaleNodePoolInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAddonInstance 获取AddonInstance详情
//
// 获取插件实例详情。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAddonInstance(request *model.ShowAddonInstanceRequest) (*model.ShowAddonInstanceResponse, error) {
	requestDef := GenReqDefForShowAddonInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAddonInstanceResponse), nil
	}
}

// ShowAddonInstanceInvoker 获取AddonInstance详情
func (c *CceClient) ShowAddonInstanceInvoker(request *model.ShowAddonInstanceRequest) *ShowAddonInstanceInvoker {
	requestDef := GenReqDefForShowAddonInstance()
	return &ShowAddonInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowChart 获取模板
//
// 获取模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowChart(request *model.ShowChartRequest) (*model.ShowChartResponse, error) {
	requestDef := GenReqDefForShowChart()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowChartResponse), nil
	}
}

// ShowChartInvoker 获取模板
func (c *CceClient) ShowChartInvoker(request *model.ShowChartRequest) *ShowChartInvoker {
	requestDef := GenReqDefForShowChart()
	return &ShowChartInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowChartValues 获取模板Values
//
// 获取模板Values
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowChartValues(request *model.ShowChartValuesRequest) (*model.ShowChartValuesResponse, error) {
	requestDef := GenReqDefForShowChartValues()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowChartValuesResponse), nil
	}
}

// ShowChartValuesInvoker 获取模板Values
func (c *CceClient) ShowChartValuesInvoker(request *model.ShowChartValuesRequest) *ShowChartValuesInvoker {
	requestDef := GenReqDefForShowChartValues()
	return &ShowChartValuesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowCluster 获取指定的集群
//
// 该API用于获取指定集群的详细信息。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowCluster(request *model.ShowClusterRequest) (*model.ShowClusterResponse, error) {
	requestDef := GenReqDefForShowCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowClusterResponse), nil
	}
}

// ShowClusterInvoker 获取指定的集群
func (c *CceClient) ShowClusterInvoker(request *model.ShowClusterRequest) *ShowClusterInvoker {
	requestDef := GenReqDefForShowCluster()
	return &ShowClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowClusterConfig 查询集群日志配置信息
//
// 获取集群组件上报的LTS的配置信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowClusterConfig(request *model.ShowClusterConfigRequest) (*model.ShowClusterConfigResponse, error) {
	requestDef := GenReqDefForShowClusterConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowClusterConfigResponse), nil
	}
}

// ShowClusterConfigInvoker 查询集群日志配置信息
func (c *CceClient) ShowClusterConfigInvoker(request *model.ShowClusterConfigRequest) *ShowClusterConfigInvoker {
	requestDef := GenReqDefForShowClusterConfig()
	return &ShowClusterConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowClusterConfigurationDetails 查询指定集群支持配置的参数列表
//
// 该API用于查询CCE服务下指定集群支持配置的参数列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowClusterConfigurationDetails(request *model.ShowClusterConfigurationDetailsRequest) (*model.ShowClusterConfigurationDetailsResponse, error) {
	requestDef := GenReqDefForShowClusterConfigurationDetails()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowClusterConfigurationDetailsResponse), nil
	}
}

// ShowClusterConfigurationDetailsInvoker 查询指定集群支持配置的参数列表
func (c *CceClient) ShowClusterConfigurationDetailsInvoker(request *model.ShowClusterConfigurationDetailsRequest) *ShowClusterConfigurationDetailsInvoker {
	requestDef := GenReqDefForShowClusterConfigurationDetails()
	return &ShowClusterConfigurationDetailsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowClusterEndpoints 获取集群访问的地址
//
// 该API用于通过集群ID获取集群访问的地址，包括PrivateIP(HA集群返回VIP)与PublicIP
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowClusterEndpoints(request *model.ShowClusterEndpointsRequest) (*model.ShowClusterEndpointsResponse, error) {
	requestDef := GenReqDefForShowClusterEndpoints()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowClusterEndpointsResponse), nil
	}
}

// ShowClusterEndpointsInvoker 获取集群访问的地址
func (c *CceClient) ShowClusterEndpointsInvoker(request *model.ShowClusterEndpointsRequest) *ShowClusterEndpointsInvoker {
	requestDef := GenReqDefForShowClusterEndpoints()
	return &ShowClusterEndpointsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowClusterSupportConfiguration 根据集群版本类型等查询集群支持的详细配置项，用于集群创建时指定
//
// 该API用于根据集群版本类型等查询集群支持的详细配置项，用于集群创建时指定。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowClusterSupportConfiguration(request *model.ShowClusterSupportConfigurationRequest) (*model.ShowClusterSupportConfigurationResponse, error) {
	requestDef := GenReqDefForShowClusterSupportConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowClusterSupportConfigurationResponse), nil
	}
}

// ShowClusterSupportConfigurationInvoker 根据集群版本类型等查询集群支持的详细配置项，用于集群创建时指定
func (c *CceClient) ShowClusterSupportConfigurationInvoker(request *model.ShowClusterSupportConfigurationRequest) *ShowClusterSupportConfigurationInvoker {
	requestDef := GenReqDefForShowClusterSupportConfiguration()
	return &ShowClusterSupportConfigurationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowClusterUpgradeInfo 获取集群升级相关信息
//
// 获取集群升级相关信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowClusterUpgradeInfo(request *model.ShowClusterUpgradeInfoRequest) (*model.ShowClusterUpgradeInfoResponse, error) {
	requestDef := GenReqDefForShowClusterUpgradeInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowClusterUpgradeInfoResponse), nil
	}
}

// ShowClusterUpgradeInfoInvoker 获取集群升级相关信息
func (c *CceClient) ShowClusterUpgradeInfoInvoker(request *model.ShowClusterUpgradeInfoRequest) *ShowClusterUpgradeInfoInvoker {
	requestDef := GenReqDefForShowClusterUpgradeInfo()
	return &ShowClusterUpgradeInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowJob 获取任务信息
//
// 该API用于获取任务信息。通过某一任务请求下发后返回的jobID来查询指定任务的进度。
// &gt; - 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径
// &gt; - 该接口通常使用场景为：
// &gt;   - 创建、删除集群时，查询相应任务的进度。
// &gt;   - 创建、删除节点时，查询相应任务的进度。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowJob(request *model.ShowJobRequest) (*model.ShowJobResponse, error) {
	requestDef := GenReqDefForShowJob()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowJobResponse), nil
	}
}

// ShowJobInvoker 获取任务信息
func (c *CceClient) ShowJobInvoker(request *model.ShowJobRequest) *ShowJobInvoker {
	requestDef := GenReqDefForShowJob()
	return &ShowJobInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowNode 获取指定的节点
//
// 该API用于通过节点ID获取指定节点的详细信息。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowNode(request *model.ShowNodeRequest) (*model.ShowNodeResponse, error) {
	requestDef := GenReqDefForShowNode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowNodeResponse), nil
	}
}

// ShowNodeInvoker 获取指定的节点
func (c *CceClient) ShowNodeInvoker(request *model.ShowNodeRequest) *ShowNodeInvoker {
	requestDef := GenReqDefForShowNode()
	return &ShowNodeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowNodePool 获取指定的节点池
//
// 该API用于获取指定节点池的详细信息。
// &gt; 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowNodePool(request *model.ShowNodePoolRequest) (*model.ShowNodePoolResponse, error) {
	requestDef := GenReqDefForShowNodePool()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowNodePoolResponse), nil
	}
}

// ShowNodePoolInvoker 获取指定的节点池
func (c *CceClient) ShowNodePoolInvoker(request *model.ShowNodePoolRequest) *ShowNodePoolInvoker {
	requestDef := GenReqDefForShowNodePool()
	return &ShowNodePoolInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowNodePoolConfigurationDetails 查询指定节点池支持配置的参数列表
//
// 该API用于查询CCE服务下指定节点池支持配置的参数列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowNodePoolConfigurationDetails(request *model.ShowNodePoolConfigurationDetailsRequest) (*model.ShowNodePoolConfigurationDetailsResponse, error) {
	requestDef := GenReqDefForShowNodePoolConfigurationDetails()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowNodePoolConfigurationDetailsResponse), nil
	}
}

// ShowNodePoolConfigurationDetailsInvoker 查询指定节点池支持配置的参数列表
func (c *CceClient) ShowNodePoolConfigurationDetailsInvoker(request *model.ShowNodePoolConfigurationDetailsRequest) *ShowNodePoolConfigurationDetailsInvoker {
	requestDef := GenReqDefForShowNodePoolConfigurationDetails()
	return &ShowNodePoolConfigurationDetailsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowNodePoolConfigurations 查询指定节点池支持配置的参数内容
//
// 该API用于查询指定节点池支持配置的参数内容。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowNodePoolConfigurations(request *model.ShowNodePoolConfigurationsRequest) (*model.ShowNodePoolConfigurationsResponse, error) {
	requestDef := GenReqDefForShowNodePoolConfigurations()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowNodePoolConfigurationsResponse), nil
	}
}

// ShowNodePoolConfigurationsInvoker 查询指定节点池支持配置的参数内容
func (c *CceClient) ShowNodePoolConfigurationsInvoker(request *model.ShowNodePoolConfigurationsRequest) *ShowNodePoolConfigurationsInvoker {
	requestDef := GenReqDefForShowNodePoolConfigurations()
	return &ShowNodePoolConfigurationsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowPartition 获取分区详情
//
// 获取分区详情
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowPartition(request *model.ShowPartitionRequest) (*model.ShowPartitionResponse, error) {
	requestDef := GenReqDefForShowPartition()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowPartitionResponse), nil
	}
}

// ShowPartitionInvoker 获取分区详情
func (c *CceClient) ShowPartitionInvoker(request *model.ShowPartitionRequest) *ShowPartitionInvoker {
	requestDef := GenReqDefForShowPartition()
	return &ShowPartitionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowPreCheck 获取集群升级前检查任务详情
//
// 获取集群升级前检查任务详情，任务ID由调用集群检查API后从响应体中uid字段获取。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowPreCheck(request *model.ShowPreCheckRequest) (*model.ShowPreCheckResponse, error) {
	requestDef := GenReqDefForShowPreCheck()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowPreCheckResponse), nil
	}
}

// ShowPreCheckInvoker 获取集群升级前检查任务详情
func (c *CceClient) ShowPreCheckInvoker(request *model.ShowPreCheckRequest) *ShowPreCheckInvoker {
	requestDef := GenReqDefForShowPreCheck()
	return &ShowPreCheckInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowQuotas 查询CCE服务下的资源配额
//
// 该API用于查询CCE服务下的资源配额。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowQuotas(request *model.ShowQuotasRequest) (*model.ShowQuotasResponse, error) {
	requestDef := GenReqDefForShowQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowQuotasResponse), nil
	}
}

// ShowQuotasInvoker 查询CCE服务下的资源配额
func (c *CceClient) ShowQuotasInvoker(request *model.ShowQuotasRequest) *ShowQuotasInvoker {
	requestDef := GenReqDefForShowQuotas()
	return &ShowQuotasInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowRelease 获取指定模板实例
//
// 获取指定模板实例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowRelease(request *model.ShowReleaseRequest) (*model.ShowReleaseResponse, error) {
	requestDef := GenReqDefForShowRelease()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowReleaseResponse), nil
	}
}

// ShowReleaseInvoker 获取指定模板实例
func (c *CceClient) ShowReleaseInvoker(request *model.ShowReleaseRequest) *ShowReleaseInvoker {
	requestDef := GenReqDefForShowRelease()
	return &ShowReleaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowReleaseHistory 查询指定模板实例历史记录
//
// 查询指定模板实例历史记录
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowReleaseHistory(request *model.ShowReleaseHistoryRequest) (*model.ShowReleaseHistoryResponse, error) {
	requestDef := GenReqDefForShowReleaseHistory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowReleaseHistoryResponse), nil
	}
}

// ShowReleaseHistoryInvoker 查询指定模板实例历史记录
func (c *CceClient) ShowReleaseHistoryInvoker(request *model.ShowReleaseHistoryRequest) *ShowReleaseHistoryInvoker {
	requestDef := GenReqDefForShowReleaseHistory()
	return &ShowReleaseHistoryInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowUpgradeClusterTask 获取集群升级任务详情
//
// 获取集群升级任务详情，任务ID由调用集群升级API后从响应体中uid字段获取。
// &gt; - 集群升级涉及多维度的组件升级操作，强烈建议统一通过CCE控制台执行交互式升级，降低集群升级过程的业务意外受损风险；
// &gt; - 当前集群升级相关接口受限开放。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowUpgradeClusterTask(request *model.ShowUpgradeClusterTaskRequest) (*model.ShowUpgradeClusterTaskResponse, error) {
	requestDef := GenReqDefForShowUpgradeClusterTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowUpgradeClusterTaskResponse), nil
	}
}

// ShowUpgradeClusterTaskInvoker 获取集群升级任务详情
func (c *CceClient) ShowUpgradeClusterTaskInvoker(request *model.ShowUpgradeClusterTaskRequest) *ShowUpgradeClusterTaskInvoker {
	requestDef := GenReqDefForShowUpgradeClusterTask()
	return &ShowUpgradeClusterTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowUpgradeWorkFlow 获取指定集群升级引导任务详情
//
// 该API用于通过升级引导任务ID获取任务的详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowUpgradeWorkFlow(request *model.ShowUpgradeWorkFlowRequest) (*model.ShowUpgradeWorkFlowResponse, error) {
	requestDef := GenReqDefForShowUpgradeWorkFlow()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowUpgradeWorkFlowResponse), nil
	}
}

// ShowUpgradeWorkFlowInvoker 获取指定集群升级引导任务详情
func (c *CceClient) ShowUpgradeWorkFlowInvoker(request *model.ShowUpgradeWorkFlowRequest) *ShowUpgradeWorkFlowInvoker {
	requestDef := GenReqDefForShowUpgradeWorkFlow()
	return &ShowUpgradeWorkFlowInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowUserChartsQuotas 获取用户模板配额
//
// 获取用户模板配额
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowUserChartsQuotas(request *model.ShowUserChartsQuotasRequest) (*model.ShowUserChartsQuotasResponse, error) {
	requestDef := GenReqDefForShowUserChartsQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowUserChartsQuotasResponse), nil
	}
}

// ShowUserChartsQuotasInvoker 获取用户模板配额
func (c *CceClient) ShowUserChartsQuotasInvoker(request *model.ShowUserChartsQuotasRequest) *ShowUserChartsQuotasInvoker {
	requestDef := GenReqDefForShowUserChartsQuotas()
	return &ShowUserChartsQuotasInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SyncNode 同步节点
//
// 该API用于同步节点。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) SyncNode(request *model.SyncNodeRequest) (*model.SyncNodeResponse, error) {
	requestDef := GenReqDefForSyncNode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SyncNodeResponse), nil
	}
}

// SyncNodeInvoker 同步节点
func (c *CceClient) SyncNodeInvoker(request *model.SyncNodeRequest) *SyncNodeInvoker {
	requestDef := GenReqDefForSyncNode()
	return &SyncNodeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SyncNodePool 同步nodePool配置到存量节点
//
// 该API用于同步nodePool配置到存量节点。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) SyncNodePool(request *model.SyncNodePoolRequest) (*model.SyncNodePoolResponse, error) {
	requestDef := GenReqDefForSyncNodePool()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SyncNodePoolResponse), nil
	}
}

// SyncNodePoolInvoker 同步nodePool配置到存量节点
func (c *CceClient) SyncNodePoolInvoker(request *model.SyncNodePoolRequest) *SyncNodePoolInvoker {
	requestDef := GenReqDefForSyncNodePool()
	return &SyncNodePoolInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAddonInstance 更新AddonInstance
//
// 更新插件实例的功能。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateAddonInstance(request *model.UpdateAddonInstanceRequest) (*model.UpdateAddonInstanceResponse, error) {
	requestDef := GenReqDefForUpdateAddonInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAddonInstanceResponse), nil
	}
}

// UpdateAddonInstanceInvoker 更新AddonInstance
func (c *CceClient) UpdateAddonInstanceInvoker(request *model.UpdateAddonInstanceRequest) *UpdateAddonInstanceInvoker {
	requestDef := GenReqDefForUpdateAddonInstance()
	return &UpdateAddonInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateChart 更新模板
//
// 更新模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateChart(request *model.UpdateChartRequest) (*model.UpdateChartResponse, error) {
	requestDef := GenReqDefForUpdateChart()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateChartResponse), nil
	}
}

// UpdateChartInvoker 更新模板
func (c *CceClient) UpdateChartInvoker(request *model.UpdateChartRequest) *UpdateChartInvoker {
	requestDef := GenReqDefForUpdateChart()
	return &UpdateChartInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateCluster 更新指定的集群
//
// 该API用于更新指定的集群。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateCluster(request *model.UpdateClusterRequest) (*model.UpdateClusterResponse, error) {
	requestDef := GenReqDefForUpdateCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateClusterResponse), nil
	}
}

// UpdateClusterInvoker 更新指定的集群
func (c *CceClient) UpdateClusterInvoker(request *model.UpdateClusterRequest) *UpdateClusterInvoker {
	requestDef := GenReqDefForUpdateCluster()
	return &UpdateClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateClusterEip 绑定、解绑集群公网apiserver地址
//
// 该API用于通过集群ID绑定、解绑集群公网apiserver地址
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateClusterEip(request *model.UpdateClusterEipRequest) (*model.UpdateClusterEipResponse, error) {
	requestDef := GenReqDefForUpdateClusterEip()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateClusterEipResponse), nil
	}
}

// UpdateClusterEipInvoker 绑定、解绑集群公网apiserver地址
func (c *CceClient) UpdateClusterEipInvoker(request *model.UpdateClusterEipRequest) *UpdateClusterEipInvoker {
	requestDef := GenReqDefForUpdateClusterEip()
	return &UpdateClusterEipInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateClusterLogConfig 配置集群日志
//
// 用户可以选择集群管理节点上哪些组件的日志上报LTS
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateClusterLogConfig(request *model.UpdateClusterLogConfigRequest) (*model.UpdateClusterLogConfigResponse, error) {
	requestDef := GenReqDefForUpdateClusterLogConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateClusterLogConfigResponse), nil
	}
}

// UpdateClusterLogConfigInvoker 配置集群日志
func (c *CceClient) UpdateClusterLogConfigInvoker(request *model.UpdateClusterLogConfigRequest) *UpdateClusterLogConfigInvoker {
	requestDef := GenReqDefForUpdateClusterLogConfig()
	return &UpdateClusterLogConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateNode 更新指定的节点
//
// 该API用于更新指定的节点。
// &gt; - 当前仅支持更新metadata下的name字段，即节点的名字。
// &gt; - 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateNode(request *model.UpdateNodeRequest) (*model.UpdateNodeResponse, error) {
	requestDef := GenReqDefForUpdateNode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateNodeResponse), nil
	}
}

// UpdateNodeInvoker 更新指定的节点
func (c *CceClient) UpdateNodeInvoker(request *model.UpdateNodeRequest) *UpdateNodeInvoker {
	requestDef := GenReqDefForUpdateNode()
	return &UpdateNodeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateNodePool 更新指定节点池
//
// 该API用于更新指定的节点池。仅支持集群在处于可用、扩容、缩容状态时调用。
//
// &gt; - 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径
//
// &gt; - 当前仅支持更新节点池名称，spec下的initialNodeCount，k8sTags，taints，login，userTags与节点池的扩缩容配置相关字段。若此次更新未设置相关值，默认更新为初始值。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateNodePool(request *model.UpdateNodePoolRequest) (*model.UpdateNodePoolResponse, error) {
	requestDef := GenReqDefForUpdateNodePool()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateNodePoolResponse), nil
	}
}

// UpdateNodePoolInvoker 更新指定节点池
func (c *CceClient) UpdateNodePoolInvoker(request *model.UpdateNodePoolRequest) *UpdateNodePoolInvoker {
	requestDef := GenReqDefForUpdateNodePool()
	return &UpdateNodePoolInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateNodePoolConfiguration 修改指定节点池配置参数的值
//
// 该API用于修改CCE服务下指定节点池配置参数的值。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateNodePoolConfiguration(request *model.UpdateNodePoolConfigurationRequest) (*model.UpdateNodePoolConfigurationResponse, error) {
	requestDef := GenReqDefForUpdateNodePoolConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateNodePoolConfigurationResponse), nil
	}
}

// UpdateNodePoolConfigurationInvoker 修改指定节点池配置参数的值
func (c *CceClient) UpdateNodePoolConfigurationInvoker(request *model.UpdateNodePoolConfigurationRequest) *UpdateNodePoolConfigurationInvoker {
	requestDef := GenReqDefForUpdateNodePoolConfiguration()
	return &UpdateNodePoolConfigurationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdatePartition 更新分区
//
// 更新分区
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdatePartition(request *model.UpdatePartitionRequest) (*model.UpdatePartitionResponse, error) {
	requestDef := GenReqDefForUpdatePartition()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePartitionResponse), nil
	}
}

// UpdatePartitionInvoker 更新分区
func (c *CceClient) UpdatePartitionInvoker(request *model.UpdatePartitionRequest) *UpdatePartitionInvoker {
	requestDef := GenReqDefForUpdatePartition()
	return &UpdatePartitionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateRelease 更新指定模板实例
//
// 更新指定模板实例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateRelease(request *model.UpdateReleaseRequest) (*model.UpdateReleaseResponse, error) {
	requestDef := GenReqDefForUpdateRelease()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateReleaseResponse), nil
	}
}

// UpdateReleaseInvoker 更新指定模板实例
func (c *CceClient) UpdateReleaseInvoker(request *model.UpdateReleaseRequest) *UpdateReleaseInvoker {
	requestDef := GenReqDefForUpdateRelease()
	return &UpdateReleaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpgradeCluster 集群升级
//
// 集群升级。
// &gt; - 集群升级涉及多维度的组件升级操作，强烈建议统一通过CCE控制台执行交互式升级，降低集群升级过程的业务意外受损风险；
// &gt; - 当前集群升级相关接口受限开放。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpgradeCluster(request *model.UpgradeClusterRequest) (*model.UpgradeClusterResponse, error) {
	requestDef := GenReqDefForUpgradeCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeClusterResponse), nil
	}
}

// UpgradeClusterInvoker 集群升级
func (c *CceClient) UpgradeClusterInvoker(request *model.UpgradeClusterRequest) *UpgradeClusterInvoker {
	requestDef := GenReqDefForUpgradeCluster()
	return &UpgradeClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpgradeNodePool 节点池升级
//
// 该API用于节点池升级。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpgradeNodePool(request *model.UpgradeNodePoolRequest) (*model.UpgradeNodePoolResponse, error) {
	requestDef := GenReqDefForUpgradeNodePool()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeNodePoolResponse), nil
	}
}

// UpgradeNodePoolInvoker 节点池升级
func (c *CceClient) UpgradeNodePoolInvoker(request *model.UpgradeNodePoolRequest) *UpgradeNodePoolInvoker {
	requestDef := GenReqDefForUpgradeNodePool()
	return &UpgradeNodePoolInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpgradeWorkFlowUpdate 更新指定集群升级引导任务状态
//
// 该API用于更新指定集群升级引导任务状态，当前仅适用于取消升级流程
// 调用该API时升级流程引导任务状态不能为进行中(running) 已完成(success) 已取消(cancel),升级子任务状态不能为running(进行中) init(已初始化) pause(任务被暂停) queue(队列中)
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpgradeWorkFlowUpdate(request *model.UpgradeWorkFlowUpdateRequest) (*model.UpgradeWorkFlowUpdateResponse, error) {
	requestDef := GenReqDefForUpgradeWorkFlowUpdate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeWorkFlowUpdateResponse), nil
	}
}

// UpgradeWorkFlowUpdateInvoker 更新指定集群升级引导任务状态
func (c *CceClient) UpgradeWorkFlowUpdateInvoker(request *model.UpgradeWorkFlowUpdateRequest) *UpgradeWorkFlowUpdateInvoker {
	requestDef := GenReqDefForUpgradeWorkFlowUpdate()
	return &UpgradeWorkFlowUpdateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UploadChart 上传模板
//
// 上传模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UploadChart(request *model.UploadChartRequest) (*model.UploadChartResponse, error) {
	requestDef := GenReqDefForUploadChart()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UploadChartResponse), nil
	}
}

// UploadChartInvoker 上传模板
func (c *CceClient) UploadChartInvoker(request *model.UploadChartRequest) *UploadChartInvoker {
	requestDef := GenReqDefForUploadChart()
	return &UploadChartInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowVersion 查询API版本信息列表
//
// 该API用于查询CCE服务当前支持的API版本信息列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowVersion(request *model.ShowVersionRequest) (*model.ShowVersionResponse, error) {
	requestDef := GenReqDefForShowVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowVersionResponse), nil
	}
}

// ShowVersionInvoker 查询API版本信息列表
func (c *CceClient) ShowVersionInvoker(request *model.ShowVersionRequest) *ShowVersionInvoker {
	requestDef := GenReqDefForShowVersion()
	return &ShowVersionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchCreateAutopilotClusterTags 批量添加指定集群的资源标签
//
// 该API用于批量添加指定集群的资源标签。
// &gt; - 每个集群支持最多20个资源标签。
// &gt; - 此接口为幂等接口：创建时，如果创建的标签已经存在（key/value均相同视为重复），默认处理成功；key相同，value不同时会覆盖原有标签。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) BatchCreateAutopilotClusterTags(request *model.BatchCreateAutopilotClusterTagsRequest) (*model.BatchCreateAutopilotClusterTagsResponse, error) {
	requestDef := GenReqDefForBatchCreateAutopilotClusterTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchCreateAutopilotClusterTagsResponse), nil
	}
}

// BatchCreateAutopilotClusterTagsInvoker 批量添加指定集群的资源标签
func (c *CceClient) BatchCreateAutopilotClusterTagsInvoker(request *model.BatchCreateAutopilotClusterTagsRequest) *BatchCreateAutopilotClusterTagsInvoker {
	requestDef := GenReqDefForBatchCreateAutopilotClusterTags()
	return &BatchCreateAutopilotClusterTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchDeleteAutopilotClusterTags 批量删除指定集群的资源标签
//
// 该API用于批量删除指定集群的资源标签。
// &gt; - 此接口为幂等接口：删除时，如果删除的标签key不存在，默认处理成功。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) BatchDeleteAutopilotClusterTags(request *model.BatchDeleteAutopilotClusterTagsRequest) (*model.BatchDeleteAutopilotClusterTagsResponse, error) {
	requestDef := GenReqDefForBatchDeleteAutopilotClusterTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchDeleteAutopilotClusterTagsResponse), nil
	}
}

// BatchDeleteAutopilotClusterTagsInvoker 批量删除指定集群的资源标签
func (c *CceClient) BatchDeleteAutopilotClusterTagsInvoker(request *model.BatchDeleteAutopilotClusterTagsRequest) *BatchDeleteAutopilotClusterTagsInvoker {
	requestDef := GenReqDefForBatchDeleteAutopilotClusterTags()
	return &BatchDeleteAutopilotClusterTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAutopilotAddonInstance 创建AddonInstance
//
// 根据提供的插件模板，安装插件实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateAutopilotAddonInstance(request *model.CreateAutopilotAddonInstanceRequest) (*model.CreateAutopilotAddonInstanceResponse, error) {
	requestDef := GenReqDefForCreateAutopilotAddonInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAutopilotAddonInstanceResponse), nil
	}
}

// CreateAutopilotAddonInstanceInvoker 创建AddonInstance
func (c *CceClient) CreateAutopilotAddonInstanceInvoker(request *model.CreateAutopilotAddonInstanceRequest) *CreateAutopilotAddonInstanceInvoker {
	requestDef := GenReqDefForCreateAutopilotAddonInstance()
	return &CreateAutopilotAddonInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAutopilotCluster 创建集群
//
// 该API用于创建一个空集群（即只有控制节点Master，没有工作节点Node）。请在调用本接口完成集群创建之后，通过[创建节点](cce_02_0242.xml)添加节点。
//
// &gt;   - 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
// &gt;   - 调用该接口创建集群时，默认不安装ICAgent，若需安装ICAgent，可在请求Body参数的annotations中加入\&quot;cluster.install.addons.external/install\&quot;:\&quot;[{\&quot;addonTemplateName\&quot;:\&quot;icagent\&quot;}]\&quot;的集群注解，将在创建集群时自动安装ICAgent。ICAgent是应用性能管理APM的采集代理，运行在应用所在的服务器上，用于实时采集探针所获取的数据，安装ICAgent是使用应用性能管理APM的前提。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateAutopilotCluster(request *model.CreateAutopilotClusterRequest) (*model.CreateAutopilotClusterResponse, error) {
	requestDef := GenReqDefForCreateAutopilotCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAutopilotClusterResponse), nil
	}
}

// CreateAutopilotClusterInvoker 创建集群
func (c *CceClient) CreateAutopilotClusterInvoker(request *model.CreateAutopilotClusterRequest) *CreateAutopilotClusterInvoker {
	requestDef := GenReqDefForCreateAutopilotCluster()
	return &CreateAutopilotClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAutopilotClusterMasterSnapshot 集群备份
//
// 集群备份
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateAutopilotClusterMasterSnapshot(request *model.CreateAutopilotClusterMasterSnapshotRequest) (*model.CreateAutopilotClusterMasterSnapshotResponse, error) {
	requestDef := GenReqDefForCreateAutopilotClusterMasterSnapshot()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAutopilotClusterMasterSnapshotResponse), nil
	}
}

// CreateAutopilotClusterMasterSnapshotInvoker 集群备份
func (c *CceClient) CreateAutopilotClusterMasterSnapshotInvoker(request *model.CreateAutopilotClusterMasterSnapshotRequest) *CreateAutopilotClusterMasterSnapshotInvoker {
	requestDef := GenReqDefForCreateAutopilotClusterMasterSnapshot()
	return &CreateAutopilotClusterMasterSnapshotInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAutopilotKubernetesClusterCert 获取集群证书
//
// 该API用于获取指定集群的证书信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateAutopilotKubernetesClusterCert(request *model.CreateAutopilotKubernetesClusterCertRequest) (*model.CreateAutopilotKubernetesClusterCertResponse, error) {
	requestDef := GenReqDefForCreateAutopilotKubernetesClusterCert()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAutopilotKubernetesClusterCertResponse), nil
	}
}

// CreateAutopilotKubernetesClusterCertInvoker 获取集群证书
func (c *CceClient) CreateAutopilotKubernetesClusterCertInvoker(request *model.CreateAutopilotKubernetesClusterCertRequest) *CreateAutopilotKubernetesClusterCertInvoker {
	requestDef := GenReqDefForCreateAutopilotKubernetesClusterCert()
	return &CreateAutopilotKubernetesClusterCertInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAutopilotPostCheck 集群升级后确认
//
// 集群升级后确认，该接口建议配合Console使用，主要用于升级步骤完成后，客户确认集群状态和业务正常后做反馈。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateAutopilotPostCheck(request *model.CreateAutopilotPostCheckRequest) (*model.CreateAutopilotPostCheckResponse, error) {
	requestDef := GenReqDefForCreateAutopilotPostCheck()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAutopilotPostCheckResponse), nil
	}
}

// CreateAutopilotPostCheckInvoker 集群升级后确认
func (c *CceClient) CreateAutopilotPostCheckInvoker(request *model.CreateAutopilotPostCheckRequest) *CreateAutopilotPostCheckInvoker {
	requestDef := GenReqDefForCreateAutopilotPostCheck()
	return &CreateAutopilotPostCheckInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAutopilotPreCheck 集群升级前检查
//
// 集群升级前检查
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateAutopilotPreCheck(request *model.CreateAutopilotPreCheckRequest) (*model.CreateAutopilotPreCheckResponse, error) {
	requestDef := GenReqDefForCreateAutopilotPreCheck()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAutopilotPreCheckResponse), nil
	}
}

// CreateAutopilotPreCheckInvoker 集群升级前检查
func (c *CceClient) CreateAutopilotPreCheckInvoker(request *model.CreateAutopilotPreCheckRequest) *CreateAutopilotPreCheckInvoker {
	requestDef := GenReqDefForCreateAutopilotPreCheck()
	return &CreateAutopilotPreCheckInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAutopilotRelease 创建模板实例
//
// 创建模板实例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateAutopilotRelease(request *model.CreateAutopilotReleaseRequest) (*model.CreateAutopilotReleaseResponse, error) {
	requestDef := GenReqDefForCreateAutopilotRelease()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAutopilotReleaseResponse), nil
	}
}

// CreateAutopilotReleaseInvoker 创建模板实例
func (c *CceClient) CreateAutopilotReleaseInvoker(request *model.CreateAutopilotReleaseRequest) *CreateAutopilotReleaseInvoker {
	requestDef := GenReqDefForCreateAutopilotRelease()
	return &CreateAutopilotReleaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAutopilotUpgradeWorkFlow 开启集群升级流程引导任务
//
// 该API用于创建一个集群升级流程引导任务。请在调用本接口完成引导任务创建之后，通过集群升级前检查开始检查任务。
// 升级流程任务用于控制集群升级任务的执行流程，执行流程为 升级前检查 &#x3D;&gt; 集群升级 &#x3D;&gt; 升级后检查。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) CreateAutopilotUpgradeWorkFlow(request *model.CreateAutopilotUpgradeWorkFlowRequest) (*model.CreateAutopilotUpgradeWorkFlowResponse, error) {
	requestDef := GenReqDefForCreateAutopilotUpgradeWorkFlow()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAutopilotUpgradeWorkFlowResponse), nil
	}
}

// CreateAutopilotUpgradeWorkFlowInvoker 开启集群升级流程引导任务
func (c *CceClient) CreateAutopilotUpgradeWorkFlowInvoker(request *model.CreateAutopilotUpgradeWorkFlowRequest) *CreateAutopilotUpgradeWorkFlowInvoker {
	requestDef := GenReqDefForCreateAutopilotUpgradeWorkFlow()
	return &CreateAutopilotUpgradeWorkFlowInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAutopilotAddonInstance 删除AddonInstance
//
// 删除插件实例的功能。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DeleteAutopilotAddonInstance(request *model.DeleteAutopilotAddonInstanceRequest) (*model.DeleteAutopilotAddonInstanceResponse, error) {
	requestDef := GenReqDefForDeleteAutopilotAddonInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAutopilotAddonInstanceResponse), nil
	}
}

// DeleteAutopilotAddonInstanceInvoker 删除AddonInstance
func (c *CceClient) DeleteAutopilotAddonInstanceInvoker(request *model.DeleteAutopilotAddonInstanceRequest) *DeleteAutopilotAddonInstanceInvoker {
	requestDef := GenReqDefForDeleteAutopilotAddonInstance()
	return &DeleteAutopilotAddonInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAutopilotChart 删除模板
//
// 删除模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DeleteAutopilotChart(request *model.DeleteAutopilotChartRequest) (*model.DeleteAutopilotChartResponse, error) {
	requestDef := GenReqDefForDeleteAutopilotChart()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAutopilotChartResponse), nil
	}
}

// DeleteAutopilotChartInvoker 删除模板
func (c *CceClient) DeleteAutopilotChartInvoker(request *model.DeleteAutopilotChartRequest) *DeleteAutopilotChartInvoker {
	requestDef := GenReqDefForDeleteAutopilotChart()
	return &DeleteAutopilotChartInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAutopilotCluster 删除集群
//
// 该API用于删除一个指定的集群。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DeleteAutopilotCluster(request *model.DeleteAutopilotClusterRequest) (*model.DeleteAutopilotClusterResponse, error) {
	requestDef := GenReqDefForDeleteAutopilotCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAutopilotClusterResponse), nil
	}
}

// DeleteAutopilotClusterInvoker 删除集群
func (c *CceClient) DeleteAutopilotClusterInvoker(request *model.DeleteAutopilotClusterRequest) *DeleteAutopilotClusterInvoker {
	requestDef := GenReqDefForDeleteAutopilotCluster()
	return &DeleteAutopilotClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAutopilotRelease 删除指定模板实例
//
// 删除指定模板实例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DeleteAutopilotRelease(request *model.DeleteAutopilotReleaseRequest) (*model.DeleteAutopilotReleaseResponse, error) {
	requestDef := GenReqDefForDeleteAutopilotRelease()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAutopilotReleaseResponse), nil
	}
}

// DeleteAutopilotReleaseInvoker 删除指定模板实例
func (c *CceClient) DeleteAutopilotReleaseInvoker(request *model.DeleteAutopilotReleaseRequest) *DeleteAutopilotReleaseInvoker {
	requestDef := GenReqDefForDeleteAutopilotRelease()
	return &DeleteAutopilotReleaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DownloadAutopilotChart 下载模板
//
// 下载模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) DownloadAutopilotChart(request *model.DownloadAutopilotChartRequest) (*model.DownloadAutopilotChartResponse, error) {
	requestDef := GenReqDefForDownloadAutopilotChart()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DownloadAutopilotChartResponse), nil
	}
}

// DownloadAutopilotChartInvoker 下载模板
func (c *CceClient) DownloadAutopilotChartInvoker(request *model.DownloadAutopilotChartRequest) *DownloadAutopilotChartInvoker {
	requestDef := GenReqDefForDownloadAutopilotChart()
	return &DownloadAutopilotChartInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutopilotAddonInstances 获取AddonInstance列表
//
// 获取集群所有已安装插件实例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAutopilotAddonInstances(request *model.ListAutopilotAddonInstancesRequest) (*model.ListAutopilotAddonInstancesResponse, error) {
	requestDef := GenReqDefForListAutopilotAddonInstances()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutopilotAddonInstancesResponse), nil
	}
}

// ListAutopilotAddonInstancesInvoker 获取AddonInstance列表
func (c *CceClient) ListAutopilotAddonInstancesInvoker(request *model.ListAutopilotAddonInstancesRequest) *ListAutopilotAddonInstancesInvoker {
	requestDef := GenReqDefForListAutopilotAddonInstances()
	return &ListAutopilotAddonInstancesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutopilotAddonTemplates 查询AddonTemplates列表
//
// 插件模板查询接口，查询插件信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAutopilotAddonTemplates(request *model.ListAutopilotAddonTemplatesRequest) (*model.ListAutopilotAddonTemplatesResponse, error) {
	requestDef := GenReqDefForListAutopilotAddonTemplates()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutopilotAddonTemplatesResponse), nil
	}
}

// ListAutopilotAddonTemplatesInvoker 查询AddonTemplates列表
func (c *CceClient) ListAutopilotAddonTemplatesInvoker(request *model.ListAutopilotAddonTemplatesRequest) *ListAutopilotAddonTemplatesInvoker {
	requestDef := GenReqDefForListAutopilotAddonTemplates()
	return &ListAutopilotAddonTemplatesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutopilotCharts 获取模板列表
//
// 获取模板列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAutopilotCharts(request *model.ListAutopilotChartsRequest) (*model.ListAutopilotChartsResponse, error) {
	requestDef := GenReqDefForListAutopilotCharts()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutopilotChartsResponse), nil
	}
}

// ListAutopilotChartsInvoker 获取模板列表
func (c *CceClient) ListAutopilotChartsInvoker(request *model.ListAutopilotChartsRequest) *ListAutopilotChartsInvoker {
	requestDef := GenReqDefForListAutopilotCharts()
	return &ListAutopilotChartsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutopilotClusterMasterSnapshotTasks 获取集群备份任务详情列表
//
// 获取集群备份任务详情列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAutopilotClusterMasterSnapshotTasks(request *model.ListAutopilotClusterMasterSnapshotTasksRequest) (*model.ListAutopilotClusterMasterSnapshotTasksResponse, error) {
	requestDef := GenReqDefForListAutopilotClusterMasterSnapshotTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutopilotClusterMasterSnapshotTasksResponse), nil
	}
}

// ListAutopilotClusterMasterSnapshotTasksInvoker 获取集群备份任务详情列表
func (c *CceClient) ListAutopilotClusterMasterSnapshotTasksInvoker(request *model.ListAutopilotClusterMasterSnapshotTasksRequest) *ListAutopilotClusterMasterSnapshotTasksInvoker {
	requestDef := GenReqDefForListAutopilotClusterMasterSnapshotTasks()
	return &ListAutopilotClusterMasterSnapshotTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutopilotClusterUpgradeFeatureGates 获取集群升级特性开关配置
//
// 获取集群升级特性开关配置
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAutopilotClusterUpgradeFeatureGates(request *model.ListAutopilotClusterUpgradeFeatureGatesRequest) (*model.ListAutopilotClusterUpgradeFeatureGatesResponse, error) {
	requestDef := GenReqDefForListAutopilotClusterUpgradeFeatureGates()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutopilotClusterUpgradeFeatureGatesResponse), nil
	}
}

// ListAutopilotClusterUpgradeFeatureGatesInvoker 获取集群升级特性开关配置
func (c *CceClient) ListAutopilotClusterUpgradeFeatureGatesInvoker(request *model.ListAutopilotClusterUpgradeFeatureGatesRequest) *ListAutopilotClusterUpgradeFeatureGatesInvoker {
	requestDef := GenReqDefForListAutopilotClusterUpgradeFeatureGates()
	return &ListAutopilotClusterUpgradeFeatureGatesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutopilotClusterUpgradePaths 获取集群升级路径
//
// 获取集群升级路径
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAutopilotClusterUpgradePaths(request *model.ListAutopilotClusterUpgradePathsRequest) (*model.ListAutopilotClusterUpgradePathsResponse, error) {
	requestDef := GenReqDefForListAutopilotClusterUpgradePaths()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutopilotClusterUpgradePathsResponse), nil
	}
}

// ListAutopilotClusterUpgradePathsInvoker 获取集群升级路径
func (c *CceClient) ListAutopilotClusterUpgradePathsInvoker(request *model.ListAutopilotClusterUpgradePathsRequest) *ListAutopilotClusterUpgradePathsInvoker {
	requestDef := GenReqDefForListAutopilotClusterUpgradePaths()
	return &ListAutopilotClusterUpgradePathsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutopilotClusters 获取指定项目下的集群
//
// 该API用于获取指定项目下所有集群的详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAutopilotClusters(request *model.ListAutopilotClustersRequest) (*model.ListAutopilotClustersResponse, error) {
	requestDef := GenReqDefForListAutopilotClusters()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutopilotClustersResponse), nil
	}
}

// ListAutopilotClustersInvoker 获取指定项目下的集群
func (c *CceClient) ListAutopilotClustersInvoker(request *model.ListAutopilotClustersRequest) *ListAutopilotClustersInvoker {
	requestDef := GenReqDefForListAutopilotClusters()
	return &ListAutopilotClustersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutopilotPreCheckTasks 获取集群升级前检查任务详情列表
//
// 获取集群升级前检查任务详情列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAutopilotPreCheckTasks(request *model.ListAutopilotPreCheckTasksRequest) (*model.ListAutopilotPreCheckTasksResponse, error) {
	requestDef := GenReqDefForListAutopilotPreCheckTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutopilotPreCheckTasksResponse), nil
	}
}

// ListAutopilotPreCheckTasksInvoker 获取集群升级前检查任务详情列表
func (c *CceClient) ListAutopilotPreCheckTasksInvoker(request *model.ListAutopilotPreCheckTasksRequest) *ListAutopilotPreCheckTasksInvoker {
	requestDef := GenReqDefForListAutopilotPreCheckTasks()
	return &ListAutopilotPreCheckTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutopilotReleases 获取模板实例列表
//
// 获取模板实例列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAutopilotReleases(request *model.ListAutopilotReleasesRequest) (*model.ListAutopilotReleasesResponse, error) {
	requestDef := GenReqDefForListAutopilotReleases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutopilotReleasesResponse), nil
	}
}

// ListAutopilotReleasesInvoker 获取模板实例列表
func (c *CceClient) ListAutopilotReleasesInvoker(request *model.ListAutopilotReleasesRequest) *ListAutopilotReleasesInvoker {
	requestDef := GenReqDefForListAutopilotReleases()
	return &ListAutopilotReleasesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutopilotUpgradeClusterTasks 获取集群升级任务详情列表
//
// 获取集群升级任务详情列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAutopilotUpgradeClusterTasks(request *model.ListAutopilotUpgradeClusterTasksRequest) (*model.ListAutopilotUpgradeClusterTasksResponse, error) {
	requestDef := GenReqDefForListAutopilotUpgradeClusterTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutopilotUpgradeClusterTasksResponse), nil
	}
}

// ListAutopilotUpgradeClusterTasksInvoker 获取集群升级任务详情列表
func (c *CceClient) ListAutopilotUpgradeClusterTasksInvoker(request *model.ListAutopilotUpgradeClusterTasksRequest) *ListAutopilotUpgradeClusterTasksInvoker {
	requestDef := GenReqDefForListAutopilotUpgradeClusterTasks()
	return &ListAutopilotUpgradeClusterTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAutopilotUpgradeWorkFlows 获取UpgradeWorkFlows列表
//
// 获取历史集群升级引导任务列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ListAutopilotUpgradeWorkFlows(request *model.ListAutopilotUpgradeWorkFlowsRequest) (*model.ListAutopilotUpgradeWorkFlowsResponse, error) {
	requestDef := GenReqDefForListAutopilotUpgradeWorkFlows()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAutopilotUpgradeWorkFlowsResponse), nil
	}
}

// ListAutopilotUpgradeWorkFlowsInvoker 获取UpgradeWorkFlows列表
func (c *CceClient) ListAutopilotUpgradeWorkFlowsInvoker(request *model.ListAutopilotUpgradeWorkFlowsRequest) *ListAutopilotUpgradeWorkFlowsInvoker {
	requestDef := GenReqDefForListAutopilotUpgradeWorkFlows()
	return &ListAutopilotUpgradeWorkFlowsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RetryAutopilotUpgradeClusterTask 重试集群升级任务
//
// 重新执行失败的集群升级任务。
// &gt; - 集群升级涉及多维度的组件升级操作，强烈建议统一通过CCE控制台执行交互式升级，降低集群升级过程的业务意外受损风险；
// &gt; - 当前集群升级相关接口受限开放。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) RetryAutopilotUpgradeClusterTask(request *model.RetryAutopilotUpgradeClusterTaskRequest) (*model.RetryAutopilotUpgradeClusterTaskResponse, error) {
	requestDef := GenReqDefForRetryAutopilotUpgradeClusterTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RetryAutopilotUpgradeClusterTaskResponse), nil
	}
}

// RetryAutopilotUpgradeClusterTaskInvoker 重试集群升级任务
func (c *CceClient) RetryAutopilotUpgradeClusterTaskInvoker(request *model.RetryAutopilotUpgradeClusterTaskRequest) *RetryAutopilotUpgradeClusterTaskInvoker {
	requestDef := GenReqDefForRetryAutopilotUpgradeClusterTask()
	return &RetryAutopilotUpgradeClusterTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RollbackAutopilotAddonInstance 回滚AddonInstance
//
// 将插件实例回滚到升级前的版本。只有在当前插件实例版本支持回滚到升级前的版本（status.isRollbackable为true），且插件实例状态为running（运行中）、available（可用）、abnormal（不可用）、upgradeFailed（升级失败）、rollbackFailed（回滚失败）时支持回滚。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) RollbackAutopilotAddonInstance(request *model.RollbackAutopilotAddonInstanceRequest) (*model.RollbackAutopilotAddonInstanceResponse, error) {
	requestDef := GenReqDefForRollbackAutopilotAddonInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RollbackAutopilotAddonInstanceResponse), nil
	}
}

// RollbackAutopilotAddonInstanceInvoker 回滚AddonInstance
func (c *CceClient) RollbackAutopilotAddonInstanceInvoker(request *model.RollbackAutopilotAddonInstanceRequest) *RollbackAutopilotAddonInstanceInvoker {
	requestDef := GenReqDefForRollbackAutopilotAddonInstance()
	return &RollbackAutopilotAddonInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotAddonInstance 获取AddonInstance详情
//
// 获取插件实例详情。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotAddonInstance(request *model.ShowAutopilotAddonInstanceRequest) (*model.ShowAutopilotAddonInstanceResponse, error) {
	requestDef := GenReqDefForShowAutopilotAddonInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotAddonInstanceResponse), nil
	}
}

// ShowAutopilotAddonInstanceInvoker 获取AddonInstance详情
func (c *CceClient) ShowAutopilotAddonInstanceInvoker(request *model.ShowAutopilotAddonInstanceRequest) *ShowAutopilotAddonInstanceInvoker {
	requestDef := GenReqDefForShowAutopilotAddonInstance()
	return &ShowAutopilotAddonInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotChart 获取模板
//
// 获取模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotChart(request *model.ShowAutopilotChartRequest) (*model.ShowAutopilotChartResponse, error) {
	requestDef := GenReqDefForShowAutopilotChart()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotChartResponse), nil
	}
}

// ShowAutopilotChartInvoker 获取模板
func (c *CceClient) ShowAutopilotChartInvoker(request *model.ShowAutopilotChartRequest) *ShowAutopilotChartInvoker {
	requestDef := GenReqDefForShowAutopilotChart()
	return &ShowAutopilotChartInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotChartValues 获取模板Values
//
// 获取模板Values
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotChartValues(request *model.ShowAutopilotChartValuesRequest) (*model.ShowAutopilotChartValuesResponse, error) {
	requestDef := GenReqDefForShowAutopilotChartValues()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotChartValuesResponse), nil
	}
}

// ShowAutopilotChartValuesInvoker 获取模板Values
func (c *CceClient) ShowAutopilotChartValuesInvoker(request *model.ShowAutopilotChartValuesRequest) *ShowAutopilotChartValuesInvoker {
	requestDef := GenReqDefForShowAutopilotChartValues()
	return &ShowAutopilotChartValuesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotCluster 获取指定的集群
//
// 该API用于获取指定集群的详细信息。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotCluster(request *model.ShowAutopilotClusterRequest) (*model.ShowAutopilotClusterResponse, error) {
	requestDef := GenReqDefForShowAutopilotCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotClusterResponse), nil
	}
}

// ShowAutopilotClusterInvoker 获取指定的集群
func (c *CceClient) ShowAutopilotClusterInvoker(request *model.ShowAutopilotClusterRequest) *ShowAutopilotClusterInvoker {
	requestDef := GenReqDefForShowAutopilotCluster()
	return &ShowAutopilotClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotClusterEndpoints 获取集群访问的地址
//
// 该API用于通过集群ID获取集群访问的地址，包括PrivateIP(HA集群返回VIP)与PublicIP
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotClusterEndpoints(request *model.ShowAutopilotClusterEndpointsRequest) (*model.ShowAutopilotClusterEndpointsResponse, error) {
	requestDef := GenReqDefForShowAutopilotClusterEndpoints()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotClusterEndpointsResponse), nil
	}
}

// ShowAutopilotClusterEndpointsInvoker 获取集群访问的地址
func (c *CceClient) ShowAutopilotClusterEndpointsInvoker(request *model.ShowAutopilotClusterEndpointsRequest) *ShowAutopilotClusterEndpointsInvoker {
	requestDef := GenReqDefForShowAutopilotClusterEndpoints()
	return &ShowAutopilotClusterEndpointsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotClusterUpgradeInfo 获取集群升级相关信息
//
// 获取集群升级相关信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotClusterUpgradeInfo(request *model.ShowAutopilotClusterUpgradeInfoRequest) (*model.ShowAutopilotClusterUpgradeInfoResponse, error) {
	requestDef := GenReqDefForShowAutopilotClusterUpgradeInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotClusterUpgradeInfoResponse), nil
	}
}

// ShowAutopilotClusterUpgradeInfoInvoker 获取集群升级相关信息
func (c *CceClient) ShowAutopilotClusterUpgradeInfoInvoker(request *model.ShowAutopilotClusterUpgradeInfoRequest) *ShowAutopilotClusterUpgradeInfoInvoker {
	requestDef := GenReqDefForShowAutopilotClusterUpgradeInfo()
	return &ShowAutopilotClusterUpgradeInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotJob 获取任务信息
//
// 该API用于获取任务信息。通过某一任务请求下发后返回的jobID来查询指定任务的进度。
// &gt; - 集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径
// &gt; - 该接口通常使用场景为：
// &gt;   - 创建、删除集群时，查询相应任务的进度。
// &gt;   - 创建、删除节点时，查询相应任务的进度。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotJob(request *model.ShowAutopilotJobRequest) (*model.ShowAutopilotJobResponse, error) {
	requestDef := GenReqDefForShowAutopilotJob()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotJobResponse), nil
	}
}

// ShowAutopilotJobInvoker 获取任务信息
func (c *CceClient) ShowAutopilotJobInvoker(request *model.ShowAutopilotJobRequest) *ShowAutopilotJobInvoker {
	requestDef := GenReqDefForShowAutopilotJob()
	return &ShowAutopilotJobInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotPreCheck 获取集群升级前检查任务详情
//
// 获取集群升级前检查任务详情，任务ID由调用集群检查API后从响应体中uid字段获取。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotPreCheck(request *model.ShowAutopilotPreCheckRequest) (*model.ShowAutopilotPreCheckResponse, error) {
	requestDef := GenReqDefForShowAutopilotPreCheck()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotPreCheckResponse), nil
	}
}

// ShowAutopilotPreCheckInvoker 获取集群升级前检查任务详情
func (c *CceClient) ShowAutopilotPreCheckInvoker(request *model.ShowAutopilotPreCheckRequest) *ShowAutopilotPreCheckInvoker {
	requestDef := GenReqDefForShowAutopilotPreCheck()
	return &ShowAutopilotPreCheckInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotQuotas 查询CCE服务下的资源配额
//
// 该API用于查询CCE服务下的资源配额。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotQuotas(request *model.ShowAutopilotQuotasRequest) (*model.ShowAutopilotQuotasResponse, error) {
	requestDef := GenReqDefForShowAutopilotQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotQuotasResponse), nil
	}
}

// ShowAutopilotQuotasInvoker 查询CCE服务下的资源配额
func (c *CceClient) ShowAutopilotQuotasInvoker(request *model.ShowAutopilotQuotasRequest) *ShowAutopilotQuotasInvoker {
	requestDef := GenReqDefForShowAutopilotQuotas()
	return &ShowAutopilotQuotasInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotRelease 获取指定模板实例
//
// 获取指定模板实例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotRelease(request *model.ShowAutopilotReleaseRequest) (*model.ShowAutopilotReleaseResponse, error) {
	requestDef := GenReqDefForShowAutopilotRelease()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotReleaseResponse), nil
	}
}

// ShowAutopilotReleaseInvoker 获取指定模板实例
func (c *CceClient) ShowAutopilotReleaseInvoker(request *model.ShowAutopilotReleaseRequest) *ShowAutopilotReleaseInvoker {
	requestDef := GenReqDefForShowAutopilotRelease()
	return &ShowAutopilotReleaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotReleaseHistory 查询指定模板实例历史记录
//
// 查询指定模板实例历史记录
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotReleaseHistory(request *model.ShowAutopilotReleaseHistoryRequest) (*model.ShowAutopilotReleaseHistoryResponse, error) {
	requestDef := GenReqDefForShowAutopilotReleaseHistory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotReleaseHistoryResponse), nil
	}
}

// ShowAutopilotReleaseHistoryInvoker 查询指定模板实例历史记录
func (c *CceClient) ShowAutopilotReleaseHistoryInvoker(request *model.ShowAutopilotReleaseHistoryRequest) *ShowAutopilotReleaseHistoryInvoker {
	requestDef := GenReqDefForShowAutopilotReleaseHistory()
	return &ShowAutopilotReleaseHistoryInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotUpgradeClusterTask 获取集群升级任务详情
//
// 获取集群升级任务详情，任务ID由调用集群升级API后从响应体中uid字段获取。
// &gt; - 集群升级涉及多维度的组件升级操作，强烈建议统一通过CCE控制台执行交互式升级，降低集群升级过程的业务意外受损风险；
// &gt; - 当前集群升级相关接口受限开放。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotUpgradeClusterTask(request *model.ShowAutopilotUpgradeClusterTaskRequest) (*model.ShowAutopilotUpgradeClusterTaskResponse, error) {
	requestDef := GenReqDefForShowAutopilotUpgradeClusterTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotUpgradeClusterTaskResponse), nil
	}
}

// ShowAutopilotUpgradeClusterTaskInvoker 获取集群升级任务详情
func (c *CceClient) ShowAutopilotUpgradeClusterTaskInvoker(request *model.ShowAutopilotUpgradeClusterTaskRequest) *ShowAutopilotUpgradeClusterTaskInvoker {
	requestDef := GenReqDefForShowAutopilotUpgradeClusterTask()
	return &ShowAutopilotUpgradeClusterTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotUpgradeWorkFlow 获取指定集群升级引导任务详情
//
// 该API用于通过升级引导任务ID获取任务的详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotUpgradeWorkFlow(request *model.ShowAutopilotUpgradeWorkFlowRequest) (*model.ShowAutopilotUpgradeWorkFlowResponse, error) {
	requestDef := GenReqDefForShowAutopilotUpgradeWorkFlow()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotUpgradeWorkFlowResponse), nil
	}
}

// ShowAutopilotUpgradeWorkFlowInvoker 获取指定集群升级引导任务详情
func (c *CceClient) ShowAutopilotUpgradeWorkFlowInvoker(request *model.ShowAutopilotUpgradeWorkFlowRequest) *ShowAutopilotUpgradeWorkFlowInvoker {
	requestDef := GenReqDefForShowAutopilotUpgradeWorkFlow()
	return &ShowAutopilotUpgradeWorkFlowInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutopilotUserChartsQuotas 获取用户模板配额
//
// 获取用户模板配额
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) ShowAutopilotUserChartsQuotas(request *model.ShowAutopilotUserChartsQuotasRequest) (*model.ShowAutopilotUserChartsQuotasResponse, error) {
	requestDef := GenReqDefForShowAutopilotUserChartsQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutopilotUserChartsQuotasResponse), nil
	}
}

// ShowAutopilotUserChartsQuotasInvoker 获取用户模板配额
func (c *CceClient) ShowAutopilotUserChartsQuotasInvoker(request *model.ShowAutopilotUserChartsQuotasRequest) *ShowAutopilotUserChartsQuotasInvoker {
	requestDef := GenReqDefForShowAutopilotUserChartsQuotas()
	return &ShowAutopilotUserChartsQuotasInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAutopilotAddonInstance 更新AddonInstance
//
// 更新插件实例的功能。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateAutopilotAddonInstance(request *model.UpdateAutopilotAddonInstanceRequest) (*model.UpdateAutopilotAddonInstanceResponse, error) {
	requestDef := GenReqDefForUpdateAutopilotAddonInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAutopilotAddonInstanceResponse), nil
	}
}

// UpdateAutopilotAddonInstanceInvoker 更新AddonInstance
func (c *CceClient) UpdateAutopilotAddonInstanceInvoker(request *model.UpdateAutopilotAddonInstanceRequest) *UpdateAutopilotAddonInstanceInvoker {
	requestDef := GenReqDefForUpdateAutopilotAddonInstance()
	return &UpdateAutopilotAddonInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAutopilotChart 更新模板
//
// 更新模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateAutopilotChart(request *model.UpdateAutopilotChartRequest) (*model.UpdateAutopilotChartResponse, error) {
	requestDef := GenReqDefForUpdateAutopilotChart()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAutopilotChartResponse), nil
	}
}

// UpdateAutopilotChartInvoker 更新模板
func (c *CceClient) UpdateAutopilotChartInvoker(request *model.UpdateAutopilotChartRequest) *UpdateAutopilotChartInvoker {
	requestDef := GenReqDefForUpdateAutopilotChart()
	return &UpdateAutopilotChartInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAutopilotCluster 更新指定的集群
//
// 该API用于更新指定的集群。
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateAutopilotCluster(request *model.UpdateAutopilotClusterRequest) (*model.UpdateAutopilotClusterResponse, error) {
	requestDef := GenReqDefForUpdateAutopilotCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAutopilotClusterResponse), nil
	}
}

// UpdateAutopilotClusterInvoker 更新指定的集群
func (c *CceClient) UpdateAutopilotClusterInvoker(request *model.UpdateAutopilotClusterRequest) *UpdateAutopilotClusterInvoker {
	requestDef := GenReqDefForUpdateAutopilotCluster()
	return &UpdateAutopilotClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAutopilotClusterEip 绑定、解绑集群公网apiserver地址
//
// 该API用于通过集群ID绑定、解绑集群公网apiserver地址
// &gt;集群管理的URL格式为：https://Endpoint/uri。其中uri为资源路径，也即API访问的路径。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateAutopilotClusterEip(request *model.UpdateAutopilotClusterEipRequest) (*model.UpdateAutopilotClusterEipResponse, error) {
	requestDef := GenReqDefForUpdateAutopilotClusterEip()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAutopilotClusterEipResponse), nil
	}
}

// UpdateAutopilotClusterEipInvoker 绑定、解绑集群公网apiserver地址
func (c *CceClient) UpdateAutopilotClusterEipInvoker(request *model.UpdateAutopilotClusterEipRequest) *UpdateAutopilotClusterEipInvoker {
	requestDef := GenReqDefForUpdateAutopilotClusterEip()
	return &UpdateAutopilotClusterEipInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAutopilotRelease 更新指定模板实例
//
// 更新指定模板实例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpdateAutopilotRelease(request *model.UpdateAutopilotReleaseRequest) (*model.UpdateAutopilotReleaseResponse, error) {
	requestDef := GenReqDefForUpdateAutopilotRelease()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAutopilotReleaseResponse), nil
	}
}

// UpdateAutopilotReleaseInvoker 更新指定模板实例
func (c *CceClient) UpdateAutopilotReleaseInvoker(request *model.UpdateAutopilotReleaseRequest) *UpdateAutopilotReleaseInvoker {
	requestDef := GenReqDefForUpdateAutopilotRelease()
	return &UpdateAutopilotReleaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpgradeAutopilotCluster 集群升级
//
// 集群升级。
// &gt; - 集群升级涉及多维度的组件升级操作，强烈建议统一通过CCE控制台执行交互式升级，降低集群升级过程的业务意外受损风险；
// &gt; - 当前集群升级相关接口受限开放。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpgradeAutopilotCluster(request *model.UpgradeAutopilotClusterRequest) (*model.UpgradeAutopilotClusterResponse, error) {
	requestDef := GenReqDefForUpgradeAutopilotCluster()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeAutopilotClusterResponse), nil
	}
}

// UpgradeAutopilotClusterInvoker 集群升级
func (c *CceClient) UpgradeAutopilotClusterInvoker(request *model.UpgradeAutopilotClusterRequest) *UpgradeAutopilotClusterInvoker {
	requestDef := GenReqDefForUpgradeAutopilotCluster()
	return &UpgradeAutopilotClusterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpgradeAutopilotWorkFlowUpdate 更新指定集群升级引导任务状态
//
// 该API用于更新指定集群升级引导任务状态，当前仅适用于取消升级流程
// 调用该API时升级流程引导任务状态不能为进行中(running) 已完成(success) 已取消(cancel),升级子任务状态不能为running(进行中) init(已初始化) pause(任务被暂停) queue(队列中)
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UpgradeAutopilotWorkFlowUpdate(request *model.UpgradeAutopilotWorkFlowUpdateRequest) (*model.UpgradeAutopilotWorkFlowUpdateResponse, error) {
	requestDef := GenReqDefForUpgradeAutopilotWorkFlowUpdate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeAutopilotWorkFlowUpdateResponse), nil
	}
}

// UpgradeAutopilotWorkFlowUpdateInvoker 更新指定集群升级引导任务状态
func (c *CceClient) UpgradeAutopilotWorkFlowUpdateInvoker(request *model.UpgradeAutopilotWorkFlowUpdateRequest) *UpgradeAutopilotWorkFlowUpdateInvoker {
	requestDef := GenReqDefForUpgradeAutopilotWorkFlowUpdate()
	return &UpgradeAutopilotWorkFlowUpdateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UploadAutopilotChart 上传模板
//
// 上传模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CceClient) UploadAutopilotChart(request *model.UploadAutopilotChartRequest) (*model.UploadAutopilotChartResponse, error) {
	requestDef := GenReqDefForUploadAutopilotChart()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UploadAutopilotChartResponse), nil
	}
}

// UploadAutopilotChartInvoker 上传模板
func (c *CceClient) UploadAutopilotChartInvoker(request *model.UploadAutopilotChartRequest) *UploadAutopilotChartInvoker {
	requestDef := GenReqDefForUploadAutopilotChart()
	return &UploadAutopilotChartInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
