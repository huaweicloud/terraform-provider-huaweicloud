---
subcategory: "Intelligent EdgeCloud (IEC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iec_bandwidths"
description: ""
---

# huaweicloud_iec_bandwidths

Use this data source to get a list of bandwidths belong to a specific IEC site.

## Example Usage

```hcl
data "huaweicloud_iec_sites" "sites_test" {}

data "huaweicloud_iec_bandwidths" "demo" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required, String) Specifies the ID of the IEC site.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `site_info` - The located information of the iec site. It contains area, province and city.

* `bandwidths` - A list of all the bandwidths found. The object is documented below.

The `bandwidths` block supports:

* `id` - The ID of the bandwidth.
* `name` - The name of the bandwidth.
* `size` - The size of the bandwidth.
* `share_type` - Whether the bandwidth is shared or exclusive.
* `charge_mode` - The charging mode of the bandwidth.
* `line` - The line name of the bandwidth.
* `status` - The status of the bandwidth.
