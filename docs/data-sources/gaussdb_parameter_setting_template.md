---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_parameter_setting_template"
description: |-
  Use this data source to query the parameter setting template for data redistribution within HuaweiCloud.
---

# huaweicloud_gaussdb_parameter_setting_template

Use this data source to query the parameter setting template for data redistribution within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_gaussdb_parameter_setting_template" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level
  region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `expansion_parameters` - The list of parameter information.  
  The [expansion_parameters](#parameter_setting_template_expansion_parameters) structure is documented below.

<a name="parameter_setting_template_expansion_parameters"></a>
The `expansion_parameters` block supports:

* `name` - The parameter name.

* `value` - The parameter value.

* `restart_required` - Whether restarting is required after modifying the parameter.

* `value_range` - The value range of the parameter.

* `type` - The parameter type. The valid values are **integer**, **boolean** and **string**.

* `description` - The parameter description.

* `risk_description` - The risk description.
