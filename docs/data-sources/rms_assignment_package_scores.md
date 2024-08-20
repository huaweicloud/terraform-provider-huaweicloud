---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_assignment_package_scores"
description: |-
  Use this data source to get the list of RMS assignment package scores.
---

# huaweicloud_rms_assignment_package_scores

Use this data source to get the list of RMS assignment package scores.

## Example Usage

```hcl
data "huaweicloud_rms_assignment_package_scores" "test" {}
```

## Argument Reference

The following arguments are supported:

* `assignment_package_name` - (Optional, String) Specifies the assignment package name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `value` - The scores of assignment packages.

  The [value](#value_struct) structure is documented below.

<a name="value_struct"></a>
The `value` block supports:

* `id` - The ID of the assignment package.

* `name` - The assignment package name.

* `score` - The score of the assignment package.
