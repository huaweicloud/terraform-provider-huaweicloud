---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_instances"
description: |-
  Use this data source to query Kafka instance list within HuaweiCloud.
---

# huaweicloud_dms_kafka_instances

Use this data source to query Kafka instance list within HuaweiCloud.

## Example Usage

### Query all instances with the keyword in the name

```hcl
variable "keyword" {}

data "huaweicloud_dms_kafka_instances" "test" {
  name        = var.keyword
  fuzzy_match = true
}
```

### Query the instance with the specified name

```hcl
variable "instance_name" {}

data "huaweicloud_dms_kafka_instances" "test" {
  name = var.instance_name
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to query the kafka instance list.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the kafka instance ID to match exactly.

* `name` - (Optional, String) Specifies the kafka instance name for data-source queries.

* `fuzzy_match` - (Optional, Bool) Specifies whether to match the instance name fuzzily, the default is a exact
  match (`flase`).

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which all instances of the list
  belong.  
  This field is only valid for enterprise users. For enterprise users, if omitted, all instances under the enterprise
  project will be queried.

* `status` - (Optional, String) Specifies the kafka instance status for data-source queries.

* `include_failure` - (Optional, Bool) Specifies whether the query results contain instances that failed to create.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The list of instances that match the filter parameters.  
  The [instances](#kafka_instances_attr) structure is documented below.

<a name="kafka_instances_attr"></a>
The `instances` block supports:

* `id` - The instance ID.

* `type` - The instance type.

* `name` - The instance name.

* `description` - The instance description.

* `availability_zones` - The list of AZ names.

* `enterprise_project_id` - The enterprise project ID to which the instance belongs.

* `product_id` - The product ID used by the instance.

* `engine_version` - The kafka engine version.

* `storage_spec_code` - The storage I/O specification.

* `storage_space` - The message storage capacity, in GB unit.

* `vpc_id` - The VPC ID to which the instance belongs.

* `network_id` - The subnet ID to which the instance belongs.

* `security_group_id` - The security group ID associated with the instance.

* `manager_user` - The username for logging in to the Kafka Manager.

* `access_user` - The access username.

* `maintain_begin` - The time at which a maintenance time window starts, the format is `HH:mm`.

* `maintain_end` - The time at which a maintenance time window ends, the format is `HH:mm`.

* `enable_public_ip` - Whether public access to the instance is enabled.

* `public_ip_ids` - The IDs of the elastic IP address (EIP).

* `security_protocol` - The protocol to use after SASL is enabled.

* `enabled_mechanisms` - The authentication mechanisms to use after SASL is enabled.

* `public_conn_addresses` - The instance public access address.
  The format of each connection address is `{IP address}:{port}`.

* `retention_policy` - The action to be taken when the memory usage reaches the disk capacity threshold.

* `dumping` - Whether to dumping is enabled.

* `enable_auto_topic` - Whether to enable automatic topic creation.

* `partition_num` - The maximum number of topics in the DMS kafka instance.

* `ssl_enable` - Whether the Kafka SASL_SSL is enabled.

* `used_storage_space` - The used message storage space, in GB unit.

* `connect_address` - The IP address for instance connection.

* `port` - The port number of the instance.

* `status` - The instance status.

* `resource_spec_code` - The resource specifications identifier.

* `user_id` - The user ID who created the instance.

* `user_name` - The username who created the instance.

* `management_connect_address` - The connection address of the Kafka manager of an instance.

* `tags` - The key/value pairs to associate with the instance.

* `cross_vpc_accesses` - Indicates the Access information of cross-VPC.  
  The [cross_vpc_accesses](#kafka_instance_cross_vpc_accesses_attr) structure is documented below.

<a name="kafka_instance_cross_vpc_accesses_attr"></a>
The `cross_vpc_accesses` block supports:

* `listener_ip` - The listener IP address.

* `advertised_ip` - The advertised IP Address.

* `port` - The port number.

* `port_id` - The port ID associated with the address.
