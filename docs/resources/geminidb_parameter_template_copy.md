---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_parameter_template_copy"
description: |-
  Manages a GeminiDB parameter template copy resource within HuaweiCloud.
---

# huaweicloud_geminidb_parameter_template_copy

Manages a GeminiDB parameter template copy resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "config_id" {}

resource "huaweicloud_geminidb_parameter_template_copy" "test" {
  config_id = var.config_id
  name      = "my-configuration-copy"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the configuration copy.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `config_id` - (Required, String, NonUpdatable) Specifies the ID of the source parameter template.

* `name` - (Required, String) Specifies the name of the copied parameter template.

* `description` - (Optional, String) Specifies the description of the copied parameter template.

* `values` - (Optional, Map) Specifies the parameter values.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `datastore_version_name` - The database version name.

* `datastore_name` - The database name.

* `mode` - The database instance mode.

* `created` - The creation time. The format is **yyyy-MM-ddTHH:mm:ssZ**.

* `updated` - The update time. The format is **yyyy-MM-ddTHH:mm:ssZ**.

* `configuration_parameters` - The list of parameter configurations.
  The [configuration_parameters](#geminidb_configuration_parameters) structure is documented below.

<a name="geminidb_configuration_parameters"></a>
The `configuration_parameters` block supports:

* `name` - The parameter name.

* `value` - The parameter value.

* `restart_required` - The parameter whether a restart is required after modifying.

* `readonly` - The parameter whether is read-only.

* `value_range` - The parameter value range.

* `type` - The parameter type. Valid values: **string**, **integer**, **boolean**, **list**, **float**.

* `description` - The parameter description.

## Import

The GeminiDB parameter template can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_geminidb_parameter_template_copy.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `config_id` and `values`.
It is generally recommended running `terraform plan` after importing a parameter template.
You can then decide if changes should be applied to the parameter template, or the resource definition should be
updated to align with the parameter template. Also you can ignore changes as below.

```hcl
resource "huaweicloud_geminidb_parameter_template_copy" "test" {
  ...

  lifecycle {
    ignore_changes = [
      config_id, values,
    ]
  }
}
