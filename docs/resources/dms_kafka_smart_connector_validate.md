---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_smart_connector_validate"
description: |-
  Use this resource to validate the connectivity between the Kafka instances within HuaweiCloud.
---

# huaweicloud_dms_kafka_smart_connector_validate

Use this resource to validate the connectivity between the Kafka instances within HuaweiCloud.

-> This resource is only a one-time action resource for validating Kafka instances connectivity. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "task_type" {}
variable "current_instance_alias" {}
variable "peer_instance_alias" {}
variable "peer_user_name" {}
variable "peer_password" {}
variable "peer_instance_id" {}

resource "huaweicloud_dms_kafka_smart_connector_validate" "test" {
  instance_id = var.instance_id
  type        = var.task_type

  task = {
    current_cluster_name          = var.current_instance_alias
    cluster_name                  = var.peer_instance_alias
    user_name                     = var.peer_user_name
    password                      = var.peer_password
    sasl_mechanism                = "SCRAM-SHA-512"
    instance_id                   = var.peer_instance_id
    security_protocol             = "SASL_SSL"
    direction                     = "push"
    sync_consumer_offsets_enabled = true
    rename_topic_enabled          = true
    replication_factor            = 3
    task_num                      = 2
    provenance_header_enabled     = true
    consumer_strategy             = "earliest"
    compression_type              = "none"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the Smart Connect is located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Kafka instance to which the
  Smart Connect belongs.

* `type` - (Optional, String, NonUpdatable) Specifies the type of the Smart Connect task.  
  The valid values are as follows:
  + **KAFKA_REPLICATOR_SOURCE**

* `task` - (Optional, List, NonUpdatable) Specifies the configuration of the Smart Connect task.  
  The [task](#kafka_smart_connector_task_validate) structure is documented below.

<a name="kafka_smart_connector_task_validate"></a>
The `task` block supports:

* `current_cluster_name` - (Optional, String, NonUpdatable) Specifies the alias of the current instance.

* `instance_id` - (Optional, String, NonUpdatable) Specifies the ID of the peer instance.

* `cluster_name` - (Optional, String, NonUpdatable) Specifies the alias of the peer instance.

* `user_name` - (Optional, String, NonUpdatable) Specifies the username of the peer instance.

* `password` - (Optional, String, NonUpdatable) Specifies the password of the peer instance.

* `sasl_mechanism` - (Optional, String, NonUpdatable) Specifies the authentication mechanism of the peer instance.  
  The valid values are as follows:
  + **SCRAM-SHA-512**
  + **PLAIN**

* `bootstrap_servers` - (Optional, String, NonUpdatable) Specifies the address of the peer instance.  
  Multiple addresses are separated by commas (,).

* `security_protocol` - (Optional, String, NonUpdatable) Specifies the authentication method of the peer instance.  
  The valid values are as follows:
  + **SASL_SSL**
  + **PLAINTEXT**
  + **SASL_PLAINTEXT**

* `direction` - (Optional, String, NonUpdatable) Specifies the synchronization direction of the Smart Connect task.  
  The valid values are as follows:
  + **push**
  + **pull**
  + **two-way**

* `sync_consumer_offsets_enabled` - (Optional, Bool, NonUpdatable) Specifies whether to synchronize consumption
  progress. Defaults to **false**.

* `replication_factor` - (Optional, Int, NonUpdatable) Specifies the number of replicas of the Smart Connect task.

* `task_num` - (Optional, Int, NonUpdatable) Specifies the number of tasks of the data replication.

* `rename_topic_enabled` - (Optional, Bool, NonUpdatable) Specifies whether to rename topic.  
  Defaults to **false**. This parameter cannot together with `topics_mapping`.

* `provenance_header_enabled` - (Optional, Bool, NonUpdatable) Specifies whether to add source header.  
  Defaults to **false**.

* `consumer_strategy` - (Optional, String, NonUpdatable) Specifies the startup offset of the Smart Connect task.  
  The valid values are as follows:
  + **latest**
  + **earliest**

* `compression_type` - (Optional, String, NonUpdatable) Specifies the compression algorithm of the Smart Connect task.  
  The valid values are as follows:
  + **none**
  + **gzip**
  + **snappy**
  + **lz4**
  + **zstd**

* `topics_mapping` - (Optional, String, NonUpdatable) Specifies the topics mapping of the Smart Connect task.  
  The format is `source_topic_name:target_topic_name`, multiple topics are separated by commas (,).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
