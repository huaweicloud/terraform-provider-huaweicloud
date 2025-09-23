---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_soc_component_actions"
description: |-
  Use this data source to get the list of SecMaster soc component actions within HuaweiCloud.
---

# huaweicloud_secmaster_soc_component_actions

Use this data source to get the list of SecMaster soc component actions within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "component_id" {}

data "huaweicloud_secmaster_soc_component_actions" "test" {
  workspace_id = var.workspace_id
  component_id = var.component_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `component_id` - (Required, String) Specifies the component ID.

* `enabled` - (Optional, Bool) Specifies whether to enable or not. The value can be **true** or **false**.
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The soc component actions list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The plugin execution function ID.

* `action_name` - The plugin execution function name.

* `action_desc` - The plugin execution function description.

* `action_type` - The action type. The value can be **action**, **connector**, or **manager**.

* `create_time` - The creation time.

* `creator_name` - The creator name.

* `can_update` - Is it upgradable. The value can be **true** or **false**.

* `action_version_id` - The plugin execution function version ID.

* `action_version_name` - The user defined action version alias.

* `action_version_number` - The internally generated action version number will only increment.

* `action_enable` - The enable or disable status.
