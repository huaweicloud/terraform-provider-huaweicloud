package v2

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/aom/v2/model"
	"net/http"
)

func GenReqDefForAddActionRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v2/{project_id}/alert/action-rules").
		WithResponse(new(model.AddActionRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddAlarmRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v2/{project_id}/alarm-rules").
		WithResponse(new(model.AddAlarmRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddEvent2alarmRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v2/{project_id}/event2alarm-rule").
		WithResponse(new(model.AddEvent2alarmRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddMetricData() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1/{project_id}/ams/report/metricdata").
		WithResponse(new(model.AddMetricDataResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddMuteRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v2/{project_id}/alert/mute-rules").
		WithResponse(new(model.AddMuteRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddOrUpdateMetricOrEventAlarmRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v4/{project_id}/alarm-rules").
		WithResponse(new(model.AddOrUpdateMetricOrEventAlarmRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ActionId").
		WithJsonTag("action_id").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("EnterpriseProjectId").
		WithJsonTag("Enterprise-Project-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddOrUpdateServiceDiscoveryRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v1/{project_id}/inv/servicediscoveryrules").
		WithResponse(new(model.AddOrUpdateServiceDiscoveryRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCountEvents() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v2/{project_id}/events/statistic").
		WithResponse(new(model.CountEventsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Type").
		WithJsonTag("type").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteActionRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v2/{project_id}/alert/action-rules").
		WithResponse(new(model.DeleteActionRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteAlarmRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v2/{project_id}/alarm-rules/{alarm_rule_id}").
		WithResponse(new(model.DeleteAlarmRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AlarmRuleId").
		WithJsonTag("alarm_rule_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteAlarmRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v2/{project_id}/alarm-rules/delete").
		WithResponse(new(model.DeleteAlarmRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteEvent2alarmRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v2/{project_id}/event2alarm-rule").
		WithResponse(new(model.DeleteEvent2alarmRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteMetricOrEventAlarmRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v4/{project_id}/alarm-rules").
		WithResponse(new(model.DeleteMetricOrEventAlarmRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteMuteRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v2/{project_id}/alert/mute-rules").
		WithResponse(new(model.DeleteMuteRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteserviceDiscoveryRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v1/{project_id}/inv/servicediscoveryrules").
		WithResponse(new(model.DeleteserviceDiscoveryRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AppRulesIds").
		WithJsonTag("appRulesIds").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListActionRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v2/{project_id}/alert/action-rules").
		WithResponse(new(model.ListActionRuleResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAgents() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1/{project_id}/{cluster_id}/{namespace}/agents").
		WithResponse(new(model.ListAgentsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClusterId").
		WithJsonTag("cluster_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Namespace").
		WithJsonTag("namespace").
		WithLocationType(def.Path))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAlarmRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v2/{project_id}/alarm-rules").
		WithResponse(new(model.ListAlarmRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListEvent2alarmRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v2/{project_id}/event2alarm-rule").
		WithResponse(new(model.ListEvent2alarmRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListEvents() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v2/{project_id}/events").
		WithResponse(new(model.ListEventsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Type").
		WithJsonTag("type").
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
		WithName("EnterpriseProjectId").
		WithJsonTag("Enterprise-Project-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListLogItems() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1/{project_id}/als/action").
		WithResponse(new(model.ListLogItemsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Type").
		WithJsonTag("type").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListMetricItems() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1/{project_id}/ams/metrics").
		WithResponse(new(model.ListMetricItemsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Type").
		WithJsonTag("type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Start").
		WithJsonTag("start").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListMetricOrEventAlarmRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v4/{project_id}/alarm-rules").
		WithResponse(new(model.ListMetricOrEventAlarmRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SortBy").
		WithJsonTag("sort_by").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("EventSource").
		WithJsonTag("event_source").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("EventSeverity").
		WithJsonTag("event_severity").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AlarmRuleStatus").
		WithJsonTag("alarm_rule_status").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AlarmRuleType").
		WithJsonTag("alarm_rule_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PromInstanceId").
		WithJsonTag("prom_instance_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("BindNotificationRuleId").
		WithJsonTag("bind_notification_rule_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RelatedCceClusters").
		WithJsonTag("related_cce_clusters").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("EnterpriseProjectId").
		WithJsonTag("Enterprise-Project-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListMuteRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v2/{project_id}/alert/mute-rules").
		WithResponse(new(model.ListMuteRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListNotifiedHistories() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v2/{project_id}/alarm-notified-histories").
		WithResponse(new(model.ListNotifiedHistoriesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("EventSn").
		WithJsonTag("event_sn").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListPermissions() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1/{project_id}/aom/auth/grant").
		WithResponse(new(model.ListPermissionsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListSample() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v2/{project_id}/samples").
		WithResponse(new(model.ListSampleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FillValue").
		WithJsonTag("fill_value").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListSeries() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v2/{project_id}/series").
		WithResponse(new(model.ListSeriesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListServiceDiscoveryRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1/{project_id}/inv/servicediscoveryrules").
		WithResponse(new(model.ListServiceDiscoveryRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForPushEvents() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v2/{project_id}/push/events").
		WithResponse(new(model.PushEventsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Action").
		WithJsonTag("action").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("XEnterprisePrjectId").
		WithJsonTag("x-enterprise-prject-id").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowActionRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v2/{project_id}/alert/action-rules/{rule_name}").
		WithResponse(new(model.ShowActionRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RuleName").
		WithJsonTag("rule_name").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAlarmRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v2/{project_id}/alarm-rules/{alarm_rule_id}").
		WithResponse(new(model.ShowAlarmRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AlarmRuleId").
		WithJsonTag("alarm_rule_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowMetricsData() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1/{project_id}/ams/metricdata").
		WithResponse(new(model.ShowMetricsDataResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FillValue").
		WithJsonTag("fillValue").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateActionRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v2/{project_id}/alert/action-rules").
		WithResponse(new(model.UpdateActionRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateAlarmRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v2/{project_id}/alarm-rules").
		WithResponse(new(model.UpdateAlarmRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateEventRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v2/{project_id}/event2alarm-rule").
		WithResponse(new(model.UpdateEventRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateMuteRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v2/{project_id}/alert/mute-rules").
		WithResponse(new(model.UpdateMuteRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreatePromInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1/{project_id}/aom/prometheus").
		WithResponse(new(model.CreatePromInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Region").
		WithJsonTag("region").
		WithLocationType(def.Header))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateRecordingRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1/{project_id}/{prometheus_instance}/aom/api/v1/rules").
		WithResponse(new(model.CreateRecordingRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PrometheusInstance").
		WithJsonTag("prometheus_instance").
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

func GenReqDefForDeletePromInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v1/{project_id}/aom/prometheus").
		WithResponse(new(model.DeletePromInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PromId").
		WithJsonTag("prom_id").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("EnterpriseProjectId").
		WithJsonTag("Enterprise-Project-Id").
		WithLocationType(def.Header))

	reqDefBuilder.WithResponseField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAccessCode() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1/{project_id}/access-code").
		WithResponse(new(model.ListAccessCodeResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListInstantQueryAomPromGet() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1/{project_id}/aom/api/v1/query").
		WithResponse(new(model.ListInstantQueryAomPromGetResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Query").
		WithJsonTag("query").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Time").
		WithJsonTag("time").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListInstantQueryAomPromPost() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1/{project_id}/aom/api/v1/query").
		WithResponse(new(model.ListInstantQueryAomPromPostResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Query").
		WithJsonTag("query").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Time").
		WithJsonTag("time").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListLabelValuesAomPromGet() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1/{project_id}/aom/api/v1/label/{label_name}/values").
		WithResponse(new(model.ListLabelValuesAomPromGetResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("LabelName").
		WithJsonTag("label_name").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListLabelsAomPromGet() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1/{project_id}/aom/api/v1/labels").
		WithResponse(new(model.ListLabelsAomPromGetResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListLabelsAomPromPost() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1/{project_id}/aom/api/v1/labels").
		WithResponse(new(model.ListLabelsAomPromPostResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListMetadataAomPromGet() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1/{project_id}/aom/api/v1/metadata").
		WithResponse(new(model.ListMetadataAomPromGetResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListPromInstance() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1/{project_id}/aom/prometheus").
		WithResponse(new(model.ListPromInstanceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PromId").
		WithJsonTag("prom_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PromType").
		WithJsonTag("prom_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("CceClusterEnable").
		WithJsonTag("cce_cluster_enable").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PromStatus").
		WithJsonTag("prom_status").
		WithLocationType(def.Query))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("EnterpriseProjectId").
		WithJsonTag("Enterprise-Project-Id").
		WithLocationType(def.Header))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListRangeQueryAomPromGet() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v1/{project_id}/aom/api/v1/query_range").
		WithResponse(new(model.ListRangeQueryAomPromGetResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Query").
		WithJsonTag("query").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Start").
		WithJsonTag("start").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("End").
		WithJsonTag("end").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Step").
		WithJsonTag("step").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListRangeQueryAomPromPost() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1/{project_id}/aom/api/v1/query_range").
		WithResponse(new(model.ListRangeQueryAomPromPostResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Query").
		WithJsonTag("query").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Start").
		WithJsonTag("start").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("End").
		WithJsonTag("end").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Step").
		WithJsonTag("step").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}
