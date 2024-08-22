---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_bandwidth"
description: ""
---

# huaweicloud_vpc_bandwidth

Manages a **Shared** Bandwidth resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_vpc_bandwidth" "bandwidth_1" {
  name = "bandwidth_1"
  size = 5
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the Shared Bandwidth.
  If omitted, the provider-level region will be used. Changing this creates a new bandwidth.

* `name` - (Required, String) Specifies the bandwidth name. The value is a string of 1 to 64 characters that
  can contain letters, digits, underscores (_), hyphens (-), and periods (.).

* `size` - (Required, Int) Specifies the size of the Shared Bandwidth.
  If `charge_mode` is **bandwidth**, the value ranges from 5 Mbit/s to 2000 Mbit/s.
  If `charge_mode` is **95peak_plus**, the value ranges from 300 Mbit/s to 2000 Mbit/s.

* `charge_mode` - (Optional, String) Specifies whether the billing is based on bandwidth or
  95th percentile bandwidth (enhanced). Possible values can be **bandwidth** and **95peak_plus**.
  The default value is **bandwidth**, and **95peak_plus** is only valid for v4 and v5 Customer.
  
-> **NOTE:** When `charging_mode` is **prePaid**, only **bandwidth** is valid, please updating `charge_mode`
  to **bandwidth** before changing to **prePaid** billing mode.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the Shared Bandwidth.
  Changing this creates a new bandwidth.

* `charging_mode` - (Optional, String) Specifies the charging mode of the Shared Bandwidth.
  The valid values are **prePaid** and **postPaid**, defaults to **postPaid**.

* `period_unit` - (Optional, String) Specifies the charging period unit of the Shared Bandwidth.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `period` - (Optional, Int) Specifies the charging period of the Shared Bandwidth.
  + If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  + If `period_unit` is set to **year**, the value ranges from `1` to `3`.

  This parameter is mandatory if `charging_mode` is set to **prePaid**.

-> **NOTE:** `period_unit`, `period` can only be updated when changing from **postPaid** to **prePaid** billing mode.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled.
  Valid values are **true** and **false**. Defaults to **false**.

* `bandwidth_type` - (Optional, String, ForceNew) Specifies the bandwidth type.
  Valid values are **share** and **edgeshare**. Default is **share**.

* `public_border_group` - (Optional, String, ForceNew) Specifies the site is center of border.
  Valid values are **center** and the name of the border site. Default is **center**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the Shared Bandwidth.

* `share_type` - Indicates whether the bandwidth is shared or dedicated.

* `status` - Indicates the bandwidth status.

* `created_at` - Indicates the bandwidth create time.

* `updated_at` - Indicates the bandwidth update time.

* `publicips` - An array of EIPs that use the bandwidth. The object includes the following:
  + `id` - The ID of the EIP or IPv6 port that uses the bandwidth.
  + `type` - The EIP type. Possible values are *5_bgp* (dynamic BGP) and *5_sbgp* (static BGP).
  + `ip_version` - The IP version, either 4 or 6.
  + `ip_address` - The IPv4 or IPv6 address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Shared Bandwidths can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpc_bandwidth.bandwidth_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```

Note that the imported state may not be identical to your resource definition, due to payment attributes missing from
the API response.
The missing attributes include: `period_unit`, `period`, `auto_renew`.
It is generally recommended running `terraform plan` after importing a Shared Bandwidth.
You can ignore changes as below.

```hcl
resource "huaweicloud_vpc_bandwidth" "bandwidth_1" {
  ...

  lifecycle {
    ignore_changes = [
      period_unit, period, auto_renew,
    ]
  }
}
```
