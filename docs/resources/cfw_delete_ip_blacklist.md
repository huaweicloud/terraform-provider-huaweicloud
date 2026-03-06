---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_delete_ip_blacklist"
description: |-
  Manages a resource to delete the imported IP blacklist within HuaweiCloud.
---

# huaweicloud_cfw_delete_ip_blacklist

Manages a resource to delete the imported IP blacklist within HuaweiCloud.

-> This resource is a one-time action resource used to delete the imported IP blacklist. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "effect_scope" {
  type = list(int)
}

resource "huaweicloud_cfw_delete_ip_blacklist" "test" {
  fw_instance_id = var.fw_instance_id
  effect_scope   = var.effect_scope
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall instance ID.

* `effect_scope` - (Optional, List, NonUpdatable) Specifies the effect scope.  
  The valid values are as follows:
  + **1**: Specify the effective scope as EIP for deletion.
  + **2**: Specify the effective range as NAT for deletion.
  + **1,2**: Effective scope for deleting EIP and NAT.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `fw_instance_id`.
