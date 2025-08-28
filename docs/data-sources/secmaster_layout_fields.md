---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_layout_fields"
description: |-
  Use this data source to get the list of layout fields.
---

# huaweicloud_secmaster_layout_fields

Use this data source to get the list of layout fields.

## Example Usage

```hcl
variable "workspace_id" {}
variable "business_code" {}

data "huaweicloud_secmaster_layout_fields" "test" {
  workspace_id  = var.workspace_id
  business_code = var.business_code
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `business_code` - (Required, String) Specifies the data class business code.
  The value can be **Incident**, **Alert**, **Resource**, **Indicator** or **Vulnerability**.

* `name` - (Optional, String) Specifies the field name

* `is_built_in` - (Optional, String) Specifies whether the field is built-in.
  The value can be **true** or **false**.

* `layout_id` - (Optional, String) Specifies the layout ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `fields` - The list of the layout fields.

  The [fields](#fields_struct) structure is documented below.

<a name="fields_struct"></a>
The `fields` block supports:

* `id` - The field ID.

* `cloud_pack_id` - The subscription package ID.

* `cloud_pack_name` - The subscription package name.

* `dataclass_id` - The data class ID.

* `cloud_pack_version` - The subscription package version.

* `field_key` - The field key.

* `name` - The field name.

* `description` - The field description.

* `en_description` - The field English description.

* `default_value` - The default value.

* `en_default_value` - The field English default value.

* `field_type` - The field type.

* `extra_json` - The additional JSON data.

* `field_tooltip` - The field tooltip.

* `en_field_tooltip` - The English field tooltip.

* `json_schema` - The JSON mode.

* `is_built_in` - Whether the field is built-in.

* `read_only` - Whether the field is read-only.

* `required` - Whether the field is required.

* `searchable` - Whether the field is searchable.

* `visible` - Whether the field is visible.

* `maintainable` - Whether the field is maintainable.

* `editable` - Whether the field is editable.

* `creatable` - Whether the field is creatable.

* `creator_id` - The creator ID.

* `creator_name` - The creator name.

* `modifier_id` - The modifier ID.

* `modifier_name` - The modifier name.

* `create_time` - The create time.

* `update_time` - The update time.

* `wizard_id` - The wizard ID bound to the field.

* `aopworkflow_id` - The workflow ID bound to the field.

* `aopworkflow_version_id` - The workflow version ID bound to the field.

* `playbook_id` - The playbook ID bound to the field.

* `playbook_version_id` - The playbook version ID bound to the field.

* `boa_version` - The BOA base version.

* `version` - The SecMaster version.
