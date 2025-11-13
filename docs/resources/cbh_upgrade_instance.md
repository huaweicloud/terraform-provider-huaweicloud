---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_upgrade_instance"
description: |-
  Manages a CBH upgrade instance resource within HuaweiCloud.
---

# huaweicloud_cbh_upgrade_instance

Manages a CBH upgrade instance resource within HuaweiCloud.

-> This resource is only a one-time action resource to upgrade a CBH instance. Deleting this resource
will not clear the corresponding upgrade instance, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "server_id" {}
variable "upgrade_time" {}

resource "huaweicloud_cbh_upgrade_instance" "test" {
  server_id    = var.server_id
  upgrade_time = var.upgrade_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to upgrade a CBH instance.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `server_id` - (Required, String, NonUpdatable) Specifies the ID of the CBH instance to upgrade.

* `upgrade_time` - (Required, Int, NonUpdatable) Specifies the upgrade time of the CBH instance, which must be `24` hours
  later than the current time. The unit is milliseconds.

* `cancel` - (Optional, String, NonUpdatable) Specifies whether to cancel the scheduled upgrade task. Once a task starts,
  it cannot be canceled. **true**: Cancel. **false**: No impact.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as `server_id`).
