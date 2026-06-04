---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_alarms"
description: |-
  Use this data source to query GaussDB alarms within HuaweiCloud.
---

# huaweicloud_gaussdb_alarms

Use this data source to query GaussDB alarms within HuaweiCloud.

## Example Usage

```hcl
variable "start_time" {}

data "huaweicloud_gaussdb_alarms" "test" {
  start_time = var.start_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level
  region will be used.

* `start_time` - (Required, String) Specifies the start time for querying alarms. The value is in the
  **yyyy-mm-ddThh:mm:ssZ** format.

* `level` - (Optional, Int) Specifies the alarm level. The valid values are as follows:
 + **1**: CRITICAL.
 + **2**: MAJOR.
 + **3**: MINOR.
 + **4**: WARNING.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `history_records` - The list of alarm history records.
  The [history_records](#history_records_struct) structure is documented below.

<a name="history_records_struct"></a>
The `history_records` block supports:

* `alarm_id` - The alarm ID.

* `name` - The alarm name.

* `status` - The alarm status.

* `alarm_type` - The alarm type.

* `level` - The alarm level.

* `instance_id` - The ID of the GaussDB instance that triggered the alarm.

* `instance_name` - The name of the GaussDB instance that triggered the alarm.

* `begin_time` - The time when the alarm was triggered.

* `update_time` - The time when the alarm was updated.
