---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_history_transaction_exported_tasks"
description: |-
  Use this data source to get the list of DAS history transaction exported tasks.
---

# huaweicloud_das_history_transaction_exported_tasks

Use this data source to get the list of DAS history transaction exported tasks.

## Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_das_history_transaction_exported_tasks" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the history transaction exported tasks are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of history transaction exported tasks.  
  The [tasks](#history_transaction_exported_tasks_tasks) structure is documented below.

<a name="history_transaction_exported_tasks_tasks"></a>
The `tasks` block supports:

* `id` - The exported task ID.

* `instance_id` - The instance ID.

* `status` - The task status.
  + **0**. Initialized.
  + **1**. Running.
  + **2**. Partially successful.
  + **3**. Successful.
  + **4**. Failed.
  + **-1**. Deleted.

* `start_time` - The start time, in RFC3339 format.

* `end_time` - The end time, in RFC3339 format.

* `created_time` - The task creation time, in RFC3339 format.

* `export_line_num` - The number of exported lines.

* `download_url` - The download URL of the exported file.
