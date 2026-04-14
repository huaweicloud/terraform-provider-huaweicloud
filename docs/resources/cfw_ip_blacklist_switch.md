---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_ip_blacklist_switch"
description: |-
  Manages a resource to enable or disable the IP blacklist feature within HuaweiCloud.
---

# huaweicloud_cfw_ip_blacklist_switch

Manages a resource to enable or disable the IP blacklist feature within HuaweiCloud.

-> 1. This resource is a one-time action resource used to enable or disable the IP blacklist feature for traffic
  filtering. Deleting this resource will not restore the previous switch state on the cloud, but will only remove the
  resource information from the tf state file.

## Example Usage

```hcl
variable "fw_instance_id" {}

resource "huaweicloud_cfw_ip_blacklist_switch" "test" {
  fw_instance_id = var.fw_instance_id
  status         = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall instance ID.

* `status` - (Required, Int, NonUpdatable) Specifies whether to enable the IP blacklist feature.  
  The valid values are as follows:
  + **0**: Disable.
  + **1**: Enable.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `fw_instance_id`.
