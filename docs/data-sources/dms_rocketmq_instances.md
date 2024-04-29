---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_instances"
description: ""
---

# huaweicloud_dms_rocketmq_instances

Use this data source to get the list of DMS RocketMQ instances.

## Example Usage

```hcl
data "huaweicloud_dms_rocketmq_instances" "test" {
  name = "rocketmq_name_test"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the DMS RocketMQ instance.

* `instance_id` - (Optional, String) Specifies the ID of the RocketMQ instance.

* `status` - (Optional, String) Specifies the status of the DMS RocketMQ instance.

* `exact_match_name` - (Optional, String) Specifies whether to search for the instance that precisely matches a
  specified instance name. Value options: **true**, **false**. Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - .

* `instances` - Indicates the list of DMS RocketMQ instances.
  The [Instance](#DmsRocketMQInstances_Instance) structure is documented below.

<a name="DmsRocketMQInstances_Instance"></a>
The `instances` block supports:

* `id` - Indicates the ID of the DMS RocketMQ instance.

* `name` - Indicates the name of the DMS RocketMQ instance.

* `status` - Indicates the status of the DMS RocketMQ instance.

* `description` - Indicates the description of the DMS RocketMQ instance.

* `type` - Indicates the DMS RocketMQ instance type.

* `specification` - Indicates the instance specification. For a cluster DMS RocketMQ instance, VM specifications
  and the number of nodes are returned.

* `engine_version` - Indicates the version of the RocketMQ engine.

* `vpc_id` - Indicates the ID of a VPC.

* `flavor_id` - Indicates a product ID.

* `security_group_id` - Indicates the ID of a security group.

* `subnet_id` - Indicates the ID of a subnet.

* `availability_zones` - Indicates the list of availability zone names, where
  instance brokers reside and which has available resources.

* `maintain_begin` - Indicates the time at which the maintenance window starts. The format is HH:mm:ss.

* `maintain_end` - Indicates the time at which the maintenance window ends. The format is HH:mm:ss.

* `storage_space` - Indicates the message storage capacity. Unit: GB.

* `used_storage_space` - Indicates the used message storage space. Unit: GB.

* `enable_publicip` - Indicates whether to enable public access.

* `publicip_id` - Indicates the ID of the EIP bound to the instance.
  Use commas (,) to separate multiple EIP IDs.
  This parameter is mandatory if public access is enabled (that is, enable_publicip is set to true).

* `publicip_address` - Indicates the public IP address.

* `ssl_enable` - Indicates whether the RocketMQ SASL_SSL is enabled. Defaults to false.

* `cross_vpc_accesses` - Indicates the Cross-VPC access information.
  The [CrossVpc](#DmsRocketMQInstances_InstanceCrossVpc) structure is documented below.

* `storage_spec_code` - Indicates the storage I/O specification.

* `ipv6_enable` - Indicates whether to support IPv6. Defaults to false.

* `node_num` - Indicates the node quantity.

* `new_spec_billing_enable` - Indicates the whether billing based on new specifications is enabled.

* `enable_acl` - Indicates whether access control is enabled.

* `broker_num` - Specifies the broker numbers. Defaults to 1.

* `namesrv_address` - Indicates the metadata address.

* `broker_address` - Indicates the service data address.

* `public_namesrv_address` - Indicates the public network metadata address.

* `public_broker_address` - Indicates the public network service data address.

* `resource_spec_code` - Indicates the resource specifications.

<a name="DmsRocketMQInstances_InstanceCrossVpc"></a>
The `cross_vpc_accesses` block supports:

* `listener_ip` - Indicates the IP of the listener.

* `advertised_ip` - Indicates the advertised IP.

* `port` - Indicates the port.

* `port_id` - Indicates the port ID associated with the address.
