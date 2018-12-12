---
layout: "huaweicloud"
page_title: "huaweicloud: huaweicloud_dcs_instance_v1"
sidebar_current: "docs-huaweicloud-resource-dcs-instance-v1"
description: |-
  Manages a DCS instance in the huaweicloud DCS Service
---

# huaweicloud\_dcs\_instance_v1

Manages a DCS instance in the huaweicloud DCS Service.

## Example Usage

### Automatically detect the correct network

```hcl
       resource "huaweicloud_networking_secgroup_v2" "secgroup_1" {
         name = "secgroup_1"
         description = "secgroup_1"
       }
       data "huaweicloud_dcs_az_v1" "az_1" {
         port = "8002"
		}
       data "huaweicloud_dcs_product_v1" "product_1" {
          spec_code = "dcs.master_standby"
		}
		resource "huaweicloud_dcs_instance_v1" "instance_1" {
		  name  = "test_dcs_instance"
          engine_version = "3.0.7"
          password = "Huawei_test"
          engine = "Redis"
          capacity = 2
          vpc_id = "1477393a-29c9-4de5-843f-18ef51257c7e"
          security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup_1.id}"
          subnet_id = "27d99e17-42f2-4751-818f-5c8c6c03ff15"
          available_zones = ["${data.huaweicloud_dcs_az_v1.az_1.id}"]
          product_id = "${data.huaweicloud_dcs_product_v1.product_1.id}"
          save_days = 1
          backup_type = "manual"
          begin_at = "00:00-01:00"
          period_type = "weekly"
          backup_at = [1]
          depends_on = ["data.huaweicloud_dcs_product_v1.product_1", "huaweicloud_networking_secgroup_v2.secgroup_1"]
		}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Indicates the name of an instance. An instance name starts with a letter,
	consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-).

* `description` - (Optional) Indicates the description of an instance. It is a character
    string containing not more than 1024 characters.

* `engine` - (Optional) Indicates a message engine. Options: rabbitmq and kafka.

* `engine_version` - (Optional) Indicates the version of a message engine.

* `capacity` - (Required) Indicates the Cache capacity. Unit: GB.
    For a DCS Redis or Memcached instance in single-node or master/standby mode, the cache
    capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB.
    For a DCS Redis instance in cluster mode, the cache capacity can be 64, 128, 256, 512,
    or 1024 GB.

* `partition_num` - (Optional) This parameter is mandatory when a Kafka instance is created.
    Indicates the maximum number of topics in a Kafka instance.
    When specification is 300 MB: 900
    When specification is 600 MB: 1800
    When specification is 1200 MB: 1800

* `access_user` - (Optional) Username used for accessing a DCS instance after password
    authentication. A username starts with a letter, consists of 1 to 64 characters,
    and supports only letters, digits, and hyphens (-).

* `password` - (Optional) Password of a DCS instance.
    The password of a DCS Redis instance must meet the following complexity requirements:

* `vpc_id` - (Required) Tenant's VPC ID. For details on how to create VPCs, see the
    Virtual Private Cloud API Reference.

* `security_group_id` - (Required) Tenant's security group ID. For details on how to
    create security groups, see the Virtual Private Cloud API Reference.

* `subnet_id` - (Required) Subnet ID. For details on how to create subnets, see the
    Virtual Private Cloud API Reference.

* `available_zones` - (Required) IDs of the AZs where cache nodes reside. For details
    on how to query AZs, see Querying AZ Information.

* `product_id` - (Required) Product ID used to differentiate DCS instance types.

* `maintain_begin` - (Optional) Indicates the time at which a maintenance time window starts.
    Format: HH:mm:ss.
    The start time and end time of a maintenance time window must indicate the time segment of
	a supported maintenance time window. For details, see section Querying Maintenance Time Windows.
    The start time must be set to 22:00, 02:00, 06:00, 10:00, 14:00, or 18:00.
    Parameters maintain_begin and maintain_end must be set in pairs. If parameter maintain_begin
	is left blank, parameter maintain_end is also blank. In this case, the system automatically
	allocates the default start time 02:00.

* `maintain_end` - (Optional) Indicates the time at which a maintenance time window ends.
    Format: HH:mm:ss.
    The start time and end time of a maintenance time window must indicate the time segment of
	a supported maintenance time window. For details, see section Querying Maintenance Time Windows.
    The end time is four hours later than the start time. For example, if the start time is 22:00,
	the end time is 02:00.
    Parameters maintain_begin and maintain_end must be set in pairs. If parameter maintain_end is left
	blank, parameter maintain_begin is also blank. In this case, the system automatically allocates
	the default end time 06:00.

* `save_days` - (Optional) Retention time. Unit: day. Range: 1–7.

* `backup_type` - (Optional) Backup type. Options:
    auto: automatic backup.
    manual: manual backup.

* `begin_at` - (Optional) Time at which backup starts. "00:00-01:00" indicates that backup
    starts at 00:00:00.

* `period_type` - (Optional) Interval at which backup is performed. Currently, only weekly
    backup is supported.

* `backup_at` - (Optional) Day in a week on which backup starts. Range: 1–7. Where: 1
    indicates Monday; 7 indicates Sunday.

## Attributes Reference

The following attributes are exported:


* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `engine` - See Argument Reference above.
* `engine_version` - See Argument Reference above.
* `capacity` - See Argument Reference above.
* `access_user` - See Argument Reference above.
* `password` - See Argument Reference above.
* `vpc_id` - See Argument Reference above.
* `vpc_name` - Indicates the name of a vpc.
* `security_group_id` - See Argument Reference above.
* `security_group_name` - Indicates the name of a security group.
* `subnet_id` - See Argument Reference above.
* `subnet_name` - Indicates the name of a subnet.
* `available_zones` - See Argument Reference above.
* `product_id` - See Argument Reference above.
* `maintain_begin` - See Argument Reference above.
* `maintain_end` - See Argument Reference above.
* `save_days` - See Argument Reference above.
* `backup_type` - See Argument Reference above.
* `begin_at` - See Argument Reference above.
* `period_type` - See Argument Reference above.
* `backup_at` - See Argument Reference above.
* `order_id` - An order ID is generated only in the monthly or yearly billing mode.
    In other billing modes, no value is returned for this parameter.
* `port` - Port of the cache node.
* `resource_spec_code` - Resource specifications.
    dcs.single_node: indicates a DCS instance in single-node mode.
    dcs.master_standby: indicates a DCS instance in master/standby mode.
    dcs.cluster: indicates a DCS instance in cluster mode.
* `used_memory` - Size of the used memory. Unit: MB.
* `internal_version` - Internal DCS version.
* `max_memory` - Overall memory size. Unit: MB.
* `user_id` - Indicates a user ID.
* `ip` - Cache node's IP address in tenant's VPC.
