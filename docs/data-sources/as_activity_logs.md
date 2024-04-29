---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_activity_logs"
description: ""
---

# huaweicloud_as_activity_logs

Use this data source to get a list of AS scaling activity logs within HuaweiCloud.

## Example Usage

```hcl
variable "group_id" {}

data "huaweicloud_as_activity_logs" "test" {
  scaling_group_id = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the AS scaling activity logs.
  If omitted, the provider-level region will be used.

* `scaling_group_id` - (Required, String) Specifies the AS scaling group ID.

* `start_time` - (Optional, String) Specifies the start time of the AS scaling activity for query. The time format is
  **yyyy-MM-ddThh:mm:ssZ**.  
  The query result is for all data with a start time greater than or equal to this value.

* `end_time` - (Optional, String) Specifies the end time of the AS scaling activity for query. The time format is
  **yyyy-MM-ddThh:mm:ssZ**.  
  The query result shows all data with an end time less than this value.

* `status` - (Optional, String) Specifies the status of the AS scaling activity for query.  
  The valid values are as follows:
  + **SUCCESS**: Scaling activity execution successfully.
  + **FAIL**: Scaling activity execution failed.
  + **DOING**: Scaling activity is currently being executed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `activity_logs` - All AS scaling activity logs that match the filter parameters.
  The [activity_logs](#as_activity_logs) structure is documented below.

<a name="as_activity_logs"></a>
The `activity_logs` block supports:

* `id` - The scaling activity log ID.

* `status` - The status of the AS scaling activity.

* `start_time` - The start time of the AS scaling activity.

* `end_time` - The end time of the AS scaling activity.

* `removed_instances` - A list of cloud server names that have completed scaling activity and are only removed from
  the elastic scaling group, separated by commas.

* `deleted_instances` - A list of cloud server names that have completed scaling activity and been removed from the
  elastic scaling group and deleted, separated by commas.

* `added_instances` - A list of cloud server names that have completed scaling activity and been added to the elastic
  scaling group, separated by commas.

* `current_instance_number` - The current number of instances of the AS scaling group.

* `desire_instance_number` - The final expected number of instances for AS scaling activity.

* `changes_instance_number` - The number of cloud servers increased or decreased during AS scaling activity.

* `description` - The description of AS scaling activity.
