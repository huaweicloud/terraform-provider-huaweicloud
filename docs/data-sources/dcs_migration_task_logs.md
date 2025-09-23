---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_migration_task_logs"
description: |-
  Use this data source to get the list of the logs of a migration task.
---

# huaweicloud_dcs_migration_task_logs

Use this data source to get the list of the logs of a migration task.

## Example Usage

```hcl
variable "task_id" {}

data "huaweicloud_dcs_migration_task_logs" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `task_id` - (Required, String) Indicates the ID of the data migration task.

* `log_level` - (Optional, String) Indicates the log level.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `migration_logs` - Indicates the log list.

  The [migration_logs](#migration_logs_struct) structure is documented below.

<a name="migration_logs_struct"></a>
The `migration_logs` block supports:

* `keyword` - Indicates the log keyword.

* `log_level` - Indicates the log level.

* `message` - Indicates the log information.

* `log_code` - Indicates the log code.

* `created_at` - Indicates the time when the migration log is generated.
