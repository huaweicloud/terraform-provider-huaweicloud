---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_overviews_classification"
description: |-
  Use this data source to get top security statistics by category.
---

# huaweicloud_waf_overviews_classification

Use this data source to get top security statistics by category.

## Example Usage

```hcl
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_waf_overviews_classification" "test" {
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

-> The parameters `from` and `to` must be used together.

* `top` - (Optional, Int) The first several results to query.
  The valid value ranges from `1` to `10`. Defaults to `5`.

* `hosts` - (Optional, String) Specifies the ID of the domain.

* `instances` - (Optional, String) Specifies the ID of the instance.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  If you want to query resources under all enterprise projects, set this parameter to **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `domain` - The attacked domain.

  The [domain](#domain_struct) structure is documented below.

* `attack_type` - The attack event distribution.

  The [attack_type](#attack_type_struct) structure is documented below.

* `ip` - The attacking source IP address.

  The [ip](#ip_struct) structure is documented below.

* `url` - The attacking URL.

  The [url](#url_struct) structure is documented below.

* `geo` - The attacking source region.

  The [geo](#geo_struct) structure is documented below.

<a name="domain_struct"></a>
The `domain` block supports:

* `items` - The attacked domain details.

  The [items](#domain_items_struct) structure is documented below.

<a name="domain_items_struct"></a>
The `items` block supports:

* `key` - The name of the domain.

* `num` - The number of times attacked.

* `web_tag` - The website name.

<a name="attack_type_struct"></a>
The `attack_type` block supports:

* `items` - The attack event details.

  The [items](#attack_type_items_struct) structure is documented below.

<a name="attack_type_items_struct"></a>
The `items` block supports:

* `key` - The attack type.

* `num` - The number of attack event.

<a name="ip_struct"></a>
The `ip` block supports:

* `items` - The IP details.

  The [items](#ip_items_struct) structure is documented below.

<a name="ip_items_struct"></a>
The `items` block supports:

* `key` - The IP address.

* `num` - The number of the attacking IP address.

<a name="url_struct"></a>
The `url` block supports:

* `items` - The URL details.

  The [items](#url_items_struct) structure is documented below.

<a name="url_items_struct"></a>
The `items` block supports:

* `key` - The URL path.

* `num` - The number of the attacking URL.

* `host` - The domain name.

<a name="geo_struct"></a>
The `geo` block supports:

* `items` - The source region details.

  The [items](#geo_items_struct) structure is documented below.

<a name="geo_items_struct"></a>
The `items` block supports:

* `key` - The source region.

* `num` - The number of the attacking source region.
