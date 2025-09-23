---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_traces"
description: |-
  Use this data source to get the list of CTS traces.
---

# huaweicloud_cts_traces

Use this data source to get the list of CTS traces.

## Example Usage

```hcl
variable "from" {}
variable "to" {}

data "huaweicloud_cts_traces" "test" {
  trace_type = "system"
  from       = var.from
  to         = var.to
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `trace_type` - (Required, String) Specifies the trace type.
  The value can be **system** (management trace) or **data** (data trace).
  The default value is **system**.

* `from` - (Required, String) Specifies the start time.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**.

* `to` - (Required, String) Specifies the end time.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**.

* `tracker_name` - (Optional, String) Specifies the tracker name.
  When **trace_type** is set to **system**, the value of this parameter is **system**.
  When **trace_type** is set to **data**, set this parameter to the name of a data tracker.

* `service_type` - (Optional, String) Specifies the cloud service type.
  This parameter is valid only when **trace_type** is set to **system**.

* `user` - (Optional, String) Specifies the user name.
  This parameter is valid only when **trace_type** is set to **system**.

* `resource_id` - (Optional, String) Specifies the cloud resource ID.
  This parameter is valid only when **trace_type** is set to **system**.

* `resource_name` - (Optional, String) Specifies the name of a resource.
  This parameter is valid only when **trace_type** is set to **system**.
  The value can contain uppercase letters.

* `resource_type` - (Optional, String) Specifies the type of a resource.
  This parameter is valid only when **trace_type** is set to **system**.

* `trace_id` - (Optional, String) Specifies the trace ID.
  If this parameter is specified, other query criteria will not take effect.
  This parameter is valid only when **trace_type** is set to **system**.

* `trace_name` - (Optional, String) Specifies the trace name.
  This parameter is valid only when **trace_type** is set to **system**.
  The value can contain uppercase letters.

* `trace_rating` - (Optional, String) Specifies the trace status.
  The value can be **normal**, **warning**, or **incident**.
  This parameter is valid only when **trace_type** is set to **system**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `traces` - The list of traces.

  The [traces](#traces_struct) structure is documented below.

<a name="traces_struct"></a>
The `traces` block supports:

* `api_version` - The version of the API called in the trace.

* `time` - The time when a trace was generated.

* `user` - The information of the user who performed the operation that triggered the trace.

  The [user](#traces_user_struct) structure is documented below.

* `resource_name` - The name of the resource on which the recorded operation was performed.

* `endpoint` - The endpoint in the details page URL of the cloud resource on which the recorded operation was performed.

* `trace_type` - The trace type.

* `record_time` - The time when a trace was recorded by CTS.

* `trace_id` - The Trace ID.

* `resource_type` - The type of the resource on which the recorded operation was performed.

* `source_ip` - The IP address of the tenant who performed the operation that triggered the trace.

* `resource_url` - The details page URL (excluding the endpoint) of the cloud resource.

* `request` - The request body of the recorded operation.

* `request_id` - The ID of the request of the recorded operation.

* `resource_id` - The ID of the cloud resource on which the recorded operation was performed.

* `trace_rating` - The trace status.

* `response` - The response body of the recorded operation.

* `code` - The returned HTTP status code of the recorded operation.

* `message` - The remarks added by other cloud services to the trace.

* `service_type` - The cloud service on which the recorded operation was performed.

* `location_info` - The information required for fault locating after a request error occurred.

* `trace_name` - The trace name.

* `read_only` - Whether a user request is read-only.

* `operation_id` - The operation ID of the trace.

<a name="traces_user_struct"></a>
The `user` block supports:

* `id` - The account ID.

* `name` - The account name.

* `domain` - The domain information of the user who performed the operation that triggered the trace.

  The [domain](#user_domain_struct) structure is documented below.

<a name="user_domain_struct"></a>
The `domain` block supports:

* `id` - The account ID.

* `name` - The account name.
