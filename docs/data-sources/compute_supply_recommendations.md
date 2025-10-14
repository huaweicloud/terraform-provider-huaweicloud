---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_supply_recommendations"
description: |-
  Use this data source to get the list of supply recommendation.
---

# huaweicloud_compute_supply_recommendations

Use this data source to get the list of supply recommendation.

## Example Usage

```hcl
variable "flavor_id" {}

data "huaweicloud_compute_supply_recommendations" "test" {
  flavor_ids = [var.flavor_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `flavor_constraint` - (Optional, List) Specifies the flavor constraint.
  The [flavor_constraint](#flavor_constraint_struct) structure is documented below.

* `flavor_ids` - (Optional, List) Specifies the flavor constraint.

* `locations` - (Optional, List) Specifies the locations.
  The [locations](#locations_struct) structure is documented below.

* `option` - (Optional, List) Specifies the option.
  The [option](#option_struct) structure is documented below.

* `strategy` - (Optional, String) Specifies the strategy. Value options: **CAPACITY** and **COST**.

<a name="flavor_constraint_struct"></a>
The `flavor_constraint` block supports:

* `architecture_type` - (Optional, List) Specifies the architecture type.

* `flavor_requirements` - (Optional, List) Specifies the flavor requirements.
  The [flavor_requirements](#flavor_requirements_struct) structure is documented below.

<a name="flavor_requirements_struct"></a>
The `flavor_requirements` block supports:

* `vcpu_count` - (Optional, List) Specifies the vcpu count.
  The [vcpu_count](#vcpu_count_struct) structure is documented below.

* `memory_mb` - (Optional, List) Specifies the memory in MByte (MB).
  The [memory_mb](#memory_mb_struct) structure is documented below.

* `cpu_manufacturers` - (Optional, List) Specifies the cpu manufacturers.

* `memory_gb_per_vcpu` - (Optional, List) Specifies the memory gb per vcpu.
  The [memory_gb_per_vcpu](#memory_gb_per_vcpu_struct) structure is documented below.

* `instance_generations` - (Optional, List) Specifies the instance generations.

<a name="vcpu_count_struct"></a>
The `vcpu_count` block supports:

* `max` - (Optional, Int) Specifies the max value. **-1** means no limit.

* `min` - (Optional, Int) Specifies the min value. **-1** means no limit.

<a name="memory_mb_struct"></a>
The `memory_mb` block supports:

* `max` - (Optional, Int) Specifies the max value. **-1** means no limit.

* `min` - (Optional, Int) Specifies the min value. **-1** means no limit.

<a name="memory_gb_per_vcpu_struct"></a>
The `memory_gb_per_vcpu` block supports:

* `max` - (Optional, Float) Specifies the max value. **-1** means no limit.

* `min` - (Optional, Float) Specifies the min value. **-1** means no limit.

<a name="locations_struct"></a>
The `locations` block supports:

* `region_id` - (Required, String) Specifies the region ID.

* `availability_zone_id` - (Optional, String) Specifies the availability zone ID.

<a name="option_struct"></a>
The `option` block supports:

* `result_granularity` - (Optional, String) Specifies the result granularity. Value options: **BY_REGION**, **BY_AZ**,
  **BY_FLAVOR**, **BY_FLAVOR_AND_REGION** and **BY_FLAVOR_AND_AZ**.

* `enable_spot` - (Optional, String) Specifies whether enable spot. Value options: **true**, **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `supply_recommendations` - Indicates the password in ciphertext.
  The [supply_recommendations](#supply_recommendations_struct) structure is documented below.

<a name="supply_recommendations_struct"></a>
The `supply_recommendations` block supports:

* `flavor_id` - Indicates the flavor ID.

* `region_id` - Indicates the region ID.

* `availability_zone_id` - Indicates the availability zone ID.

* `score` - Indicates the score.
