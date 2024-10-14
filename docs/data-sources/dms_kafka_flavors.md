---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_flavors"
description: ""
---

# huaweicloud_dms_kafka_flavors

Use this data source to get the list of available flavor details within HuaweiCloud.

## Example Usage

### Query the list of kafka flavors for cluster type

```hcl
data "huaweicloud_dms_kafka_flavors" "test" {
  type = "cluster"
}
```

### Query the kafka flavor details of the specified ID

```hcl
data "huaweicloud_dms_kafka_flavors" "test" {
  flavor_id = "c6.2u4g.cluster"
}
```

### Query list of kafka flavors that available in the availability zone list

```hcl
variable "az1" {}
variable "az2" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  availability_zones = [
    var.az1,
    var.az2,
  ]
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the dms kafka flavors.
  If omitted, the provider-level region will be used.

* `flavor_id` - (Optional, String) Specifies the DMS flavor ID, e.g. **c6.2u4g.cluster**.

* `storage_spec_code` - (Optional, String) Specifies the disk IO encoding.
  + **dms.physical.storage.high.v2**: Type of the disk that uses high I/O.
  + **dms.physical.storage.ultra.v2**: Type of the disk that uses ultra-high I/O.

* `type` - (Optional, String) Specifies flavor type. The valid values are **single**, **cluster** and **cluster.small**.

* `arch_type` - (Optional, String) Specifies the type of CPU architecture, e.g. **X86**.

* `availability_zones` - (Optional, List) Specifies the list of availability zones with available resources.

* `charging_mode` - (Optional, String) Specifies the flavor billing mode.
  The valid values are **prePaid** and **postPaid**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - The supported flavor versions.

* `flavors` - The list of flavor details.
  The [object](#dms_kafka_flavors) structure is documented below.

<a name="dms_kafka_flavors"></a>
The `flavors` block supports:

* `id` - The flavor ID.

* `type` - The flavor type.

* `vm_specification` - The underlying VM specification.

* `arch_types` - The list of supported CPU architectures.

* `charging_modes` - The list of supported billing modes.

* `ios` - The list of supported disk IO types.
  The [object](#dms_kafka_flavor_ios) structure is documented below.

* `support_features` - The list of features supported by the current specification.
  The [object](#dms_kafka_flavor_support_features) structure is documented below.

* `properties` - The properties of the current specification.
  The [object](#dms_kafka_flavor_properties) structure is documented below.

<a name="dms_kafka_flavor_ios"></a>
The `ios` block supports:

* `storage_spec_code` - The disk IO encoding.

* `type` - The disk type.

* `availability_zones` - The list of availability zones with available resources.

* `unavailability_zones` - The list of unavailability zones with available resources.

<a name="dms_kafka_flavor_support_features"></a>
The `support_features` block supports:

* `name` - The function name, e.g. **connector_obs**.

* `properties` - The function property details.
  The [object](#dms_kafka_flavor_support_feature_properties) structure is documented below.

<a name="dms_kafka_flavor_support_feature_properties"></a>
The `properties` block supports:

* `max_task` - The maximum number of tasks for the dump function.

* `min_task` - The minimum number of tasks for the dump function.

* `max_node` - The maximum number of nodes for the dump function.

* `min_node` - The minimum number of nodes for the dump function.

<a name="dms_kafka_flavor_properties"></a>
The `properties` block supports:

* `max_broker` - The maximum number of brokers.

* `min_broker` - The minimum number of brokers.

* `max_bandwidth_per_broker` - The maximum bandwidth per broker.

* `max_consumer_per_broker` - The maximum number of consumers per broker.

* `max_partition_per_broker` - The maximum number of partitions per broker.

* `max_tps_per_broker` - The maximum TPS per broker.

* `max_storage_per_node` - The maximum storage per node. The unit is GB.

* `min_storage_per_node` - The minimum storage per node. The unit is GB.

* `flavor_alias` - The flavor ID alias.
