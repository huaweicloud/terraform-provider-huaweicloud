---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_connections"
description: |-
  Use this data source to get the list of EG connections within HuaweiCloud.
---

# huaweicloud_eg_connections

Use this data source to get the list of EG connections within HuaweiCloud.

## Example Usage

### Query all connections

```hcl
data "huaweicloud_eg_connections" "test" {}
```

### Query connections by connection name

```hcl
variable "connection_name" {}

data "huaweicloud_eg_connections" "test" {
  name = var.connection_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the connections are located.  
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the exact name of the connection to be queried.

* `fuzzy_name` - (Optional, String) Specifies the name of the connection to be queried for fuzzy matching.

* `sort` - (Optional, String) Specifies the sorting method for query results.  
  The format is `field:order`, where `field` is the field name and `order` is `ASC` or `DESC`. e.g. `name:ASC`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - All connections that match the filter parameters.  
  The [connections](#data_attr_connections) structure is documented below.

<a name="data_attr_connections"></a>
The `connections` block supports:

* `id` - The ID of the connection.

* `name` - The name of the connection.

* `description` - The description of the connection.

* `vpc_id` - The ID of the VPC to which the connection belongs.

* `subnet_id` - The ID of the subnet to which the connection belongs.

* `type` - The type of the connection.
  + **WEBHOOK**
  + **KAFKA**

* `status` - The status of the connection.

* `kafka_detail` - The Kafka detail information for the connection.  
  The [kafka_detail](#data_connections_kafka_detail) structure is documented below.

* `agency` - The user delegation name used by the private network connection.

* `flavor` - The flavor information of the connection.
  The [flavor](#data_connections_flavor) structure is documented below.

* `created_time` - The creation time of the connection, in UTC format.

* `updated_time` - The latest update time of the connection, in UTC format.

* `error_info` - The error information of the connection.  
  The [error_info](#data_connections_error_info) structure is documented below.

<a name="data_connections_kafka_detail"></a>
The `kafka_detail` block supports:

* `instance_id` - The ID of the Kafka instance.

* `connect_address` - The connection address of the Kafka instance.

* `security_protocol` - The security protocol used for the connection.

* `enable_sasl_ssl` - Whether SASL_SSL is enabled for the Kafka instance.

* `user_name` - The username of the Kafka instance.

* `acks` - The number of confirmation signals the producer needs to receive to consider the message sent successfully.

* `address` - The connection address of Kafka instance.

<a name="data_connections_flavor"></a>
The `flavor` block supports:

* `name` - The name of the flavor.

* `concurrency_type` - The concurrency type of the flavor.

* `concurrency` - The concurrency value of the flavor.

* `bandwidth_type` - The bandwidth type of the flavor.

<a name="data_connections_error_info"></a>
The `error_info` block supports:

* `error_code` - The error code.

* `error_detail` - The detailed error information.

* `error_msg` - The error message.
