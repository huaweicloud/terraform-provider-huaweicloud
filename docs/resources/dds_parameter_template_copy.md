---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_parameter_template_copy"
description: |-
  Manages a DDS parameter template copy resource within HuaweiCloud.
---

# huaweicloud_dds_parameter_template_copy

Manages a DDS parameter template copy resource within HuaweiCloud.

## Example Usage

```hcl
variable "configuration_id" {}
variable "name" {}

resource "huaweicloud_dds_parameter_template_copy" "test" {
  configuration_id = var.configuration_id
  name             = var.name
  description      = "test copy"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `configuration_id` - (Required, String, ForceNew) Specifies the parameter template ID.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of replicated parameter template.
  The parameter template name can contain **1** to **64** characters. It can contain only letters, digits, hyphens (-),
  underscores (_), and periods (.).
  Changing this creates a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of replicated parameter template.
  The value is left blank by default. The description must consist of a maximum of **256** characters and cannot contain
  the carriage return character or the following special characters: >!<"&'=
  Changing this creates a new resource.

* `parameter_values` - (Optional, Map) Specifies the mapping between parameter names and parameter values.
  You can customize parameter values based on the parameters in the default parameter template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `node_version` - Indicates the database version.

* `parameters` - Indicates the parameters defined by users based on the default parameter templates.
  The [Parameter](#DdsParameterTemplate_Parameter) structure is documented below.

* `created_at` - The create time of the parameter template.

* `updated_at` - The update time of the parameter template.

<a name="DdsParameterTemplate_Parameter"></a>
The `Parameter` block supports:

* `name` - Indicates the parameter name.

* `value` - Indicates the parameter value.

* `description` - Indicates the parameter description.

* `type` - Indicates the parameter type. The value can be integer, string, boolean, float, or list.

* `value_range` - Indicates the value range.

* `restart_required` - Indicates whether the instance needs to be restarted.
  + If the value is **true**, restart is required.
  + If the value is **false**, restart is not required.

* `readonly` - Indicates whether the parameter is read-only.
  + If the value is **true**, the parameter is read-only.
  + If the value is **false**, the parameter is not read-only.

## Import

The DDS copyed parameter template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dds_parameter_template_copy.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `configuration_id`.
It is generally recommended running `terraform plan` after importing an template.
You can then decide if changes should be applied to the template, or the resource definition should be updated to
align with the template. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dds_parameter_template_copy" "test" {
    ...

  lifecycle {
    ignore_changes = [
      configuration_id,
    ]
  }
}
```
