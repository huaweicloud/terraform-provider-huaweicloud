---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_layout_wizards"
description: |-
  Use this data source to get the list of layout wizards under a specified layout.
---

# huaweicloud_secmaster_layout_wizards

Use this data source to get the list of layout wizards under a specified layout.

## Example Usage

```hcl
variable "workspace_id" {}
variable "layout_id" {}

data "huaweicloud_secmaster_layout_wizards" "test" {
  workspace_id = var.workspace_id
  layout_id    = var.layout_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `layout_id` - (Required, String) Specifies the layout ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of layout wizards.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The layout wizard ID.

* `name` - The layout wizard name.

* `en_name` - The layout wizard English name.

* `description` - The layout wizard description.

* `en_description` - The layout wizard English description.

* `wizard_json` - The layout wizard related information.

* `creator_id` - The creator ID.

* `create_time` - The creation time.

* `update_time` - The update time.

* `project_id` - The project ID.

* `workspace_id` - The workspace ID.

* `is_binding` - Whether bound the button.

* `binding_button` - The binding button information.
  The [binding_button](#data_binding_button_struct) structure is documented below.

* `is_built_in` - Whether the page is a system page.

* `boa_version` - The BOA base version.

* `version` - The Secmaster version.

<a name="data_binding_button_struct"></a>
The `binding_button` block supports:

* `button_id` - The button ID.

* `button_name` - The button name.
