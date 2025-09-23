---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_overviews_attack_ip"
description: |-
  Use this data source to query the IP address of an attack source.
---

# huaweicloud_waf_overviews_attack_ip

Use this data source to query the IP address of an attack source.

## Example Usage

```hcl
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_waf_overviews_attack_ip" "test" {
  from = var.start_time
  to   = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `from` - (Required, Int) Specifies the query start time.
  The format is 13-digit timestamp in millisecond.

* `to` - (Required, Int) Specifies the query end time.
  The format is 13-digit timestamp in millisecond.

* `top` - (Optional, Int) Specifies the first several results to query.
  The valid value ranges from `1` to `10`.

* `hosts` - (Optional, String) Specifies the ID of the domain.

* `instances` - (Optional, String) Specifies the ID of the instance.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The detail information.
  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `key` - The type.
  The value can be **ACCESS**, **CRAWLER**, **ATTACK**, **WEB_ATTACK**, **PRECISE** or **CC**.

* `num` - The quantity.
