---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_eip"
description: ""
---

# huaweicloud_vpc_eip

Manages an EIP resource within HuaweiCloud.

## Example Usage

### Create an EIP with Dedicated Bandwidth

```hcl
var "bandwidth_name" {}

resource "huaweicloud_vpc_eip" "dedicated" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = var.bandwidth_name
    size        = 10
    charge_mode = "traffic"
  }
}
```

### Create an EIP with Shared Bandwidth

```hcl
var "bandwidth_name" {}

resource "huaweicloud_vpc_bandwidth" "test" {
  name = var.bandwidth_name
  size = 5
}

resource "huaweicloud_vpc_eip" "shared" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type = "WHOLE"
    id         = huaweicloud_vpc_bandwidth.test.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the EIP resource.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `publicip` - (Required, List) Specifies the EIP configuration.  
  The [object](#vpc_eip_publicip) structure is documented below.

* `bandwidth` - (Required, List) Specifies the bandwidth configuration.  
  The [object](#vpc_eip_bandwidth) structure is documented below.

* `name` - (Optional, String) Specifies the name of the EIP.  
  The name can contain `1` to `64` characters, including English letters, Chinese characters, digits, underscores (_),
  hyphens (-), and periods (.).

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the EIP belongs.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the EIP.

* `charging_mode` - (Optional, String) Specifies the charging mode of the EIP.  
  The valid values are **prePaid** and **postPaid**, defaults to **postPaid**.

-> **NOTE:** Please update the `charge_mode` of `bandwidth` to **bandwidth** before changing to **prePaid** billing mode.

* `period_unit` - (Optional, String) Specifies the charging period unit of the EIP.  
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `period` - (Optional, Int) Specifies the charging period of the EIP.
  + If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  + If `period_unit` is set to **year**, the value ranges from `1` to `3`.

  This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled.  
  Valid values are **true** and **false**. Defaults to **false**.

-> **NOTE:** `period_unit`, `period` and `auto_renew` can only be updated when changing to **prePaid** billing mode.

<a name="vpc_eip_publicip"></a>
The `publicip` block supports:

* `type` - (Optional, String, ForceNew) Specifies the EIP type. Possible values are **5_bgp** (dynamic BGP)
  and **5_sbgp** (static BGP), the default value is **5_bgp**. Changing this will create a new resource.

* `ip_address` - (Optional, String, ForceNew) Specifies the EIP address to be assigned.  
  The value must be a valid **IPv4** address in the available IP address range.
  The system automatically assigns an EIP if you do not specify it. Changing this will create a new resource.

* `ip_version` - (Optional, Int) Specifies the IP version, either `4` (default) or `6`.

<a name="vpc_eip_bandwidth"></a>
The `bandwidth` block supports:

* `share_type` - (Required, String, ForceNew) Specifies whether the bandwidth is dedicated or shared.  
  Changing this will create a new resource. Possible values are as follows:
  + **PER**: Dedicated bandwidth
  + **WHOLE**: Shared bandwidth

* `name` - (Optional, String) Specifies the bandwidth name.  
  The name can contain `1` to `64` characters, including letters, digits, underscores (_), hyphens (-), and periods (.).
  This parameter is mandatory when `share_type` is set to **PER**.

* `size` - (Optional, Int) The bandwidth size.  
  The value ranges from `1` to `300` Mbit/s. This parameter is mandatory when `share_type` is set to **PER**.

* `id` - (Optional, String, ForceNew) The shared bandwidth ID.  
  This parameter is mandatory when `share_type` is set to **WHOLE**. Changing this will create a new resource.

* `charge_mode` - (Optional, String) Specifies whether the bandwidth is billed by traffic or by bandwidth
  size. The value can be **traffic** or **bandwidth**. If the `charging_mode` is **prePaid**, only **bandwidth** is valid.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `address` - The IPv4 address of the EIP.
* `ipv6_address` - The IPv6 address of the EIP.
* `private_ip` - The private IP address bound to the EIP.
* `port_id` - The port ID which the EIP associated with.
* `status` - The status of EIP.
* `created_at` - The create time of EIP.
* `updated_at` - The update time of EIP.
* `associate_type` - The associate type of EIP. Values are **PORT**, **NATGW**, **ELB**, **ELBV1** and **VPN**.
* `associate_id` - The associate id of EIP.
* `instance_type` - The instance type to which the port belongs. Return when `associate_type` is **PORT**.
* `instance_id` - The instance id to which the port belongs. Return when `associate_type` is **PORT**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

EIPs can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpc_eip.test 2c7f39f3-702b-48d1-940c-b50384177ee1
```
