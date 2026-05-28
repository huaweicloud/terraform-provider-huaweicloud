---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_pt_appliction_records"
description: |-
  Use this data source to get the application records of a GeminiDB parameter template.
---

# huaweicloud_geminidb_pt_appliction_records

Use this data source to get the application records of a GeminiDB parameter template.

## Example Usage

```hcl
variable "config_id" {}

data "huaweicloud_geminidb_pt_appliction_records" "test" {
  config_id = var.config_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the applied histories.
  If omitted, the provider-level region will be used.

* `config_id` - (Required, String) Specifies the ID of the parameter template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `histories` - The list of parameter template application histories.
  The [histories](#geminidb_applied_histories) structure is documented below.

<a name="geminidb_applied_histories"></a>
The `histories` block supports:

* `instance_id` - The instance ID.

* `instance_name` - The instance name.

* `applied_at` - The applied time in the format **yyyy-MM-ddTHH:mm:ssZ**.

* `apply_result` - The application result. Valid values are:
  + **SUCCESS**: Application successful.
  + **APPLYING**: Application in progress.
  + **FAILED**: Application failed.

* `failure_reason` - The failure reason (if the application failed).
