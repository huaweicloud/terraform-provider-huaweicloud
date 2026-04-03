---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_pops"
description: |-
  Use this data source to get the list of access points within HuaweiCloud.
---

# huaweicloud_ga_pops

Use this data source to get the list of access points within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_ga_pops" "test" {}
```

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `pops` - The list of access points.

  The [pops](#pops_struct) structure is documented below.

<a name="pops_struct"></a>
The `pops` block supports:

* `id` - The access point ID.
