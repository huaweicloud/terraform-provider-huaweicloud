---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_data_standard_template_customs"
description: |-
  Use this data source to query DataArts Architecture data standard template customs within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_data_standard_template_customs

Use this data source to query DataArts Architecture data standard template customs within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_data_standard_template_customs" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the data standard template customs are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the data standard template customs belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `customs` - The list of data standard template customs that matched filter parameters.  
  The [customs](#dataarts_architecture_data_standard_template_customs) structure is documented below.

<a name="dataarts_architecture_data_standard_template_customs"></a>
The `customs` block supports:

* `id` - The ID of the data standard template custom, in UUID format.

* `fd_name` - The field name of data standard template custom.

* `fd_name_en` - The field english name of data standard template custom.

* `description` - The description of data standard template custom.

* `actived` - Whether the data standard template custom field is visible.

* `required` - Whether the data standard template custom field is required.

* `searchable` - Whether the data standard template custom field is searchable.

* `optional_values` - Valid range for the custom field of the data standard template.

* `create_time` - The creation time of the data standard template custom, in RFC3339 format.

* `update_time` - The latest update time of the data standard template custom, in RFC3339 format.

* `create_by` - The creator of the data standard template custom.

* `update_by` - The last editor of the data standard template custom.
