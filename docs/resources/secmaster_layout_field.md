---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_layout_field"
description: |-
  Manages a SecMaster layout field resource within HuaweiCloud.
---

# huaweicloud_secmaster_layout_field

Manages a SecMaster layout field resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

resource "huaweicloud_secmaster_layout_field" "test" {
  workspace_id = var.workspace_id
  name         = "test-field"
  field_key    = "test_field"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SecMaster layout field resource. If omitted,
  the provider-level region will be used. Changing this will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.

* `name` - (Required, String) Specifies the name of the layout field.

* `field_key` - (Required, String) Specifies the key of the layout field.

* `description` - (Optional, String) Specifies the description of the layout field.

* `layout_id` - (Optional, String) Specifies the layout ID this field belongs to.

* `wizard_id` - (Optional, String) Specifies the ID of the bound page.

* `aopworkflow_id` - (Optional, String) Specifies the ID of the bound workflow.

* `aopworkflow_version_id` - (Optional, String) Specifies the version ID of the bound workflow.

* `playbook_id` - (Optional, String) Specifies the ID of the bound playbook.

* `playbook_version_id` - (Optional, String) Specifies the version ID of the bound playbook.

* `default_value` - (Optional, String) Specifies the default value of the field.

* `field_type` - (Optional, String) Specifies the type of the layout field, such as **shorttext**, **radio**, **grid**, etc.

* `display_type` - (Optional, String) Specifies the display type of the field.

* `extra_json` - (Optional, String) Specifies the additional JSON data.

* `field_tooltip` - (Optional, String) Specifies the tooltip text of the field.

* `json_schema` - (Optional, String) Specifies the JSON schema of the field.

* `readonly` - (Optional, Bool) Specifies whether the field is read-only. Defaults to **false**.

* `required` - (Optional, Bool) Specifies whether the field is required. Defaults to **false**.

* `searchable` - (Optional, Bool) Specifies whether the field is searchable. Defaults to **false**.

* `visible` - (Optional, Bool) Specifies whether the field is visible. Defaults to **false**.

* `maintainable` - (Optional, Bool) Specifies whether the field is maintainable. Defaults to **false**.

* `editable` - (Optional, Bool) Specifies whether the field is editable. Defaults to **false**.

* `creatable` - (Optional, Bool) Specifies whether the field is creatable. Defaults to **false**.

* `boa_version` - (Optional, String) Specifies the BOA base version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also the layout field ID).

* `cloud_pack_id` - The subscription package ID.

* `cloud_pack_name` - The subscription package name.

* `dataclass_id` - The data class ID.

* `cloud_pack_version` - The subscription package version.

* `en_description` - The English description of the field.

* `en_default_value` - The English default value of the field.

* `en_field_tooltip` - The English tooltip text of the field.

* `is_built_in` - Whether the field is built-in.

* `creator_id` - The creator ID.

* `creator_name` - The creator name.

* `modifier_id` - The modifier ID.

* `modifier_name` - The modifier name.

* `create_time` - The creation time.

* `update_time` - The update time.

* `version` - The version of the SecMaster.

## Import

Layout field resource can be imported using the `workspace_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_secmaster_layout_field.test <workspace_id>/<id>
```
