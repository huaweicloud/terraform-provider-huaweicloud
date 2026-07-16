---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_parameter_template_save"
description: |-
  Manages a resource to save instance parameters as a parameter template within HuaweiCloud.
---

# huaweicloud_gaussdb_parameter_template_save

Manages a resource to save instance parameters as a parameter template within HuaweiCloud.

## Example Usage

```hcl
variable "config_id" {}
variable "name" {}

resource "huaweicloud_gaussdb_parameter_template_save" "test" {
  config_id = var.config_id
  name      = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `config_id` - (Required, String, NonUpdatable) Specifies the ID of the parameter template to be saved.

* `name` - (Required, String, NonUpdatable) Specifies the name of the saved parameter template.
  The name must be unique.
  The template name can contain `1` to `64` characters. Only letters (case-sensitive), digits, hyphens (-),
  underscores (_), and periods (.) are allowed.

* `description` - (Optional, String) Specifies the parameter template description.
  The parameter can contain a maximum of `256` characters. Carriage return characters and the following special
  characters are not allowed: **>!<"&'=**.

* `values` - (Optional, Map) Specifies the template parameters and corresponding to values need to be modified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `engine_version` - Indicates the engine version.

* `instance_mode` - Indicates the deployment model.
  + **ha**: Centralized deployment.
  + **independent**: Independent.

* `created_at` - Indicates the creation time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `updated_at` - Indicates the modification time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `configuration_parameters` - Indicates the list of the template parameters.
  The [configuration_parameters](#configuration_parameters_struct) structure is documented below.

<a name="configuration_parameters_struct"></a>
The `configuration_parameters` block supports:

* `name` - Indicates the parameter name.

* `value` - Indicates the parameter value.

* `need_restart` - Whether the instance needs to be rebooted.

* `readonly` - Whether the parameter is read-only.

* `value_range` - Indicates the parameter value range.

* `data_type` - Indicates the data type.
  The value can be **string**, **integer**, **boolean**, **list**, **float** or **all**.

* `description` - Indicates the parameter description.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_parameter_template_save.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `config_id`, `values`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_gaussdb_parameter_template_save" "test" {
  ...

  lifecycle {
    ignore_changes = [
      config_id, values,
    ]
  }
}
```
