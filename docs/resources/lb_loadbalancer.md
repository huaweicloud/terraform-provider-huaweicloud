---
subcategory: "Elastic Load Balance (ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lb_loadbalancer"
description: ""
---

# huaweicloud_lb_loadbalancer

Manages an ELB loadbalancer resource within HuaweiCloud.

## Example Usage

### Basic Loadbalancer

```hcl
variable "ipv4_subnet_id" {}

resource "huaweicloud_lb_loadbalancer" "lb_1" {
  vip_subnet_id = var.ipv4_subnet_id

  tags = {
    key = "value"
  }
}
```

### Loadbalancer With EIP

```hcl
variable "ipv4_subnet_id" {}

resource "huaweicloud_lb_loadbalancer" "lb_1" {
  vip_subnet_id = var.ipv4_subnet_id
}

resource "huaweicloud_vpc_eip_associate" "eip_1" {
  public_ip = "1.2.3.4"
  port_id   = huaweicloud_lb_loadbalancer.lb_1.vip_port_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the loadbalancer resource. If omitted, the
  provider-level region will be used. Changing this creates a new loadbalancer.

* `name` - (Optional, String) Human-readable name for the loadbalancer. Does not have to be unique.

* `description` - (Optional, String) Human-readable description for the loadbalancer.

* `vip_subnet_id` - (Required, String, ForceNew) The **IPv4 subnet ID** of the subnet where the load balancer works.
  Changing this creates a new loadbalancer.

* `vip_address` - (Optional, String, ForceNew) The ip address of the load balancer. Changing this creates a new
  loadbalancer.

* `tags` - (Optional, Map) The key/value pairs to associate with the loadbalancer.

* `enterprise_project_id` - (Optional, String) The enterprise project id of the loadbalancer.

* `protection_status` - (Optional, String) Specifies whether modification protection is enabled. Value options:
  + **nonProtection**: No protection.
  + **consoleProtection**: Console modification protection.

  Defaults to **nonProtection**.

* `protection_reason` - (Optional, String) Specifies the reason to enable modification protection. Only valid when
  `protection_status` is **consoleProtection**.

* `charging_mode` - (Optional, String) Specifies the charging mode of the loadbalancer.  
  The valid values are **prePaid** and **postPaid**, defaults to **postPaid**.

* `period_unit` - (Optional, String) Specifies the charging period unit of the loadbalancer.  
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `period` - (Optional, Int) Specifies the charging period of the loadbalancer.
  + If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  + If `period_unit` is set to **year**, the value ranges from `1` to `3`.

  This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled.  
  Valid values are **true** and **false**. Defaults to **false**.

-> **NOTE:** `period_unit`, `period` and `auto_renew` can only be updated when changing to **prePaid** billing mode.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `vip_port_id` - The Port ID of the Load Balancer IP.

* `public_ip` - The EIP address that is associated to the Load Balancer instance.

* `charge_mode` - Indicates how the load balancer will be billed.

* `frozen_scene` - Indicates the scenario where the load balancer is frozen.

* `created_at` - The create time of the load balancer.

* `updated_at` - The update time of the load balancer.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

Load balancers can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_lb_loadbalancer.test 3e3632db-36c6-4b28-a92e-e72e6562daa6
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `period_unit`, `period`, `auto_renew`.
It is generally recommended running `terraform plan` after importing a loadbalancer.
You can then decide if changes should be applied to the loadbalancer, or the resource definition should be updated to
align with the loadbalancer. Also you can ignore changes as below.

```hcl
resource "huaweicloud_lb_loadbalancer" "test" {
  ...

  lifecycle {
    ignore_changes = [
      period_unit, period, auto_renew,
    ]
  }
}
```
