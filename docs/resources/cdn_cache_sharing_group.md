---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_cache_sharing_group"
description: |-
  Manages a cache sharing group resource within HuaweiCloud.
---

# huaweicloud_cdn_cache_sharing_group

Manages a cache sharing group resource within HuaweiCloud.

## Example Usage

```hcl
variable "cache_sharing_group_name" {}
variable "primary_domain_name" {}
variable "share_cache_records" {
  type = list(string)
}

resource "huaweicloud_cdn_cache_sharing_group" "test" {
  name           = var.cache_sharing_group_name
  primary_domain = var.primary_domain_name

  dynamic "share_cache_records" {
    for_each = var.share_cache_records

    content {
      domain_name = share_cache_records.value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, NonUpdatable) Specifies the name of the cache sharing group.

* `primary_domain` - (Required, String, NonUpdatable) Specifies the primary domain name of the cache sharing group.

* `share_cache_records` - (Required, List) Specifies the list of associated domain names of the cache sharing group.  
  The [share_cache_records](#cdn_share_cache_records) structure is documented below.

<a name="cdn_share_cache_records"></a>
The `share_cache_records` block supports:

* `domain_name` - (Required, String) Specifies the associated domain name.

-> The primary domain will not be automatically join the cache sharing, users need to explicitly declare in the
   `share_cache_records` list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - The creation time of the cache sharing group, in RFC3339 format.

## Import

The cache sharing group can be imported using the `id` or `name`, e.g.

```bash
$ terraform import huaweicloud_cdn_cache_sharing_group.test <id>
```

or

```bash
$ terraform import huaweicloud_cdn_cache_sharing_group.test <name>
```
