---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_data_class_fields"
description: |-
  Use this data source to get the list of SecMaster data class fields.
---

# huaweicloud_secmaster_data_class_fields

Use this data source to get the list of SecMaster data class fields.

## Example Usage

```hcl
variable "workspace_id" {}
variable "data_class_id" {}

data "huaweicloud_secmaster_data_class_fields" "test" {
  workspace_id  = var.workspace_id
  data_class_id = var.data_class_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `data_class_id` - (Required, String) Specifies the data class ID.

* `name` - (Optional, String) Specifies the field name.

* `is_built_in` - (Optional, String) Specifies whether it is built in SecMaster. The value can be **true** or **false**.

* `mapping` - (Optional, String) Specifies whether to display in other places other the classification and mapping module.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `fields` - The field list.

  The [fields](#fields_struct) structure is documented below.

<a name="fields_struct"></a>
The `fields` block supports:

* `id` - The field ID.

* `name` - The field name.

* `type` - The field type, such as **short text**, **radio** and **grid**.

* `business_code` - The business code of the field.

* `business_type` - The associated service.

* `description` - The field description.

* `data_class_name` - The data class name.

* `subscribed_version` - The subscribed version.

* `mapping` - Whether to display in other places other the classification and mapping module.

* `io_type` - The input and output types.

* `field_key` - The field key.

* `extra_json` - The additional JSON.

* `default_value` - The default value.

* `is_built_in` - Whether the field is built in SecMaster.

* `business_id` - The ID of associated service.

* `used_by` - Which services are used by.

* `target_api` - The target API.

* `json_schema` - The JSON mode.

* `field_tooltip` - The tool tip.

* `display_type` - The display type.

* `required` - Whether the field is required.

* `case_sensitive` - Whether the field is case sensitive.

* `editabled` - Whether the field can be edited.

* `visible` - Whether the field is visible.

* `maintainabled` - Whether the field can be maintained.

* `searchabled` - Whether the field is searchable mode.

* `read_only` - Whether the field is read-only.

* `creatabled` - Whether the field can be created.

* `created_at` - The create time.

* `updated_at` - The update time.

* `creator` - The creator.

* `modifier` - The modifier.

* `creator_id` - The creator ID.

* `modifier_id` - The modifier ID.
