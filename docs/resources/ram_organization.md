---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_organization"
description: "Manages a RAM organization resource within HuaweiCloud."
---

# huaweicloud_ram_organization

Manages a RAM organization resource within HuaweiCloud.

-> Only organization administrators can enable or disable sharing with the organization. After enabling sharing with
organizations, resources can be shared to organizations and organizational units, and the resource users do not need to
accept the sharing to take effect. Destroying resources does not change the current enabled state.

## Example Usage

```hcl
resource "huaweicloud_ram_organization" "test" {
  enabled = true
} 
```

## Argument Reference

The following arguments are supported:

* `enabled` - (Required, Bool) Specifies whether sharing with organizations is enabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
