---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_logical_sessions"
description: |-
  Use this data source to get the logical sessions of a DDM instance.
---

# huaweicloud_ddm_logical_sessions

Use this data source to get the logical sessions of a DDM instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_ddm_logical_sessions" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the DDM instance ID.

* `keyword` - (Optional, String) Specifies the keyword filtered by the session result.
  It is a fuzzy match field and can contain a maximum of 255 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `logical_processes` - Indicates the logical sessions.

  The [logical_processes](#logical_processes_struct) structure is documented below.

<a name="logical_processes_struct"></a>
The `logical_processes` block supports:

* `id` - Indicates the session ID

* `user` - Indicates the current user.

* `host` - Indicates the IP address and port number.

* `db` - Indicates the database name.

* `state` - Indicates the status of the SQL statement.

* `command` - Indicates the connection status.
  Generally, the value can be **sleep**, **query**, or **connect**.

* `info` - Indicates the SQL statement that is being executed.

* `time` - Indicates the duration of a connection, in seconds.
