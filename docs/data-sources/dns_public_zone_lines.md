---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_public_zone_lines"
description: |-
  Use this data source to get the list of DNS public zone lines.
---

# huaweicloud_dns_public_zone_lines

Use this data source to get the list of DNS public zone lines.

## Example Usage

```hcl
variable "zone_id" {}

data "huaweicloud_dns_public_zone_lines" "test" {
  zone_id = var.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `zone_id` - (Required, String) Specifies the zone ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `lines` - Indicates the lines.

  The [lines](#lines_struct) structure is documented below.

<a name="lines_struct"></a>
The `lines` block supports:

* `id` - Indicates the line ID.

* `line` - Indicates the line name.

* `create_time` - Indicates the creation time. Format is **yyyy-MM-dd'T'HH:mm:ss.SSS**.
