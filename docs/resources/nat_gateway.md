---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_gateway"
description: |-
  Manages a gateway resource of the **public** NAT within HuaweiCloud.
---

# huaweicloud_nat_gateway

Manages a gateway resource of the **public** NAT within HuaweiCloud.

## Example Usage

### Creating a postpaid NAT gateway

```hcl
variable "gateway_name" {}
variable "vpc_id" {}
variable "network_id" {}
variable "gateway_specification" {}

resource "huaweicloud_nat_gateway" "test" {
  name        = var.gateway_name
  description = "test for terraform"
  spec        = var.gateway_specification
  vpc_id      = var.vpc_id
  subnet_id   = var.network_id
}
```

### Creating a prepaid NAT gateway

```hcl
variable "gateway_name" {}
variable "vpc_id" {}
variable "network_id" {}
variable "gateway_specification" {}

resource "huaweicloud_nat_gateway" "test" {
  name        = var.gateway_name
  description = "test for terraform"
  spec        = var.gateway_specification
  vpc_id      = var.vpc_id
  subnet_id   = var.network_id

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the NAT gateway is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC to which the NAT gateway belongs.  
  Changing this will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the subnet ID of the downstream interface (the next hop of the
  DVR) of the NAT gateway.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the NAT gateway name.  
  The valid length is limited from `1` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.

* `spec` - (Required, String) Specifies the specification of the NAT gateway. The valid values are as follows:
  + **1**: Small type, which supports up to `10,000` SNAT connections.
  + **2**: Medium type, which supports up to `50,000` SNAT connections.
  + **3**: Large type, which supports up to `200,000` SNAT connections.
  + **4**: Extra-large type, which supports up to `1,000,000` SNAT connections.

* `description` - (Optional, String) Specifies the description of the NAT gateway, which contain maximum of `512`
  characters, and angle brackets (<) and (>) are not allowed.

-> Fields `name`, `spec` and `description` only support editing for pay-per-use billing mode NAT gateways.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the NAT gateway.
  The valid values are as follows:
  + **prePaid**: the yearly/monthly billing mode.
  + **postPaid**: the pay-per-use billing mode.

  Defaults to **postPaid**. Changing this will create a new resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the NAT gateway.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this will create a new resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the NAT gateway.
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this will create a new resource.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled. This parameter is only valid when
  `charging_mode` is set to **prePaid**. Valid values are **true** and **false**. Defaults to **false**.

* `ngport_ip_address` - (Optional, String, ForceNew) Specifies the IP address used for the NG port of the NAT gateway.
  The IP address must be one of the IP addresses of the VPC subnet associated with the NAT gateway.
  If not spacified, it will be automatically allocated.
  Changing this will creates a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of the NAT gateway.  
  Changing this will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the NAT gateway.

* `session_conf` - (Optional, List) Specifies the session configuration of the NAT gateway.
  The [session_conf](#nat_session_conf) structure is documented below.

<a name="nat_session_conf"></a>
The `session_conf` block supports:

* `tcp_session_expire_time` - (Optional, Int) Specifies the TCP session expiration time, in seconds.
  The valid value from `40` to `7,200`, default value is `900`.

* `udp_session_expire_time` - (Optional, Int) Specifies the UDP session expiration time, in seconds.
  The valid value from `40` to `7,200`, default value is `300`.

* `icmp_session_expire_time` - (Optional, Int) Specifies the ICMP session expiration time, in seconds.
  The valid value from `10` to `7,200`, default value is `10`.

* `tcp_time_wait_time` - (Optional, Int) Specifies the duration of TIME_WAIT state when TCP connection is closed,
  in seconds. The valid value from `0` to `1,800`, default value is `5`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The current status of the NAT gateway.

* `created_at` - The creation time of the NAT gateway.

* `billing_info` - The order information of the NAT gateway.
  When the `charging_mode` is set to **prePaid**, this parameter is available.

* `dnat_rules_limit` - The maximum number of DNAT rules on the NAT gateway. Defaults to `200`.

* `snat_rule_public_ip_limit` - The maximum number of SNAT rules on the NAT gateway. Defaults to `20`.

* `pps_max` - The number of packets that the NAT gateway can receive or send per second.

* `bps_max` - The bandwidth that the NAT gateway can receive or send per second, unit is MB.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

NAT gateways can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_nat_gateway.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `charging_mode`, `period_unit`,
`period` and `auto_renew`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_nat_gateway" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      charging_mode, period_unit, period, auto_renew,
    ]
  }
}
```
