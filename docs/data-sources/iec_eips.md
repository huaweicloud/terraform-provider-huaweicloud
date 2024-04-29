---
subcategory: "Intelligent EdgeCloud (IEC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iec_eips"
description: ""
---

# huaweicloud_iec_eips

Use this data source to get a list of EIPs belong to a specific IEC site.

## Example Usage

```hcl
data "huaweicloud_iec_sites" "sites_test" {}

data "huaweicloud_iec_eips" "site" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required, String) Specifies the ID of the IEC site.

* `port_id` - (Optional, String) Specifies the ID of the port.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `site_info` - The located information of the iec site. It contains area, province and city.

* `eips` - A list of all the EIPs found. The object is documented below.

The `eips` block supports:

* `id` - The ID of elastic IP.
* `status` - The status elastic IP.
* `ip_version` - The version of elastic IP address.
* `public_ip` - The address of elastic IP.
* `private_ip` - The address of private IP.
* `bandwidth_id` - The ID of bandwidth.
* `bandwidth_name` - The name of bandwidth.
* `bandwidth_size` - The size of bandwidth.
* `bandwidth_share_type` - Whether the bandwidth is shared or exclusive.
