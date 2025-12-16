---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_instance"
description: ""
---

# huaweicloud_dms_kafka_instance

Manage DMS Kafka instance resources within HuaweiCloud.

## Example Usage

### Create a Kafka instance using flavor ID

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "access_password" {}

variable "availability_zones" {
  default = ["your_availability_zones_a", "your_availability_zones_b", "your_availability_zones_c"]
}
variable "flavor_id" {
  default = "your_flavor_id, such: c6.2u4g.cluster"
}
variable "storage_spec_code" {
  default = "your_storage_spec_code, such: dms.physical.storage.ultra.v2"
}

# Query flavor information based on flavorID and storage I/O specification.
# Make sure the flavors are available in the availability zone.
data "huaweicloud_dms_kafka_flavors" "test" {
  type               = "cluster"
  flavor_id          = var.flavor_id
  availability_zones = var.availability_zones
  storage_spec_code  = var.storage_spec_code
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "kafka_test"
  vpc_id            = var.vpc_id
  network_id        = var.subnet_id
  security_group_id = var.security_group_id

  flavor_id          = data.huaweicloud_dms_kafka_flavors.test.flavor_id
  storage_spec_code  = data.huaweicloud_dms_kafka_flavors.test.flavors[0].ios[0].storage_spec_code
  availability_zones = var.availability_zones
  engine_version     = "2.7"
  storage_space      = 600
  broker_num         = 3

  ssl_enable  = true
  access_user = "user"
  password    = var.access_password

  parameters {
    name  = "min.insync.replicas"
    value = "2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DMS Kafka instances. If omitted, the
  provider-level region will be used. Changing this creates a new instance resource.

* `name` - (Required, String) Specifies the name of the DMS Kafka instance. An instance name starts with a letter,
  consists of 4 to 64 characters, and supports only letters, digits, hyphens (-) and underscores (_).

* `flavor_id` - (Optional, String) Specifies the Kafka [flavor ID](https://support.huaweicloud.com/intl/en-us/productdesc-kafka/Kafka-specification.html),
  e.g. **c6.2u4g.cluster**. This parameter and `product_id` are alternative.

  -> It is recommended to use `flavor_id` if the region supports it.

* `product_id` - (Optional, String) Specifies a product ID, which includes bandwidth, partition, broker and default
  storage capacity.

  -> **NOTE:** Change this to change the bandwidth, partition and broker of the Kafka instances. Please note that the
  broker changes may cause storage capacity changes. So, if you specify the value of `storage_space`, you need to
  manually modify the value of `storage_space` after changing the `product_id`.

* `engine_version` - (Required, String, ForceNew) Specifies the version of the Kafka engine,
  such as 1.1.0, 2.3.0, 2.7 or other supported versions. Changing this creates a new instance resource.

* `storage_spec_code` - (Required, String, ForceNew) Specifies the storage I/O specification.
  If the instance is created with `flavor_id`, the valid values are as follows:
  + **dms.physical.storage.high.v2**: Type of the disk that uses high I/O.
  + **dms.physical.storage.ultra.v2**: Type of the disk that uses ultra-high I/O.

  If the instance is created with `product_id`, the valid values are as follows:
  + **dms.physical.storage.high**: Type of the disk that uses high I/O.
    The corresponding bandwidths are **100MB** and **300MB**.
  + **dms.physical.storage.ultra**: Type of the disk that uses ultra-high I/O.
    The corresponding bandwidths are **100MB**, **300MB**, **600MB** and **1,200MB**.

  Changing this creates a new instance resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC. Changing this creates a new instance resource.

* `network_id` - (Required, String, ForceNew) Specifies the ID of a subnet. Changing this creates a new instance
  resource.

* `security_group_id` - (Required, String) Specifies the ID of a security group.

* `availability_zones` - (Required, List, ForceNew) The names of the AZ where the Kafka instances reside.
  The parameter value can not be left blank or an empty array. Changing this creates a new instance resource.

  -> **NOTE:** Deploy one availability zone or at least three availability zones. Do not select two availability zones.
  Deploy to more availability zones, the better the reliability and SLA coverage.
  [Learn more](https://support.huaweicloud.com/intl/en-us/kafka_faq/kafka-faq-200426002.html)

  ~> The parameter behavior of `availability_zones` has been changed from `list` to `set`.

* `ipv6_enable` - (Optional, Bool, ForceNew) Specifies whether to enable IPv6. Defaults to **false**.
  Changing this creates a new instance resource.

* `arch_type` - (Optional, String, ForceNew) Specifies the CPU architecture. Valid value is **X86**.
  Changing this creates a new instance resource.

* `storage_space` - (Optional, Int) Specifies the message storage capacity, the unit is GB.
  The storage spaces corresponding to the product IDs are as follows:
  + **c6.2u4g.cluster** (100MB bandwidth): `300` to `300,000` GB
  + **c6.4u8g.cluster** (300MB bandwidth): `300` to `600,000` GB
  + **c6.8u16g.cluster** (600MB bandwidth): `300` to `900,000` GB
  + **c6.12u12g.cluster**: `300` to `900,000` GB
  + **c6.16u32g.cluster** (1,200MB bandwidth): `300` to `900,000` GB

  It is required when creating an instance with `flavor_id`.

* `broker_num` - (Optional, Int) Specifies the broker numbers.
  It is required when creating an instance with `flavor_id`.

* `new_tenant_ips` - (Optional, List) Specifies the IPv4 private IP addresses for the new brokers.

  -> The number of specified IP addresses must be less than or equal to the number of new brokers.

* `access_user` - (Optional, String) Specifies the username of SASL_SSL user. A username consists of 4
  to 64 characters and supports only letters, digits, and hyphens (-). Changing this creates a new instance resource.

  -> This parameter can be modified only when SSL is enabled for the first time using the `port_protocol` parameter.
     This parameter cannot be modified after encrypted access is enabled.

* `password` - (Optional, String) Specifies the password of SASL_SSL user. A password must meet the following
  complexity requirements: Must be 8 to 32 characters long. Must contain at least 2 of the following character types:
  lowercase letters, uppercase letters, digits, and special characters (`~!@#$%^&*()-_=+\\|[{}]:'",<.>/?).

  -> **NOTE:** `access_user` and `password` is mandatory and available when `ssl_enable` is **true**.

* `security_protocol` - (Optional, String, ForceNew) Specifies the protocol to use after SASL is enabled. Value options:
  + **SASL_SSL**: Data is encrypted with SSL certificates for high-security transmission.
  + **SASL_PLAINTEXT**: Data is transmitted in plaintext with username and password authentication. This protocol only
    uses the SCRAM-SHA-512 mechanism and delivers high performance.
  
  Defaults to **SASL_SSL**. Changing this creates a new instance resource.

  -> If `port_protocol` is used to set the private network access security protocol and the public network access
  security protocol, this parameter is invalid.

* `enabled_mechanisms` - (Optional, List) Specifies the authentication mechanisms to use after SASL is
  enabled. Value options:
  + **PLAIN**: Simple username and password verification.
  + **SCRAM-SHA-512**: User credential verification, which is more secure than **PLAIN**.
  
  Defaults to [**PLAIN**]. Changing this creates a new instance resource.

  -> This parameter can be modified only when SSL is enabled for the first time using the `port_protocol` parameter.
     This parameter cannot be modified after encrypted access is enabled.

* `description` - (Optional, String) Specifies the description of the DMS Kafka instance. It is a character string
  containing not more than 1,024 characters.

* `maintain_begin` - (Optional, String) Specifies the time at which a maintenance time window starts. Format: HH:mm. The
  start time and end time of a maintenance time window must indicate the time segment of a supported maintenance time
  window. The start time must be set to 22:00, 02:00, 06:00, 10:00, 14:00, or 18:00. Parameters `maintain_begin`
  and `maintain_end` must be set in pairs. If parameter `maintain_begin` is left blank, parameter `maintain_end` is also
  blank. In this case, the system automatically allocates the default start time 02:00.

* `maintain_end` - (Optional, String) Specifies the time at which a maintenance time window ends. Format: HH:mm. The
  start time and end time of a maintenance time window must indicate the time segment of a supported maintenance time
  window. The end time is four hours later than the start time. For example, if the start time is 22:00, the end time is
  02:00. Parameters `maintain_begin`
  and `maintain_end` must be set in pairs. If parameter `maintain_end` is left blank, parameter
  `maintain_begin` is also blank. In this case, the system automatically allocates the default end time 06:00.

* `public_ip_ids` - (Optional, List) Specifies the IDs of the elastic IP address (EIP)
  bound to the DMS Kafka instance. Changing this creates a new instance resource.
  + If the instance is created with `flavor_id`, the total number of public IPs is equal to `broker_num`.
  + If the instance is created with `product_id`, the total number of public IPs must provide as follows:

  | Bandwidth | Total number of public IPs |
  | ---- | ---- |
  | 100MB | 3 |
  | 300MB | 3 |
  | 600MB | 4 |
  | 1,200MB | 8 |

  -> Only support to **add** public IP nums when `broker_num` adding and the instance is **created** with public IP
  **enabled** and using `flavor_id`.

  ~> The parameter behavior of `public_ip_ids` has been changed from `list` to `set`.

* `retention_policy` - (Optional, String) Specifies the action to be taken when the memory usage reaches the disk
  capacity threshold. The valid values are as follows:
  + **time_base**: Automatically delete the earliest messages.
  + **produce_reject**: Stop producing new messages.

* `dumping` - (Optional, Bool, ForceNew) Specifies whether to enable  message dumping(smart connect).
  Changing this creates a new instance resource.

* `enable_auto_topic` - (Optional, Bool) Specifies whether to enable automatic topic creation. If automatic
  topic creation is enabled, a topic will be automatically created with 3 partitions and 3 replicas when a message is
  produced to or consumed from a topic that does not exist.
  The default value is false.

* `parameters` - (Optional, List) Specifies the array of one or more parameters to be set to the Kafka instance after
  launched. The [parameters](#dms_parameters) structure is documented below.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the Kafka instance.

* `ssl_enable` - (Optional, Bool, ForceNew) Specifies whether the Kafka SASL_SSL is enabled.  
  Defaults to **false**.  
  Changing this creates a new resource.  
  When both `port_protocol` and `ssl_enable` parameters are set, `port_protocol` takes precedence.

* `vpc_client_plain` - (Optional, Bool, ForceNew) Specifies whether the intra-VPC plaintext access is enabled.
  Defaults to **false**. Changing this creates a new resource.

* `tags` - (Optional, Map) The key/value pairs to associate with the DMS Kafka instance.

* `cross_vpc_accesses` - (Optional, List) Specifies the cross-VPC access information.
  The [object](#dms_cross_vpc_accesses) structure is documented below.

* `port_protocol` - (Optional, List) Specifies the port protocol information.  
  The [object](#kafka_instance_port_protocol) structure is documented below.

* `disk_encrypted_enable` - (Optional, Bool, ForceNew) Specifies whether to enable disk encryption.  
  Defaults to **false**.  
  Changing this creates a new instance resource.

* `disk_encrypted_key` - (Optional, String, ForceNew) Specifies the key ID of the disk encryption.  
  Changing this creates a new instance resource.  
  This parameter is required when `disk_encrypted_enable` is set to **true**.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the instance. Valid values are *prePaid*
  and *postPaid*, defaults to *postPaid*. Changing this creates a new resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the instance.
  Valid values are *month* and *year*. This parameter is mandatory if `charging_mode` is set to *prePaid*.
  Changing this creates a new resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the instance. If `period_unit` is set to *month*
  , the value ranges from 1 to 9. If `period_unit` is set to *year*, the value ranges from 1 to 3. This parameter is
  mandatory if `charging_mode` is set to *prePaid*. Changing this creates a new resource.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled. Valid values are "true" and "false".

<a name="dms_cross_vpc_accesses"></a>
The `cross_vpc_accesses` block supports:

* `advertised_ip` - (Optional, String) The advertised IP Address or domain name.

<a name="dms_parameters"></a>
The `parameters` block supports:

* `name` - (Required, String) Specifies the parameter name. Static parameter needs to restart the instance to take effect.

* `value` - (Required, String) Specifies the parameter value.

<a name="kafka_instance_port_protocol"></a>
The `port_protocol` block supports:

* `private_plain_enable` - (Optional, Bool) Specifies whether the private plaintext access is enabled.

  -> The private plaintext access and private SSL access cannot be disabled at the same time.

* `private_sasl_ssl_enable` - (Optional, Bool) Specifies whether the private SASL SSL access is enabled.  
  This parameter and `private_sasl_plaintext_enable` cannot be set to `true` at the same time.

* `private_sasl_plaintext_enable` - (Optional, Bool) Specifies whether the private SASL plaintext access is enabled.

* `public_plain_enable` - (Optional, Bool) Specifies whether the public plaintext access is enabled.

* `public_sasl_ssl_enable` - (Optional, Bool) Specifies whether the public SASL SSL access is enabled.
  This parameter and `public_sasl_plaintext_enable` cannot be set to **true** at the same time.

* `public_sasl_plaintext_enable` - (Optional, Bool) Specifies whether the public SASL plaintext access is enabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `engine` - Indicates the message engine.
* `partition_num` - Indicates the number of partitions in Kafka instance.
* `used_storage_space` - Indicates the used message storage space. Unit: GB
* `port` - Indicates the port number of the DMS Kafka instance.
* `status` - Indicates the status of the DMS Kafka instance.
* `enable_public_ip` - Indicates whether public access to the DMS Kafka instance is enabled.
* `resource_spec_code` - Indicates a resource specifications identifier.
* `type` - Indicates the DMS Kafka instance type.
* `user_id` - Indicates the ID of the user who created the DMS Kafka instance
* `user_name` - Indicates the name of the user who created the DMS Kafka instance
* `connect_address` - Indicates the IP address of the DMS Kafka instance.
* `cross_vpc_accesses` - Indicates the Access information of cross-VPC.
  The [cross_vpc_accesses](#attr_cross_vpc_accesses) structure is documented below.
* `charging_mode` - Indicates the charging mode of the instance.
* `public_ip_address` - Indicates the public IP addresses list of the instance.
* `extend_times` - Indicates the extend times. If the value exceeds `20`, disk expansion is no longer allowed.
* `connector_id` - Indicates the connector ID.
* `connector_node_num` - Indicates the number of connector node.
* `storage_resource_id` - Indicates the storage resource ID.
* `storage_type` - Indicates the storage type.
* `ipv6_connect_addresses` - Indicates the IPv6 connect addresses list.
* `created_at` - Indicates the create time.
* `cert_replaced` - Indicates whether the certificate can be replaced.
* `is_logical_volume` - Indicates whether the instance is a new instance.
* `message_query_inst_enable` - Indicates whether message query is enabled.
* `node_num` - Indicates the node quantity.
* `pod_connect_address` - Indicates the connection address on the tenant side.
* `public_bandwidth` - Indicates the public network access bandwidth.
* `ssl_two_way_enable` - Indicates whether to enable two-way authentication.

* `port_protocol` - Indicates instance connection address. The structure is documented below.
  The [port_protocol](#dms_instance_port_protocol_attr) structure is documented below.

<a name="attr_cross_vpc_accesses"></a>
The `cross_vpc_accesses` block supports:

* `listener_ip` - The listener IP address.
* `port` - The port number.
* `port_id` - The port ID associated with the address.

<a name="dms_instance_port_protocol_attr"></a>
The `port_protocols` block supports:

* `private_plain_address` - The private plain address.

* `private_plain_domain_name` - The private plain domain name.

* `private_sasl_ssl_address` - The private sasl ssl address.

* `private_sasl_ssl_domain_name` - The private sasl ssl domain name.

* `private_sasl_plaintext_address` - The private sasl plaintext address.

* `private_sasl_plaintext_domain_name` - The private sasl plaintext domain name.

* `public_plain_address` - The public plain address.

* `public_plain_domain_name` - The public plain domain name.

* `public_sasl_ssl_address` - The public sasl ssl address.

* `public_sasl_ssl_domain_name` - The public sasl ssl domain name.

* `public_sasl_plaintext_address` - The public sasl plaintext address.

* `public_sasl_plaintext_domain_name` - The public sasl plaintext domain name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 50 minutes.
* `update` - Default is 50 minutes.
* `delete` - Default is 15 minutes.

## Import

DMS Kafka instance can be imported using the instance id, e.g.

```
 $ terraform import huaweicloud_dms_kafka_instance.instance_1 8d3c7938-dc47-4937-a30f-c80de381c5e3
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`password`, `security_protocol`, `enabled_mechanisms`, `arch_type` and `new_tenant_ips`.
It is generally recommended running `terraform plan` after importing
a DMS Kafka instance. You can then decide if changes should be applied to the instance, or the resource definition
should be updated to align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dms_kafka_instance" "instance_1" {
    ...

  lifecycle {
    ignore_changes = [
      password, security_protocol, enabled_mechanisms, arch_type, new_tenant_ips
    ]
  }
}
```
