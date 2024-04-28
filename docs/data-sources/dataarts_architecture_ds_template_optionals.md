---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_ds_template_optionals"
description: ""
---

# huaweicloud_dataarts_architecture_ds_template_optionals

Use this data source to get the list of DataArts Architecture data standard template optionals.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_ds_template_optionals" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID of DataArts Architecture.

* `fd_name` - (Optional, String) Specifies the name of the optional field.

* `required` - (Optional, Bool) Specifies whether the field is required. Defaults to **false**.

* `searchable` - (Optional, Bool) Specifies whether the field is search supported. Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `optional_fields` - Indicates the list of DataArts Architecture data standard template optional fields.
  The [optional_fields](#TemplateOptionalFields_OptionalField) structure is documented below.

<a name="TemplateOptionalFields_OptionalField"></a>
The `optional_fields` block supports:

* `fd_name` - Indicates the name of the optional field.

* `description` - Indicates the description of the field.

* `description_en` - Indicates the English description of the field.

* `required` - Indicates whether the field is required.

* `searchable` - Indicates whether the field is search supported.
