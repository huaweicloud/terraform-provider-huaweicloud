---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_register_kafka_instance"
description: |-
  Use this resource to register the Kafka instance to LTS within HuaweiCloud.
---


# huaweicloud_lts_register_kafka_instance

Use this resource to register the Kafka instance to LTS within HuaweiCloud.

-> 1. Before registering a Kafka instance, please configure the inbound rules of the security group to which the Kafka
  instance belongs. Please refer to the [document](https://support.huaweicloud.com/intl/en-us/usermanual-lts/lts_04_0043.html).
  <br>2. The same Kafka instance can only be registered once.
  <br>3. This resource is only a one-time action resource for registering the Kafka instance to LTS. Deleting this resource
  will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "kafka_instance_id" {}
variable "kafka_instance_name" {}
variable "kafka_instance_access_user" {}
variable "kafka_instance_access_password" {}

resource "huaweicloud_lts_register_kafka_instance" "test" {
  instance_id = var.kafka_instance_id
  kafka_name  = var.kafka_instance_name

  connect_info {
    user_name = var.kafka_instance_access_user
    pwd       = var.kafka_instance_access_password
  } 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Kafka instance to be registered to the LTS.

* `kafka_name` - (Required, String, NonUpdatable) Specifies the name of the Kafka instance to be registered to the LTS.

* `connect_info` - (Optional, List, NonUpdatable) Specifies the connection information of the Kafka instance to be
  registered to the LTS.  
  The [connect_info](#register_kafka_to_lts_connect_info) structure is documented below.  
  This parameter is available and required only when the registered Kafka instance is encrypted access.

<a name="register_kafka_to_lts_connect_info"></a>
The `connect_info` block supports:

* `user_name` - (Optional, String, NonUpdatable) Specifies the name of the SASL_SSL user of the Kafka instance.

* `pwd` - (Optional, String, NonUpdatable) Specifies the password of the SASL_SSL user of the Kafka instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
