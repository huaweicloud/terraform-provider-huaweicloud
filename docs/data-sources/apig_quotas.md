---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_quotas"
description: |-
  Use this data source to query the quotas within HuaweiCloud.
---

# huaweicloud_apig_quotas

Use this data source to query the quotas within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_apig_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the quotas.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of quotas.  
  The [quotas](#apig_quotas_attr) structure is documented below.

<a name="apig_quotas_attr"></a>
The `quotas` block supports:

* `id` - The ID of the quota.

* `name` - The name of the quota.

* `value` - The value of the quota.

* `description` - The description of the quota.

* `created_at` - The creation time of the quota, in RFC3339 format.
