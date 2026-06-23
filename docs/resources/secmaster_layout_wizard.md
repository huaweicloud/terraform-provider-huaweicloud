---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_layout_wizard"
description: |-
  Manages a SecMaster layout wizard resource within HuaweiCloud.
---

# huaweicloud_secmaster_layout_wizard

Manages a SecMaster layout wizard resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "layout_id" {}
variable "name" {}

resource "huaweicloud_secmaster_layout_wizard" "test" {
  workspace_id = var.workspace_id
  layout_id    = var.layout_id
  name         = var.name
  description  = "test description"
  boa_version  = "v3"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.

* `layout_id` - (Required, String, NonUpdatable) Specifies the layout ID.

* `name` - (Required, String) Specifies the wizard name.

* `description` - (Optional, String) Specifies the description.

* `wizard_json` - (Optional, String) Specifies the layout wizard information.

* `is_binding` - (Optional, String) Specifies whether the button is bound.  
  The valid values are as follows:
  + **true**
  + **false**

* `binding_button` - (Optional, List) Specifies the binding buttons.  
  The [binding_button](#layout_wizard_binding_button) block is supported.

* `boa_version` - (Optional, String) Specifies the BOA version.

<a name="layout_wizard_binding_button"></a>
The `binding_button` block supports:

* `button_id` - (Required, String) Specifies the button ID.

* `button_name` - (Optional, String) Specifies the button name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (the wizard ID).

* `create_time` - The creation time.

* `creator_id` - The creator ID.

* `en_description` - The English description.

* `en_name` - The English name.

* `update_time` - The update time.

* `is_built_in` - Whether the wizard is built-in.

* `version` - The SecMaster version.

## Import

The layout wizard can be imported using the `workspace_id` and `id`, separated by slashes (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_layout_wizard.test <workspace_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `layout_id`.
It is generally recommended running `terraform plan` after importing a layout wizard.
You can then decide if changes should be applied to the layout wizard, or the resource definition should be updated to
align with the layout wizard. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_secmaster_layout_wizard" "test" {
  ...

  lifecycle {
    ignore_changes = [
      layout_id,
    ]
  }
}
```
