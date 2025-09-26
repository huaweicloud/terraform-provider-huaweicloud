---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_loadbalancer"
description: |-
  Manages a Dedicated load balancer resource within HuaweiCloud.
---

# huaweicloud_elb_loadbalancer

Manages a Dedicated load balancer resource within HuaweiCloud.

## Example Usage

### Basic Loadbalancer

```hcl
variable "vpc_id" {}
variable "ipv4_subnet_id" {}
variable "l4_flavor_id" {}
variable "l7_flavor_id" {}
variable "eps_id" {}

resource "huaweicloud_elb_loadbalancer" "basic" {
  name        = "basic"
  description = "basic example"

  vpc_id         = var.vpc_id
  ipv4_subnet_id = var.ipv4_subnet_id

  l4_flavor_id = var.l4_flavor_id
  l7_flavor_id = var.l7_flavor_id

  availability_zone = [
    "cn-north-4a",
    "cn-north-4b",
  ]

  enterprise_project_id = var.eps_id
}
```

### Loadbalancer With Existing EIP

```hcl
variable "vpc_id" {}
variable "ipv4_subnet_id" {}
variable "ipv6_network_id" {}
variable "ipv6_bandwidth_id" {}
variable "l4_flavor_id" {}
variable "l7_flavor_id" {}
variable "eps_id" {}
variable "eip_id" {}

resource "huaweicloud_elb_loadbalancer" "basic" {
  name        = "basic"
  description = "basic example"

  vpc_id            = var.vpc_id
  ipv6_network_id   = var.ipv6_network_id
  ipv6_bandwidth_id = var.ipv6_bandwidth_id
  ipv4_subnet_id    = var.ipv4_subnet_id

  l4_flavor_id = var.l4_flavor_id
  l7_flavor_id = var.l7_flavor_id

  availability_zone = [
    "cn-north-4a",
    "cn-north-4b",
  ]

  enterprise_project_id = var.eps_id

  ipv4_eip_id = var.eip_id
}
```

### Loadbalancer With EIP

```hcl
variable "vpc_id" {}
variable "ipv4_subnet_id" {}
variable "ipv6_network_id" {}
variable "ipv6_bandwidth_id" {}
variable "l4_flavor_id" {}
variable "l7_flavor_id" {}

resource "huaweicloud_elb_loadbalancer" "basic" {
  name        = "basic"
  description = "basic example"

  vpc_id            = var.vpc_id
  ipv6_network_id   = var.ipv6_network_id
  ipv6_bandwidth_id = var.ipv6_bandwidth_id
  ipv4_subnet_id    = var.ipv4_subnet_id

  l4_flavor_id = var.l4_flavor_id
  l7_flavor_id = var.l7_flavor_id

  availability_zone = [
    "cn-north-4a",
    "cn-north-4b",
  ]
}
```

### Loadbalancer With gateway

