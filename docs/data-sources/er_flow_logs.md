---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_flow_logs"
description: ""
---

# huaweicloud_er_flow_logs

Use this data source to get the list of flow logs.

## Example Usage

```hcl
variable "instance_id" {}
variable "resource_id" {}

data "huaweicloud_er_flow_logs" "test" {
  instance_id = var.instance_id
  resource_id = var.resource_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the flow logs are located.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the ER instance to which the flow logs belong.

* `resource_type` - (Optional, String) Specifies the type of the flow logs.
  The valid values are as follows:
  + **attachment**: The flow logs type are attachment.

* `resource_id` - (Optional, String) Specifies the ID of the attachment to which the flow logs belong.

* `flow_log_id` - (Optional, String) Specifies the ID of the flow log.

* `name` - (Optional, String) Specifies the name of the flow log.

* `status` - (Optional, String) Specifies the status of the flow logs.

* `enabled` - (Optional, String) Specifies the switch status of the flow log.
  The value can be **true** and **false**.

* `log_group_id` - (Optional, String) Specifies the ID of the log group to which the flow logs belong.

* `log_stream_id` - (Optional, String) Specifies the ID of the log stream to which the flow logs belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flow_logs` - The list ot the flow logs.
  The [flow_logs](#flowLogs) structure is documented below.

<a name="flowLogs"></a>
The `flow_logs` block supports:

* `id` - The ID of the flow log.

* `name` - The name of the flow log.

* `description` - The description of the flow log.

* `resource_type` - The type of the flow log.

* `resource_id` - The ID of the attachment to which the flow log belongs.

* `log_group_id` - The ID of the log group to which the flow log belongs.

* `log_stream_id` - The ID of the log stream to which the flow log belongs.

* `log_store_type` - The storage type of the flow log.

* `created_at` - The creation time of the flow log.

* `updated_at` - The latest update time of the flow log.

* `status` - The status of the flow log.

* `enabled` - The switch of the flow log.
