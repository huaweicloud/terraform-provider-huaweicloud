---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_quotas"
description: |-
  Use this data source to get the list of Organizations quotas within HuaweiCloud.
---

# huaweicloud_organizations_quotas

Use this data source to get the list of Organizations quotas within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_organizations_quotas" "test"{}
```

## Argument Reference

The following arguments are supported:

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - Indicates the list of the resource quotas.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - Indicates the quota information.

  The [resources](#quotas_resources_struct) structure is documented below.

<a name="quotas_resources_struct"></a>
The `resources` block supports:

* `type` - Indicates the quota type.  
  + **account**
  + **organizational_unit**
  + **policy**
  + **close_account_count_in_30_days**
  + **service_control_policy_per_entity**

* `quota` - Indicates the number of quotas.

* `min` - Indicates the minimum quota.

* `max` - Indicates the maximum quota.

* `used` - Indicates the used quantity.
