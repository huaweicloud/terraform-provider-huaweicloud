---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_instances"
description: ""
---

# huaweicloud_dms_rabbitmq_instances

Use this data source to get the list of RabbitMQ instances.

## Example Usage

```hcl
var "instance_id" {}
var "name" {}
var "enterprise_project_id" {}
var "engine_version" {}
var "flavor_id" {}
var "type" {}

data "huaweicloud_dms_rabbitmq_instances" "test" {
  instance_id           = var.instance_id
  name                  = var.name
  exact_match_name      = "true" 
  enterprise_project_id = var.enterprise_project_id
  status                = "RUNNING" 
  engine_version        = var.engine_version
  flavor_id             = var.flavor_id
  type                  = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `engine_version` - (Optional, String) Specifies the version of the RabbitMQ engine.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the RabbitMQ instance belongs.

* `exact_match_name` - (Optional, String) Specifies whether to search for the instance that precisely matches
  a specified instance name. Value options: **true**, **false**. Defaults to **false**.

* `flavor_id` - (Optional, String) Specifies the flavor ID of the RabbitMQ instance.

* `instance_id` - (Optional, String) Specifies the ID of the RabbitMQ instance.

* `name` - (Optional, String) Specifies the name of the RabbitMQ instance.

* `status` - (Optional, String) Specifies the status of the RabbitMQ instance. Value options: **CREATING** **RUNNING**,
  **FAULTY**,  **RESTARTING**, **STARTING**, **CHANGING**, **CHANGE FAILED**, **FROZEN**, **FREEZING**, **UPGRADING**,
  **EXTENDING**, **ROLLING BACK**.

* `type` - (Optional, String) Specifies the RabbitMQ instance type. Value options: **cluster**, **single**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `instances` - Indicates the list of RabbitMQ instances.
The [instances](#DMS_rabbitmq_instances) structure is documented below.

<a name="DMS_rabbitmq_instances"></a>
The `instances` block supports:

* `id` - Indicates the ID of the RabbitMQ instance.

* `access_user` - Indicates the name of the user accessing the RabbitMQ instance.

* `availability_zones` - Indicates the list of the availability zone names.

* `broker_num` - Indicates the number of the brokers.

* `charging_mode` - Indicates the billing mode. The value can be: **prePaid** or **postPaid**.

* `connect_address` - Indicates the IP address of the RabbitMQ instance.

* `created_at` - Indicates the creation time of the RabbitMQ instance.

* `description` - Indicates the description of the RabbitMQ instance.

* `engine` - Indicates the message engine type. The value is `rabbitmq`.

* `engine_version` - Indicates the version of the RabbitMQ engine.

* `enterprise_project_id` - Indicates the enterprise project ID to which the RabbitMQ instance belongs.

* `extend_times` - Indicates the number of disk expansion times. If the value exceeds 20, disk expansion is
  no longer allowed.

* `flavor_id` - Indicates the flavor ID of the RabbitMQ instance.

* `is_logical_volume` - Indicates whether the instance is a new instance. This parameter is used to
  distinguish old instances from new instances during instance capacity expansion.

* `maintain_begin` - Indicates the time at which the maintenance window starts. The format is HH:mm:ss.

* `maintain_end` - Indicates the time at which the maintenance window ends. The format is HH:mm:ss.

* `management_connect_address` - Indicates the management address of the RabbitMQ instance.

* `name` - Indicates the name of the RabbitMQ instance.

* `port` - Indicates the port of the RabbitMQ instance.

* `security_group_id` - Indicates the ID of a security group.

* `security_group_name` - Indicates the name of a security group.

* `specification` - Indicates the instance specification.

* `ssl_enable` - Indicates whether the RabbitMQ SASL_SSL is enabled. The value can be: **true** or **false**.

* `status` - Indicates the status of the RabbitMQ instance. The value can be: **CREATING** **RUNNING**, **FAULTY**,
  **RESTARTING**, **STARTING**, **CHANGING**, **CHANGE FAILED**, **FROZEN**, **FREEZING**, **UPGRADING**, **EXTENDING**
  or **ROLLING BACK**.

* `storage_resource_id` - Indicates the ID of the storage resource.

* `storage_space` - Indicates the message storage space in GB.

* `storage_spec_code` - Indicates the storage I/O specification.

* `subnet_id` - Indicates the ID of a subnet.

* `tags` - Indicates the key/value pairs of tags associated with the RabbitMQ instance.

* `type` - Indicates the RabbitMQ instance type. The value can be: **cluster** or **single**.

* `used_storage_space` - Indicates the used message storage space in GB.

* `user_name` - Indicates the name of the user creating the RabbitMQ instance.

* `vpc_id` - Indicates the ID of a VPC.

* `vpc_name` - Indicates the name of a VPC.

* `disk_encrypted_enable` - Whether the disk is encrypted.

* `disk_encrypted_key` - The key of the disk encryption.
