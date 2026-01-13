---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_types"
description: |-
  Use this data source to get the list of HSS event types within HuaweiCloud.
---

# huaweicloud_hss_event_types

Use this data source to get the list of HSS event types within HuaweiCloud.

## Example Usage

```hcl
variable "category" {}

data "huaweicloud_hss_event_types" "test" {
  category = var.category
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `category` - (Required, String) Specifies the event category. Valid values are:
  + **host**: Host security event.
  + **container**: Container security event.
  + **serverless**: Serverless scenario security event.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `begin_time` - (Optional, Int) Specifies the start time, 13-digit timestamp. Must be less than or equal to `end_time`.
  If `end_time` is not passed, the current time will be queried by default.

* `end_time` - (Optional, Int) Specifies the end time, 13-digit timestamp. Must be greater than or equal to
  `begin_time`. If `begin_time` is not passed, the query will start from timestamp 0 by default.

* `last_days` - (Optional, Int) Specifies the number of days for query time range.
  Mutually exclusive with custom query time `begin_time` and `end_time`.

* `host_name` - (Optional, String) Specifies the server name.

* `host_id` - (Optional, String) Specifies the server ID.

* `private_ip` - (Optional, String) Specifies the server private IP.

* `public_ip` - (Optional, String) Specifies the server public IP.

* `container_name` - (Optional, String) Specifies the container instance name.

* `handle_status` - (Optional, String) Specifies the handle status. Valid values are:
  + **unhandled**: Not handled.
  + **handled**: Handled.

* `severity` - (Optional, String) Specifies the threat level. Valid values are:
  + **Security**: Security.
  + **Low**: Low risk.
  + **Medium**: Medium risk.
  + **High**: High risk.
  + **Critical**: Critical.

* `severity_list` - (Optional, List) Specifies the threat level list. Valid values are the same as `severity`.

* `attack_tag` - (Optional, String) Specifies the attack identifier. Valid values are:
  + **attack_success**: Attack succeeded.
  + **attack_attempt**: Attack attempted.
  + **attack_blocked**: Attack blocked.
  + **abnormal_behavior**: Abnormal behavior.
  + **collapsible_host**: Host compromised.
  + **system_vulnerability**: System vulnerability.

* `asset_value` - (Optional, String) Specifies the asset importance. Valid values are:
  + **important**: Important asset.
  + **common**: Common asset.
  + **test**: Test asset.

* `tag_list` - (Optional, List) Specifies the event tag list, for example: `["热点事件"]`.

* `att_ck` - (Optional, String) Specifies the ATT&CK attack level. Valid values are:
  + **Reconnaissance**
  + **Initial Access**
  + **Execution**
  + **Persistence**
  + **Privilege Escalation**
  + **Defense Evasion**
  + **Credential Access**
  + **Command and Control**
  + **Impact**

* `event_name` - (Optional, String) Specifies the alarm name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `data_list` - The list of event type details.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `event_type_name` - The names corresponding to major categories. For its values,
  see [API docs](https://support.huaweicloud.com/api-hss2.0/ListEventType.html#ListEventType__response_EventTypeResponseInfo).

* `event_type_num` - The total number of event categories.

* `event_type_list` - The list of names and quantities corresponding to subcategories.

  The [event_type_list](#event_type_list_struct) structure is documented below.

<a name="event_type_list_struct"></a>
The `event_type_list` block supports:

* `event_type` - The event type. For its values,
  see [API docs](https://support.huaweicloud.com/api-hss2.0/ListEventType.html#ListEventType__response_EventTypeResponseInfo).

* `event_num` - The number of events.

* `status` - The status. Valid values are:
  + **locked**
  + **unlocked**
