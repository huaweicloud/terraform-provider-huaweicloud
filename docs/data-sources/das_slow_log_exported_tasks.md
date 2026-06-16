---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_slow_log_exported_tasks"
description: |-
  Use this data source to get the list of DAS slow log exported tasks.
---

# huaweicloud_das_slow_log_exported_tasks

Use this data source to get the list of DAS slow log exported tasks.

## Example Usage

### Filter by export type

```hcl
variable "instance_id" {}

data "huaweicloud_das_slow_log_exported_tasks" "test" {
  instance_id = var.instance_id
  export_type = "slowsqldetails"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the slow log exported tasks are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `export_type` - (Optional, String) Specifies the export type.  
  The valid values are as follows:
  + **slowsql**
  + **slowsqldetails**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of slow log exported tasks that matched filter parameters.  
  The [tasks](#slow_log_exported_tasks_attr) structure is documented below.

<a name="slow_log_exported_tasks_attr"></a>
The `tasks` block supports:

* `id` - The task ID.

* `instance_id` - The instance ID.

* `status` - The task status.

* `created_time` - The task creation time, in RFC3339 format.

* `start_time` - The task start time, in RFC3339 format.

* `end_time` - The task end time, in RFC3339 format.

* `download_url` - The download URL.
