---
subcategory: Content Delivery Network (CDN)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domain_tags"
description: |-
  Use this data source to get a list of CDN domain tags within HuaweiCloud.
---

# huaweicloud_cdn_domain_tags

Use this data source to get a list of CDN domain tags within HuaweiCloud.

## Example Usage

```hcl
variable "cdn_domain_id" {}

data "huaweicloud_cdn_domain_tags" "test" {
  resource_id = var.cdn_domain_id
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String) Specifies the ID of the domain to query tags.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of domain tags that matched filter parameters.  
  The [tags](#cdn_domain_tags) structure is documented below.

<a name="cdn_domain_tags"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
