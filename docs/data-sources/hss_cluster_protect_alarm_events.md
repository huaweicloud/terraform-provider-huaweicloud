---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_cluster_protect_alarm_events"
description: |-
  Use this data source to get the list of HSS cluster protect all alarm events within HuaweiCloud.
---

# huaweicloud_hss_cluster_protect_alarm_events

Use this data source to get the list of HSS cluster protect all alarm events within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_cluster_protect_alarm_events" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `last_update_time` - The last update time.

* `data_list` - The list of cluster protect alarm events.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `action` - The blocking action.

* `cluster_name` - The cluster name.

* `cluster_id` - The cluster ID.

* `event_name` - The event name.

* `event_class_id` - The event unique identifier.

* `event_id` - The event ID.

* `event_type` - The event type.

* `event_content` - The event content.

* `handle_status` - The handling status.  
  The valid values are as follows:
  + **unhandled**: Unhandled.
  + **handled**: Handled.

* `create_time` - The creation time.

* `update_time` - The update time.

* `project_id` - The project ID.

* `enterprise_project_id` - The enterprise project ID.

* `policy_name` - The policy name.

* `policy_id` - The policy ID.

* `resource_info` - The event resource information.

  The [resource_info](#resource_info_struct) structure is documented below.

<a name="resource_info_struct"></a>
The `resource_info` block supports:

* `enforcement_action` - The enforcement action.

* `group` - The group.

* `message` - The message.

* `name` - The name.

* `namespace` - The namespace.

* `version` - The version.

* `kind` - The resource type.

* `resource_name` - The resource name.
