---
subcategory: "Distributed Message Service (DMS)"
---

# huaweicloud_dms_kafka_instance

Manage DMS Kafka instance resources within HuaweiCloud.

## Example Usage

```hcl
variable vpc_id {}
variable subnet_id {}
variable security_group_id {}

data "huaweicloud_availability_zones" "zones" {}
data "huaweicloud_dms_product" "product_1" {
  engine            = "kafka"
  instance_type     = "cluster"
  version           = "2.3.0"
  bandwidth         = "100MB"
  storage_spec_code = "dms.physical.storage.ultra"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name               = "instance_1"
  product_id         = data.huaweicloud_dms_product.product_1.id
  engine_version     = data.huaweicloud_dms_product.product_1.version
  storage_spec_code  = data.huaweicloud_dms_product.product_1.storage_spec_code
  availability_zones = [
    data.huaweicloud_availability_zones.zones.names[0],
    data.huaweicloud_availability_zones.zones.names[1]
  ]
  
  vpc_id            = var.vpc_id
  network_id        = var.subnet_id
  security_group_id = var.security_group_id

  access_user = "user"
  password    = "Kafkatest@123"

  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DMS kafka instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new instance resource.

* `name` - (Required, String) Specifies the name of the DMS kafka instance. An instance name starts with a letter,
  consists of 4 to 64 characters, and supports only letters, digits, hyphens (-) and underscores (_).

* `product_id` - (Required, String) Specifies a product ID, which includes bandwidth, partition, broker and default
  storage capacity.

-> **NOTE:** Change this to change the bandwidth, partition and broker of the Kafka instances. Please note that the
broker changes may cause storage capacity changes. So, if you specify the value of `storage_space`, you need to manually
modify the value of `storage_space` after changing the `product_id`.

* `engine_version` - (Required, String, ForceNew) Specifies the version of the kafka engine. Valid values are "1.1.0"
  and "2.3.0". Changing this creates a new instance resource.

* `storage_spec_code` - (Required, String, ForceNew) Specifies the storage I/O specification. Value range:
  + When bandwidth is 100MB: dms.physical.storage.high or dms.physical.storage.ultra
  + When bandwidth is 300MB: dms.physical.storage.high or dms.physical.storage.ultra
  + When bandwidth is 600MB: dms.physical.storage.ultra
  + When bandwidth is 1,200MB: dms.physical.storage.ultra

  Changing this creates a new instance resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC. Changing this creates a new instance resource.

* `network_id` - (Required, String, ForceNew) Specifies the ID of a subnet. Changing this creates a new instance
  resource.

* `security_group_id` - (Required, String) Specifies the ID of a security group.

* `availability_zones` - (Required, List, ForceNew) The codes of the AZ where the Kafka instances reside.
  The parameter value can not be left blank or an empty array. Changing this creates a new instance resource.

* `manager_user` - (Required, String, ForceNew) Specifies the username for logging in to the Kafka Manager. The username
  consists of 4 to 64 characters and can contain letters, digits, hyphens (-), and underscores (_). Changing this
  creates a new instance resource.

* `manager_password` - (Required, String, ForceNew) Specifies the password for logging in to the Kafka Manager. The
  password must meet the following complexity requirements: Must be 8 to 32 characters long. Must contain at least 2 of
  the following character types: lowercase letters, uppercase letters, digits, and special characters (`~!@#$%^&*()-_
  =+\\|[{}]:'",<.>/?). Changing this creates a new instance resource.

* `storage_space` - (Optional, Int) Specifies the message storage capacity, the unit is GB. Value range:
  + When bandwidth is 100MB: 600–90000 GB
  + When bandwidth is 300MB: 1200–90000 GB
  + When bandwidth is 600MB: 2400–90000 GB
  + When bandwidth is 1200MB: 4800–90000 GB

  The storage capacity of the product used by default.

* `access_user` - (Optional, String, ForceNew) Specifies a username. A username consists of 4 to 64 characters and
  supports only letters, digits, and hyphens (-). Changing this creates a new instance resource.

* `password` - (Optional, String, ForceNew) Specifies the password of the DMS kafka instance. A password must meet the
  following complexity requirements: Must be 8 to 32 characters long. Must contain at least 2 of the following character
  types: lowercase letters, uppercase letters, digits, and special characters (`~!@#$%^&*()-_=+\\|[{}]:'",<.>/?).
  Changing this creates a new instance resource.

  -> **NOTE:** If `access_user` and `password` are specified, Kafka SASL_SSL will be automatically enabled.

* `description` - (Optional, String) Specifies the description of the DMS kafka instance. It is a character string
  containing not more than 1024 characters.

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

* `public_ip_ids` - (Optional, List, ForceNew) Specifies the IDs of the elastic IP address (EIP)
  bound to the DMS kafka instance. The num of IDs needed ranges:
  + When bandwidth is 100MB: 3
  + When bandwidth is 300MB: 3
  + When bandwidth is 600MB: 4
  + When bandwidth is 1200MB: 8

  Changing this creates a new instance resource.

* `retention_policy` - (Optional, String) Specifies the action to be taken when the memory usage reaches the disk
  capacity threshold. Value range:
  + `time_base`: Automatically delete the earliest messages.
  + `produce_reject`: Stop producing new messages.

* `dumping` - (Optional, Bool, ForceNew) Specifies whether to enable dumping. Changing this creates a new instance
  resource.

* `enable_auto_topic` - (Optional, Bool, ForceNew) Specifies whether to enable automatic topic creation. If automatic
  topic creation is enabled, a topic will be automatically created with 3 partitions and 3 replicas when a message is
  produced to or consumed from a topic that does not exist. Changing this creates a new instance resource.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the kafka instance.

* `tags` - (Optional, Map) The key/value pairs to associate with the DMS kafka instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `engine` - Indicates the message engine.
* `bandwidth` - The Bandwidth of a Kafka instance, that is the maximum amount of data transferred per unit time.
  Unit: MB.
* `partition_num` - Indicates the maximum number of topics in the DMS kafka instance.
* `used_storage_space` - Indicates the used message storage space. Unit: GB
* `port` - Indicates the port number of the DMS kafka instance.
* `status` - Indicates the status of the DMS kafka instance.
* `ssl_enable` - Indicates whether the Kafka SASL_SSL is enabled.
* `enable_public_ip` - Indicates whether public access to the DMS kafka instance is enabled.
* `resource_spec_code` - Indicates a resource specifications identifier.
* `type` - Indicates the DMS kafka instance type.
* `user_id` - Indicates the ID of the user who created the DMS kafka instance
* `user_name` - Indicates the name of the user who created the DMS kafka instance
* `connect_address` - Indicates the IP address of the DMS kafka instance.
* `manegement_connect_address` - Indicates the connection address of the Kafka Manager of a Kafka instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minute.
* `update` - Default is 20 minute.
* `delete` - Default is 15 minute.

## Import

DMS kafka instance can be imported using the instance id, e.g.

```
 $ terraform import huaweicloud_dms_kafka_instance.instance_1 8d3c7938-dc47-4937-a30f-c80de381c5e3
```

Note that the imported state may not be identical to your resource definition, due to some attrubutes missing from the
API response, security or some other reason. The missing attributes include:
`password`, `manager_password` and `public_ip_ids`. It is generally recommended running `terraform plan` after importing
a DMS kafka instance. You can then decide if changes should be applied to the instance, or the resource definition
should be updated to align with the instance. Also you can ignore changes as below.

```
resource "huaweicloud_dms_kafka_instance" "instance_1" {
    ...

  lifecycle {
    ignore_changes = [
      password, manager_password,
    ]
  }
}
```
