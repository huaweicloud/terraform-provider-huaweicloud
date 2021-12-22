---
subcategory: "Elastic IP (EIP)"
---

# huaweicloud_vpc_eip

Manages an EIP resource within HuaweiCloud. This is an alternative to `huaweicloud_vpc_eip_v1`

## Example Usage

### EIP with Dedicated Bandwidth

```hcl
resource "huaweicloud_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type  = "PER"
    name        = "test"
    size        = 10
    charge_mode = "traffic"
  }
}
```

### EIP with Shared Bandwidth

```hcl
resource "huaweicloud_vpc_bandwidth" "bandwidth_1" {
  name = "bandwidth_1"
  size = 5
}

resource "huaweicloud_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type = "WHOLE"
    id         = huaweicloud_vpc_bandwidth.bandwidth_1.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the EIP resource. If omitted, the provider-level
  region will be used. Changing this creates a new resource.

* `publicip` - (Required, List) The elastic IP address object.

* `bandwidth` - (Required, List) The bandwidth object.

* `name` - (Optional, String) Specifies the name of the elastic IP. The value can contain 1 to 64 characters,
  including letters, digits, underscores (_), hyphens (-), and periods (.).

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the elastic IP.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the elastic IP. Changing this
  creates a new eip.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the elastic IP. Valid values are
  *prePaid* and *postPaid*, defaults to *postPaid*. Changing this creates a new eip.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the elastic IP. Valid values are
  *month* and *year*. This parameter is mandatory if `charging_mode` is set to *prePaid*. Changing this creates a new
  eip.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the elastic IP. If `period_unit` is set to
  *month*, the value ranges from 1 to 9. If `period_unit` is set to *year*, the value ranges from 1 to 3. This parameter
  is mandatory if `charging_mode` is set to *prePaid*. Changing this creates a new resource.

* `auto_renew` - (Optional, String, ForceNew) Specifies whether auto renew is enabled. Valid values are "true" and "
  false". Changing this creates a new resource.

The `publicip` block supports:

* `type` - (Optional, String, ForceNew) Specifies the EIP type. Possible values are *5_bgp* (dynamic BGP)
  and *5_sbgp* (static BGP), the default value is *5_bgp*. Changing this creates a new resource.

* `ip_address` - (Optional, String, ForceNew) Specifies the EIP to be assigned. The value must be a valid **IPv4**
  address in the available IP address range. The system automatically assigns an EIP if you do not specify it.
  Changing this creates a new resource.

* `ip_version` - (Optional, Int) Specifies the IP version, either 4 (default) or 6.

* `port_id` - (Optional, String) The port id which this EIP will associate with. If the value is "" or not
  specified, the EIP will be in unbind state.

The `bandwidth` block supports:

* `share_type` - (Required, String, ForceNew) Whether the bandwidth is dedicated or shared. Changing this creates a new
  resource. Possible values are as follows:
  + *PER*: Dedicated bandwidth
  + *WHOLE*: Shared bandwidth

* `name` - (Optional, String) The bandwidth name, which is a string of 1 to 64 characters that contain letters, digits,
  underscores (_), and hyphens (-). This parameter is mandatory when `share_type` is set to *PER*.

* `size` - (Optional, Int) The bandwidth size. The value ranges from 1 to 300 Mbit/s. This parameter is mandatory
  when `share_type` is set to *PER*.

* `id` - (Optional, String, ForceNew) The shared bandwidth id. This parameter is mandatory when
  `share_type` is set to *WHOLE*. Changing this creates a new resource.

* `charge_mode` - (Optional, String, ForceNew) Specifies whether the bandwidth is billed by traffic or by bandwidth
  size. The value can be *traffic* or *bandwidth*. Changing this creates a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `address` - The IPv4 address of the EIP.
* `ipv6_address` - The IPv6 address of the EIP.
* `private_ip` - The private IP address bound to the EIP.
* `status` - The status of EIP.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

EIPs can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_vpc_eip.eip_1 2c7f39f3-702b-48d1-940c-b50384177ee1
```
