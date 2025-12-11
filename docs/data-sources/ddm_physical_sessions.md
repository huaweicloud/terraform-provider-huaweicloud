---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_physical_sessions"
description: |-
  Use this data source to get the physical sessions of a DDM instance.
---

# huaweicloud_ddm_physical_sessions

Use this data source to get the physical sessions of a DDM instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_ddm_physical_sessions" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the associated RDS instance.

* `keyword` - (Optional, String) Specifies the Keyword filtered by the sessions result. It can contain a maximum of 255
  characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `physical_processes` - Indicates the physical session information list.

  The [physical_processes](#physical_processes_struct) structure is documented below.

<a name="physical_processes_struct"></a>
The `physical_processes` block supports:

* `id` - Indicates the session ID.

* `user` - Indicates the current user.

* `host` - Indicates the IP address and port number.

* `db` - Indicates the database name.

* `state` - Indicates the status of the SQL statement being executed.

* `command` - Indicates the connection status.
  Generally, the value can be **sleep**, **query** or **connect**.

* `info` - Indicates the SQL statement that is being executed.

* `time` - Indicates the duration of a connection, in seconds.

* `trx_executed_time` - Indicates the duration of a transaction, in seconds.
