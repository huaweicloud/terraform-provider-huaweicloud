package v3

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
	"net/http"
)

func GenReqDefForAddSecurityGroups() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/ports/{port_id}/insert-security-groups").
		WithResponse(new(model.AddSecurityGroupsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PortId").
		WithJsonTag("port_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddSourcesToTrafficMirrorSession() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-sessions/{traffic_mirror_session_id}/add-sources").
		WithResponse(new(model.AddSourcesToTrafficMirrorSessionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorSessionId").
		WithJsonTag("traffic_mirror_session_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchCreatePortTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/ports/{port_id}/tags/create").
		WithResponse(new(model.BatchCreatePortTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PortId").
		WithJsonTag("port_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchCreateSecurityGroupRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/vpc/security-groups/{security_group_id}/security-group-rules/batch-create").
		WithResponse(new(model.BatchCreateSecurityGroupRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SecurityGroupId").
		WithJsonTag("security_group_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchCreateSubNetworkInterface() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/vpc/sub-network-interfaces/batch-create").
		WithResponse(new(model.BatchCreateSubNetworkInterfaceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchDeletePortTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/ports/{port_id}/tags/delete").
		WithResponse(new(model.BatchDeletePortTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PortId").
		WithJsonTag("port_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCountPortsByTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/ports/resource-instances/count").
		WithResponse(new(model.CountPortsByTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreatePortTag() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/ports/{port_id}/tags").
		WithResponse(new(model.CreatePortTagResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PortId").
		WithJsonTag("port_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateSecurityGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/vpc/security-groups").
		WithResponse(new(model.CreateSecurityGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateSecurityGroupRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/vpc/security-group-rules").
		WithResponse(new(model.CreateSecurityGroupRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateSubNetworkInterface() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/vpc/sub-network-interfaces").
		WithResponse(new(model.CreateSubNetworkInterfaceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateTrafficMirrorFilter() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-filters").
		WithResponse(new(model.CreateTrafficMirrorFilterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateTrafficMirrorFilterRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-filter-rules").
		WithResponse(new(model.CreateTrafficMirrorFilterRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateTrafficMirrorSession() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-sessions").
		WithResponse(new(model.CreateTrafficMirrorSessionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeletePortTag() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/ports/{port_id}/tags/{tag_key}").
		WithResponse(new(model.DeletePortTagResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PortId").
		WithJsonTag("port_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TagKey").
		WithJsonTag("tag_key").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteSecurityGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/vpc/security-groups/{security_group_id}").
		WithResponse(new(model.DeleteSecurityGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SecurityGroupId").
		WithJsonTag("security_group_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteSecurityGroupRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/vpc/security-group-rules/{security_group_rule_id}").
		WithResponse(new(model.DeleteSecurityGroupRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SecurityGroupRuleId").
		WithJsonTag("security_group_rule_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteSubNetworkInterface() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/vpc/sub-network-interfaces/{sub_network_interface_id}").
		WithResponse(new(model.DeleteSubNetworkInterfaceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SubNetworkInterfaceId").
		WithJsonTag("sub_network_interface_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteTrafficMirrorFilter() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-filters/{traffic_mirror_filter_id}").
		WithResponse(new(model.DeleteTrafficMirrorFilterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorFilterId").
		WithJsonTag("traffic_mirror_filter_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteTrafficMirrorFilterRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-filter-rules/{traffic_mirror_filter_rule_id}").
		WithResponse(new(model.DeleteTrafficMirrorFilterRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorFilterRuleId").
		WithJsonTag("traffic_mirror_filter_rule_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteTrafficMirrorSession() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-sessions/{traffic_mirror_session_id}").
		WithResponse(new(model.DeleteTrafficMirrorSessionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorSessionId").
		WithJsonTag("traffic_mirror_session_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListPortTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/ports/tags").
		WithResponse(new(model.ListPortTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListPortsByTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/ports/resource-instances/filter").
		WithResponse(new(model.ListPortsByTagsResponse)).
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

func GenReqDefForListSecurityGroupRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/security-group-rules").
		WithResponse(new(model.ListSecurityGroupRulesResponse)).
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
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SecurityGroupId").
		WithJsonTag("security_group_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Protocol").
		WithJsonTag("protocol").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Description").
		WithJsonTag("description").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RemoteGroupId").
		WithJsonTag("remote_group_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Direction").
		WithJsonTag("direction").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Action").
		WithJsonTag("action").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("RemoteIpPrefix").
		WithJsonTag("remote_ip_prefix").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListSecurityGroups() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/security-groups").
		WithResponse(new(model.ListSecurityGroupsResponse)).
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
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Description").
		WithJsonTag("description").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("EnterpriseProjectId").
		WithJsonTag("enterprise_project_id").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListSubNetworkInterfaces() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/sub-network-interfaces").
		WithResponse(new(model.ListSubNetworkInterfacesResponse)).
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
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("VirsubnetId").
		WithJsonTag("virsubnet_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PrivateIpAddress").
		WithJsonTag("private_ip_address").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("MacAddress").
		WithJsonTag("mac_address").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("VpcId").
		WithJsonTag("vpc_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Description").
		WithJsonTag("description").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ParentId").
		WithJsonTag("parent_id").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListTrafficMirrorFilterRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-filter-rules").
		WithResponse(new(model.ListTrafficMirrorFilterRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Description").
		WithJsonTag("description").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorFilterId").
		WithJsonTag("traffic_mirror_filter_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Direction").
		WithJsonTag("direction").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Protocol").
		WithJsonTag("protocol").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SourceCidrBlock").
		WithJsonTag("source_cidr_block").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DestinationCidrBlock").
		WithJsonTag("destination_cidr_block").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SourcePortRange").
		WithJsonTag("source_port_range").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("DestinationPortRange").
		WithJsonTag("destination_port_range").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Action").
		WithJsonTag("action").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Priority").
		WithJsonTag("priority").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListTrafficMirrorFilters() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-filters").
		WithResponse(new(model.ListTrafficMirrorFiltersResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Description").
		WithJsonTag("description").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("CreatedAt").
		WithJsonTag("created_at").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("UpdatedAt").
		WithJsonTag("updated_at").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListTrafficMirrorSessions() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-sessions").
		WithResponse(new(model.ListTrafficMirrorSessionsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Description").
		WithJsonTag("description").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorFilterId").
		WithJsonTag("traffic_mirror_filter_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorTargetId").
		WithJsonTag("traffic_mirror_target_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorTargetType").
		WithJsonTag("traffic_mirror_target_type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("VirtualNetworkId").
		WithJsonTag("virtual_network_id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PacketLength").
		WithJsonTag("packet_length").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Priority").
		WithJsonTag("priority").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Enabled").
		WithJsonTag("enabled").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Type").
		WithJsonTag("type").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("CreatedAt").
		WithJsonTag("created_at").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("UpdatedAt").
		WithJsonTag("updated_at").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Marker").
		WithJsonTag("marker").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForMigrateSubNetworkInterface() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/sub-network-interfaces/migrate").
		WithResponse(new(model.MigrateSubNetworkInterfaceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForRemoveSecurityGroups() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/ports/{port_id}/remove-security-groups").
		WithResponse(new(model.RemoveSecurityGroupsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PortId").
		WithJsonTag("port_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForRemoveSourcesFromTrafficMirrorSession() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-sessions/{traffic_mirror_session_id}/remove-sources").
		WithResponse(new(model.RemoveSourcesFromTrafficMirrorSessionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorSessionId").
		WithJsonTag("traffic_mirror_session_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowPortTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/ports/{port_id}/tags").
		WithResponse(new(model.ShowPortTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("PortId").
		WithJsonTag("port_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowSecurityGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/security-groups/{security_group_id}").
		WithResponse(new(model.ShowSecurityGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SecurityGroupId").
		WithJsonTag("security_group_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowSecurityGroupRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/security-group-rules/{security_group_rule_id}").
		WithResponse(new(model.ShowSecurityGroupRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SecurityGroupRuleId").
		WithJsonTag("security_group_rule_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowSubNetworkInterface() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/sub-network-interfaces/{sub_network_interface_id}").
		WithResponse(new(model.ShowSubNetworkInterfaceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SubNetworkInterfaceId").
		WithJsonTag("sub_network_interface_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowSubNetworkInterfacesQuantity() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/sub-network-interfaces/count").
		WithResponse(new(model.ShowSubNetworkInterfacesQuantityResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowTrafficMirrorFilter() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-filters/{traffic_mirror_filter_id}").
		WithResponse(new(model.ShowTrafficMirrorFilterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorFilterId").
		WithJsonTag("traffic_mirror_filter_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowTrafficMirrorFilterRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-filter-rules/{traffic_mirror_filter_rule_id}").
		WithResponse(new(model.ShowTrafficMirrorFilterRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorFilterRuleId").
		WithJsonTag("traffic_mirror_filter_rule_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowTrafficMirrorSession() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-sessions/{traffic_mirror_session_id}").
		WithResponse(new(model.ShowTrafficMirrorSessionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorSessionId").
		WithJsonTag("traffic_mirror_session_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateSecurityGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/security-groups/{security_group_id}").
		WithResponse(new(model.UpdateSecurityGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SecurityGroupId").
		WithJsonTag("security_group_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateSubNetworkInterface() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/sub-network-interfaces/{sub_network_interface_id}").
		WithResponse(new(model.UpdateSubNetworkInterfaceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("SubNetworkInterfaceId").
		WithJsonTag("sub_network_interface_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateTrafficMirrorFilter() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-filters/{traffic_mirror_filter_id}").
		WithResponse(new(model.UpdateTrafficMirrorFilterResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorFilterId").
		WithJsonTag("traffic_mirror_filter_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateTrafficMirrorFilterRule() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-filter-rules/{traffic_mirror_filter_rule_id}").
		WithResponse(new(model.UpdateTrafficMirrorFilterRuleResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorFilterRuleId").
		WithJsonTag("traffic_mirror_filter_rule_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateTrafficMirrorSession() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/traffic-mirror-sessions/{traffic_mirror_session_id}").
		WithResponse(new(model.UpdateTrafficMirrorSessionResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TrafficMirrorSessionId").
		WithJsonTag("traffic_mirror_session_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddFirewallRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/firewalls/{firewall_id}/insert-rules").
		WithResponse(new(model.AddFirewallRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAssociateSubnetFirewall() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/firewalls/{firewall_id}/associate-subnets").
		WithResponse(new(model.AssociateSubnetFirewallResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchCreateFirewallTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/firewalls/{firewall_id}/tags/create").
		WithResponse(new(model.BatchCreateFirewallTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchDeleteFirewallTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/firewalls/{firewall_id}/tags/delete").
		WithResponse(new(model.BatchDeleteFirewallTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCountFirewallsByTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/firewalls/resource-instances/count").
		WithResponse(new(model.CountFirewallsByTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateFirewall() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/vpc/firewalls").
		WithResponse(new(model.CreateFirewallResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateFirewallTag() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/firewalls/{firewall_id}/tags").
		WithResponse(new(model.CreateFirewallTagResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteFirewall() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/vpc/firewalls/{firewall_id}").
		WithResponse(new(model.DeleteFirewallResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteFirewallTag() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/firewalls/{firewall_id}/tags/{tag_key}").
		WithResponse(new(model.DeleteFirewallTagResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TagKey").
		WithJsonTag("tag_key").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDisassociateSubnetFirewall() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/firewalls/{firewall_id}/disassociate-subnets").
		WithResponse(new(model.DisassociateSubnetFirewallResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListFirewall() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/firewalls").
		WithResponse(new(model.ListFirewallResponse)).
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
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Status").
		WithJsonTag("status").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AdminStateUp").
		WithJsonTag("admin_state_up").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("EnterpriseProjectId").
		WithJsonTag("enterprise_project_id").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListFirewallTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/firewalls/tags").
		WithResponse(new(model.ListFirewallTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Limit").
		WithJsonTag("limit").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Offset").
		WithJsonTag("offset").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListFirewallsByTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/firewalls/resource-instances/filter").
		WithResponse(new(model.ListFirewallsByTagsResponse)).
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

func GenReqDefForRemoveFirewallRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/firewalls/{firewall_id}/remove-rules").
		WithResponse(new(model.RemoveFirewallRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowFirewall() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/firewalls/{firewall_id}").
		WithResponse(new(model.ShowFirewallResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowFirewallTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/firewalls/{firewall_id}/tags").
		WithResponse(new(model.ShowFirewallTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateFirewall() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/firewalls/{firewall_id}").
		WithResponse(new(model.UpdateFirewallResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateFirewallRules() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/firewalls/{firewall_id}/update-rules").
		WithResponse(new(model.UpdateFirewallRulesResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("FirewallId").
		WithJsonTag("firewall_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddClouddcnSubnetsTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/clouddcn-subnets/{resource_id}/tags").
		WithResponse(new(model.AddClouddcnSubnetsTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceId").
		WithJsonTag("resource_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchCreateClouddcnSubnetsTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/clouddcn-subnets/{resource_id}/tags/create").
		WithResponse(new(model.BatchCreateClouddcnSubnetsTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceId").
		WithJsonTag("resource_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForBatchDeleteClouddcnSubnetsTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/clouddcn-subnets/{resource_id}/tags/delete").
		WithResponse(new(model.BatchDeleteClouddcnSubnetsTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceId").
		WithJsonTag("resource_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateClouddcnSubnet() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/vpc/clouddcn-subnets").
		WithResponse(new(model.CreateClouddcnSubnetResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteClouddcnSubnet() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/vpc/clouddcn-subnets/{clouddcn_subnet_id}").
		WithResponse(new(model.DeleteClouddcnSubnetResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClouddcnSubnetId").
		WithJsonTag("clouddcn_subnet_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteClouddcnSubnetsTag() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/clouddcn-subnets/{resource_id}/tags/{tag_key}").
		WithResponse(new(model.DeleteClouddcnSubnetsTagResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceId").
		WithJsonTag("resource_id").
		WithLocationType(def.Path))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("TagKey").
		WithJsonTag("tag_key").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListClouddcnSubnets() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/clouddcn-subnets").
		WithResponse(new(model.ListClouddcnSubnetsResponse)).
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
		WithName("VpcId").
		WithJsonTag("vpc_id").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListClouddcnSubnetsCountFilterTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/clouddcn-subnets/resource-instances/count").
		WithResponse(new(model.ListClouddcnSubnetsCountFilterTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListClouddcnSubnetsFilterTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/clouddcn-subnets/resource-instances/filter").
		WithResponse(new(model.ListClouddcnSubnetsFilterTagsResponse)).
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

func GenReqDefForListClouddcnSubnetsTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/clouddcn-subnets/tags").
		WithResponse(new(model.ListClouddcnSubnetsTagsResponse)).
		WithContentType("application/json")

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowClouddcnSubnet() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/clouddcn-subnets/{clouddcn_subnet_id}").
		WithResponse(new(model.ShowClouddcnSubnetResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClouddcnSubnetId").
		WithJsonTag("clouddcn_subnet_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowClouddcnSubnetsTags() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/clouddcn-subnets/{resource_id}/tags").
		WithResponse(new(model.ShowClouddcnSubnetsTagsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ResourceId").
		WithJsonTag("resource_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateClouddcnSubnet() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/clouddcn-subnets/{clouddcn_subnet_id}").
		WithResponse(new(model.UpdateClouddcnSubnetResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("ClouddcnSubnetId").
		WithJsonTag("clouddcn_subnet_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateAddressGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/vpc/address-groups").
		WithResponse(new(model.CreateAddressGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteAddressGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/vpc/address-groups/{address_group_id}").
		WithResponse(new(model.DeleteAddressGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AddressGroupId").
		WithJsonTag("address_group_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteIpAddressGroupForce() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/vpc/address-groups/{address_group_id}/force").
		WithResponse(new(model.DeleteIpAddressGroupForceResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AddressGroupId").
		WithJsonTag("address_group_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListAddressGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/address-groups").
		WithResponse(new(model.ListAddressGroupResponse)).
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
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("IpVersion").
		WithJsonTag("ip_version").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Description").
		WithJsonTag("description").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("EnterpriseProjectId").
		WithJsonTag("enterprise_project_id").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowAddressGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/address-groups/{address_group_id}").
		WithResponse(new(model.ShowAddressGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AddressGroupId").
		WithJsonTag("address_group_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateAddressGroup() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/address-groups/{address_group_id}").
		WithResponse(new(model.UpdateAddressGroupResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("AddressGroupId").
		WithJsonTag("address_group_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForAddVpcExtendCidr() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/vpcs/{vpc_id}/add-extend-cidr").
		WithResponse(new(model.AddVpcExtendCidrResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("VpcId").
		WithJsonTag("vpc_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForCreateVpc() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v3/{project_id}/vpc/vpcs").
		WithResponse(new(model.CreateVpcResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForDeleteVpc() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodDelete).
		WithPath("/v3/{project_id}/vpc/vpcs/{vpc_id}").
		WithResponse(new(model.DeleteVpcResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("VpcId").
		WithJsonTag("vpc_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForListVpcs() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/vpcs").
		WithResponse(new(model.ListVpcsResponse)).
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
		WithName("Id").
		WithJsonTag("id").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Name").
		WithJsonTag("name").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Description").
		WithJsonTag("description").
		WithLocationType(def.Query))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Cidr").
		WithJsonTag("cidr").
		WithLocationType(def.Query))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForRemoveVpcExtendCidr() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/vpcs/{vpc_id}/remove-extend-cidr").
		WithResponse(new(model.RemoveVpcExtendCidrResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("VpcId").
		WithJsonTag("vpc_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForShowVpc() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{project_id}/vpc/vpcs/{vpc_id}").
		WithResponse(new(model.ShowVpcResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("VpcId").
		WithJsonTag("vpc_id").
		WithLocationType(def.Path))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func GenReqDefForUpdateVpc() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPut).
		WithPath("/v3/{project_id}/vpc/vpcs/{vpc_id}").
		WithResponse(new(model.UpdateVpcResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("VpcId").
		WithJsonTag("vpc_id").
		WithLocationType(def.Path))

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}
