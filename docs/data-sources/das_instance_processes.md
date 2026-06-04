---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_instance_processes"
description: |-
  Use this data source to query DAS instance processes within HuaweiCloud.
---

# huaweicloud_das_instance_processes

Use this data source to query DAS instance processes within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "db_user_id" {}

data "huaweicloud_das_instance_processes" "test" {
  instance_id = var.instance_id
  db_user_id  = var.db_user_id
}
```

### Filter by user

```hcl
variable "instance_id" {}
variable "db_user_id" {}
variable "db_user_name" {}

data "huaweicloud_das_instance_processes" "test" {
  instance_id  = var.instance_id
  db_user_id   = var.db_user_id
  db_user_name = var.db_user_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the instance processes are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the instance to which the processes belong.

* `db_user_id` - (Required, String) Specifies the ID of the database user.

* `db_user_name` - (Optional, String) Specifies the name of the database user.

* `db_name` - (Optional, String) Specifies the name of the database.

* `node_id` - (Optional, String) Specifies the ID of the instance node.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `processes` - The list of processes that matched filter parameters.  
  The [processes](#das_instance_processes_attr) structure is documented below.

<a name="das_instance_processes_attr"></a>
The processes block supports:

* `id` - The ID of the process.

* `db_user_name` - The name of the database user.

* `host` - The host of the process.

* `db_name` - The name of the database.

* `command` - The command being executed.

* `time` - The duration of the process, in seconds.

* `state` - The state of the process.

* `sql` - The SQL statement being executed.

* `trx_executed_time` - The duration of the transaction, in seconds.
