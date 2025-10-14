---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_quotas"
description: |-
  Use this data source to get the CDN resource quotas within HuaweiCloud.
---

# huaweicloud_cdn_quotas

Use this data source to get the CDN resource quotas within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cdn_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the resource quotas are located.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of the resource quotas that matched filter parameters.  
  The [quotas](#cdn_quotas) structure is documented below.

<a name="cdn_quotas"></a>
The `quotas` block supports:

* `limit` - The limit of the resource quota.

* `type` - The type of the resource quota.

* `used` - The used capacity of the resource quota.

* `user_domain_id` - The domain ID of the user to which the resource quota belong.
