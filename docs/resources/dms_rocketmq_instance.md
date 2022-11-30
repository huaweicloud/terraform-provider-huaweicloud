---
subcategory: "Distributed Message Service (DMS)"
---

# huaweicloud_dms_rocketmq_instance

Manage DMS RocketMQ instance resources within HuaweiCloud.

## Example Usage

```HCL
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "availability_zones" {
  type = list(string)
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name                = "rocketmq_name_test"
  description         = "this is a rocketmq instance"
  engine_version      = "4.8.0"
  storage_space       = 300
  vpc_id              = var.vpc_id
  subnet_id           = var.subnet_id
  security_group_id   = var.security_group_id
  availability_zones  = var.availability_zones
  flavor_id          = "c6.4u8g.cluster"
  storage_spec_code   = "dms.physical.storage.high.v2"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the DMS RocketMQ instance.
  An instance name starts with a letter, consists of 4 to 64 characters, and can contain only letters,
  digits, underscores (_), and hyphens (-).

* `engine_version` - (Required, String, ForceNew) Specifies the version of the RocketMQ engine. Value: 4.8.0.
  Changing this parameter will create a new resource.

* `storage_space` - (Required, Int, ForceNew) Specifies the message storage capacity, Unit: GB.
  Value range: 300-3000.
  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC.
  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of a subnet.
  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String) Specifies the ID of a security group.

* `availability_zones` - (Required, List, ForceNew) Specifies the list of availability zone names, where
  instance brokers reside and which has available resources.

  Changing this parameter will create a new resource.

* `flavor_id` - (Required, String, ForceNew) Specifies a product ID. The options are as follows:
  + **c6.4u8g.cluster**: maximum number of topics on each broker: 4000; maximum number of consumer groups
    on each broker: 4000
  + **c6.8u16g.cluster**: maximum number of topics on each broker: 8000; maximum number of consumer groups
    on each broker: 8000
  + **c6.12u24g.cluster**: maximum number of topics on each broker: 12,000; maximum number of consumer groups
    on each broker: 12,000
  + **c6.16u32g.cluster**: maximum number of topics on each broker: 16,000; maximum number of consumer groups
    on each broker: 16,000
  Changing this parameter will create a new resource.

* `storage_spec_code` - (Required, String, ForceNew) Specifies the storage I/O specification.
  The options are as follows:
  + **dms.physical.storage.high.v2**: high I/O disk
  + **dms.physical.storage.ultra.v2**: ultra-high I/O disk
  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the DMS RocketMQ instance.
  The description can contain a maximum of 1024 characters.

* `ssl_enable` - (Optional, Bool, ForceNew) Specifies whether the RocketMQ SASL_SSL is enabled. Defaults to false.
  Changing this parameter will create a new resource.

* `ipv6_enable` - (Optional, Bool, ForceNew) Specifies whether to support IPv6. Defaults to false.
  Changing this parameter will create a new resource.

* `enable_publicip` - (Optional, Bool, ForceNew) Specifies whether to enable public access.
  By default, public access is disabled.
  Changing this parameter will create a new resource.

* `publicip_id` - (Optional, String, ForceNew) Specifies the ID of the EIP bound to the instance.
  Use commas (,) to separate multiple EIP IDs.
  This parameter is mandatory if public access is enabled (that is, enable_publicip is set to true).
  Changing this parameter will create a new resource.

* `broker_num` - (Optional, Int, ForceNew) Specifies the broker numbers. Defaults to 1.
  Changing this parameter will create a new resource.

* `retention_policy` - (Optional, Bool) Specifies the ACL access control.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `status` - Indicates the status of the DMS RocketMQ instance.

* `type` - Indicates the DMS RocketMQ instance type. Value: cluster.

* `specification` - Indicates the instance specification. For a cluster DMS RocketMQ instance, VM specifications
  and the number of nodes are returned.

* `maintain_begin` - Indicates the time at which the maintenance window starts. The format is HH:mm:ss.

* `maintain_end` - Indicates the time at which the maintenance window ends. The format is HH:mm:ss.

* `used_storage_space` - Indicates the used message storage space. Unit: GB.

* `publicip_address` - Indicates the public IP address.

* `cross_vpc_info` - Indicates the Cross-VPC access information.

* `node_num` - Indicates the node quantity.

* `new_spec_billing_enable` - Indicates whether billing based on new specifications is enabled.

* `enable_acl` - Indicates whether access control is enabled.

* `namesrv_address` - Indicates the metadata address.

* `broker_address` - Indicates the service data address.

* `public_namesrv_address` - Indicates the public network metadata address.

* `public_broker_address` - Indicates the public network service data address.

* `resource_spec_code` - Indicates the resource specifications.

## Import

The rocketmq instance can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_dms_rocketmq_instance.test 8d3c7938-dc47-4937-a30f-c80de381c5e3
```
