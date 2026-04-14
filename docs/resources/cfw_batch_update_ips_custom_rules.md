---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_batch_update_ips_custom_rules"
description: |-
  Manages a resource to batch update IPS custom rules action within HuaweiCloud.
---

# huaweicloud_cfw_batch_update_ips_custom_rules

Manages a resource to batch update IPS custom rules action within HuaweiCloud.

-> This resource is a one-time action resource used to batch update IPS custom rule actions.
Deleting this resource will not revert the actions of the IPS rules, but will only remove the resource
information from the tfstate file.

## Example Usage

```hcl
variable "fw_instance_id" {} 
variable "action_type" {} 
variable "ips_ids" { 
  type = list(string) 
}

resource "huaweicloud_cfw_batch_update_ips_custom_rules" "test" {
  fw_instance_id = var.fw_instance_id 
  action_type    = var.action_type 
  ips_ids        = var.ips_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall instance ID.

* `action_type` - (Required, Int, NonUpdatable) Specifies the action type for IPS custom rules.
  The values are as follows:
  + **0**: Log only.
  + **1**: Reset/Intercept.

* `ips_ids` - (Required, List, NonUpdatable) Specifies the IDs of IPS custom rules.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `fw_instance_id`.
