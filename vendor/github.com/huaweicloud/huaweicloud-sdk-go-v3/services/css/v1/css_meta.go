package v1

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v1/model"
	"net/http"
)

func GenReqDefForAddIndependentNode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/type/{type}/independent").
		WithResponse(new(model.AddIndependentNodeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Type").
		WithJsonTag("type").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForChangeMode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/mode/change").
		WithResponse(new(model.ChangeModeResponse)).
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

func GenReqDefForChangeSecurityGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/sg/change").
		WithResponse(new(model.ChangeSecurityGroupResponse)).
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

func GenReqDefForCreateAiOps() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/ai-ops").
		WithResponse(new(model.CreateAiOpsResponse)).
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

func GenReqDefForCreateAutoCreatePolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/policy").
		WithResponse(new(model.CreateAutoCreatePolicyResponse)).
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

func GenReqDefForCreateBindPublic() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/public/open").
		WithResponse(new(model.CreateBindPublicResponse)).
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

func GenReqDefForCreateCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters").
		WithResponse(new(model.CreateClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateClustersTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/{resource_type}/{cluster_id}/tags").
		WithResponse(new(model.CreateClustersTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceType").
		WithJsonTag("resource_type").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateElbListener() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/es-listeners").
		WithResponse(new(model.CreateElbListenerResponse)).
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

func GenReqDefForCreateLoadIkThesaurus() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/thesaurus").
		WithResponse(new(model.CreateLoadIkThesaurusResponse)).
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

func GenReqDefForCreateLogBackup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/logs/collect").
		WithResponse(new(model.CreateLogBackupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateSnapshot() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/index_snapshot").
		WithResponse(new(model.CreateSnapshotResponse)).
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

func GenReqDefForDeleteAiOps() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/ai-ops/{aiops_id}").
		WithResponse(new(model.DeleteAiOpsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AiopsId").
		WithJsonTag("aiops_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}").
		WithResponse(new(model.DeleteClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteClustersTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v1.0/{project_id}/{resource_type}/{cluster_id}/tags/{key}").
		WithResponse(new(model.DeleteClustersTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceType").
		WithJsonTag("resource_type").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Key").
		WithJsonTag("key").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteIkThesaurus() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/thesaurus").
		WithResponse(new(model.DeleteIkThesaurusResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteSnapshot() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/{snapshot_id}").
		WithResponse(new(model.DeleteSnapshotResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SnapshotId").
		WithJsonTag("snapshot_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDownloadCert() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/cer/download").
		WithResponse(new(model.DownloadCertResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForEnableOrDisableElb() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/loadbalancers/es-switch").
		WithResponse(new(model.EnableOrDisableElbResponse)).
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

func GenReqDefForListAiOps() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/ai-ops").
		WithResponse(new(model.ListAiOpsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Start").
		WithJsonTag("start").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListClustersDetails() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters").
		WithResponse(new(model.ListClustersDetailsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Start").
		WithJsonTag("start").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DatastoreType").
		WithJsonTag("datastoreType").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListClustersTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/{resource_type}/tags").
		WithResponse(new(model.ListClustersTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceType").
		WithJsonTag("resource_type").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListElbCerts() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/elb/certificates").
		WithResponse(new(model.ListElbCertsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListElbs() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/loadbalancers").
		WithResponse(new(model.ListElbsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListFlavors() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/es-flavors").
		WithResponse(new(model.ListFlavorsResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListImages() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/target/{upgrade_type}/images").
		WithResponse(new(model.ListImagesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("UpgradeType").
		WithJsonTag("upgrade_type").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListLogsJob() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/logs/records").
		WithResponse(new(model.ListLogsJobResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Start").
		WithJsonTag("start").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListSmnTopics() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/domains/{domain_id}/ai-ops/smn-topics").
		WithResponse(new(model.ListSmnTopicsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DomainId").
		WithJsonTag("domain_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListSnapshots() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/index_snapshots").
		WithResponse(new(model.ListSnapshotsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListYmls() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/ymls/template").
		WithResponse(new(model.ListYmlsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListYmlsJob() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/ymls/joblists").
		WithResponse(new(model.ListYmlsJobResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Start").
		WithJsonTag("start").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForResetPassword() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/password/reset").
		WithResponse(new(model.ResetPasswordResponse)).
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

func GenReqDefForRestartCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/restart").
		WithResponse(new(model.RestartClusterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForRestoreSnapshot() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/{snapshot_id}/restore").
		WithResponse(new(model.RestoreSnapshotResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SnapshotId").
		WithJsonTag("snapshot_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForRetryUpgradeTask() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/upgrade/{action_id}/retry").
		WithResponse(new(model.RetryUpgradeTaskResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ActionId").
		WithJsonTag("action_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RetryMode").
		WithJsonTag("retry_mode").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAutoCreatePolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/policy").
		WithResponse(new(model.ShowAutoCreatePolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowClusterDetail() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}").
		WithResponse(new(model.ShowClusterDetailResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowClusterTag() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/{resource_type}/{cluster_id}/tags").
		WithResponse(new(model.ShowClusterTagResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceType").
		WithJsonTag("resource_type").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowElbDetail() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/es-listeners").
		WithResponse(new(model.ShowElbDetailResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowGetLogSetting() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/logs/settings").
		WithResponse(new(model.ShowGetLogSettingResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowIkThesaurus() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/thesaurus").
		WithResponse(new(model.ShowIkThesaurusResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowLogBackup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/logs/search").
		WithResponse(new(model.ShowLogBackupResponse)).
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

func GenReqDefForShowVpcepConnection() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/connections").
		WithResponse(new(model.ShowVpcepConnectionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Start").
		WithJsonTag("start").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForStartAutoSetting() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/auto_setting").
		WithResponse(new(model.StartAutoSettingResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForStartLogAutoBackupPolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/logs/policy/update").
		WithResponse(new(model.StartLogAutoBackupPolicyResponse)).
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

func GenReqDefForStartLogs() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/logs/open").
		WithResponse(new(model.StartLogsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Action").
		WithJsonTag("action").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForStartPublicWhitelist() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/public/whitelist/update").
		WithResponse(new(model.StartPublicWhitelistResponse)).
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

func GenReqDefForStartTargetClusterConnectivityTest() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/logs/connectivity").
		WithResponse(new(model.StartTargetClusterConnectivityTestResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForStartVpecp() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/open").
		WithResponse(new(model.StartVpecpResponse)).
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

func GenReqDefForStopLogAutoBackupPolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/logs/policy/close").
		WithResponse(new(model.StopLogAutoBackupPolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForStopLogs() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/logs/close").
		WithResponse(new(model.StopLogsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForStopPublicWhitelist() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/public/whitelist/close").
		WithResponse(new(model.StopPublicWhitelistResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForStopSnapshot() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/index_snapshots").
		WithResponse(new(model.StopSnapshotResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForStopVpecp() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/close").
		WithResponse(new(model.StopVpecpResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateAzByInstanceType() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/inst-type/{inst_type}/azmigrate").
		WithResponse(new(model.UpdateAzByInstanceTypeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstType").
		WithJsonTag("inst_type").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateBatchClustersTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/{resource_type}/{cluster_id}/tags/action").
		WithResponse(new(model.UpdateBatchClustersTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceType").
		WithJsonTag("resource_type").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateClusterName() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/changename").
		WithResponse(new(model.UpdateClusterNameResponse)).
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

func GenReqDefForUpdateEsListener() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/es-listeners/{listener_id}").
		WithResponse(new(model.UpdateEsListenerResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ListenerId").
		WithJsonTag("listener_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateExtendCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/extend").
		WithResponse(new(model.UpdateExtendClusterResponse)).
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

func GenReqDefForUpdateExtendInstanceStorage() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/role_extend").
		WithResponse(new(model.UpdateExtendInstanceStorageResponse)).
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

func GenReqDefForUpdateFlavor() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/flavor").
		WithResponse(new(model.UpdateFlavorResponse)).
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

func GenReqDefForUpdateFlavorByType() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/{types}/flavor").
		WithResponse(new(model.UpdateFlavorByTypeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Types").
		WithJsonTag("types").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/instance/{instance_id}/replace").
		WithResponse(new(model.UpdateInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("instance_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateLogSetting() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/logs/settings").
		WithResponse(new(model.UpdateLogSettingResponse)).
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

func GenReqDefForUpdateOndemandClusterToPeriod() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/cluster/{cluster_id}/period").
		WithResponse(new(model.UpdateOndemandClusterToPeriodResponse)).
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

func GenReqDefForUpdatePublicBandWidth() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/public/bandwidth").
		WithResponse(new(model.UpdatePublicBandWidthResponse)).
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

func GenReqDefForUpdateShrinkCluster() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/extend/{project_id}/clusters/{cluster_id}/role/shrink").
		WithResponse(new(model.UpdateShrinkClusterResponse)).
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

func GenReqDefForUpdateShrinkNodes() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/node/offline").
		WithResponse(new(model.UpdateShrinkNodesResponse)).
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

func GenReqDefForUpdateSnapshotSetting() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/setting").
		WithResponse(new(model.UpdateSnapshotSettingResponse)).
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

func GenReqDefForUpdateUnbindPublic() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/public/close").
		WithResponse(new(model.UpdateUnbindPublicResponse)).
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

func GenReqDefForUpdateVpcepConnection() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/connections").
		WithResponse(new(model.UpdateVpcepConnectionResponse)).
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

func GenReqDefForUpdateVpcepWhitelist() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/permissions").
		WithResponse(new(model.UpdateVpcepWhitelistResponse)).
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

func GenReqDefForUpdateYmls() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/ymls/update").
		WithResponse(new(model.UpdateYmlsResponse)).
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

func GenReqDefForUpgradeCore() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/inst-type/{inst_type}/image/upgrade").
		WithResponse(new(model.UpgradeCoreResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstType").
		WithJsonTag("inst_type").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpgradeDetail() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/upgrade/detail").
		WithResponse(new(model.UpgradeDetailResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Start").
		WithJsonTag("start").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ActionMode").
		WithJsonTag("action_mode").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForStartKibanaPublic() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/publickibana/open").
		WithResponse(new(model.StartKibanaPublicResponse)).
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

func GenReqDefForStopPublicKibanaWhitelist() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/publickibana/whitelist/close").
		WithResponse(new(model.StopPublicKibanaWhitelistResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateAlterKibana() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/publickibana/bandwidth").
		WithResponse(new(model.UpdateAlterKibanaResponse)).
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

func GenReqDefForUpdateCloseKibana() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/publickibana/close").
		WithResponse(new(model.UpdateCloseKibanaResponse)).
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

func GenReqDefForUpdatePublicKibanaWhitelist() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/publickibana/whitelist/update").
		WithResponse(new(model.UpdatePublicKibanaWhitelistResponse)).
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

func GenReqDefForAddFavorite() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/lgsconf/favorite").
		WithResponse(new(model.AddFavoriteResponse)).
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

func GenReqDefForCreateCnf() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/lgsconf/submit").
		WithResponse(new(model.CreateCnfResponse)).
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

func GenReqDefForDeleteConf() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/lgsconf/delete").
		WithResponse(new(model.DeleteConfResponse)).
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

func GenReqDefForDeleteConfig() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v2.0/{project_id}/clusters/{cluster_id}/lgsconf/delete").
		WithResponse(new(model.DeleteConfigResponse)).
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

func GenReqDefForDeleteTemplate() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v1.0/{project_id}/lgsconf/deletetemplate").
		WithResponse(new(model.DeleteTemplateResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListActions() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/lgsconf/listactions").
		WithResponse(new(model.ListActionsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListCerts() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/certs").
		WithResponse(new(model.ListCertsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Start").
		WithJsonTag("start").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("CertsType").
		WithJsonTag("certsType").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListConfs() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/lgsconf/listconfs").
		WithResponse(new(model.ListConfsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListPipelines() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/lgsconf/listpipelines").
		WithResponse(new(model.ListPipelinesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListTemplates() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/lgsconf/template").
		WithResponse(new(model.ListTemplatesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Type").
		WithJsonTag("type").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowGetConfDetail() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/lgsconf/confdetail").
		WithResponse(new(model.ShowGetConfDetailResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForStartConnectivityTest() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/checkconnection").
		WithResponse(new(model.StartConnectivityTestResponse)).
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

func GenReqDefForStartPipeline() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/lgsconf/start").
		WithResponse(new(model.StartPipelineResponse)).
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

func GenReqDefForStopHotPipeline() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/lgsconf/hot-stop").
		WithResponse(new(model.StopHotPipelineResponse)).
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

func GenReqDefForStopPipeline() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/lgsconf/stop").
		WithResponse(new(model.StopPipelineResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateCnf() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1.0/{project_id}/clusters/{cluster_id}/lgsconf/update").
		WithResponse(new(model.UpdateCnfResponse)).
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
