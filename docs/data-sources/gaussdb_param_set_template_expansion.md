---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_param_set_template_expansion"
description: |-
  Use this data source to query the parameter setting template for scale-out optimization of GaussDB within HuaweiCloud.
---

# huaweicloud_gaussdb_param_set_template_expansion

Use this data source to query the parameter setting template for scale-out optimization of GaussDB within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_gaussdb_param_set_template_expansion" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the parameter setting template for expansion.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `param_set_template_expansion` - The list of expansion optimization parameters.
  The [param_set_template_expansion](#param_set_template_expansion_attr) structure is documented below.

<a name="param_set_template_expansion_attr"></a>
The `param_set_template_expansion` block supports:

* `name` - The parameter name.

* `value` - The parameter value.

* `restart_required` - Indicates whether a restart is required after modifying the parameter.

* `value_range` - The valid value range of the parameter.

* `type` - The parameter type.
  Valid values are **integer**, **boolean**, and **string**.

* `description` - The parameter description.
