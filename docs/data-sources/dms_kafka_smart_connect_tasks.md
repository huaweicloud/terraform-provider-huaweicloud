---
subcategory: "Distributed Message Service (DMS)"
---

# huaweicloud_dms_kafka_smart_connect_tasks

Use this data source to get the list of DMS kafka smart connect task.

## Example Usage

```hcl
var "connector_id" {}

data "huaweicloud_dms_kafka_smart_connect_tasks" "test" {
  connector_id     = var.connector_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `connector_id` - (Required, String) Specifies the connector id of the kafka instance.

* `task_id` - (Optional, String) Specifies the id of the smart connect task.

* `task_name` - (Optional, String) Specifies the name of the smart connect task.

* `destination_type` - (Optional, String) Specifies the destination type of the smart connect task.

* `status` - (Optional, String) Specifies the status of the smart connect task.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of smart connector task.
  The [tasks](#DMS_kafka_smart_connect_tasks) structure is documented below.

<a name="DMS_kafka_smart_connect_tasks"></a>
The `groups` block supports:

* `task_id` - Indicates the id of the smart connector task.

* `task_name` - Indicates the name of the smart connector task.

* `destination_type` - Indicates the destination type of the smart connector task.

* `created_at` - Indicates the creation time of the smart connector task.

* `status` - Indicates the status of the smart connector task. The value can be **RUNNING**, **PAUSED**.

* `topics` - Indicates the topic names string or topic regular expression.

