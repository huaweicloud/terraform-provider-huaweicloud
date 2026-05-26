---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_minor_version_upgrade"
description: |-
  Use this resource to upgrade the minor version or proxy version of a DCS instance within HuaweiCloud.
---

# huaweicloud_dcs_instance_minor_version_upgrade

Use this resource to upgrade the minor version or proxy version of a DCS instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dcs_instance_minor_version_upgrade" "test" {
  instance_id          = var.instance_id
  proxy_minor_version  = "latest"
  engine_minor_version = "latest"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to upgrade the instance.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance to be upgraded.

* `engine_minor_version` - (Optional, String, NonUpdatable) Specifies the target engine minor version.
  Set to **latest** to upgrade to the latest version.

* `proxy_minor_version` - (Optional, String, NonUpdatable) Specifies the target proxy node minor version.
  Set to **latest** to upgrade to the latest version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
