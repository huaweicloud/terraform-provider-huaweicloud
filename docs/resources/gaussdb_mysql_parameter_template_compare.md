---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_parameter_template_compare"
description: |-
  Manages a GaussDB MySQL parameter template compare resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_parameter_template_compare

Manages a GaussDB MySQL parameter template compare resource within HuaweiCloud.

## Example Usage

```hcl
variable "source_configuration_id" {}
variable "target_configuration_id" {}

resource "huaweicloud_gaussdb_mysql_parameter_template_compare" "test" {
  source_configuration_id = var.source_configuration_id
  target_configuration_id = var.target_configuration_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `source_configuration_id` - (Required, String, ForceNew) Specifies the ID of the source parameter template to be
  compared. Changing this parameter will create a new resource.

* `target_configuration_id` - (Required, String, ForceNew) Specifies the ID of the destination parameter template to be
  compared. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `<source_configuration_id>/<target_configuration_id>`.

* `differences` - Indicates the differences between parameters.
  The [differences](#differences_struct) structure is documented below.

<a name="differences_struct"></a>
The `differences` block supports:

* `parameter_name` -  Indicates the parameter name.

* `source_value` -  Indicates the parameter value in the source parameter template.

* `target_value` -  Indicates the parameter value in the destination parameter template.
