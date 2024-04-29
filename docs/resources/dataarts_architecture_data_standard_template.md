---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_data_standard_template"
description: ""
---

# huaweicloud_dataarts_architecture_data_standard_template

Manages DataArts Architecture data standard template resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

resource "huaweicloud_dataarts_architecture_data_standard_template" "test"{
  workspace_id = var.workspace_id

  optional_fields {
    fd_name    = "dataLength"
    required   = true
    searchable = true
  }

  optional_fields {
    fd_name    = "dqcRule"
    required   = true
    searchable = false
  }

  custom_fields {
    fd_name         = "custom_field_test"
    optional_values = "value1;value2;value3"
    required        = true
    searchable      = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID of DataArts Architecture. Changing this
  parameter will create a new resource.

* `optional_fields` - (Optional, List) Specifies the optional fields of the data standard template to be activated.
  The [optional_fields](#DataStandardTemplate_OptionalField) structure is documented below.

* `custom_fields` - (Optional, List) Specifies the custom fields of the data standard template to be added.
  The [custom_fields](#DataStandardTemplate_CustomField) structure is documented below.

<a name="DataStandardTemplate_OptionalField"></a>
The `optional_fields` block supports:

* `fd_name` - (Required, String) Specifies the name of the field.

* `required` - (Optional, Bool) Specifies whether the field is required.

* `searchable` - (Optional, Bool) Specifies whether the field is search supported.

<a name="DataStandardTemplate_CustomField"></a>
The `custom_fields` block supports:

* `fd_name` - (Required, String) Specifies the name of the field.

* `optional_values` - (Optional, String) Specifies the optional values of the field. Multiple values are separated by
  semicolons (;).

* `required` - (Optional, Bool) Specifies whether the field is required.

* `searchable` - (Optional, Bool) Specifies whether the field is search supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `optional_fields` - Indicates the optional fields of the data standard template to be activated.

  The [optional_fields](#DataStandardTemplate_OptionalField) structure is documented below.

* `custom_fields` - Indicates the custom fields of the data standard template to be added.

  The [custom_fields](#DataStandardTemplate_CustomField) structure is documented below.

<a name="DataStandardTemplate_OptionalField"></a>
The `optional_fields` block supports:

* `id` - Indicates the ID of the optional field.

* `created_by` - Indicates the name of creator.

* `updated_by` - Indicates the name of updater.

* `created_at` - Indicates the creation time of the field.

* `updated_at` - Indicates the latest update time of the field.

<a name="DataStandardTemplate_CustomField"></a>
The `custom_fields` block supports:

* `id` - Indicates the ID of the custom field.

* `created_by` - Indicates the name of creator.

* `updated_by` - Indicates the name of updater.

* `created_at` - Indicates the creation time of the field.

* `updated_at` - Indicates the latest update time of the field.

## Import

The DataArts architecture data standard template can be imported using the `workspace_id`, e.g.

```bash
$ terraform import huaweicloud_dataarts_architecture_data_standard_template.test <workspace_id>
```
