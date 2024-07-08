---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_quotas"
description: |-
  Using this data source to query the list of available resource quotas for FunctionGraph service within HuaweiCloud.
---

# huaweicloud_fgs_quotas

Using this data source to query the list of available resource quotas for FunctionGraph service within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_fgs_quotas" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - All quotas that match the filter parameters.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `type` - The resource type corresponding to quota.
  + **fgs_func_scale_down_timeout**: Release time of idle function instances in FunctionGraph `v1`.
  + **fgs_func_occurs**: Indicates instance quota for functions in FunctionGraph `v1` and reserved instance quota for
    functions in FunctionGraph `v2`.
  + **fgs_func_pat_idle_time**: Release time of idle PAT in VPC function.
  + **fgs_func_num**: User function quantity quota.
  + **fgs_func_code_size**: Total code size quota of user functions.
  + **fgs_workflow_num**: Function flow quantity quota.
  + **fgs_on_demand_instance_limit**: Maximum number of instances per function in FunctionGraph `v2`.
  + **fgs_func_qos_limit**: Instance quantity quota of user functions.

* `unit` - The unit of usage.
  
  -> If the resource type is **fgs_func_code_size**, the unit is `MB`. In other scenarios, there is no unit.

* `limit` - The number of available quota.

* `used` - The number of quota used.
