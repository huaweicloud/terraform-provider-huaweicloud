---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_background_tasks"
description: |-
  Use this data source to get the list of Kafka background tasks.
---

# huaweicloud_dms_kafka_background_tasks

Use this data source to get the list of Kafka background tasks.

## Example Usage

```hcl
variable "instance_id" {}
variable "begin_time" {}
variable "end_time" {}

data "huaweicloud_dms_kafka_background_tasks" "test" {
  instance_id = var.instance_id
  begin_time  = var.begin_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `begin_time` - (Optional, String) Specifies the time of task where the query starts. The format is **YYYYMMDDHHmmss**.

* `end_time` - (Optional, String) Specifies the time of task where the query ends. The format is **YYYYMMDDHHmmss**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the task list.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - Indicates the task ID.

* `name` - Indicates the task name.

* `params` - Indicates the task parameters.

* `status` - Indicates the task status.

* `user_id` - Indicates the user ID.

* `user_name` - Indicates the username.

* `created_at` - Indicates the start time.

* `updated_at` - Indicates the end time.
