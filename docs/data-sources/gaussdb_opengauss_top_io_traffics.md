---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_top_io_traffics"
description: |-
  Use this data source to get the top I/O statistics of instance database processes and return the results associated with session information.
---

# huaweicloud_gaussdb_opengauss_top_io_traffics

Use this data source to get the top I/O statistics of instance database processes and return the results associated with
session information.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "component_id" {}

data "huaweicloud_gaussdb_opengauss_top_io_traffics" "test" {
  instance_id  = var.instance_id
  node_id      = var.node_id
  component_id = var.component_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the GaussDB OpenGauss instance ID.

* `node_id` - (Required, String) Specifies the node ID.
  It must be a DN node of a non-log role or a CN node, and the node status must be normal.

* `component_id` - (Required, String) Specifies the component ID.
  It must be a CN or a DN component of a non-log role.

* `top_io_num` - (Optional, Int) Specifies the number of top I/O threads to be queried in the database process.

* `sort_condition` - (Optional, String) Specifies the top I/O sorting condition.
  Value options: **read**, **write**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `top_io_infos` - Indicates the I/O information.

  The [top_io_infos](#top_io_infos_struct) structure is documented below.

<a name="top_io_infos_struct"></a>
The `top_io_infos` block supports:

* `thread_id` - Indicates the thread ID.

* `thread_type` - Indicates the thread type.
  The value can be **worker** or **background**.

* `disk_read_rate` - Indicates the rate of reading data from disks, in KB/s.

* `disk_write_rate` - Indicates the rate of writing data to the disk, in KB/s.

* `session_id` - Indicates the session ID.

* `unique_sql_id` - Indicates the SQL ID.

* `database_name` - Indicates the database name.

* `client_ip` - Indicates the IP address of the client.

* `user_name` - Indicates the username.

* `state` - Indicates the status.

* `sql_start` - Indicates the start time of the statement.
