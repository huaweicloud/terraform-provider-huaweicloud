---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_sessions"
description: |-
  Use this data source to query the current sessions of a TaurusDB HTAP instance within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_sessions

Use this data source to query the current sessions of a TaurusDB HTAP instance within HuaweiCloud.

## Example Usage

```hcl
variable "htap_instance_id" {}

data "huaweicloud_taurusdb_htap_sessions" "test" {
  instance_id = var.htap_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HTAP sessions.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the HTAP instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `process_list` - The list of sessions of the HTAP instance.
  The [process_list](#taurusdb_htap_process_list_attr) structure is documented below.

<a name="taurusdb_htap_process_list_attr"></a>
The `process_list` block supports:

* `id` - The session ID.

* `user` - The session username.

* `host` - The session host.

* `state` - The session status.

* `database` - The database corresponding to the session.

* `sql_statement` - The SQL statement executed by the session.

* `duration` - The session duration, in seconds.

* `command` - The session command type.
