---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_terminal_binding"
description: ""
---

# huaweicloud_workspace_terminal_binding

Manages the terminal bindings between MAC addresses and desktops within HuaweiCloud.

-> Only one resource can be created in a region.

## Example Usage

### Allow a machine to remote the desktops

```hcl
variable "desktop_name" {}

resource "huaweicloud_workspace_terminal_binding" "test" {
  enabled               = true
  disabled_after_delete = true

  bindings {
    desktop_name = var.desktop_name
    mac          = "FA-16-3E-E2-3A-1D"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the desktops (to be bound to the MAC address) are located.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `bindings` - (Required, List) Specifies the terminal bindings configuration between MAC addresses and desktops.
  The [blacklist](#terminal_bindings_args) structure is documented below.

* `enabled` - (Optional, Bool) Specifies whether bindings are available.

* `disabled_after_delete` - (Optional, Bool) Specifies whether disabled the binding function before destroy resource.
  Defaults to **true**.

<a name="terminal_bindings_args"></a>
The `bindings` block supports:

* `mac` - (Required, String) Specifies the MAC address.

* `desktop_name` - (Required, String) Specifies the desktop name.

* `description` - (Optional, String) Specifies the binding description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `bindings` - The terminal bindings configuration between MAC addresses and desktops.
  The [blacklist](#terminal_bindings_attr) structure is documented below.

<a name="terminal_bindings_attr"></a>
The `bindings` block supports:

* `id` - The ID of the binding policy.

## Import

Bindings can be imported using the resource `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_terminal_binding.test <id>
```

Also you can using any UUID string to replace this ID in the import phase.

Note that the imported state may not be identical to your resource definition, because of parameter
`disabled_after_delete` is not a remote parameter.  
It is generally recommended running `terraform plan` after importing resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated.
