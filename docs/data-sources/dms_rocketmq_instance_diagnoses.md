---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_instance_diagnoses"
description: |-
  Use this data source to query the list of RocketMQ instance diagnosis reports.
---

# huaweicloud_dms_rocketmq_instance_diagnoses

Use this data source to query the list of RocketMQ instance diagnosis reports.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_rocketmq_instance_diagnoses" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the diagnosis reports are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `reports` - The list of the diagnosis reports.  
  The [reports](#instance_diagnoses_reports) structure is documented below.

<a name="instance_diagnoses_reports"></a>
The `reports` block supports:

* `report_id` - The ID of the diagnosis report.

* `group_name` - The name of the consumer group.

* `status` - The status of the report.

* `created_at` - The creation time of the report.

* `abnormal_item_sum` - The number of abnormal items.

* `faulted_node_sum` - The number of faulted nodes.
