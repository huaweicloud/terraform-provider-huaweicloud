---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_share_cache_group"
description: |-
  Manages a share cache group resource within HuaweiCloud.
---

# huaweicloud_cdn_share_cache_group

Manages a share cache group resource within HuaweiCloud.

## Example Usage

```hcl
variable "share_cache_group_name" {}
variable "primary_domain_name" {}
variable "share_cache_records" {
  type = list(string)
}

resource "huaweicloud_cdn_share_cache_group" "test" {
  name           = var.share_cache_group_name
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

* `region` - (Optional, String, ForceNew) Specifies the region where the share cache group is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the share cache group.

* `primary_domain` - (Required, String, NonUpdatable) Specifies the primary domain name.

* `share_cache_records` - (Required, List) Specifies the list of associated domain names.  
  The [share_cache_records](#cdn_share_cache_group_share_cache_records) structure is documented below.

<a name="cdn_share_cache_group_share_cache_records"></a>
The `share_cache_records` block supports:

* `domain_name` - (Required, String) Specifies the associated domain name.

-> The primary domain will not be automatically join the cache sharing, users need to explicitly declare in the
   `share_cache_records` list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - The creation time of the share cache group, in RFC3339 format.

## Import

The share cache group can be imported using the `id` or `group_name`, e.g.

```bash
$ terraform import huaweicloud_cdn_share_cache_group.test <id>
```

or

```bash
$ terraform import huaweicloud_cdn_share_cache_group.test <group_name>
```
