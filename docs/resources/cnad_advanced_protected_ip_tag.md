---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_protected_ip_tag"
description: |-
  Manages a CNAD advanced protected IP tag resource within HuaweiCloud.
---

# huaweicloud_cnad_advanced_protected_ip_tag

Manages a CNAD advanced protected IP tag resource within HuaweiCloud.

-> This resource is a one-time action resource for setting tag on protected IP. Deleting this resource will not
remove the tag from the protected IP, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "protected_ip_id" {}
variable "ip_tag" {}

resource "huaweicloud_cnad_advanced_protected_ip_tag" "test" {
  protected_ip_id = var.protected_ip_id
  tag             = var.ip_tag
}
```

## Argument Reference

The following arguments are supported:

* `protected_ip_id` - (Required, String, NonUpdatable) Specifies the ID of the protected IP.

* `tag` - (Required, String, NonUpdatable) Specifies the tag to be set on the protected IP.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
