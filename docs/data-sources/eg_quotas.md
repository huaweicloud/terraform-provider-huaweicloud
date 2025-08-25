---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_quotas"
description: |-
  Use this data source to query EG quota of the current tenant within HuaweiCloud.
---

# huaweicloud_eg_quotas

Use this data source to query EG quota of the current tenant within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_eg_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the quota of the resource type to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of resource quotas that matched filter parameters.  
  The [quotas](#eg_quotas_attr) structure is documented below.

<a name="eg_quotas_attr"></a>
The `quotas` block supports:

* `name` - The quota name.

* `type` - The quota type.

* `quota` - The quota of current tenant.

* `used` - The quota used by the current tenant.

* `max` - The maximum quota.

* `min` - The minimum quota.
