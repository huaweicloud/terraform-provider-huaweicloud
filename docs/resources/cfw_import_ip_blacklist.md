---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_import_ip_blacklist"
description: |-
  Manages a resource to import the IP blacklist within HuaweiCloud.
---

# huaweicloud_cfw_import_ip_blacklist

Manages a resource to import the IP blacklist within HuaweiCloud.

-> This resource is a one-time action resource used to import the IP blacklist. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

### Incremental Import

```hcl
variable "fw_instance_id" {}
variable "ip_blacklist" {}
variable "effect_scope" {
  type = list(int)
}

resource "huaweicloud_cfw_import_ip_blacklist" "test" {
  fw_instance_id = var.fw_instance_id
  add_type       = 0
  ip_blacklist   = var.ip_blacklist
  effect_scope   = var.effect_scope
}
```

### Full Import

```hcl
variable "fw_instance_id" {}
variable "ip_blacklist" {}
variable "effect_scope" {
  type = list(int)
}

resource "huaweicloud_cfw_import_ip_blacklist" "test" {
  fw_instance_id = var.fw_instance_id
  add_type       = 1
  ip_blacklist   = var.ip_blacklist
  effect_scope   = var.effect_scope
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall instance ID.

* `add_type` - (Required, Int, NonUpdatable) Specifies the import type.  
  The valid values are as follows:
  + **0**: Incremental import (appending to the existing list).
  + **1**: Full import (overwrite the existing one).

* `ip_blacklist` - (Required, String, NonUpdatable) Specifies a list of IP addresses. Currently, different IP addresses
  can be separated by delimiters such as ",", " ", ";, "\r\n", "\n", and "\t".

* `effect_scope` - (Required, List, NonUpdatable) Specifies the effect scope.  
  The valid values are as follows:
  + **1**: The scope of effect is EIP.
  + **2**: The scope of effect is NAT.
  + **1,2**: The scope of effect is EIP and NAT.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `fw_instance_id`.
