---
subcategory: "Distributed Message Service (DMS)"
---

# huaweicloud_dms_rabbitmq_instance

Manage DMS RabbitMQ instance resources within HuaweiCloud.

## Example Usage

### Basic Instance

```hcl
variable vpc_id {}
variable subnet_id {}
variable security_group_id {}

data "huaweicloud_availability_zones" "zones" {}

data "huaweicloud_dms_product" "test" {
  engine        = "rabbitmq"
  instance_type = "cluster"
  version       = "3.7.17"
  node_num      = 3
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name              = "instance_1"
  product_id        = data.huaweicloud_dms_product.test.id
  engine_version    = data.huaweicloud_dms_product.test.version
  storage_spec_code = data.huaweicloud_dms_product.test.storage_spec_code

  vpc_id             = var.vpc_id
  network_id         = var.subnet_id
  security_group_id  = var.security_group_id
  availability_zones = [
    data.huaweicloud_availability_zones.zones.names[0]
  ]

  access_user = "user"
  password    = "Rabbitmqtest@123"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DMS RabbitMQ instance resource. If omitted,
  the provider-level region will be used. Changing this creates a new instance resource.

* `name` - (Required, String) Specifies the name of the DMS RabbitMQ instance. An instance name starts with a letter,
  consists of 4 to 64 characters, and supports only letters, digits, hyphens (-) and underscores (_).

* `product_id` - (Required, String) Specifies a product ID. Changing this creates a new instance resource.

  -> **NOTE:** Change this will change number of nodes and storage capacity. If you specify the value of
  `storage_space`, you need to manually modify the value of `storage_space` after changing the `product_id`.

* `engine_version` - (Optional, String, ForceNew) Specifies the version of the RabbitMQ engine. Default to "3.7.17".
  Changing this creates a new instance resource.

* `storage_spec_code` - (Required, String, ForceNew) Specifies the storage I/O specification.
  Valid values are **dms.physical.storage.high** and **dms.physical.storage.ultra**.
  Changing this creates a new instance resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC. Changing this creates a new instance resource.

* `network_id` - (Required, String, ForceNew) Specifies the ID of a subnet. Changing this creates a new instance
  resource.

* `security_group_id` - (Required, String) Specifies the ID of a security group.

* `availability_zones` - (Required, List, ForceNew) Specifies the names of an AZ.
  The parameter value can not be left blank or an empty array.
  Changing this creates a new instance resource.

  ~> The parameter behavior of `availability_zones` has been changed from `list` to `set`.

* `access_user` - (Required, String, ForceNew) Specifies a username. A username consists of 4 to 64 characters and
  supports only letters, digits, and hyphens (-). Changing this creates a new instance resource.

* `password` - (Required, String, ForceNew) Specifies the password of the DMS RabbitMQ instance. A password must meet
  the following complexity requirements: Must be 8 to 32 characters long. Must contain at least 2 of the following
  character types: lowercase letters, uppercase letters, digits,
  and special characters (`~!@#$%^&*()-_=+\\|[{}]:'",<.>/?).
  Changing this creates a new instance resource.

* `storage_space` - (Optional, Int, ForceNew) Specifies the message storage space, unit is GB. Value range:
  + Single-node RabbitMQ instance: 100–90,000 GB
  + Cluster RabbitMQ instance: 100 GB x number of nodes to 90,000 GB, 200 GB x number of nodes to 90,000 GB,
    and 300 GB x number of nodes to 90,000 GB

  The storage capacity of the product used by default. Changing this creates a new instance resource.

* `description` - (Optional, String) Specifies the description of the DMS RabbitMQ instance.
  It is a character string containing not more than 1,024 characters.

* `maintain_begin` - (Optional, String) Specifies the time at which a maintenance time window starts. Format: HH:mm.
  The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance
  time window.
  The start time must be set to 22:00, 02:00, 06:00, 10:00, 14:00, or 18:00. Parameters `maintain_begin`
  and `maintain_end` must be set in pairs. If parameter `maintain_begin` is left blank, parameter `maintain_end` is also
  blank. In this case, the system automatically allocates the default start time 02:00.

* `maintain_end` - (Optional, String) Specifies the time at which a maintenance time window ends. Format: HH:mm.
  The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance
  time window. The end time is four hours later than the start time.
  For example, if the start time is 22:00, the end time is 02:00.
  Parameters `maintain_begin` and `maintain_end` must be set in pairs.
  If parameter `maintain_end` is left  blank, parameter `maintain_begin` is also blank.
  In this case, the system automatically allocates the default end time 06:00.

* `ssl_enable` - (Optional, String, ForceNew) Specifies whether to enable public access for the DMS RabbitMQ instance.
  Changing this creates a new instance resource.

* `public_ip_id` - (Optional, String) Specifies the ID of the elastic IP address (EIP)
  bound to the DMS RabbitMQ instance.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the RabbitMQ instance.

* `tags` - (Optional, Map) The key/value pairs to associate with the DMS RabbitMQ instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `engine` - Indicates the message engine.
* `specification` - Indicates the instance specification. For a single-node DMS RabbitMQ instance, VM specifications are
  returned. For a cluster DMS RabbitMQ instance, VM specifications and the number of nodes are returned.
* `used_storage_space` - Indicates the used message storage space. Unit: GB
* `port` - Indicates the port number of the DMS RabbitMQ instance.
* `status` - Indicates the status of the DMS RabbitMQ instance.
* `enable_public_ip` - Indicates whether public access to the DMS RabbitMQ instance is enabled.
* `resource_spec_code` - Indicates a resource specifications identifier.
* `type` - Indicates the DMS RabbitMQ instance type.
* `user_id` - Indicates the ID of the user who created the DMS RabbitMQ instance
* `user_name` - Indicates the name of the user who created the DMS RabbitMQ instance
* `connect_address` - Indicates the IP address of the DMS RabbitMQ instance.
* `manegement_connect_address` - Indicates the management address of the DMS RabbitMQ instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 50 minute.
* `update` - Default is 50 minute.
* `delete` - Default is 15 minute.

## Import

DMS RabbitMQ instance can be imported using the instance id, e.g.

```
 $ terraform import huaweicloud_dms_rabbitmq_instance.instance_1 8d3c7938-dc47-4937-a30f-c80de381c5e3
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`password`. It is generally recommended running `terraform plan` after importing a DMS RabbitMQ instance. You can then
decide if changes should be applied to the instance, or the resource definition should be updated to align with the
instance. Also you can ignore changes as below.

```
resource "huaweicloud_dms_rabbitmq_instance" "instance_1" {
    ...

  lifecycle {
    ignore_changes = [
      password
    ]
  }
}
```
