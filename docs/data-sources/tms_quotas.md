---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_quotas"
description: |-
  Use this data source to get the list of tag quotas.
---

# huaweicloud_tms_quotas

Use this data source to get the list of tag quotas.

## Example Usage

```hcl
data "huaweicloud_tms_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - Indicates the list of quotas.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `quota_key` - Indicates the quota key.

* `quota_limit` - Indicates the quota value.

* `used` - Indicates the quota used.

* `unit` - Indicates the unit.
