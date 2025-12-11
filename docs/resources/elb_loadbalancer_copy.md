---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_loadbalancer_copy"
description: |-
  Manages a Dedicated load balancer copy resource within HuaweiCloud.
---

# huaweicloud_elb_loadbalancer_copy

Manages a Dedicated load balancer copy resource within HuaweiCloud.

## Example Usage

```hcl
variable "loadbalancer_id" {}
variable "ipv4_subnet_id" {}
variable "ipv6_network_id" {}
variable "backend_subnet_id" {}
variable "l4_flavor_id" {}
variable "l7_flavor_id" {}
variable "ipv6_bandwidth_id" {}

resource "huaweicloud_elb_loadbalancer_copy" "test" {
  loadbalancer_id   = var.loadbalancer_id
  name              = "test_elb_name"
  availability_zone = ["cn-north-4a", "cn-north-4b"]
  ipv4_subnet_id    = var.ipv4_subnet_id
  ipv4_address      = "192.168.0.216"
  ipv6_network_id   = var.ipv6_network_id
  ipv6_address      = "2407:c080:1200:2a02:34e6:8059:ce7f:1add"
  backend_subnets   = [var.backend_subnet_id]
  l4_flavor_id      = var.l4_flavor_id
  l7_flavor_id      = var.l7_flavor_id
  reuse_pool        = true

  description                = "test elb description"
  ipv6_bandwidth_id          = var.ipv6_bandwidth_id
  cross_vpc_backend          = "true"
  protection_status          = "consoleProtection"
  protection_reason          = "test protection reason"
  deletion_protection_enable = "true"
  waf_failure_action         = "discard"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the load balancer resource. If omitted, the
  provider-level region will be used. Changing this creates a new load balancer.

* `loadbalancer_id` - (Required, String, NonUpdatable) Specifies the source load balancer ID.

* `name` - (Optional, String) Specifies the load balancer name.

* `availability_zone` - (Optional, List) Specifies the list of AZ names.

  -> **NOTE:** Removing an AZ may disconnect existing connections. Exercise caution when performing this
  operation.

* `ipv4_subnet_id` - (Optional, String) Specifies the ID of the IPv4 subnet where the load balancer works. If it is not
  specified, the IPv4 subnet of the original load balancer is used. The subnets where the original and new load balancers
  work must be in the same VPC.

* `ipv4_address` - (Optional, String) Specifies the private IPv4 address of the load balancer.

* `ipv6_network_id` - (Optional, String) Specifies the ID of the IPv6 subnet where the new load balancer works. If it is
  not specified, the IPv6 subnet of the original load balancer is used. The subnets where the original and new load
  balancers work must be in the same VPC.

* `ipv6_address` - (Optional, String) Specifies the private IPv6 address of the load balancer.

* `backend_subnets` - (Optional, List) Specifies the ID of the backend subnet of the load balancer. If it is not specified,
  the backend subnet of the original load balancer is used. The subnets where the original and new load balancers work must
  be in the same VPC.

* `l4_flavor_id` - (Optional, String) Specifies the Layer 4 specifications of the new load balancer. If it is not specified,
  the Layer 4 specifications of the original load balancer are used.

* `l7_flavor_id` - (Optional, String) Specifies the Layer 7 specifications of the new load balancer. If it is not specified,
  the Layer 7 specifications of the original load balancer are used.

* `enterprise_project_id` - (Optional, String) The enterprise project ID of the load balancer.

* `reuse_pool` - (Optional, String, NonUpdatable) Specifies whether to reuse the backend server group and backend server
  ID of the original load balancer.
  + If it is set to **true**, the backend server group of the original load balancer will be used.
  + If no backend server group is selected, a new backend server group is created by default.
  + It is invalid when `enterprise_project_id` is set to another enterprise project.

* `ipv6_bandwidth_id` - (Optional, String) Specifies the ipv6 bandwidth ID. Only support shared bandwidth.

* `description` - (Optional, String) Specifies the description of the load balancer.

* `cross_vpc_backend` - (Optional, String) Specifies whether to add backend servers that are not in the load balancer's
  VPC. Can only be **true** when updating. Value options: **true**, **false**.

* `protection_status` - (Optional, String) Specifies the protection status for update. Value options:
  + **nonProtection**: No protection.
  + **consoleProtection**: Console modification protection.

* `protection_reason` - (Optional, String) Specifies the reason for update protection. Only valid when `protection_status`
  is **consoleProtection**.

* `deletion_protection_enable` - (Optional, String) Specifies whether to enable deletion protection for the load balancer.
  Value options:
  + **true**: Enable deletion protection.
  + **false**: Disable deletion protection.

* `waf_failure_action` - (Optional, String) Specifies traffic distributing policies when the WAF is faulty.
  Value options:
  + **discard**: Traffic will not be distributed.
  + **forward**: Traffic will be distributed to the default backend servers.

* `force_delete` - (Optional, Bool, NonUpdatable) Specifies whether to forcibly delete the load balancer, remove the load
  balancer, listeners, unbind associated pools. Defaults to **false**.

* `tags` - (Optional, Map) The key/value pairs to associate with the load balancer.

* `charging_mode` - (Optional, String) Specifies the charging mode of the ELB load balancer.
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.

* `period_unit` - (Optional, String) Specifies the charging period unit of the ELB load balancer.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `period` - (Optional, Int) Specifies the charging period of the ELB load balancer.
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled. Valid values are **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the load balancer.

* `loadbalancer_type` - Indicates the type of the load balancer.

* `vpc_id` - Indicates the ID of the VPC where the load balancer resides.

* `elb_virsubnet_type` - Indicates the type of the subnet on the downstream plane. The value can be:
  + **ipv4**: IPv4 subnet
  + **dualstack**: subnet that supports IPv4/IPv6 dual stack

* `frozen_scene` - Indicates the scenario where the load balancer is frozen. Multiple values are separated using
  commas (,).
  The value can be:
  + **POLICE**: The load balancer is frozen due to security reasons.
  + **ILLEGAL**: The load balancer is frozen due to violation of laws and regulations.
  + **VERIFY**: Your account has not completed real-name authentication.
  + **PARTNER**: The load balancer is frozen by the partner.
  + **ARREAR**: Your account is in arrears.

* `operating_status` - Indicates the operating status of the load balancer. The value can be:
  + **ONLINE**: indicates that the load balancer is running normally.
  + **FROZEN**: indicates that the load balancer is frozen.

* `public_border_group` - Indicates the AZ group to which the load balancer belongs.

* `charge_mode` - Indicates the billing mode. The value can be:
  + **flavor**: Billed by the specifications you will select.
  + **lcu**: Billed by LCU usage.

* `ipv4_port_id` - Indicates the ID of the port bound to the private IPv4 address of the load balancer.

* `gw_flavor_id` - Indicates the flavor ID of the gateway load balancer.

* `created_at` - Indicates the time when the load balancer was created, in RFC3339 format.

* `updated_at` - Indicates the time when the load balancer was updated, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

The ELB load balancer copy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_loadbalancer_copy.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `loadbalancer_id`, `ipv6_bandwidth_id`,
`deletion_protection_enable`, `reuse_pool`,  `period_unit`, `period`, `auto_renew` and `force_delete`. It is generally
recommended running `terraform plan` after importing a load balancer copy. You can then decide if changes should be applied
to the load balancer copy, or the resource definition should be updated to align with the load balancer. Also you can
ignore changes as below.

```hcl
resource "huaweicloud_elb_loadbalancer_copy" "test" {
    ...
  lifecycle {
    ignore_changes = [
      loadbalancer_id, ipv6_bandwidth_id, deletion_protection_enable, reuse_pool, period_unit, period, auto_renew,
      force_delete,
    ]
  }
}
```
