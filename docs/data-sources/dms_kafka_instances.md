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

* `storage_type` - TThe storage type of the instance.

* `storage_resource_id` - The storage resource ID of the instance.

* `vpc_id` - The VPC ID to which the instance belongs.

* `vpc_name` - The VPC name to which the instance belongs.

* `vpc_client_plain` - Whether the intra-VPC plaintext access is enabled.

* `network_id` - The subnet ID to which the instance belongs.

* `subnet_name` - The subnet name to which the instance belongs.

* `subnet_cidr` - The CIDR of the subnet to which the instance belongs.

* `security_group_id` - The security group ID associated with the instance.

* `security_group_name` - The security group name associated with the instance.

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

* `connector_id` - The ID of the dump task.

* `connector_node_num` - The number of dump nodes.

* `enable_auto_topic` - Whether to enable automatic topic creation.

* `partition_num` - The maximum number of topics in the DMS kafka instance.

* `ssl_enable` - Whether the Kafka SASL_SSL is enabled.

* `used_storage_space` - The used message storage space, in GB unit.

* `connect_address` - The IP address for instance connection.

* `port` - The port number of the instance.

* `status` - The instance status.

* `resource_spec_code` - The resource specifications identifier.

* `specification` - The specification of the instance.

* `user_id` - The user ID who created the instance.

* `user_name` - The username who created the instance.

* `management_connect_address` - The connection address of the Kafka manager of an instance.

* `tags` - The key/value pairs to associate with the instance.

* `cross_vpc_accesses` - Indicates the Access information of cross-VPC.  
  The [cross_vpc_accesses](#kafka_instance_cross_vpc_accesses_attr) structure is documented below.

* `broker_num` - The number of brokers in the instance.

* `ces_version` - The CES version corresponding to the instance.

* `charging_mode` - The charging mode of the instance.

* `extend_times` - The extend times.

* `ipv6_enable` - Whether the IPv6 is enabled.

* `ipv6_connect_addresses` - The IPv6 connect addresses of the instance.

* `is_logical_volume` - Whether the expansion is new instance.
  + **true**: The new instance, allowing disk dynamic expansion without restart.
  + **false**: The old instance.

* `message_query_inst_enable` - Whether message query is enabled.

* `enable_log_collection` - Whether log collection is enabled.

* `new_auth_cert` - Whether the new auth cert is enabled.

* `new_spec_billing_enable` - Whether the new billing specification is enabled.

* `node_num` - The number of nodes in the instance.

* `order_id` - The order ID of the instance.

* `pod_connect_address` - The connection address on the tenant side.

* `port_protocol` - The port protocol of the instance.  
  The [port_protocol](#kafka_instances_port_protocol_attr) structure is documented below.

* `support_features` - The support features of the instance.

* `ssl_two_way_enable` - Whether the SSL two-way authentication is enabled.

* `public_bandwidth` - The public bandwidth of the instance.

* `public_boundwidth` - The public boundwidth of the instance.

* `created_at` - The creation time of the instance, in RFC3339 format.

<a name="kafka_instance_cross_vpc_accesses_attr"></a>
The `cross_vpc_accesses` block supports:

* `listener_ip` - The listener IP address.

* `advertised_ip` - The advertised IP Address.

* `port` - The port number.

* `port_id` - The port ID associated with the address.

<a name="kafka_instances_port_protocol_attr"></a>
The `port_protocol` block supports:

* `private_plain_enable` - Whether private plaintext access is enabled.

* `private_sasl_ssl_enable` - Whether private SASL SSL access is enabled.

* `private_sasl_plaintext_enable` - Whether private SASL plaintext access is enabled.

* `public_plain_enable` - Whether public plaintext access is enabled.

* `public_sasl_ssl_enable` - Whether public SASL SSL access is enabled.

* `public_sasl_plaintext_enable` - Whether public SASL plaintext access is enabled.

* `private_plain_address` - The private plain address.

* `private_sasl_ssl_address` - The private SASL SSL address.

* `private_sasl_plaintext_address` - The private SASL plaintext address.

* `public_plain_address` - The public plain address.

* `public_sasl_ssl_address` - The public SASL SSL address.

* `public_sasl_plaintext_address` - The public SASL plaintext address.

* `private_plain_domain_name` - The private plain domain name.

* `private_sasl_ssl_domain_name` - The private SASL SSL domain name.

* `private_sasl_plaintext_domain_name` - The private SASL plaintext domain name.

* `public_plain_domain_name` - The public plain domain name.

* `public_sasl_ssl_domain_name` - The public SASL SSL domain name.

* `public_sasl_plaintext_domain_name` - The public SASL plaintext domain name.
