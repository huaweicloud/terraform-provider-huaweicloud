---
subcategory: "Intelligent EdgeCloud (IEC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iec_eip"
description: ""
---

# huaweicloud_iec_eip

Manages a eip resource within HuaweiCloud IEC.

## Example Usage

```hcl
data "huaweicloud_iec_sites" "iec_sites" {}

resource "huaweicloud_iec_eip" "eip_test" {
  site_id = data.huaweicloud_iec_sites.iec_sites.sites[0].id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `site_id` - (Required, String, ForceNew) Specifies the ID of IEC service site. Changing this parameter creates a new
  resource.

* `line_id` - (Optional, String, ForceNew) Specifies the line ID of IEC service site.
  Changing this parameter creates a new resource.

* `port_id` - (Optional, String) Specifies the port ID which this eip will associate with.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `ip_version` - The version of elastic IP address.
* `status` - The status of elastic IP.
* `public_ip` - The address of elastic IP.
* `private_ip` - The address of private IP.
* `bandwidth_id` - The id of bandwidth.
* `bandwidth_name` - The name of bandwidth.
* `bandwidth_size` - The size of bandwidth.
* `bandwidth_share_type` - Whether the bandwidth is shared or exclusive.
* `site_info` - The located information of the IEC site. It contains area, province and city.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 3 minutes.

## Import

IEC EIPs can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iec_eip.eip_test b5ad19d1-57d1-48fd-aab7-1378f9bee169
```
