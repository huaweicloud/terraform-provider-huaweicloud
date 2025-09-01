---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_instance_batches"
description: |-
  Use this data source to get the list of COC instance batches.
---

# huaweicloud_coc_instance_batches

Use this data source to get the list of COC instance batches.

## Example Usage

```hcl
variable "resource_id" {}
variable "region_id" {}
variable "host_name" {}
variable "fixed_ip" {}
variable "zone_id" {}

data "huaweicloud_coc_instance_batches" "test" {
  batch_strategy = "AUTO_BATCH"
  target_instances {
    resource_id        = var.resource_id
    cloud_service_name = "ecs"
    region_id          = var.region_id
    type               = "CLOUDSERVER"
    properties {
      host_name = var.hostname
      fixed_ip  = var.fixed_ip
      region_id = var.region_id
      zone_id   = var.zone_id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `batch_strategy` - (Required, String) Specifies the batching strategy, only automatic batching is supported.
  The value is **AUTO_BATCH**.

* `target_instances` - (Required, List) Specifies the target host instance.

  The [target_instances](#target_instances_struct) structure is documented below.

<a name="target_instances_struct"></a>
The `target_instances` block supports:

* `resource_id` - (Required, String) Specifies the unique ID of the instance.

* `cloud_service_name` - (Required, String) Specifies the resource provider: ECS. For a single script ticket, the
  provider is the same for each instance.

* `region_id` - (Required, String) Specifies the resource region ID.

* `type` - (Required, String) Specifies the resource type under the resource provider. If not specified, it defaults to
  **CLOUDSERVER**.

* `custom_attributes` - (Optional, List) Specifies the users can customize five attributes in the key_value format.

  The [custom_attributes](#target_instances_custom_attributes_struct) structure is documented below.

* `properties` - (Optional, List) Specifies the additional host properties.

  The [properties](#target_instances_properties_struct) structure is documented below.

<a name="target_instances_custom_attributes_struct"></a>
The `custom_attributes` block supports:

* `key` - (Required, String) Specifies the custom attribute key.

* `value` - (Required, String) Specifies the custom attribute value.

<a name="target_instances_properties_struct"></a>
The `properties` block supports:

* `host_name` - (Required, String) Specifies the host name.

* `fixed_ip` - (Required, String) Specifies the fixed IP.

* `region_id` - (Required, String) Specifies the region ID.

* `zone_id` - (Required, String) Specifies the availability zone.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the batch results.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `batch_index` - Indicates the batch ID.

* `target_instances` - Indicates the machines in the current batch.

  The [target_instances](#data_target_instances_struct) structure is documented below.

<a name="data_target_instances_struct"></a>
The `target_instances` block supports:

* `resource_id` - Indicates the unique ID of the instance.

* `cloud_service_name` - Indicates the resource provider: ECS. For a single script ticket, the provider is the same
  for each instance.

* `region_id` - Indicates the resource region ID.

* `type` - Indicates the resource type under the resource provider. If not specified, it defaults to **CLOUDSERVER**.

* `custom_attributes` - Indicates the users can customize five attributes in the key_value format.

  The [custom_attributes](#target_instances_custom_attributes_struct) structure is documented below.

* `properties` - Indicates the additional host properties.

  The [properties](#target_instances_properties_struct) structure is documented below.

<a name="target_instances_custom_attributes_struct"></a>
The `custom_attributes` block supports:

* `key` - Indicates the custom attribute key.

* `value` - Indicates the custom attribute value.

<a name="target_instances_properties_struct"></a>
The `properties` block supports:

* `host_name` - Indicates the host name.

* `fixed_ip` - Indicates the fixed IP.

* `region_id` - Indicates the region ID.

* `zone_id` - Indicates the availability zone.
