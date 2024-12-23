---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_parameter_template"
description: |-
  Manages a GaussDB OpenGauss parameter template resource within HuaweiCloud.
---

# huaweicloud_gaussdb_opengauss_parameter_template

Manages a GaussDB OpenGauss parameter template resource within HuaweiCloud.

## Example Usage

### create parameter template

```hcl
resource "huaweicloud_gaussdb_opengauss_parameter_template" "test" {
  name           = "test_gaussdb_opengauss_parameter_template"
  engine_version = "8.201"
  instance_mode  = "independent"

  parameters {
    name  = "audit_system_object"
    value = "100"
  }

  parameters {
    name  = "cms:enable_finishredo_retrieve"
    value = "on"
  }
}
```

### replica parameter template from existed configuration

```hcl
variable "source_configuration_id" {}

resource "huaweicloud_gaussdb_opengauss_parameter_template" "test" {
  name                    = "test_copy_from_configuration"
  source_configuration_id = var.source_configuration_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the parameter template, which must be unique. The template
  name can contain up to **64** characters. It can contain only letters (case-sensitive), digits, hyphens (-),
  underscores (_), and periods (.).

  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the Parameter template description. This parameter is left blank
  by default. Up to **256** characters are displayed. Carriage return characters or special characters (>!<"&'=) are not
  allowed.

  Changing this parameter will create a new resource.

* `engine_version` - (Optional, String, ForceNew) Specifies the DB engine version.

  Changing this parameter will create a new resource.

  -> **NOTE:** It is mandatory when `instance_mode` is specified, and can not be specified when `source_configuration_id`
  is specified.

* `instance_mode` - (Optional, String, ForceNew) Specifies the deployment model.

  Changing this parameter will create a new resource.

  -> **NOTE:** It is mandatory when `engine_version` is specified, and can not be specified when `source_configuration_id`
  is specified.

* `parameters` - (Optional, List, ForceNew) Specifies the list of the template parameters.
  The [parameters](#parameters_struct) structure is documented below.

  Changing this parameter will create a new resource.

  -> **NOTE:** It can not be specified when `source_configuration_id` is specified.

* `source_configuration_id` - (Optional, String, ForceNew) Specifies the source parameter template ID.

  Changing this parameter will create a new resource.

  -> **NOTE:** It can not be specified when `engine_version`, `instance_mode` or `parameters` are specified.

  -> **NOTE:** Exactly one of `engine_version` and `source_configuration_id` must be provided.

<a name="parameters_struct"></a>
The `parameters` block supports:

* `name` - (Required, String, ForceNew) Specifies the name of a specific parameter.

* `value` - (Required, String, ForceNew) Specifies the value of a specific parameter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - Indicates the creation time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `updated_at` - Indicates the modification time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `parameters` - Indicates the list of the template parameters.
  The [parameters](#parameters_struct) structure is documented below.

<a name="parameters_struct"></a>
The `parameters` block supports:

* `need_restart` - Indicates whether the instance needs to be rebooted.

* `readonly` - Indicates whether the parameter is read-only.

* `value_range` - Indicates the parameter value range.

* `data_type` - Indicates the data type. The value can be **string**, **integer**, **boolean**, **list**, **all**,
  or **float**.

* `description` - Indicates the parameter description.

## Import

The GaussDB OpenGauss parameter template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_opengauss_parameter_template.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `source_configuration_id` and `parameters`.
It is generally recommended running `terraform plan` after importing a GaussDB OpenGauss parameter template. You can then
decide if changes should be applied to the GaussDB OpenGauss parameter template, or the resource definition should be
updated to align with the GaussDB OpenGauss parameter template. Also you can ignore changes as below.

```hcl
resource "huaweicloud_gaussdb_opengauss_parameter_template" "test" {
  ...

  lifecycle {
    ignore_changes = [
      source_configuration_id, parameters,
    ]
  }
}
```
