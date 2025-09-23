---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_protected_instance_add_nic"
description: |-
  Using this resource to add a network interface card (NIC) to a protected instance in SDRS within HuaweiCloud.
---

# huaweicloud_sdrs_protected_instance_add_nic

Using this resource to add a network interface card (NIC) to a protected instance in SDRS within HuaweiCloud.

-> This is a one-time action resource to add a NIC to a protected instance. Deleting this resource will
not change the current NIC configuration, but will only remove the resource information from the tfstate file.

-> Using this resource may cause unexpected changes to the ECS security group used to protect the instance.
Before using this resource, use `lifecycle` to ignore unexpected changes to the `security_group_ids` field in
resource `huaweicloud_compute_instance`. The following restrictions apply before using this resource:
<br/>1. Status of the protection group must be **available** or **protected**.
<br/>2. Status of the protected instance must be **available** or **protected**.
<br/>3. The subnet of the NIC to be added must belong to the same VPC of the protected group and protected instance.

## Example Usage

```hcl
variable "protected_instance_id" {}
variable "subnet_id" {}
variable "ip_address" {}
variable "security_groups" {
  type = list(string)
}

resource "huaweicloud_sdrs_protected_instance_add_nic" "test" {
  protected_instance_id = var.protected_instance_id
  subnet_id             = var.subnet_id
  ip_address            = var.ip_address
  security_groups       = var.security_groups
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to execute the request.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `protected_instance_id` - (Required, String, NonUpdatable) Specifies the ID of the protected instance to add the NIC to.

* `subnet_id` - (Required, String, NonUpdatable) Specifies the ID of the subnet to which the NIC will be attached.
  It is network ID of the subnet, which is the same as the `neutron_network_id` value.

* `security_groups` - (Optional, List, NonUpdatable) Specifies the security groups to associate with the NIC.
  Defaults to the system default security group.
  The [security_groups](#security_groups_struct) structure is documented below.

* `ip_address` - (Optional, String, NonUpdatable) Specifies the IP address to assign to the NIC.
  If not specified, an available IP will be automatically assigned.

<a name="security_groups_struct"></a>
The `security_groups` block supports:

* `id` - (Required, String, NonUpdatable) Specifies the ID of the security group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
