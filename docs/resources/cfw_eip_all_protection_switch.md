---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_eip_all_protection_switch"
description: |-
  Manages a resource to switch CFW all EIP protection (kill switch / restore) within HuaweiCloud.
---

# huaweicloud_cfw_eip_all_protection_switch

Manages a resource to switch CFW all EIP protection (one-click bypass / restore) within HuaweiCloud.

-> This resource is a one-time action resource used to switch CFW EIP all protection. Deleting this resource will not
  clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "fw_instance_id" {}

resource "huaweicloud_cfw_eip_all_protection_switch" "restore" {
  fw_instance_id     = var.fw_instance_id
  bypass_operation   = 0
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall instance ID.

* `bypass_operation` - (Required, Int, NonUpdatable) Specifies the protection operation type.  
  The valid values are as follows:
  + **1**: One-click turn off protection (bypass).
  + **0**: One-click restore protection.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `fw_instance_id`.

* `object_id` - The protected object ID.

* `fail_reason` - The reason when the bypass or restore operation fails.
