---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_proxy_upgrade"
description: |
  Use this resource to upgrade the kernel version of a TaurusDB proxy instance within HuaweiCloud.
---

# huaweicloud_taurusdb_proxy_upgrade

Use this resource to upgrade the kernel version of a TaurusDB proxy instance within HuaweiCloud.

-> This resource is a one-time action resource for upgrading kernel version of a TaurusDB proxy instance. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from
the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "proxy_id" {}
variable "source_version" {}
variable "target_version" {}

resource "huaweicloud_taurusdb_proxy_upgrade" "test" {
  instance_id    = var.instance_id
  proxy_id       = var.proxy_id
  source_version = var.source_version
  target_version = var.target_version
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the provider-level
  region will be used.

* `instance_id` - (Required, String, NoneUpdatable) Specifies the ID of the TaurusDB instance.

* `proxy_id` - (Required, String, NoneUpdatable) Specifies the ID of the proxy instance.

* `source_version` - (Required, String, NoneUpdatable) Specifies the source kernel version of the proxy.

* `target_version` - (Required, String, NoneUpdatable) Specifies the target kernel version of the proxy.
  The target kernel version must be later than or equal to the source kernel version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
