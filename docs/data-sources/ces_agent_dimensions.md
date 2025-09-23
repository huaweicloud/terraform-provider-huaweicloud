---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_agent_dimensions"
description: |-
  Use this data source to get the list of CES agent dimensions.
---

# huaweicloud_ces_agent_dimensions

Use this data source to get the list of CES agent dimensions.

## Example Usage

```hcl
variable "instance_id" {}
variable "dim_name" {}

data "huaweicloud_ces_agent_dimensions" "test" {
  instance_id = var.instance_id
  dim_name    = var.dim_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `dim_name` - (Required, String) Specifies the dimension name.
  The valid values are **mount_point**, **disk**, **proc**, **gpu** and **raid**.

* `dim_value` - (Optional, String) Specifies the dimension value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `dimensions` - The dimension information list.

  The [dimensions](#dimensions_struct) structure is documented below.

<a name="dimensions_struct"></a>
The `dimensions` block supports:

* `name` - The dimension name.

* `value` - The dimension value.

* `origin_value` - The actual dimension value.
