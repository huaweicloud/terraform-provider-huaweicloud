---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_export_ip_blacklist"
description: |-
  Manages a resource to export the IP blacklist within HuaweiCloud.
---

# huaweicloud_cfw_export_ip_blacklist

Manages a resource to export the IP blacklist within HuaweiCloud.

-> 1. This resource is a one-time action resource used to export the IP blacklist. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tf state file.
  <br/>2. Before using this resource, it is necessary to ensure that the IP blacklist has been imported into the current
  firewall instance, otherwise executing the resource will result in an error.

## Example Usage

```hcl
variable "fw_instance_id" {}

resource "huaweicloud_cfw_export_ip_blacklist" "test" {
  fw_instance_id = var.fw_instance_id
  name           = "ip-blacklist-eip.txt"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall instance ID.

* `name` - (Required, String, NonUpdatable) Specifies the IP blacklist name.  
  The valid values are as follows:
  + **ip-blacklist-eip.txt**: Export IP blacklist with effect scope EIP.
  + **ip-blacklist-nat.txt**: Export IP blacklist with effect scope NAT.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `fw_instance_id`.

* `data` - The exported IP blacklist data.
