---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_eip

Manages a eip resource within HuaweiCloud IEC.

## Example Usage

```hcl
data "huaweicloud_iec_sites" "iec_sites" {}

resource "huaweicloud_iec_eip" "eip_test" {
  site_id    = data.huaweicloud_iec_sites.iec_sites.sites[0].id
  ip_version = 4
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required, String, ForceNew) Specifies the id of iec sevice site. Changing this parameter creates a new iec eip resource.

* `ip_version` - (Optional, Int) The elastic IP address object. IEC services only support IPv4 now.

* `port_id` - (Optional, String) The elastic IP address object.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `status` - Specifies the [status](https://support.huaweicloud.com/intl/en-us/api-eip/eip_api_0002.html#eip_api_0002__en-us_topic_0201534285_table3035698) 
    of iec eip.

* `public_ip` - Indicates the elastic IP address object.

* `private_ip` - Indicates the private IP address.

* `bandwitch_name` - The name of bandwidth object.

* `bandwitch_size` - The size of bandwidth object.

* `bandwitch_share_type` - Whether the bandwidth is shared or exclusive. 

* `site_info` - Specifies the information of iec subnet site.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 3 minute.

## Import

IEC EIPs can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_iec_eip.eip_1 b5ad19d1-57d1-48fd-aab7-1378f9bee169
```
