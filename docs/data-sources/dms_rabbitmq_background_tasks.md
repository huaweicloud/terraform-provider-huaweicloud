---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_background_tasks"
description: |-
  Use this data source to get the list of DMS RabbitMQ background tasks.
---

# huaweicloud_dms_rabbitmq_background_tasks

Use this data source to get the list of DMS RabbitMQ background tasks.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_rabbitmq_background_tasks" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `begin_time` - (Optional, String) Specifies the time of task where the query starts. The format is YYYYMMDDHHmmss.

* `end_time` - (Optional, String) Specifies the time of task where the query ends. The format is YYYYMMDDHHmmss.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the task list.
  The [tasks](#attrblock--tasks) structure is documented below.

<a name="attrblock--tasks"></a>
The `tasks` block supports:

* `created_at` - Indicates the start time.

* `id` - Indicates the task ID.

* `name` - Indicates the task name.

* `params` - Indicates the task parameters.

* `status` - Indicates the task status.

* `updated_at` - Indicates the end time.

* `user_id` - Indicates the user ID.

* `user_name` - Indicates the username.
