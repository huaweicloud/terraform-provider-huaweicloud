---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_bigkey_analyses"
description: ""
---

# huaweicloud_dcs_bigkey_analyses

Use this data source to get the list of DCS big key analyses.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_bigkey_analyses" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `analysis_id` - (Optional, String) Specifies the  ID of the big key analysis.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `scan_type` - (Optional, String) Specifies the mode of the big key analysis.
  Value options: **manual**, **auto**.

* `status` - (Optional, String) Specifies the status of the big key analysis.
  Value options: **waiting**, **running**, **success**, **failed**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - Indicates the list of big key analysis records.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `id` - Indicates the id of the big key analysis.

* `scan_type` - Indicates the mode of the big key analysis

* `status` - Indicates the status of the big key analysis.

* `created_at` - Indicates the creation time of the big key analysis. The value is in UTC format.

* `started_at` - Indicates the time when the big key analysis started. The value is in UTC format.

* `finished_at` - Indicates the time when the big key analysis ended. The value is in UTC format.
