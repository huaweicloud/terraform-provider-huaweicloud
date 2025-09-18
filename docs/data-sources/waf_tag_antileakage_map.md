---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_tag_antileakage_map"
description: |-
  Use this data source to query details about sensitive information options.
---

# huaweicloud_waf_tag_antileakage_map

Use this data source to query details about sensitive information options.

## Example Usage

```hcl
variable "lang" {}

data "huaweicloud_waf_tag_antileakage_map" "test" {
  lang = var.lang
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `lang` - (Required, String) Specifies the language.
  The value can be **cn** or **en**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `leakagemap` - The content type and return code of sensitive information.

  The [leakagemap](leakagemap_struct) structure is documented below.

* `locale` - The description of each type in the sensitive information option.

  The [locale](locale_struct) structure is documented below.

<a name="leakagemap_struct"></a>
The `leakagemap` block supports:

* `sensitive` - The content type of sensitive information.

* `code` - The return code of sensitive information.

<a name="locale_struct"></a>
The `locale` block supports:

* `code` - The response code interception, which is used to capture and process specific HTTP response codes.

* `id_card` - The ID card number, which is the unique code for identifying an individual.

* `sensitive` - The sensitive information filtering, which is used to detect and process sensitive information.

* `phone` - The phone number, which is used for contact.

* `responsecode` - The various response codes involved in the options.

* `email` - The Email address, which is used for electronic communication.
