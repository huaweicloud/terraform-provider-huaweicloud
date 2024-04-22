---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_global_connection_bandwidth_spec_codes"
description: |-
  Use this data source to get the list of CC line specification of global connection bandwidths.
---

# huaweicloud_cc_global_connection_bandwidth_spec_codes

Use this data source to get the list of CC line specification of global connection bandwidths.

## Example Usage

```hcl
variable "spec_id" {}

data "huaweicloud_cc_global_connection_bandwidth_spec_codes" "test" {
  spec_id = var.spec_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `spec_id` - (Optional, String) Line specification ID.

* `local_area` - (Optional, String) Local access point included in the line specification.

* `remote_area` - (Optional, String) Remote access point included in the line specification.

* `level` - (Optional, String) Line grade.
  The valid values are as follows:
  + **Pt**: platinum.
  + **Au**: gold.
  + **Ag**: silver.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `spec_codes` - The line specification list.

  The [spec_codes](#spec_codes_struct) structure is documented below.

<a name="spec_codes_struct"></a>
The `spec_codes` block supports:

* `id` - Line specification ID.

* `local_area` - Local access point.

* `remote_area` - Remote access point.

* `created_at` - Time when the resource was created.

* `updated_at` - Time when the resource was updated.

* `name_zh` - Line specification in Chinese.

* `name_en` - Line specification in English.

* `level` - Line grade.

* `sku` - Product code of specific global connection bandwidth line specifications.

* `size` - Minimum bandwidth for sale, in Mbit/s.
