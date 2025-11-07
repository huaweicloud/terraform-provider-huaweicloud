---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_rasp_events"
description: |-
  Use this data source to get the list of application protection events.
---

# huaweicloud_hss_rasp_events

Use this data source to get the list of application protection events.

## Example Usage

```hcl
variable "host_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_hss_rasp_events" "test" {
  host_id    = var.host_id
  start_time = var.start_time
  end_time   = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Required, String) Specifies the host ID.

* `start_time` - (Required, Int) Specifies the query start time.
  The format is 13-digit timestamp in millisecond.

* `end_time` - (Required, Int) Specifies the query end time.
  The format is 13-digit timestamp in millisecond.

* `app_type` - (Optional, String) Specifies the application type.
  The value can be **java**.

* `severity` - (Optional, String) Specifies the alarm severity.
  The valid values are as follows:
  + **0**: Info.
  + **1**: Low level alarm.
  + **2**: Medium level alarm.
  + **3**: High level alarm.
  + **4**: Critical.

* `attack_tag` - (Optional, String) Specifies the attack tag.
  The valid values are as follows:
  + **Attack Success**
  + **Attack Attempt**
  + **Attack Blocked**
  + **Abnormal Behavior**
  + **Collapsible Host**
  + **System Vulnerability**

* `protect_status` - (Optional, String) Specifies the protection status.
  The valid values are as follows:
  + **closed**
  + **opened**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of application protection events.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_name` - The host name.

* `private_ip` - The host private IP address.

* `event_name` - The alarm name.

* `severity` - The alarm severity.

* `req_src_ip` - The source IP address.

* `app_stack` - The application call stack information.

* `attack_input_name` - The attack affiliated field.

* `attack_input_value` - The attack payload content.

* `query_string` - The query string.

* `req_headers` - The web request header information.

* `req_method` - The web request method.

* `req_params` - The web request parameters.

* `req_path` - The web request path.

* `req_protocol` - The web request protocol.

* `req_url` - The web request URL.

* `attack_tag` - The attack tag.

* `chk_probe` - The probe identification.

* `chk_rule` - The check rule identification.

* `chk_rule_desc` - The check rule description.

* `exist_bug` - Whether the application exist a bug.
