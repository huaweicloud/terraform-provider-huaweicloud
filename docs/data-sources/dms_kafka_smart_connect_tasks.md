---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_smart_connect_tasks"
description: ""
---

# huaweicloud_dms_kafka_smart_connect_tasks

Use this data source to get the list of DMS kafka smart connect tasks.

## Example Usage

```hcl
var "connector_id" {}
var "task_id" {}

data "huaweicloud_dms_kafka_smart_connect_tasks" "test" {
  connector_id     = var.connector_id
  task_id          = var.task_id
  task_name        = "test_task"
  destination_type = "OBS"
  status           = "RUNNING"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `connector_id` - (Required, String) Specifies the connector ID of the kafka instance.

* `task_id` - (Optional, String) Specifies the ID of the smart connect task.

* `task_name` - (Optional, String) Specifies the name of the smart connect task.

* `destination_type` - (Optional, String) Specifies the destination type of the smart connect task.

* `status` - (Optional, String) Specifies the status of the smart connect task. Value options: **RUNNING**, **PAUSED**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of smart connect tasks.
  The [tasks](#DMS_kafka_smart_connect_tasks) structure is documented below.

<a name="DMS_kafka_smart_connect_tasks"></a>
The `tasks` block supports:

* `id` - Indicates the ID of the smart connect task.

* `task_name` - Indicates the name of the smart connect task.

* `destination_type` - Indicates the destination type of the smart connect task.

* `created_at` - Indicates the creation time of the smart connect task.

* `status` - Indicates the status of the smart connect task. The value can be: **RUNNING**, **PAUSED**.

* `topics` - Indicates the topic names separated by commas or the topic regular expression of the smart connect task.
