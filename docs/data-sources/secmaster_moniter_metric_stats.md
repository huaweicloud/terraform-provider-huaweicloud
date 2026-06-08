---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_moniter_metric_stats"
description: |-
  Use this data source to get the SecMaster moniter metric statistics within HuaweiCloud.
---

# huaweicloud_secmaster_moniter_metric_stats

Use this data source to get the SecMaster moniter metric statistics within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "dataspace_id" {}
variable "pipe_id" {}

data "huaweicloud_secmaster_moniter_metric_stats" "test" {
  workspace_id    = var.workspace_id
  dataspace_id    = var.dataspace_id
  pipe_id         = var.pipe_id
  start_timestamp = 1780313858887
  end_timestamp   = 1780918658887
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `dataspace_id` - (Required, String) Specifies the dataspace ID.

* `pipe_id` - (Required, String) Specifies the pipe ID.

* `start_timestamp` - (Required, Int) Specifies the start timestamp.

* `end_timestamp` - (Required, Int) Specifies the end timestamp.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `results` - The metric statistics list.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `average_msg_bytes` - The average message bytes.

* `subscribe_msgs` - The number of subscribed messages.
