---
subcategory: "Elastic IP (EIP)"
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

* `size` - (Required, Int) Specifies the size of the Shared Bandwidth. The value ranges from 5 Mbit/s to 2000 Mbit/s.

* `charge_mode` - (Optional, String, ForceNew) Specifies whether the billing is based on bandwidth or
  95th percentile bandwidth (enhanced). Possible values can be **bandwidth** and **95peak_plus**.
  The default value is **bandwidth**, and **95peak_plus** is only valid for v4 and v5 Customer.
  Changing this creates a new bandwidth.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the Shared Bandwidth.
  Changing this creates a new bandwidth.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the Shared Bandwidth.

* `share_type` - Indicates whether the bandwidth is shared or dedicated.

* `bandwidth_type` - Indicates the bandwidth type.

* `status` - Indicates the bandwidth status.

* `publicips` - An array of EIPs that use the bandwidth. The object includes the following:
  + `id` - The ID of the EIP or IPv6 port that uses the bandwidth.
  + `type` - The EIP type. Possible values are *5_bgp* (dynamic BGP) and *5_sbgp* (static BGP).
  + `ip_version` - The IP version, either 4 or 6.
  + `ip_address` - The IPv4 or IPv6 address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

Shared Bandwidths can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_vpc_bandwidth.bandwidth_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
