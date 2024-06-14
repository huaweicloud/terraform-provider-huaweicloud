---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_extend_flavors"
description: |-
  Use this data source to get the list of Kafka instance extend flavors.
---

# huaweicloud_dms_kafka_extend_flavors

Use this data source to get the list of Kafka instance extend flavors.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_extend_flavors" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `type` - (Optional, String) Specifies the flavor type.

* `charging_mode` - (Optional, String) Specifies the flavor billing mode. The valid values are **prePaid** and **postPaid**.

* `arch_type` - (Optional, String) Specifies the type of CPU architecture, e.g. **X86**.

* `storage_spec_code` - (Optional, String) Specifies the disk IO encoding, e.g. **dms.physical.storage.high.v2**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - Indicates the versions supported by the message engine.

* `flavors` - Indicates the flavor information for specification modification.

  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `id` - Indicates the flavor ID.

* `vm_specification` - Indicates the ECS flavor used by the flavor.

* `type` - Indicates the instance type.

* `charging_modes` - Indicates the supported billing modes.

* `arch_types` - Indicates the supported CPU architectures.

* `ios` - Indicates the disk I/O information.

  The [ios](#flavors_ios_struct) structure is documented below.

* `properties` - Indicates the properties.

  The [properties](#flavors_properties_struct) structure is documented below.

* `support_features` - Indicates the supported features.

  The [support_features](#flavors_support_features_struct) structure is documented below.

* `available_zones` - Indicates the AZs where there are available resources.

* `unavailable_zones` - Indicates the AZs where resources are sold out.

<a name="flavors_ios_struct"></a>
The `ios` block supports:

* `storage_spec_code` - Indicates the storage I/O specification.

* `type` - Indicates the I/O type.

* `available_zones` - Indicates the AZs where there are available resources.

* `unavailable_zones` - Indicates the AZs where resources are sold out.

<a name="flavors_properties_struct"></a>
The `properties` block supports:

* `max_partition_per_broker` - Indicates the maximum number of partitions of each broker.

* `max_bandwidth_per_broker` - Indicates the maximum bandwidth of each broker.

* `max_tps_per_broker` - Indicates the maximum TPS of each broker.

* `max_consumer_per_broker` - Indicates the maximum number of consumers of each broker.

* `min_broker` - Indicates the minimum number of brokers.

* `max_broker` - Indicates the maximum number of brokers.

* `min_storage_per_node` - Indicates the minimum storage space of each broker. Unit: GB.

* `max_storage_per_node` - Indicates the maximum storage space of each broker. Unit: GB.

* `flavor_alias` - Indicates the alias of **flavor_id**.

<a name="flavors_support_features_struct"></a>
The `support_features` block supports:

* `name` - Indicates the feature name.

* `properties` - Indicates the key-value pair of a feature.
