---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_offline_key_analyses"
description: |-
  Use this data source to query the offline key analyses.
---

# huaweicloud_dcs_offline_key_analyses

Use this data source to query the offline key analyses

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_offline_key_analyses" "test" {
  instance_id = var.instance_id
}
```

### Filter Tasks by Status

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_offline_key_analyses" "test" {
  instance_id = var.instance_id
  status      = "success"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the offline key analyses.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `status` - (Optional, String) Specifies the status of the offline key analysis task to filter.
  The valid values are as follows:
  + **waiting**: Task is waiting.
  + **running**: Task is running.
  + **success**: Task completed successfully.
  + **failed**: Task failed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of offline key analysis task records.
  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `id` - The task execution record ID.

* `status` - The status of the offline key analysis task.
  The valid values are as follows:
  + **waiting**: Task is waiting.
  + **running**: Task is running.
  + **success**: Task completed successfully.
  + **failed**: Task failed.

* `created_at` - The time when the analysis task was created.

* `started_at` - The time when the analysis task started.

* `finished_at` - The time when the analysis task finished.
