package v3

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"
	"net/http"
)

func GenReqDefForBatchCreateResourceTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/{resource_type}/{resource_id}/tags/create").
		WithResponse(new(model.BatchCreateResourceTagsResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceId").
		WithJsonTag("resource_id").
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

func GenReqDefForBatchDeleteResourceTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/{resource_type}/{resource_id}/tags/delete").
		WithResponse(new(model.BatchDeleteResourceTagsResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceId").
		WithJsonTag("resource_id").
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

func GenReqDefForCheckObsBuckets() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{domain_id}/checkbucket").
		WithResponse(new(model.CheckObsBucketsResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DomainId").
		WithJsonTag("domain_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateNotification() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/notifications").
		WithResponse(new(model.CreateNotificationResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateTracker() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/tracker").
		WithResponse(new(model.CreateTrackerResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteNotification() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/notifications").
		WithResponse(new(model.DeleteNotificationResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NotificationId").
		WithJsonTag("notification_id").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteTracker() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/trackers").
		WithResponse(new(model.DeleteTrackerResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrackerName").
		WithJsonTag("tracker_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrackerType").
		WithJsonTag("tracker_type").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListNotifications() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/notifications/{notification_type}").
		WithResponse(new(model.ListNotificationsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NotificationType").
		WithJsonTag("notification_type").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("NotificationName").
		WithJsonTag("notification_name").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListOperations() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/operations").
		WithResponse(new(model.ListOperationsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ServiceType").
		WithJsonTag("service_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceType").
		WithJsonTag("resource_type").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListQuotas() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/quotas").
		WithResponse(new(model.ListQuotasResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListTraceResources() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{domain_id}/resources").
		WithResponse(new(model.ListTraceResourcesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DomainId").
		WithJsonTag("domain_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListTraces() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/traces").
		WithResponse(new(model.ListTracesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TraceType").
		WithJsonTag("trace_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("From").
		WithJsonTag("from").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Next").
		WithJsonTag("next").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("To").
		WithJsonTag("to").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrackerName").
		WithJsonTag("tracker_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ServiceType").
		WithJsonTag("service_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("User").
		WithJsonTag("user").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceId").
		WithJsonTag("resource_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceName").
		WithJsonTag("resource_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceType").
		WithJsonTag("resource_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TraceId").
		WithJsonTag("trace_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TraceName").
		WithJsonTag("trace_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TraceRating").
		WithJsonTag("trace_rating").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AccessKeyId").
		WithJsonTag("access_key_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("EnterpriseProjectId").
		WithJsonTag("enterprise_project_id").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListTrackers() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/trackers").
		WithResponse(new(model.ListTrackersResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrackerName").
		WithJsonTag("tracker_name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrackerType").
		WithJsonTag("tracker_type").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListUserResources() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/user-resources").
		WithResponse(new(model.ListUserResourcesResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateNotification() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/notifications").
		WithResponse(new(model.UpdateNotificationResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateTracker() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/tracker").
		WithResponse(new(model.UpdateTrackerResponse)).
		WithContentType("application/json;charset=UTF-8")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}
