---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_cache_sharing_groups"
description: |-
  Use this data source to get a list of cache sharing groups within HuaweiCloud.
---

# huaweicloud_cdn_cache_sharing_groups

Use this data source to get a list of cache sharing groups within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cdn_cache_sharing_groups" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the cache sharing groups are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - The list of the cache sharing groups.  
  The [groups](#cdn_share_cache_groups_attr) structure is documented below.

<a name="cdn_share_cache_groups_attr"></a>
The `groups` block supports:

* `id` - The ID of the cache sharing group.

* `group_name` - The name of the cache sharing group.

* `primary_domain` - The primary domain name of the cache sharing group.

* `share_cache_records` - The list of associated domain names of the cache sharing group.  
  The [share_cache_records](#cdn_share_cache_records_attr) structure is documented below.

* `create_time` - The creation time of the cache sharing group, in RFC3339 format.

<a name="cdn_share_cache_records_attr"></a>
The `share_cache_records` block supports:

* `domain_name` - The associated domain name.
