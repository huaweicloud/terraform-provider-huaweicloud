---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_flow_logs"
description: |-
  Use this data source to get the list of VPC flow logs.
---

# huaweicloud_vpc_flow_logs

Use this data source to get the list of VPC flow logs.

## Example Usage

```hcl
variable "flow_log_name" {}

data "huaweicloud_vpc_flow_logs" "basic" {
  name = var.flow_log_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the VPC flow log name.
  The value can contain no more than 64 characters,
  including letters, digits, underscores (_), hyphens (-), and periods (.).

* `flow_log_id` - (Optional, String) Specifies the VPC flow log ID.

* `resource_type` - (Optional, String) Specifies the resource type for which that the logs to be collected.
  The value can be: **port**, **network,** and **vpc**.

* `resource_id` - (Optional, String) Specifies the resource ID for which that the logs to be collected.

* `log_group_id` - (Optional, String) Specifies the LTS log group ID.

* `log_stream_id` - (Optional, String) Specifies the LTS log stream ID.

* `traffic_type` - (Optional, String) Specifies the type of traffic to log.
  The value can be: **all**, **accept** and **reject**.

* `status` - (Optional, String) Specifies the status of the flow log.
  The value can be **ACTIVE**, **DOWN** or **ERROR**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flow_logs` - The list of VPC flow logs.

  The [flow_logs](#flow_logs_struct) structure is documented below.

<a name="flow_logs_struct"></a>
The `flow_logs` block supports:

* `name` - The VPC flow log name.

* `id` - The ID of a VPC flow log

* `description` - The VPC flow log description.

* `resource_type` - The resource type for which that the logs to be collected.

* `resource_id` - The resource ID for which that the logs to be collected.

* `log_group_id` - The LTS log group ID.

* `log_stream_id` - The LTS log stream ID.

* `traffic_type` - The type of traffic to log.

* `enabled` - Whether to enable the VPC flow log.

* `status` - The VPC flow log status.

* `created_at` - The time when the resource is created.

* `updated_at` - The time when the resource is last updated.
