---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_intelligent_session_kill_history"
description: |-
  Use this data source to get the intelligent session killing history of a DB instance.
---

# huaweicloud_rds_intelligent_session_kill_history

Use this data source to get the intelligent session killing history of a DB instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_intelligent_session_kill_history" "test" {
  instance_id = var.instance_id
  start_time  = "2026-03-25 15:04:05"
  end_time    = "2026-03-28 15:04:05"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `start_time` - (Optional, String) Specifies the query start time. For example, **2026-03-15 15:04:05**.

* `end_time` - (Optional, String) Specifies the query end time. For example, **2026-03-18 15:04:05**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `history` - Indicates the list of intelligent session killing history.

  The [history](#history_struct) structure is documented below.

<a name="history_struct"></a>
The `history` block supports:

* `task_id` - Indicates the task ID.

* `start_time` - Indicates the start time for the killing operation.

* `end_time` - Indicates the end time for the killing operation.

* `download_link` - Indicates the link for downloading the operation history.
