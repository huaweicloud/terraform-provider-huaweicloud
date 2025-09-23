---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_bandwidth_package_lines"
description: |-
  Use this data source to get the list of CC bandwidth lines.
---

# huaweicloud_cc_bandwidth_package_lines

Use this data source to get the list of CC bandwidth lines.

## Example Usage

```hcl
data "huaweicloud_cc_bandwidth_package_lines" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bandwidth_package_lines` - Indicates the list of bandwidth package lines.

  The [bandwidth_package_lines](#bandwidth_package_lines_struct) structure is documented below.

<a name="bandwidth_package_lines_struct"></a>
The `bandwidth_package_lines` block supports:

* `local_region_id` - Indicates the local region ID.

* `remote_region_id` - Indicates the remote region ID.

* `local_site_code` - Indicates the local site code.

* `remote_site_code` - Indicates the remote site code.

* `support_levels` - Indicates the list of supported classes.

* `spec_codes` - Indicates the offering code list.

  The [spec_codes](#bandwidth_package_lines_spec_codes_struct) structure is documented below.

<a name="bandwidth_package_lines_spec_codes_struct"></a>
The `spec_codes` block supports:

* `level` - Indicates the bandwidth package class.

* `spec_code` - Indicates the specification code of the bandwidth package.

* `name_cn` - Indicates the Chinese instance name.

* `name_en` - Indicates the English instance name.

* `max_bandwidth` - Indicates the maximum bandwidth.

* `min_bandwidth` - Indicates the minimum bandwidth.

* `support_billing_modes` - Indicates the billing mode.
