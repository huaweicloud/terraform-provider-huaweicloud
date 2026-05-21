---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_restorable_time_window"
description: |-
  Use this data source to get a list of restorable time window for a GeminiDB instance.
---

# huaweicloud_geminidb_restorable_time_window

Use this data source to get a list of restorable time window for a GeminiDB instance.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_restorable_time_window" "test" {
  instance_id = var.instance_id
}
```

### With Time Range

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_restorable_time_window" "test" {
  instance_id = var.instance_id
  start_time  = "2024-01-01T00:00:00+0800"
  end_time    = "2024-12-31T23:59:59+0800"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the restorable time window.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GeminiDB instance.

* `start_time` - (Optional, String) Specifies the start time for the query.
  The format is **yyyy-mm-ddThh:mm:ssZ**.
  Defaults to one day before the current query time.

* `end_time` - (Optional, String) Specifies the end time for the query.
  The format is **yyyy-mm-ddThh:mm:ssZ**.
  Defaults to the current query time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `restorable_time_periods` - The list of restorable time periods.
  The [restorable_time_periods](#geminidb_restorable_time_periods) structure is documented below.

<a name="geminidb_restorable_time_periods"></a>
The `restorable_time_periods` block supports:

* `start_time` - The start time of the restorable time period.
  This is a UNIX timestamp in milliseconds, in UTC timezone.

* `end_time` - The end time of the restorable time period.
  This is a UNIX timestamp in milliseconds, in UTC timezone.
