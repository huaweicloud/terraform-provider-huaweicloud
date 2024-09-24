package v5

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
	"net/http"
)

func GenReqDefForCreateAccessCode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/auth/accesscode").
		WithResponse(new(model.CreateAccessCodeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddQueue() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/amqp-queues").
		WithResponse(new(model.AddQueueResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchShowQueue() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/amqp-queues").
		WithResponse(new(model.BatchShowQueueResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("QueueName").
		WithJsonTag("queue_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteQueue() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/amqp-queues/{queue_id}").
		WithResponse(new(model.DeleteQueueResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("QueueId").
		WithJsonTag("queue_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowQueue() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/amqp-queues/{queue_id}").
		WithResponse(new(model.ShowQueueResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("QueueId").
		WithJsonTag("queue_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddApplication() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/apps").
		WithResponse(new(model.AddApplicationResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteApplication() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/apps/{app_id}").
		WithResponse(new(model.DeleteApplicationResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowApplication() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/apps/{app_id}").
		WithResponse(new(model.ShowApplicationResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowApplications() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/apps").
		WithResponse(new(model.ShowApplicationsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DefaultApp").
		WithJsonTag("default_app").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateApplication() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/apps/{app_id}").
		WithResponse(new(model.UpdateApplicationResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateAsyncCommand() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/async-commands").
		WithResponse(new(model.CreateAsyncCommandResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAsyncDeviceCommand() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/async-commands/{command_id}").
		WithResponse(new(model.ShowAsyncDeviceCommandResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("CommandId").
		WithJsonTag("command_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateRoutingBacklogPolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/routing-rule/backlog-policy").
		WithResponse(new(model.CreateRoutingBacklogPolicyResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteRoutingBacklogPolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}").
		WithResponse(new(model.DeleteRoutingBacklogPolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyId").
		WithJsonTag("policy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListRoutingBacklogPolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/routing-rule/backlog-policy").
		WithResponse(new(model.ListRoutingBacklogPolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyName").
		WithJsonTag("policy_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowRoutingBacklogPolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}").
		WithResponse(new(model.ShowRoutingBacklogPolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyId").
		WithJsonTag("policy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateRoutingBacklogPolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}").
		WithResponse(new(model.UpdateRoutingBacklogPolicyResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyId").
		WithJsonTag("policy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateBatchTask() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/batchtasks").
		WithResponse(new(model.CreateBatchTaskResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteBatchTask() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/batchtasks/{task_id}").
		WithResponse(new(model.DeleteBatchTaskResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TaskId").
		WithJsonTag("task_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListBatchTasks() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/batchtasks").
		WithResponse(new(model.ListBatchTasksResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TaskType").
		WithJsonTag("task_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Status").
		WithJsonTag("status").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForRetryBatchTask() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/batchtasks/{task_id}/retry").
		WithResponse(new(model.RetryBatchTaskResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TaskId").
		WithJsonTag("task_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowBatchTask() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/batchtasks/{task_id}").
		WithResponse(new(model.ShowBatchTaskResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TaskId").
		WithJsonTag("task_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TaskDetailStatus").
		WithJsonTag("task_detail_status").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Target").
		WithJsonTag("target").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForStopBatchTask() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/batchtasks/{task_id}/stop").
		WithResponse(new(model.StopBatchTaskResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TaskId").
		WithJsonTag("task_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteBatchTaskFile() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/batchtask-files/{file_id}").
		WithResponse(new(model.DeleteBatchTaskFileResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FileId").
		WithJsonTag("file_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListBatchTaskFiles() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/batchtask-files").
		WithResponse(new(model.ListBatchTaskFilesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUploadBatchTaskFile() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/batchtask-files").
		WithResponse(new(model.UploadBatchTaskFileResponse)).
		WithContentType("multipart/form-data")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddBridge() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/bridges").
		WithResponse(new(model.AddBridgeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteBridge() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/bridges/{bridge_id}").
		WithResponse(new(model.DeleteBridgeResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("BridgeId").
		WithJsonTag("bridge_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListBridges() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/bridges").
		WithResponse(new(model.ListBridgesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForResetBridgeSecret() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/bridges/{bridge_id}/reset-secret").
		WithResponse(new(model.ResetBridgeSecretResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("BridgeId").
		WithJsonTag("bridge_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBroadcastMessage() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/broadcast-messages").
		WithResponse(new(model.BroadcastMessageResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddCertificate() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/certificates").
		WithResponse(new(model.AddCertificateResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SpAuthToken").
		WithJsonTag("Sp-Auth-Token").
		WithLocationType(def.Header))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("StageAuthToken").
		WithJsonTag("Stage-Auth-Token").
		WithLocationType(def.Header))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCheckCertificate() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/certificates/{certificate_id}/action").
		WithResponse(new(model.CheckCertificateResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("CertificateId").
		WithJsonTag("certificate_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ActionId").
		WithJsonTag("action_id").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SpAuthToken").
		WithJsonTag("Sp-Auth-Token").
		WithLocationType(def.Header))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("StageAuthToken").
		WithJsonTag("Stage-Auth-Token").
		WithLocationType(def.Header))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteCertificate() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/certificates/{certificate_id}").
		WithResponse(new(model.DeleteCertificateResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("CertificateId").
		WithJsonTag("certificate_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SpAuthToken").
		WithJsonTag("Sp-Auth-Token").
		WithLocationType(def.Header))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("StageAuthToken").
		WithJsonTag("Stage-Auth-Token").
		WithLocationType(def.Header))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListCertificates() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/certificates").
		WithResponse(new(model.ListCertificatesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SpAuthToken").
		WithJsonTag("Sp-Auth-Token").
		WithLocationType(def.Header))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("StageAuthToken").
		WithJsonTag("Stage-Auth-Token").
		WithLocationType(def.Header))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateCertificate() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/certificates/{certificate_id}").
		WithResponse(new(model.UpdateCertificateResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("CertificateId").
		WithJsonTag("certificate_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateCommand() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/commands").
		WithResponse(new(model.CreateCommandResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateDeviceAuthorizer() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/device-authorizers").
		WithResponse(new(model.CreateDeviceAuthorizerResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteDeviceAuthorizer() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/device-authorizers/{authorizer_id}").
		WithResponse(new(model.DeleteDeviceAuthorizerResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AuthorizerId").
		WithJsonTag("authorizer_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListDeviceAuthorizers() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/device-authorizers").
		WithResponse(new(model.ListDeviceAuthorizersResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AuthorizerName").
		WithJsonTag("authorizer_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowDeviceAuthorizer() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/device-authorizers/{authorizer_id}").
		WithResponse(new(model.ShowDeviceAuthorizerResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AuthorizerId").
		WithJsonTag("authorizer_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateDeviceAuthorizer() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/device-authorizers/{authorizer_id}").
		WithResponse(new(model.UpdateDeviceAuthorizerResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AuthorizerId").
		WithJsonTag("authorizer_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddDeviceGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/device-group").
		WithResponse(new(model.AddDeviceGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateOrDeleteDeviceInGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/device-group/{group_id}/action").
		WithResponse(new(model.CreateOrDeleteDeviceInGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("GroupId").
		WithJsonTag("group_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ActionId").
		WithJsonTag("action_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteDeviceGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/device-group/{group_id}").
		WithResponse(new(model.DeleteDeviceGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("GroupId").
		WithJsonTag("group_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListDeviceGroups() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/device-group").
		WithResponse(new(model.ListDeviceGroupsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("LastModifiedTime").
		WithJsonTag("last_modified_time").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("GroupType").
		WithJsonTag("group_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowDeviceGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/device-group/{group_id}").
		WithResponse(new(model.ShowDeviceGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("GroupId").
		WithJsonTag("group_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowDevicesInGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/device-group/{group_id}/devices").
		WithResponse(new(model.ShowDevicesInGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("GroupId").
		WithJsonTag("group_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateDeviceGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/device-group/{group_id}").
		WithResponse(new(model.UpdateDeviceGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("GroupId").
		WithJsonTag("group_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddDevice() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/devices").
		WithResponse(new(model.AddDeviceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteDevice() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/devices/{device_id}").
		WithResponse(new(model.DeleteDeviceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForFreezeDevice() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/freeze").
		WithResponse(new(model.FreezeDeviceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListDeviceGroupsByDevice() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/list-device-group").
		WithResponse(new(model.ListDeviceGroupsByDeviceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListDevices() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/devices").
		WithResponse(new(model.ListDevicesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ProductId").
		WithJsonTag("product_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("GatewayId").
		WithJsonTag("gateway_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("IsCascadeQuery").
		WithJsonTag("is_cascade_query").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NodeId").
		WithJsonTag("node_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceName").
		WithJsonTag("device_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("StartTime").
		WithJsonTag("start_time").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("EndTime").
		WithJsonTag("end_time").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForResetDeviceSecret() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/action").
		WithResponse(new(model.ResetDeviceSecretResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ActionId").
		WithJsonTag("action_id").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForResetFingerprint() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/reset-fingerprint").
		WithResponse(new(model.ResetFingerprintResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForSearchDevices() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/search/query-devices").
		WithResponse(new(model.SearchDevicesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowDevice() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/devices/{device_id}").
		WithResponse(new(model.ShowDeviceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUnfreezeDevice() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/unfreeze").
		WithResponse(new(model.UnfreezeDeviceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateDevice() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/devices/{device_id}").
		WithResponse(new(model.UpdateDeviceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateDeviceProxy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/device-proxies").
		WithResponse(new(model.CreateDeviceProxyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteDeviceProxy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/device-proxies/{proxy_id}").
		WithResponse(new(model.DeleteDeviceProxyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ProxyId").
		WithJsonTag("proxy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListDeviceProxies() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/device-proxies").
		WithResponse(new(model.ListDeviceProxiesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ProxyName").
		WithJsonTag("proxy_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowDeviceProxy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/device-proxies/{proxy_id}").
		WithResponse(new(model.ShowDeviceProxyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ProxyId").
		WithJsonTag("proxy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateDeviceProxy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/device-proxies/{proxy_id}").
		WithResponse(new(model.UpdateDeviceProxyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ProxyId").
		WithJsonTag("proxy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowDeviceShadow() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/shadow").
		WithResponse(new(model.ShowDeviceShadowResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateDeviceShadowDesiredData() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/shadow").
		WithResponse(new(model.UpdateDeviceShadowDesiredDataResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateRoutingFlowControlPolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/routing-rule/flowcontrol-policy").
		WithResponse(new(model.CreateRoutingFlowControlPolicyResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteRoutingFlowControlPolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}").
		WithResponse(new(model.DeleteRoutingFlowControlPolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyId").
		WithJsonTag("policy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListRoutingFlowControlPolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/routing-rule/flowcontrol-policy").
		WithResponse(new(model.ListRoutingFlowControlPolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Scope").
		WithJsonTag("scope").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ScopeValue").
		WithJsonTag("scope_value").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyName").
		WithJsonTag("policy_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowRoutingFlowControlPolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}").
		WithResponse(new(model.ShowRoutingFlowControlPolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyId").
		WithJsonTag("policy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateRoutingFlowControlPolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}").
		WithResponse(new(model.UpdateRoutingFlowControlPolicyResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyId").
		WithJsonTag("policy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateMessage() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/messages").
		WithResponse(new(model.CreateMessageResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListDeviceMessages() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/messages").
		WithResponse(new(model.ListDeviceMessagesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowDeviceMessage() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/messages/{message_id}").
		WithResponse(new(model.ShowDeviceMessageResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("MessageId").
		WithJsonTag("message_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateOtaPackage() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/ota-upgrades/packages").
		WithResponse(new(model.CreateOtaPackageResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteOtaPackage() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/ota-upgrades/packages/{package_id}").
		WithResponse(new(model.DeleteOtaPackageResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PackageId").
		WithJsonTag("package_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListOtaPackageInfo() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/ota-upgrades/packages").
		WithResponse(new(model.ListOtaPackageInfoResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PackageType").
		WithJsonTag("package_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ProductId").
		WithJsonTag("product_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Version").
		WithJsonTag("version").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowOtaPackage() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/ota-upgrades/packages/{package_id}").
		WithResponse(new(model.ShowOtaPackageResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PackageId").
		WithJsonTag("package_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBindDevicePolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/device-policies/{policy_id}/bind").
		WithResponse(new(model.BindDevicePolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyId").
		WithJsonTag("policy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateDevicePolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/device-policies").
		WithResponse(new(model.CreateDevicePolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteDevicePolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/device-policies/{policy_id}").
		WithResponse(new(model.DeleteDevicePolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyId").
		WithJsonTag("policy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListDevicePolicies() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/device-policies").
		WithResponse(new(model.ListDevicePoliciesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyName").
		WithJsonTag("policy_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowDevicePolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/device-policies/{policy_id}").
		WithResponse(new(model.ShowDevicePolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyId").
		WithJsonTag("policy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowTargetsInDevicePolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/device-policies/{policy_id}/list-targets").
		WithResponse(new(model.ShowTargetsInDevicePolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyId").
		WithJsonTag("policy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUnbindDevicePolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/device-policies/{policy_id}/unbind").
		WithResponse(new(model.UnbindDevicePolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyId").
		WithJsonTag("policy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateDevicePolicy() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/device-policies/{policy_id}").
		WithResponse(new(model.UpdateDevicePolicyResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PolicyId").
		WithJsonTag("policy_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateProduct() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/products").
		WithResponse(new(model.CreateProductResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteProduct() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/products/{product_id}").
		WithResponse(new(model.DeleteProductResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ProductId").
		WithJsonTag("product_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListProducts() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/products").
		WithResponse(new(model.ListProductsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ProductName").
		WithJsonTag("product_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowProduct() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/products/{product_id}").
		WithResponse(new(model.ShowProductResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ProductId").
		WithJsonTag("product_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateProduct() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/products/{product_id}").
		WithResponse(new(model.UpdateProductResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ProductId").
		WithJsonTag("product_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListProperties() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/properties").
		WithResponse(new(model.ListPropertiesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ServiceId").
		WithJsonTag("service_id").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateProperties() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/devices/{device_id}/properties").
		WithResponse(new(model.UpdatePropertiesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateProvisioningTemplate() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/provisioning-templates").
		WithResponse(new(model.CreateProvisioningTemplateResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteProvisioningTemplate() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/provisioning-templates/{template_id}").
		WithResponse(new(model.DeleteProvisioningTemplateResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TemplateId").
		WithJsonTag("template_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListProvisioningTemplates() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/provisioning-templates").
		WithResponse(new(model.ListProvisioningTemplatesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TemplateName").
		WithJsonTag("template_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowProvisioningTemplate() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/provisioning-templates/{template_id}").
		WithResponse(new(model.ShowProvisioningTemplateResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TemplateId").
		WithJsonTag("template_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateProvisioningTemplate() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/provisioning-templates/{template_id}").
		WithResponse(new(model.UpdateProvisioningTemplateResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TemplateId").
		WithJsonTag("template_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateRoutingRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/routing-rule/rules").
		WithResponse(new(model.CreateRoutingRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateRuleAction() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/routing-rule/actions").
		WithResponse(new(model.CreateRuleActionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteRoutingRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/routing-rule/rules/{rule_id}").
		WithResponse(new(model.DeleteRoutingRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RuleId").
		WithJsonTag("rule_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteRuleAction() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/routing-rule/actions/{action_id}").
		WithResponse(new(model.DeleteRuleActionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ActionId").
		WithJsonTag("action_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListRoutingRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/routing-rule/rules").
		WithResponse(new(model.ListRoutingRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Resource").
		WithJsonTag("resource").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Event").
		WithJsonTag("event").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppType").
		WithJsonTag("app_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RuleName").
		WithJsonTag("rule_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Active").
		WithJsonTag("active").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListRuleActions() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/routing-rule/actions").
		WithResponse(new(model.ListRuleActionsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RuleId").
		WithJsonTag("rule_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Channel").
		WithJsonTag("channel").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppType").
		WithJsonTag("app_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowRoutingRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/routing-rule/rules/{rule_id}").
		WithResponse(new(model.ShowRoutingRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RuleId").
		WithJsonTag("rule_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowRuleAction() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/routing-rule/actions/{action_id}").
		WithResponse(new(model.ShowRuleActionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ActionId").
		WithJsonTag("action_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateRoutingRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/routing-rule/rules/{rule_id}").
		WithResponse(new(model.UpdateRoutingRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RuleId").
		WithJsonTag("rule_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateRuleAction() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/routing-rule/actions/{action_id}").
		WithResponse(new(model.UpdateRuleActionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ActionId").
		WithJsonTag("action_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForChangeRuleStatus() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/rules/{rule_id}/status").
		WithResponse(new(model.ChangeRuleStatusResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RuleId").
		WithJsonTag("rule_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/rules").
		WithResponse(new(model.CreateRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/rules/{rule_id}").
		WithResponse(new(model.DeleteRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RuleId").
		WithJsonTag("rule_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/rules").
		WithResponse(new(model.ListRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppId").
		WithJsonTag("app_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RuleType").
		WithJsonTag("rule_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/rules/{rule_id}").
		WithResponse(new(model.ShowRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RuleId").
		WithJsonTag("rule_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/rules/{rule_id}").
		WithResponse(new(model.UpdateRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RuleId").
		WithJsonTag("rule_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListResourcesByTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/tags/query-resources").
		WithResponse(new(model.ListResourcesByTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForTagDevice() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/tags/bind-resource").
		WithResponse(new(model.TagDeviceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUntagDevice() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/tags/unbind-resource").
		WithResponse(new(model.UntagDeviceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddTunnel() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v5/iot/{project_id}/tunnels").
		WithResponse(new(model.AddTunnelResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCloseDeviceTunnel() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v5/iot/{project_id}/tunnels/{tunnel_id}").
		WithResponse(new(model.CloseDeviceTunnelResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TunnelId").
		WithJsonTag("tunnel_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteDeviceTunnel() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v5/iot/{project_id}/tunnels/{tunnel_id}").
		WithResponse(new(model.DeleteDeviceTunnelResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TunnelId").
		WithJsonTag("tunnel_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListDeviceTunnels() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/tunnels").
		WithResponse(new(model.ListDeviceTunnelsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DeviceId").
		WithJsonTag("device_id").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowDeviceTunnel() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v5/iot/{project_id}/tunnels/{tunnel_id}").
		WithResponse(new(model.ShowDeviceTunnelResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TunnelId").
		WithJsonTag("tunnel_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("InstanceId").
		WithJsonTag("Instance-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}
