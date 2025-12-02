---
subcategory: "Dedicated Host (DeH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_deh_quotas"
description: |-
  Use this data source to get the list of DeH quotas.
---

# huaweicloud_deh_quotas

Use this data source to get the list of DeH quotas.

## Example Usage

```hcl
data "huaweicloud_deh_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource` - (Optional, String) Indicates the quota resource type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quota_set` - Indicates the quotas information.
  The [quota_set](#attrblock--quota_set) structure is documented below.

<a name="attrblock--quota_set"></a>
The `quota_set` block supports:

* `hard_limit` - Indicates the existing quota.

* `resource` - Indicates the quota resource type.

* `used` - Indicates the number of the used instances.
