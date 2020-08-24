---
layout: "huaweicloud"
page_title: "huaweicloud: huaweicloud_dcs_instance"
sidebar_current: "docs-huaweicloud-resource-dcs-instance"
description: |-
  Manages a DCS instance in the huaweicloud DCS Service
---

# huaweicloud\_dcs\_instance

Manages a DCS instance in the huaweicloud DCS Service.
This is an alternative to `huaweicloud_dcs_instance_v1`

## Example Usage

### Automatically detect the correct network

```hcl
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "secgroup_1"
}
data "huaweicloud_dcs_az" "az_1" {
  port = "8002"
  code = "cn-north-1a"
}

resource "huaweicloud_vpc" "vpc_4" {
  name = "terraform_provider_vpc4"
  cidr= "192.168.0.0/16"
}
resource "huaweicloud_vpc_subnet" "subnet_1" {
  name = "huaweicloud_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = huaweicloud_vpc.vpc_4.id
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name              = "test_dcs_instance"
  engine_version    = "3.0"
  password          = "Huawei_test"
  engine            = "Redis"
  capacity          = 2
  vpc_id            = huaweicloud_vpc.vpc_4.id
  security_group_id = "${huaweicloud_networking_secgroup.secgroup_1.id}"
  subnet_id         = huaweicloud_vpc_subnet.subnet_1.id
  available_zones   = ["${data.huaweicloud_dcs_az.az_1.id}"]
  product_id        = "dcs.master_standby-h"
  save_days         = 1
  backup_type       = "manual"
  begin_at          = "00:00-01:00"
  period_type       = "weekly"
  backup_at         = [1]
  depends_on        = ["huaweicloud_networking_secgroup.secgroup_1"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Indicates the name of an instance. It starts with English characters 
    and can only be composed of English letters, numbers, underscores and underscores. 
    When creating a single instance, the name is a string of 4 to 64 bits in length. 
    When creating instances in batches, the length of the name is a string of 4 to 56 characters, 
    and the format of the instance name is "custom name-n", where n starts from 000 and increases in sequence.
    For example, if you create two instances in batches and the custom name is dcs_demo, 
    the names of the two instances are dcs_demo-000 and dcs_demo-001.

* `description` - (Optional) Indicates the description of an instance. It is a character
    string containing not more than 1024 characters.

* `engine` - (Required) Indicates a cache engine. Options: Redis and Memcached. Changing this
    creates a new instance.

* `engine_version` - (Optional) Indicates the version of a message engine.When the cache engine is Redis, 
    the value is 3.0, 4.0 or 5.0. 
    Changing this creates a new instance.

* `capacity` - (Required) Indicates the Cache capacity. Unit: GB.
    Redis3.0: Stand-alone and active/standby type instance values: 2, 4, 8, 16, 32, 64. 
    Proxy cluster instance specifications support 64, 128, 256, 512, and 1024.

    Redis4.0 and Redis5.0: Stand-alone and active/standby type instance 
    values: 0.125, 0.25, 0.5, 1, 2, 4, 8, 16, 32, 64. Cluster instance specifications 
    support 24, 32, 48, 64, 96, 128, 192, 256, 384, 512, 768, 1024.

    Memcached: Stand-alone and active/standby type instance values: 2, 4, 8, 16, 32, 64.
    Changing this creates a new instance.

* `access_user` - (Optional) Username used for accessing a DCS instance after password
    authentication. A username starts with a letter, consists of 1 to 64 characters,
    and supports only letters, digits, and hyphens (-).
    - When the cache engine is Memcached, this parameter is optional.
    - When the cache engine is Redis, this parameter does not need to be set.
    Changing this creates a new instance.

* `password` - (Optional) Password of a DCS instance.
    The password of a DCS Redis instance must meet the following complexity requirements:
    - Enter a string of 8 to 32 bits in length.
    - The new password cannot be the same as the old password.
    - Must contain three combinations of the following four characters: Lower case letters,
        uppercase letter, digital, Special characters include (`~!@#$%^&*()-_=+|[{}]:'",<.>/?).
    Changing this creates a new instance.

* `vpc_id` - (Required) Specifies the id of the VPC.
    Changing this creates a new instance.

* `security_group_id` - (Required) Tenant's security group ID.

* `subnet_id` - (Required) Specifies the id of the subnet.
    Changing this creates a new instance.

* `available_zones` - (Required) IDs of the AZs where cache nodes reside.
    If you are creating active/standby, Proxy cluster, and Cluster cluster instances to support 
    cross-zone deployment, you can specify the standby zone for the standby node. When specifying 
    availability zones for nodes, separate them with commas.
    Changing this creates a new instance.

* `product_id` - (Required) Product ID or Names used to differentiate DCS instance types.
    Changing this creates a new instance.

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

* `save_days` - (Optional) Retention time. Unit: day. Range: 1–7. Changing this creates a new instance.

* `backup_type` - (Optional) Backup type. Options:
    auto: automatic backup.
    manual: manual backup.
    Changing this creates a new instance.

* `begin_at` - (Optional) Time at which backup starts. "00:00-01:00" indicates that backup
    starts at 00:00:00. Changing this creates a new instance.

* `period_type` - (Optional) Interval at which backup is performed. Currently, only weekly
    backup is supported. Changing this creates a new instance.

* `backup_at` - (Optional) Day in a week on which backup starts. Range: 1–7. Where: 1
    indicates Monday; 7 indicates Sunday. Changing this creates a new instance.

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
