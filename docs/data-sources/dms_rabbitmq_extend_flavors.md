---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_extend_flavors"
description: |-
  Use this data source to get the list of RabbitMQ extend available flavor details within HuaweiCloud.
---

# huaweicloud_dms_rabbitmq_extend_flavors

Use this data source to get the list of RabbitMQ extend available flavor details within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_rabbitmq_extend_flavors" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the DMS RabbitMQ instance ID.

* `type` - (Optional, String) Specifies flavor type. The valid values are **single** and **cluster**.

* `charging_mode` - (Optional, String) Specifies the flavor billing mode.
  The valid values are **prePaid** and **postPaid**.

* `arch_type` - (Optional, String) Specifies the type of CPU architecture, e.g. **X86**.

* `storage_spec_code` - (Optional, String) Specifies the disk IO encoding, e.g. **dms.physical.storage.high.v2**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - Indicates the supported flavor versions.

* `flavors` - Indicates the list of flavor details.

  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `id` - Indicates the flavor ID.

* `type` - Indicates the flavor type.

* `charging_modes` - Indicates the list of supported billing modes.

* `vm_specification` - Indicates the underlying VM specification.

* `arch_types` - Indicates the list of supported CPU architectures.

* `ios` - Indicates the list of supported disk IO types.

  The [ios](#flavors_ios_struct) structure is documented below.

* `properties` - Indicates the properties of the current specification.

  The [properties](#flavors_properties_struct) structure is documented below.

* `support_features` - Indicates the list of features supported by the current specification.

  The [support_features](#flavors_support_features_struct) structure is documented below.

<a name="flavors_ios_struct"></a>
The `ios` block supports:

* `type` - Indicates the disk type.

* `storage_spec_code` - Indicates the disk IO encoding.

* `available_zones` - Indicates the list of availability zones with available resources.

* `unavailable_zones` - Indicates the list of unavailability zones with available resources.

<a name="flavors_properties_struct"></a>
The `properties` block supports:

* `min_broker` - Indicates the minimum number of brokers.

* `max_broker` - Indicates the maximum number of brokers.

* `min_storage_per_node` - Indicates the minimum storage per node. The unit is GB.

* `max_storage_per_node` - Indicates the maximum storage per node. The unit is GB.

* `max_queue_per_broker` - Indicates the maximum number of queues.

* `max_connection_per_broker` - Indicates the maximum number of connections.

* `step_length` - Indicates the step length.

* `flavor_alias` - Indicates the alias of **flavor_id**.

<a name="flavors_support_features_struct"></a>
The `support_features` block supports:

* `name` - Indicates the feature name.

* `properties` - Indicates the property details.
