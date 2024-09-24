package v3

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cce/v3/model"
	"net/http"
)

func GenReqDefForAddNode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/add").
		WithResponse(new(model.AddNodeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddNodesToNodePool() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}/nodes/add").
		WithResponse(new(model.AddNodesToNodePoolResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodepoolId").
		WithJsonTag("nodepool_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAwakeCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/awake").
		WithResponse(new(model.AwakeClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchCreateClusterTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/tags/create").
		WithResponse(new(model.BatchCreateClusterTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchDeleteClusterTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/tags/delete").
		WithResponse(new(model.BatchDeleteClusterTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchSyncNodes() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/sync").
		WithResponse(new(model.BatchSyncNodesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForContinueUpgradeClusterTask() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade/continue").
		WithResponse(new(model.ContinueUpgradeClusterTaskResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateAddonInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/addons").
		WithResponse(new(model.CreateAddonInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateCloudPersistentVolumeClaims() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v1/namespaces/{namespace}/cloudpersistentvolumeclaims").
		WithResponse(new(model.CreateCloudPersistentVolumeClaimsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("XClusterID").
		WithJsonTag("X-Cluster-ID").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters").
		WithResponse(new(model.CreateClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateClusterMasterSnapshot() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3.1/projects/{project_id}/clusters/{cluster_id}/operation/snapshot").
		WithResponse(new(model.CreateClusterMasterSnapshotResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateKubernetesClusterCert() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/clustercert").
		WithResponse(new(model.CreateKubernetesClusterCertResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("PortID").
		WithJsonTag("Port-ID").
		WithKindName("string").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateNode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodes").
		WithResponse(new(model.CreateNodeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodepoolScaleUp").
		WithJsonTag("nodepoolScaleUp").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateNodePool() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools").
		WithResponse(new(model.CreateNodePoolResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreatePartition() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/partitions").
		WithResponse(new(model.CreatePartitionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreatePostCheck() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/postcheck").
		WithResponse(new(model.CreatePostCheckResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreatePreCheck() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/precheck").
		WithResponse(new(model.CreatePreCheckResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateRelease() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/cce/cam/v3/clusters/{cluster_id}/releases").
		WithResponse(new(model.CreateReleaseResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateUpgradeWorkFlow() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgradeworkflows").
		WithResponse(new(model.CreateUpgradeWorkFlowResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteAddonInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/api/v3/addons/{id}").
		WithResponse(new(model.DeleteAddonInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Query))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteChart() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v2/charts/{chart_id}").
		WithResponse(new(model.DeleteChartResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ChartId").
		WithJsonTag("chart_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteCloudPersistentVolumeClaims() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/api/v1/namespaces/{namespace}/cloudpersistentvolumeclaims/{name}").
		WithResponse(new(model.DeleteCloudPersistentVolumeClaimsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteVolume").
		WithJsonTag("deleteVolume").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("StorageType").
		WithJsonTag("storageType").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("XClusterID").
		WithJsonTag("X-Cluster-ID").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}").
		WithResponse(new(model.DeleteClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteEfs").
		WithJsonTag("delete_efs").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteEni").
		WithJsonTag("delete_eni").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteEvs").
		WithJsonTag("delete_evs").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteNet").
		WithJsonTag("delete_net").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteObs").
		WithJsonTag("delete_obs").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteSfs").
		WithJsonTag("delete_sfs").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteSfs30").
		WithJsonTag("delete_sfs30").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("LtsReclaimPolicy").
		WithJsonTag("lts_reclaim_policy").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Tobedeleted").
		WithJsonTag("tobedeleted").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("OndemandNodePolicy").
		WithJsonTag("ondemand_node_policy").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PeriodicNodePolicy").
		WithJsonTag("periodic_node_policy").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteNode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/{node_id}").
		WithResponse(new(model.DeleteNodeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodeId").
		WithJsonTag("node_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodepoolScaleDown").
		WithJsonTag("nodepoolScaleDown").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteNodePool() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}").
		WithResponse(new(model.DeleteNodePoolResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodepoolId").
		WithJsonTag("nodepool_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteRelease() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}").
		WithResponse(new(model.DeleteReleaseResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDownloadChart() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v2/charts/{chart_id}/archive").
		WithResponse(new(model.DownloadChartResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ChartId").
		WithJsonTag("chart_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForHibernateCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/hibernate").
		WithResponse(new(model.HibernateClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAddonInstances() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/addons").
		WithResponse(new(model.ListAddonInstancesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AddonTemplateName").
		WithJsonTag("addon_template_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAddonTemplates() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/addontemplates").
		WithResponse(new(model.ListAddonTemplatesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AddonTemplateName").
		WithJsonTag("addon_template_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("BaseUpdateAddonVersion").
		WithJsonTag("base_update_addon_version").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Newest").
		WithJsonTag("newest").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Version").
		WithJsonTag("version").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListCharts() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v2/charts").
		WithResponse(new(model.ListChartsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListClusterMasterSnapshotTasks() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3.1/projects/{project_id}/clusters/{cluster_id}/operation/snapshot/tasks").
		WithResponse(new(model.ListClusterMasterSnapshotTasksResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListClusterUpgradeFeatureGates() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/clusterupgradefeaturegates").
		WithResponse(new(model.ListClusterUpgradeFeatureGatesResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListClusterUpgradePaths() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/clusterupgradepaths").
		WithResponse(new(model.ListClusterUpgradePathsResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListClusters() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters").
		WithResponse(new(model.ListClustersResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Detail").
		WithJsonTag("detail").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Status").
		WithJsonTag("status").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Type").
		WithJsonTag("type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Version").
		WithJsonTag("version").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListNodePools() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools").
		WithResponse(new(model.ListNodePoolsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ShowDefaultNodePool").
		WithJsonTag("showDefaultNodePool").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListNodes() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodes").
		WithResponse(new(model.ListNodesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListPartitions() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/partitions").
		WithResponse(new(model.ListPartitionsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListPreCheckTasks() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/precheck/tasks").
		WithResponse(new(model.ListPreCheckTasksResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListReleases() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/cce/cam/v3/clusters/{cluster_id}/releases").
		WithResponse(new(model.ListReleasesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ChartId").
		WithJsonTag("chart_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Query))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListUpgradeClusterTasks() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade/tasks").
		WithResponse(new(model.ListUpgradeClusterTasksResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListUpgradeWorkFlows() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgradeworkflows").
		WithResponse(new(model.ListUpgradeWorkFlowsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForMigrateNode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/operation/migrateto/{target_cluster_id}").
		WithResponse(new(model.MigrateNodeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TargetClusterId").
		WithJsonTag("target_cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForPauseUpgradeClusterTask() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade/pause").
		WithResponse(new(model.PauseUpgradeClusterTaskResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForRemoveNode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/operation/remove").
		WithResponse(new(model.RemoveNodeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForResetNode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/reset").
		WithResponse(new(model.ResetNodeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForResizeCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/resize").
		WithResponse(new(model.ResizeClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForRetryUpgradeClusterTask() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade/retry").
		WithResponse(new(model.RetryUpgradeClusterTaskResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForRevokeKubernetesClusterCert() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/clustercertrevoke").
		WithResponse(new(model.RevokeKubernetesClusterCertResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForRollbackAddonInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/addons/{id}/operation/rollback").
		WithResponse(new(model.RollbackAddonInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForScaleNodePool() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}/operation/scale").
		WithResponse(new(model.ScaleNodePoolResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodepoolId").
		WithJsonTag("nodepool_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAddonInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/addons/{id}").
		WithResponse(new(model.ShowAddonInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowChart() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v2/charts/{chart_id}").
		WithResponse(new(model.ShowChartResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ChartId").
		WithJsonTag("chart_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowChartValues() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v2/charts/{chart_id}/values").
		WithResponse(new(model.ShowChartValuesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ChartId").
		WithJsonTag("chart_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}").
		WithResponse(new(model.ShowClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Detail").
		WithJsonTag("detail").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowClusterConfig() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/cluster/{cluster_id}/log-configs").
		WithResponse(new(model.ShowClusterConfigResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Type").
		WithJsonTag("type").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowClusterConfigurationDetails() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/configuration/detail").
		WithResponse(new(model.ShowClusterConfigurationDetailsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowClusterEndpoints() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/openapi").
		WithResponse(new(model.ShowClusterEndpointsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowClusterSupportConfiguration() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/clusters/configuration/detail").
		WithResponse(new(model.ShowClusterSupportConfigurationResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterType").
		WithJsonTag("cluster_type").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterVersion").
		WithJsonTag("cluster_version").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NetworkMode").
		WithJsonTag("network_mode").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowClusterUpgradeInfo() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/upgradeinfo").
		WithResponse(new(model.ShowClusterUpgradeInfoResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowJob() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/jobs/{job_id}").
		WithResponse(new(model.ShowJobResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("JobId").
		WithJsonTag("job_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowNode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/{node_id}").
		WithResponse(new(model.ShowNodeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodeId").
		WithJsonTag("node_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowNodePool() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}").
		WithResponse(new(model.ShowNodePoolResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodepoolId").
		WithJsonTag("nodepool_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowNodePoolConfigurationDetails() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}/configuration/detail").
		WithResponse(new(model.ShowNodePoolConfigurationDetailsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodepoolId").
		WithJsonTag("nodepool_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowNodePoolConfigurations() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}/configuration").
		WithResponse(new(model.ShowNodePoolConfigurationsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodepoolId").
		WithJsonTag("nodepool_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowPartition() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/partitions/{partition_name}").
		WithResponse(new(model.ShowPartitionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PartitionName").
		WithJsonTag("partition_name").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowPreCheck() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/precheck/tasks/{task_id}").
		WithResponse(new(model.ShowPreCheckResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TaskId").
		WithJsonTag("task_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowQuotas() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/quotas").
		WithResponse(new(model.ShowQuotasResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowRelease() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}").
		WithResponse(new(model.ShowReleaseResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowReleaseHistory() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}/history").
		WithResponse(new(model.ShowReleaseHistoryResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowUpgradeClusterTask() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade/tasks/{task_id}").
		WithResponse(new(model.ShowUpgradeClusterTaskResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TaskId").
		WithJsonTag("task_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowUpgradeWorkFlow() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgradeworkflows/{upgrade_workflow_id}").
		WithResponse(new(model.ShowUpgradeWorkFlowResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("UpgradeWorkflowId").
		WithJsonTag("upgrade_workflow_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowUserChartsQuotas() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v2/charts/{project_id}/quotas").
		WithResponse(new(model.ShowUserChartsQuotasResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForSyncNode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/api/v2/projects/{project_id}/clusters/{cluster_id}/nodes/{node_id}/sync").
		WithResponse(new(model.SyncNodeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodeId").
		WithJsonTag("node_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForSyncNodePool() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/api/v3.1/projects/{project_id}/clusters/{cluster_id}/nodepool/{nodepool_id}/sync").
		WithResponse(new(model.SyncNodePoolResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodepoolId").
		WithJsonTag("nodepool_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateAddonInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/api/v3/addons/{id}").
		WithResponse(new(model.UpdateAddonInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateChart() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v2/charts/{chart_id}").
		WithResponse(new(model.UpdateChartResponse)).
		WithContentType("multipart/form-data")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ChartId").
		WithJsonTag("chart_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}").
		WithResponse(new(model.UpdateClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateClusterEip() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/mastereip").
		WithResponse(new(model.UpdateClusterEipResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateClusterLogConfig() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/api/v3/projects/{project_id}/cluster/{cluster_id}/log-configs").
		WithResponse(new(model.UpdateClusterLogConfigResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateNode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/{node_id}").
		WithResponse(new(model.UpdateNodeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodeId").
		WithJsonTag("node_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateNodePool() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}").
		WithResponse(new(model.UpdateNodePoolResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodepoolId").
		WithJsonTag("nodepool_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateNodePoolConfiguration() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}/configuration").
		WithResponse(new(model.UpdateNodePoolConfigurationResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodepoolId").
		WithJsonTag("nodepool_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdatePartition() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/partitions/{partition_name}").
		WithResponse(new(model.UpdatePartitionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PartitionName").
		WithJsonTag("partition_name").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateRelease() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}").
		WithResponse(new(model.UpdateReleaseResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpgradeCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade").
		WithResponse(new(model.UpgradeClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpgradeNodePool() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}/operation/upgrade").
		WithResponse(new(model.UpgradeNodePoolResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodepoolId").
		WithJsonTag("nodepool_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpgradeWorkFlowUpdate() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPatch).
		WithPath("/api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgradeworkflows/{upgrade_workflow_id}").
		WithResponse(new(model.UpgradeWorkFlowUpdateResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("UpgradeWorkflowId").
		WithJsonTag("upgrade_workflow_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUploadChart() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v2/charts").
		WithResponse(new(model.UploadChartResponse)).
		WithContentType("multipart/form-data")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowVersion() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/").
		WithResponse(new(model.ShowVersionResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchCreateAutopilotClusterTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/tags/create").
		WithResponse(new(model.BatchCreateAutopilotClusterTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchDeleteAutopilotClusterTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/tags/delete").
		WithResponse(new(model.BatchDeleteAutopilotClusterTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateAutopilotAddonInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v3/addons").
		WithResponse(new(model.CreateAutopilotAddonInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateAutopilotCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v3/projects/{project_id}/clusters").
		WithResponse(new(model.CreateAutopilotClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateAutopilotClusterMasterSnapshot() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v3.1/projects/{project_id}/clusters/{cluster_id}/operation/snapshot").
		WithResponse(new(model.CreateAutopilotClusterMasterSnapshotResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateAutopilotKubernetesClusterCert() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/clustercert").
		WithResponse(new(model.CreateAutopilotKubernetesClusterCertResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("PortID").
		WithJsonTag("Port-ID").
		WithKindName("string").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateAutopilotPostCheck() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/postcheck").
		WithResponse(new(model.CreateAutopilotPostCheckResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateAutopilotPreCheck() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/precheck").
		WithResponse(new(model.CreateAutopilotPreCheckResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateAutopilotRelease() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/cam/v3/clusters/{cluster_id}/releases").
		WithResponse(new(model.CreateAutopilotReleaseResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateAutopilotUpgradeWorkFlow() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgradeworkflows").
		WithResponse(new(model.CreateAutopilotUpgradeWorkFlowResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteAutopilotAddonInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/autopilot/v3/addons/{id}").
		WithResponse(new(model.DeleteAutopilotAddonInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Query))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteAutopilotChart() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/autopilot/v2/charts/{chart_id}").
		WithResponse(new(model.DeleteAutopilotChartResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ChartId").
		WithJsonTag("chart_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteAutopilotCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}").
		WithResponse(new(model.DeleteAutopilotClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteEfs").
		WithJsonTag("delete_efs").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteEni").
		WithJsonTag("delete_eni").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteNet").
		WithJsonTag("delete_net").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteObs").
		WithJsonTag("delete_obs").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeleteSfs30").
		WithJsonTag("delete_sfs30").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("LtsReclaimPolicy").
		WithJsonTag("lts_reclaim_policy").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteAutopilotRelease() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}").
		WithResponse(new(model.DeleteAutopilotReleaseResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDownloadAutopilotChart() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v2/charts/{chart_id}/archive").
		WithResponse(new(model.DownloadAutopilotChartResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ChartId").
		WithJsonTag("chart_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAutopilotAddonInstances() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/addons").
		WithResponse(new(model.ListAutopilotAddonInstancesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AddonTemplateName").
		WithJsonTag("addon_template_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAutopilotAddonTemplates() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/addontemplates").
		WithResponse(new(model.ListAutopilotAddonTemplatesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AddonTemplateName").
		WithJsonTag("addon_template_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("BaseUpdateAddonVersion").
		WithJsonTag("base_update_addon_version").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Newest").
		WithJsonTag("newest").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Version").
		WithJsonTag("version").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAutopilotCharts() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v2/charts").
		WithResponse(new(model.ListAutopilotChartsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAutopilotClusterMasterSnapshotTasks() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3.1/projects/{project_id}/clusters/{cluster_id}/operation/snapshot/tasks").
		WithResponse(new(model.ListAutopilotClusterMasterSnapshotTasksResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAutopilotClusterUpgradeFeatureGates() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/clusterupgradefeaturegates").
		WithResponse(new(model.ListAutopilotClusterUpgradeFeatureGatesResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAutopilotClusterUpgradePaths() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/clusterupgradepaths").
		WithResponse(new(model.ListAutopilotClusterUpgradePathsResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAutopilotClusters() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/projects/{project_id}/clusters").
		WithResponse(new(model.ListAutopilotClustersResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Detail").
		WithJsonTag("detail").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Status").
		WithJsonTag("status").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Type").
		WithJsonTag("type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Version").
		WithJsonTag("version").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAutopilotPreCheckTasks() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/precheck/tasks").
		WithResponse(new(model.ListAutopilotPreCheckTasksResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAutopilotReleases() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/cam/v3/clusters/{cluster_id}/releases").
		WithResponse(new(model.ListAutopilotReleasesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ChartId").
		WithJsonTag("chart_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Query))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAutopilotUpgradeClusterTasks() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade/tasks").
		WithResponse(new(model.ListAutopilotUpgradeClusterTasksResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAutopilotUpgradeWorkFlows() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgradeworkflows").
		WithResponse(new(model.ListAutopilotUpgradeWorkFlowsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForRetryAutopilotUpgradeClusterTask() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade/retry").
		WithResponse(new(model.RetryAutopilotUpgradeClusterTaskResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForRollbackAutopilotAddonInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v3/addons/{id}/operation/rollback").
		WithResponse(new(model.RollbackAutopilotAddonInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotAddonInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/addons/{id}").
		WithResponse(new(model.ShowAutopilotAddonInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotChart() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v2/charts/{chart_id}").
		WithResponse(new(model.ShowAutopilotChartResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ChartId").
		WithJsonTag("chart_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotChartValues() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v2/charts/{chart_id}/values").
		WithResponse(new(model.ShowAutopilotChartValuesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ChartId").
		WithJsonTag("chart_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}").
		WithResponse(new(model.ShowAutopilotClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Detail").
		WithJsonTag("detail").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotClusterEndpoints() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/openapi").
		WithResponse(new(model.ShowAutopilotClusterEndpointsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotClusterUpgradeInfo() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/upgradeinfo").
		WithResponse(new(model.ShowAutopilotClusterUpgradeInfoResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotJob() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/projects/{project_id}/jobs/{job_id}").
		WithResponse(new(model.ShowAutopilotJobResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("JobId").
		WithJsonTag("job_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotPreCheck() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/precheck/tasks/{task_id}").
		WithResponse(new(model.ShowAutopilotPreCheckResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TaskId").
		WithJsonTag("task_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotQuotas() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/projects/{project_id}/quotas").
		WithResponse(new(model.ShowAutopilotQuotasResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotRelease() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}").
		WithResponse(new(model.ShowAutopilotReleaseResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotReleaseHistory() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}/history").
		WithResponse(new(model.ShowAutopilotReleaseHistoryResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotUpgradeClusterTask() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade/tasks/{task_id}").
		WithResponse(new(model.ShowAutopilotUpgradeClusterTaskResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TaskId").
		WithJsonTag("task_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotUpgradeWorkFlow() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgradeworkflows/{upgrade_workflow_id}").
		WithResponse(new(model.ShowAutopilotUpgradeWorkFlowResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("UpgradeWorkflowId").
		WithJsonTag("upgrade_workflow_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutopilotUserChartsQuotas() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/autopilot/v2/charts/{project_id}/quotas").
		WithResponse(new(model.ShowAutopilotUserChartsQuotasResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateAutopilotAddonInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/autopilot/v3/addons/{id}").
		WithResponse(new(model.UpdateAutopilotAddonInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateAutopilotChart() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/autopilot/v2/charts/{chart_id}").
		WithResponse(new(model.UpdateAutopilotChartResponse)).
		WithContentType("multipart/form-data")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ChartId").
		WithJsonTag("chart_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateAutopilotCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}").
		WithResponse(new(model.UpdateAutopilotClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ErrorStatus").
		WithJsonTag("errorStatus").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateAutopilotClusterEip() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/mastereip").
		WithResponse(new(model.UpdateAutopilotClusterEipResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateAutopilotRelease() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}").
		WithResponse(new(model.UpdateAutopilotReleaseResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpgradeAutopilotCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade").
		WithResponse(new(model.UpgradeAutopilotClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpgradeAutopilotWorkFlowUpdate() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPatch).
		WithPath("/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgradeworkflows/{upgrade_workflow_id}").
		WithResponse(new(model.UpgradeAutopilotWorkFlowUpdateResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("UpgradeWorkflowId").
		WithJsonTag("upgrade_workflow_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUploadAutopilotChart() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/autopilot/v2/charts").
		WithResponse(new(model.UploadAutopilotChartResponse)).
		WithContentType("multipart/form-data")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}