```hcl
variable "vpc_id" {}
variable "ipv4_subnet_id" {}
variable "ipv6_network_id" {}

resource "huaweicloud_elb_loadbalancer" "basic" {
  name              = "basic"
  description       = "basic example"
  loadbalancer_type = "gateway"

  vpc_id          = var.vpc_id
  ipv4_subnet_id  = var.ipv4_subnet_id
  ipv6_network_id = var.ipv6_network_id

  availability_zone = ["cn-north-4a"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the load balancer resource. If omitted, the
  provider-level region will be used. Changing this creates a new load balancer.

* `name` - (Required, String) Human-readable name for the load balancer.

* `availability_zone` - (Required, List) Specifies the list of AZ names.

  -> **NOTE:** Removing an AZ may disconnect existing connections. Exercise caution when performing this
  operation.

* `loadbalancer_type` - (Optional, String, ForceNew) Specifies the type of the load balancer. Value options:
  + **gateway**: indicates a gateway load balancer.
  + Keep empty(default) indicates other types of load balancers.

  -> **NOTE:** 1. `vip_address` and `ipv6_vip_address` are not supported by gateway load balancers.
  <br/> 2. `vip_subnet_cidr_id` and `ipv6_subnet_cidr_id` cannot be left blank at the same time.
  <br/> 3. You cannot bind an EIP to gateway load balancers.

* `description` - (Optional, String) Human-readable description for the load balancer.

* `cross_vpc_backend` - (Optional, Bool) Enable this if you want to associate the IP addresses of backend servers with
  your load balancer. Can only be true when updating. Defaults to **false**.

* `vpc_id` - (Optional, String, ForceNew) The vpc on which to create the load balancer. Changing this creates a new
  load balancer.

* `ipv4_subnet_id` - (Optional, String) The **IPv4 subnet ID** of the subnet on which to allocate the load balancer
  ipv4 address.

* `ipv6_network_id` - (Optional, String) The **ID** of the subnet on which to allocate the load balancer ipv6 address.

* `ipv6_bandwidth_id` - (Optional, String) The ipv6 bandwidth id. Only support shared bandwidth.

* `ipv4_address` - (Optional, String) The ipv4 address of the load balancer.

* `ipv6_address` - (Optional, String) The ipv6 address of the Load Balancer.

* `ipv4_eip_id` - (Optional, String, ForceNew) The ID of the EIP. Changing this parameter will create a new resource.

  -> **NOTE:** If the ipv4_eip_id parameter is configured, you do not need to configure the bandwidth parameters:
  `iptype`, `bandwidth_charge_mode`, `bandwidth_size`, `share_type` and `bandwidth_id`.

* `iptype` - (Optional, String, ForceNew) Elastic IP type. Changing this parameter will create a new resource.

* `bandwidth_charge_mode` - (Optional, String, ForceNew) Bandwidth billing type. Value options:
  + **bandwidth**: Billed by bandwidth.
  + **traffic**: Billed by traffic.
  
  It is mandatory when `iptype` is set and `bandwidth_id` is empty.
  Changing this parameter will create a new resource.

* `sharetype` - (Optional, String, ForceNew) Bandwidth sharing type. Value options:
  + **PER**: Dedicated bandwidth.
  + **WHOLE**: Shared bandwidth.
  
  It is mandatory when `iptype` is set and `bandwidth_id` is empty.
  Changing this parameter will create a new resource.

* `bandwidth_size` - (Optional, Int, ForceNew) Bandwidth size. It is mandatory when `iptype` is set and `bandwidth_id`
  is empty. Changing this parameter will create a new resource.

* `bandwidth_id` - (Optional, String, ForceNew) Bandwidth ID of the shared bandwidth. It is mandatory when `sharetype`
  is **WHOLE**. Changing this parameter will create a new resource.

  -> **NOTE:** If the `bandwidth_id` parameter is configured, you can not configure the parameters:
  `bandwidth_charge_mode`, `bandwidth_size`.

* `l4_flavor_id` - (Optional, String) The L4 flavor id of the load balancer.

* `l7_flavor_id` - (Optional, String) The L7 flavor id of the load balancer.

* `backend_subnets` - (Optional, List) The IDs of subnets on the downstream plane.
  + If this parameter is not specified, select subnets as follows:
      - If IPv6 is enabled for a load balancer, the ID of subnet specified in `ipv6_network_id` will be used.
      - If IPv4 is enabled for a load balancer, the ID of subnet specified in `ipv4_subnet_id` will be used.
      - If only public network is available for a load balancer, the ID of any subnet in the VPC where the load balancer
        resides will be used. Subnets with more IP addresses are preferred.
  + If there is more than one subnet, the first subnet in the list will be used, and the subnets must be in the VPC
    where the load balancer resides.

* `protection_status` - (Optional, String) The protection status for update. Value options:
  + **nonProtection**: No protection.
  + **consoleProtection**: Console modification protection.

  Defaults to **nonProtection**.

* `protection_reason` - (Optional, String) The reason for update protection. Only valid when `protection_status` is
  **consoleProtection**.

* `tags` - (Optional, Map) The key/value pairs to associate with the load balancer.

* `enterprise_project_id` - (Optional, String) The enterprise project id of the load balancer.

* `charging_mode` - (Optional, String) Specifies the charging mode of the ELB load balancer.
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.

* `period_unit` - (Optional, String) Specifies the charging period unit of the ELB load balancer.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `period` - (Optional, Int) Specifies the charging period of the ELB load balancer.
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled. Valid values are **true** and **false**.

-> **NOTE:** `period_unit`, `period` and `auto_renew` can only be updated when `charging_mode` changed to **prePaid**
  billing mode.

* `force_delete` - (Optional, Bool) Specifies whether to forcibly delete the load balancer, remove the load balancer,
  listeners, unbind associated pools. Defaults to **false**.

* `deletion_protection_enable` - (Optional, Bool) Specifies whether to enable deletion protection
  for the load balancer. Value options:
  + **true**: Enable deletion protection.
  + **false**: Disable deletion protection.

* `waf_failure_action` - (Optional, String) Specifies traffic distributing policies when the WAF is faulty.
  Value options:
  + **discard**: Traffic will not be distributed.
  + **forward**: Traffic will be distributed to the default backend servers.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the load balancer.

* `gw_flavor_id` - The flavor ID of the gateway load balancer.

* `ipv4_port_id` - The ID of the port bound to the private IPv4 address of the load balancer.

* `ipv4_eip` - The ipv4 eip address of the load balancer.

* `ipv6_eip` - The ipv6 eip address of the load balancer.

* `ipv6_eip_id` - The ipv6 eip id of the load balancer.

* `ipv6_eip_id` - The type of the subnet on the downstream plane. The value can be:
  + **ipv4**: IPv4 subnet
  + **dualstack**: subnet that supports IPv4/IPv6 dual stack

* `elb_virsubnet_type` - The type of the subnet on the downstream plane. The value can be:
  + **ipv4**: IPv4 subnet
  + **dualstack**: subnet that supports IPv4/IPv6 dual stack

* `frozen_scene` - The scenario where the load balancer is frozen. Multiple values are separated using commas (,).
  The value can be:
  + **POLICE**: The load balancer is frozen due to security reasons.
  + **ILLEGAL**: The load balancer is frozen due to violation of laws and regulations.
  + **VERIFY**: Your account has not completed real-name authentication.
  + **PARTNER**: The load balancer is frozen by the partner.
  + **ARREAR**: Your account is in arrears.

* `operating_status` - The operating status of the load balancer. The value can be:
  + **ONLINE**: indicates that the load balancer is running normally.
  + **FROZEN**: indicates that the load balancer is frozen.

* `public_border_group` - The AZ group to which the load balancer belongs.

* `charge_mode` - Indicates the billing mode. The value can be one of the following:
  + **flavor**: Billed by the specifications you will select.
  + **lcu**: Billed by LCU usage.

* `guaranteed` - Indicates whether the load balancer is a dedicated load balancer.
  The value can be one of the following:
  + **false**: The load balancer is a shared load balancer.
  + **true**: The load balancer is a dedicated load balancer.

* `created_at` - Indicates the time when the load balancer was created, in RFC3339 format.

* `updated_at` - Indicates the time when the load balancer was updated, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

ELB load balancer can be imported using the ID, e.g.

```bash
$ terraform import huaweicloud_elb_loadbalancer.loadbalancer_1 <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `ipv6_bandwidth_id`, `iptype`,
`bandwidth_charge_mode`, `sharetype`,  `bandwidth_size`, `bandwidth_id`, `force_delete`
and `deletion_protection_enable`. It is generally recommended running `terraform plan` after importing a load balancer.
You can then decide if changes should be applied to the load balancer, or the resource
definition should be updated to align with the load balancer. Also you can ignore changes as below.

```hcl
resource "huaweicloud_elb_loadbalancer" "loadbalancer_1" {
    ...
  lifecycle {
    ignore_changes = [
      ipv6_bandwidth_id, iptype, bandwidth_charge_mode, sharetype, bandwidth_size, bandwidth_id, force_delete,
      deletion_protection_enable,
    ]
  }
}
```
