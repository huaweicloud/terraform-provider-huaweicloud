---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_rts_waiting_event_statistics"
description: |-
  Use this data source to query the session wait event statistics of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_rts_waiting_event_statistics

Use this data source to query the session wait event statistics of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_rts_waiting_event_statistics" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the session wait event statistics.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `wait_events_info` - The list of wait event statistics.
  The [wait_events_info](#session_wait_events_info_attr) structure is documented below.

<a name="session_wait_events_info_attr"></a>
The `wait_events_info` block supports:

* `node_name` - The component name.

* `event_name` - The wait event name.

* `count` - The wait event count.
