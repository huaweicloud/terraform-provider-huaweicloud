---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_instance"
description: ""
---

# huaweicloud_dms_rocketmq_instance

Manage DMS RocketMQ instance resources within HuaweiCloud.

## Example Usage

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "availability_zones" {
  type = list(string)
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name               = "rocketmq_name_test"
  description        = "this is a rocketmq instance"
  engine_version     = "4.8.0"
  storage_space      = 300
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  security_group_id  = var.security_group_id
  availability_zones = var.availability_zones
  flavor_id          = "c6.4u8g.cluster"
  storage_spec_code  = "dms.physical.storage.high.v2"
  broker_num         = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the DMS RocketMQ instance.
  An instance name starts with a letter, consists of `4` to `64` characters, and can contain only letters,
  digits, underscores (_), and hyphens (-).

* `engine_version` - (Required, String, ForceNew) Specifies the version of the RocketMQ engine.
  Valid values are **4.8.0** and **5.x**.
  Changing this parameter will create a new resource.

* `storage_space` - (Required, Int) Specifies the message storage capacity, Unit: GB.
  When `engine_version` is **4.8.0**, value ranges from `300` to `30,000`.
  When `engine_version` is **5.x**, value ranges from `200` to `60,000`.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC.
  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of a subnet.
  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String) Specifies the ID of a security group.

* `availability_zones` - (Required, List, ForceNew) Specifies the list of availability zone names, where
  instance brokers reside and which has available resources.

  Changing this parameter will create a new resource.

* `flavor_id` - (Required, String) Specifies the flavor ID.

* `storage_spec_code` - (Required, String, ForceNew) Specifies the storage I/O specification.
  The options are as follows:
  + **dms.physical.storage.high.v2**: high I/O disk
  + **dms.physical.storage.ultra.v2**: ultra-high I/O disk
  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the DMS RocketMQ instance.
  The description can contain a maximum of `1,024` characters.

* `ssl_enable` - (Optional, Bool, ForceNew) Specifies whether the RocketMQ SASL_SSL is enabled. Defaults to **false**.
  Changing this parameter will create a new resource.  
  If this parameter is set to **true**, `tls_mode` can be omitted or must be set to **SSL**.
  If this parameter is set to **false**, `tls_mode` cannot be set to **SSL**.

* `ipv6_enable` - (Optional, Bool, ForceNew) Specifies whether to support IPv6. Defaults to **false**.
  Changing this parameter will create a new resource.

* `enable_publicip` - (Optional, Bool) Specifies whether to enable public access. By default, public access is disabled.

* `publicip_id` - (Optional, String) Specifies the ID of the EIP bound to the instance. Use commas (,) to separate
  multiple EIP IDs. It is mandatory if `enable_publicip` is **true** and should be empty when `enable_publicip` is **false**.

* `broker_num` - (Optional, Int) Specifies the broker numbers. It's **required** when instance architecture is
  **cluster**. Defaults to `1` when instance architecture is **single node**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the instance.

* `enable_acl` - (Optional, Bool) Specifies whether access control is enabled.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the instance. Valid values are *prePaid*
  and *postPaid*, defaults to *postPaid*. Changing this creates a new resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the instance.
  Valid values are *month* and *year*. This parameter is mandatory if `charging_mode` is set to *prePaid*.
  Changing this creates a new resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the instance. If `period_unit` is set to *month*
  , the value ranges from 1 to 9. If `period_unit` is set to *year*, the value ranges from `1` to `3`. This parameter is
  mandatory if `charging_mode` is set to *prePaid*. Changing this creates a new resource.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled. Valid values are "true" and "false".

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the instance.

* `configs` - (Optional, List) Specifies the instance configs.
  The [configs](#dms_configs) structure is documented below.

* `tls_mode` - (Optional, String) Specifies TLS mode of the instance.  
  The valid values are as follows:
  + **PLAINTEXT**
  + **SSL**
  + **PERMISSIVE**

<a name="dms_configs"></a>
The `configs` block supports:

* `name` - (Required, String) Specifies the config name.

* `value` - (Required, String) Specifies the config value.

-> When `name` is **fileReservedTime**, `value` ranges from `1` to `720` and unit is **hour**. Defaults to `48`.

## Attribute Reference

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

* `namesrv_address` - Indicates the metadata address.

* `broker_address` - Indicates the service data address.

* `public_namesrv_address` - Indicates the public network metadata address.

* `public_broker_address` - Indicates the public network service data address.

* `resource_spec_code` - Indicates the resource specifications.

* `cross_vpc_accesses` - Indicates the Access information of cross-VPC. The structure is documented below.

The `cross_vpc_accesses` block supports:

* `advertised_ip` - The advertised IP Address or domain name.
* `listener_ip` - The listener IP address.
* `port` - The port number.
* `port_id` - The port ID associated with the address.

## Import

The rocketmq instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dms_rocketmq_instance.test 8d3c7938-dc47-4937-a30f-c80de381c5e3
```
