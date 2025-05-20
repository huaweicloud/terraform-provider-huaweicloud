---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_script_order_batch_details"
description: |-
  Use this data source to get the instance list of COC script orders in a batch.
---

# huaweicloud_coc_script_order_batch_details

Use this data source to get the instance list of COC script orders in a batch.

## Example Usage

```hcl
variable "execute_uuid" {}

data "huaweicloud_coc_script_order_batch_details" "test" {
  execute_uuid = var.execute_uuid
  batch_index  = 1
}
```

## Argument Reference

The following arguments are supported:

* `batch_index` - (Required, Int) Specifies the batch index.

* `execute_uuid` - (Required, String) Specifies the execution ID of the script order.

* `status` - (Optional, String) Specifies instance execution status.
  Values can be as follows:
  + **READY**: The operation is to be performed.
  + **PROCESSING**: The operation is in progress.
  + **ABNORMAL**: Abnormal.
  + **CANCELED**: Canceled.
  + **FINISHED**: Success.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `execute_instances` - Indicates a list of executed instances.

  The [execute_instances](#data_execute_instances_struct) structure is documented below.

<a name="data_execute_instances_struct"></a>
The `execute_instances` block supports:

* `id` - Indicates the primary key ID.

* `target_instance` - Indicates the destination instance.

  The [target_instance](#execute_instances_target_instance_struct) structure is documented below.

* `gmt_created` - Indicates the creation time.

* `gmt_finished` - Indicates the completion time.

* `execute_costs` - Indicates the time consumed in seconds.

* `status` - Indicates the instance execution status.
  Values can be as follows:
  + **READY**: The operation is to be performed.
  + **PROCESSING**: The operation is in progress.
  + **ABNORMAL**: Abnormal.
  + **CANCELED**: Canceled.
  + **FINISHED**: Success.

* `message` - Indicates the instance execution logs.

<a name="execute_instances_target_instance_struct"></a>
The `target_instance` block supports:

* `resource_id` - Indicates the unique ID of an instance.

* `provider` - Indicates the resource provider.

* `region_id` - Indicates the ID of the region to which the host belongs.

* `type` - Indicates the resource type of the resource provider.
  Values can be as follows:
  + **CLOUDSERVER**: Indicates the cloud server type.

* `custom_attributes` - Indicates that users are allowed to customize five attributes in the key/value format.

  The [custom_attributes](#target_instance_custom_attributes_struct) structure is documented below.

* `agent_sn` - Indicates the agent management ID.

* `agent_status` - Indicates the agent management status.

* `properties` - Indicates the additional attributes of the host.

  The [properties](#target_instance_properties_struct) structure is documented below.

<a name="target_instance_custom_attributes_struct"></a>
The `custom_attributes` block supports:

* `key` - Indicates the user-defined attribute key.

* `value` - Indicates the value of a custom property.

<a name="target_instance_properties_struct"></a>
The `properties` block supports:

* `host_name` - Indicates the host name.

* `fixed_ip` - Indicates the internal network IP address.

* `floating_ip` - Indicates the elastic IP address.

* `region_id` - Indicates the regions.

* `zone_id` - Indicates the availability zone.

* `application` - Indicates the CMDB application.

* `group` - Indicates the CMDB group.

* `project_id` - Indicates the project ID of an instance.
