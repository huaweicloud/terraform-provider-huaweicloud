---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_tag_ip_reputation_map"
description: |-
  Use this data source to query details about threat intelligence control protection options.
---

# huaweicloud_waf_tag_ip_reputation_map

Use this data source to query details about threat intelligence control protection options.

## Example Usage

```hcl
variable "lang" {}
variable "type" {}

data "huaweicloud_waf_tag_ip_reputation_map" "test" {
  lang = var.lang
  type = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `lang` - (Required, String) Specifies the language.
  The value can be **cn** or **en**.

* `type` - (Required, String) Specifies the language.
  The value only can be **idc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ip_reputation_map` - The content type of the threat intelligence control protection options.

  The [ip_reputation_map](reputaion_struct) structure is documented below.

* `locale` - The description of each option in the threat intelligence control protection options.

<a name="reputaion_struct"></a>
The `ip_reputation_map` block supports:

* `idc` - The types of content controlled by threat intelligence.
