---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_eip_alarm_whitelist"
description: |-
  Manages a CFW EIP alarm whitelist resource within HuaweiCloud.
---

# huaweicloud_cfw_eip_alarm_whitelist

Manages a CFW EIP alarm whitelist resource within HuaweiCloud.

-> Currently, this resource does not support the delete function. Deleting this resource will not clear the
  corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "eip_id" {}
variable "public_ip" {}

resource "huaweicloud_cfw_eip_alarm_whitelist" "test" {
  fw_instance_id = var.fw_instance_id
  eip_id         = var.eip_id
  public_ip      = var.public_ip
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall instance ID.

* `eip_id` - (Required, String, NonUpdatable) Specifies the EIP ID.  

* `public_ip` - (Required, String, NonUpdatable) Specifies the IPv4 address.  

* `object_id` - (Optional, String, NonUpdatable) Specifies the protected object ID.

* `public_ipv6` - (Optional, String, NonUpdatable) Specifies the IPv6 address.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `fw_instance_id`.

* `device_name` - The device name.

* `type` -  The EIP whitelist flag.  
  The valid values are as follows:
  + **1**: The EIP is in the whitelist.
  + **0**: The EIP is not in the whitelist.

## Import

The EIP alarm whitelist resource can be imported using `id` and `public_ip`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cfw_eip_alarm_whitelist.test <id>/<public_ip>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `object_id`, `enterprise_project_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_cfw_eip_alarm_whitelist" "test" { 
  ...
  
  lifecycle {
    ignore_changes = [
      object_id, enterprise_project_id,
    ]
  }
}
```
