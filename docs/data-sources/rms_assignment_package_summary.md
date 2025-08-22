---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_assignment_package_summary"
description: |-
  Use this data source to get the compliance results of conformance packages.
---

# huaweicloud_rms_assignment_package_summary

Use this data source to get the compliance results of conformance packages.

## Example Usage

```hcl
data "huaweicloud_rms_assignment_package_summary" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `conformance_pack_name` - (Optional, String) Specifies the conformance package name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `value` - Indicates the summary of the compliance results of conformance packages.

  The [value](#value_struct) structure is documented below.

<a name="value_struct"></a>
The `value` block supports:

* `id` - Indicates the ID of a conformance package.

* `name` - Indicates the conformance package name.

* `compliance` - Indicates the compliance result of a conformance package.
