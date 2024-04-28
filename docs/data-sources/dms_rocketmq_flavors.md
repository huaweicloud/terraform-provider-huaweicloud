---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_flavors"
description: ""
---

# huaweicloud_dms_rocketmq_flavors

Use this data source to get the list of RocketMQ flavors within HuaweiCloud.

## Example Usage

```hcl
variable "az1" {}
variable "az2" {}

data "huaweicloud_dms_rocketmq_flavors" "test" {
  availability_zones = [
    var.az1, var.az2,
  ]
  arch_type          = "X86"
  charging_mode      = "prePaid"
  type               = "cluster"
  flavor_id          = "c6.2u4g.cluster"
  storage_spec_code  = "dms.physical.storage.high.v2"  
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `arch_type` - (Optional, String) Specifies the type of CPU architecture, e.g. **X86**.

* `availability_zones` - (Optional, List) Specifies the list of availability zone names.

* `charging_mode` - (Optional, String) Specifies the billing mode of the flavor.
  Value options: **prePaid** and **postPaid**.

* `flavor_id` - (Optional, String) Specifies the ID of the flavor, e.g. **c6.2u4g.cluster**.

* `storage_spec_code` - (Optional, String) Specifies the disk IO encoding.
  Value options:
  + **dms.physical.storage.high.v2**: Type of the disk that uses high I/O.
  + **dms.physical.storage.ultra.v2**: Type of the disk that uses ultra-high I/O.

* `type` - (Optional, String) Specifies the type of the flavor. Value options: **single** and **cluster**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - Indicates the list of flavor versions.

* `flavors` - Indicates the list of flavors.
  The [flavors](#DMS_rocketmq_flavors) structure is documented below.

<a name="DMS_rocketmq_flavors"></a>
The `flavors` block supports:

* `id` - Indicates the ID of the flavor.

* `arch_types` - Indicates the list of the types of CPU architecture.

* `charging_modes` - Indicates the list of the billing modes.

* `ios` - Indicates the list of disk IO types.
  The [ios](#DMS_rocketmq_flavor_ios) structure is documented below.

* `support_features` - Indicates the list of features supported by the current specification.
  The [support_features](#DMS_rocketmq_flavor_support_features) structure is documented below.

* `type` - Indicates the type of the flavor. The value can be: **single** and **cluster**.

* `properties` - Indicates the list of the properties of the current specification.
  The [properties](#DMS_rocketmq_flavor_properties) structure is documented below.

* `vm_specification` - Indicates the underlying VM specification, e.g. **c6.large.2**

<a name="DMS_rocketmq_flavor_ios"></a>
The `ios` block supports:

* `storage_spec_code` - Indicates the disk IO encoding.

* `type` - Indicates the disk type.

* `availability_zones` - Indicates the list of availability zone names.

* `unavailability_zones` - Indicates the list of unavailability zone names.

<a name="DMS_rocketmq_flavor_support_features"></a>
The `support_features` block supports:

* `name` - Indicates the function name, e.g. **connector_obs**.

* `properties` - Indicates the list of the function property details.
  The [properties](#DMS_rocketmq_flavor_support_feature_properties) structure is documented below.

<a name="DMS_rocketmq_flavor_properties"></a>
The `properties` block supports:

* `max_broker` - Indicates the maximum number of brokers.

* `min_broker` - Indicates the minimum number of brokers.

* `max_bandwidth_per_broker` - Indicates the maximum bandwidth per broker.

* `max_consumer_per_broker` - Indicates the maximum number of consumers per broker.

* `max_partition_per_broker` - Indicates the maximum number of partitions per broker.

* `max_tps_per_broker` - Indicates the maximum TPS per broker.

* `max_storage_per_node` - Indicates the maximum storage per node. The unit is GB.

* `min_storage_per_node` - Indicates the minimum storage per node. The unit is GB.

* `flavor_alias` - Indicates the alias of the flavor.

<a name="DMS_rocketmq_flavor_support_feature_properties"></a>
The `properties` block supports:

* `max_task` - Indicates the maximum number of tasks for the dump function.

* `min_task` - Indicates the minimum number of tasks for the dump function.

* `max_node` - Indicates the maximum number of nodes for the dump function.

* `min_node` - Indicates the minimum number of nodes for the dump function.
