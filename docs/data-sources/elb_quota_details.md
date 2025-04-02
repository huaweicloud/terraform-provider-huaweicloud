---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_quota_details"
description: |-
  Use this data source to get the current quotas and used quotas of resources related to a ELB in a specific project.
---

# huaweicloud_elb_quota_details

Use this data source to get the current quotas and used quotas of resources related to a ELB in a specific project.

## Example Usage

```hcl
data "huaweicloud_elb_quota_details" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `quota_key` - (Optional, List) Specifies the resource type.
  Value options: **loadbalancer**, **listener**, **ipgroup**, **pool**, **member**, **healthmonitor**, **l7policy**,
  **certificate**, **security_policy**, **listeners_per_loadbalancer**, **listeners_per_pool**, **members_per_pool**,
  **condition_per_policy**, **ipgroup_bindings**, **ipgroup_max_length**, **ipgroups_per_listener**, **pools_per_l7policy**
  or **l7policies_per_listener**.
  
  Multiple values can be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - Indicates the resource quotas.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `quota_key` - Indicates the resource type.

* `quota_limit` - Indicates the total quota.
  + If the value is greater than or equal to **0**, it indicates the current quota.
  + **-1** indicates that the quota is not limited.

* `used` - Indicates the used quota.

* `unit` - Indicates the quota unit.
  The value can only be **count**.
