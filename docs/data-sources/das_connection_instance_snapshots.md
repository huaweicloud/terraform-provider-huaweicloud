---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_connection_instance_snapshots"
description: |-
  Use this data source to get the list of DAS connection instance snapshots.
---

# huaweicloud_das_connection_instance_snapshots

Use this data source to get the list of DAS connection instance snapshots.

-> This data source only supports to query snapshots of **MySQL** instances.

## Example Usage

```hcl
variable "user_id" {}

data "huaweicloud_das_connection_instance_snapshots" "test" {
  user_id    = var.user_id
  module     = 1
  start_time = "2026-05-20T00:00:00+08:00"
  end_time   = "2026-05-26T00:00:00+08:00"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the connection instance snapshots are located.  
  If omitted, the provider-level region will be used.

* `user_id` - (Required, String) Specifies the database user ID.

* `module` - (Required, Int) Specifies the lock snapshot type.  
  The valid values are as follows:
  + **0**: MDL lock snapshot.
  + **1**: InnoDB lock snapshot.
  + **2**: Recent deadlocks.

* `start_time` - (Required, String) Specifies the start time of the query range, in RFC3339 format.  

* `end_time` - (Required, String) Specifies the end time of the query range, in RFC3339 format.

-> 1.The earliest `start_time` and `end_time` is `7` days earlier than the current time.  
2.The `end_time` must be greater than the `start_time`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `snapshots` - The list of lock snapshots.  
  The [snapshots](#connection_instance_snapshots_attr) structure is documented below.

<a name="connection_instance_snapshots_attr"></a>
The `snapshots` block supports:

* `id` - The snapshot ID.

* `status` - The snapshot status.
  + **0**: Waiting.
  + **1**: Running.
  + **2**: Failed.
  + **3**: Successful.

* `created_time` - The snapshot creation time, in RFC3339 format.

* `find_lock` - Whether a lock was found.
  + **0**: No.
  + **1**: Yes.
