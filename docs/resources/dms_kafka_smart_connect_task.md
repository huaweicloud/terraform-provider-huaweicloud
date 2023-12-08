---
subcategory: "Distributed Message Service (DMS)"
---

# huaweicloud_dms_kafka_smart_connect_task

Manage DMS kafka smart connect task resource within HuaweiCloud.

## Example Usage

```hcl
variable "connector_id" {}

resource "huaweicloud_dms_kafka_smart_connect_task" "test" {
  connector_id     = var.connector_id
  source_type      = "BLOB"
  task_name        = "task12"
  destination_type = "OBS"

  obs_destination_descriptor {   
    topics                = "topic-test"
    consumer_strategy     = "latest"
    destination_file_type = "TEXT"
    access_key            = ""
    secret_key            = ""
    obs_bucket_name       = "afdafd"
    obs_path              = "afdasfd"
    partition_format      = "yyyy/MM/dd/HH/mm"
    record_delimiter      = "\n"
    deliver_time_interval = 300
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `connector_id` - (Required, String, ForceNew) Specifies the connector_id of the kafka instance.

  Changing this parameter will create a new resource.

* `source_type` - (Required, String, ForceNew) Specifies the source type of the smart connect task.
  Defaults to **BLOB**.

  Changing this parameter will create a new resource.

* `task_name` - (Required, String, ForceNew) Specifies the name of the smart connect task.

  Changing this parameter will create a new resource.

* `destination_type` - (Required, String, ForceNew) Specifies the destination type of the smart connect task.
  Defaults to **OBS**.

  Changing this parameter will create a new resource.

* `obs_destination_descriptor` - The list of obs destination descriptor.
  The [obs_destination_descriptor](#smart_connect_obs_destination_descriptor) structure is documented below.

<a name="smart_connect_obs_destination_descriptor"></a>
  The `obs_destination_descriptor` block supports:

* `access_key` - (Required, String, ForceNew) Specifies the access key of the HuaweiCloud.

* `secret_key` - (Required, String, ForceNew) Specifies the secret_key key of the HuaweiCloud.

* `consumer_strategy` - (Required, String, ForceNew) Specifies the consumer strategy of the smart connect task.
  The value can be : **latest**, **earliest**. Defaults to **latest**.

* `deliver_time_interval` - (Required, Int, ForceNew) Specifies the deliver time interval of the smart connect task.
  The value should be between 30 and 900.

* `obs_bucket_name` - (Required, String, ForceNew) Specifies the obs bucket name of the smart connect task.

* `partiton_format` - (Required, String, ForceNew) Specifies the time directory format of the smart connect task.
  The value can be : **yyyy**, **yyyy/MM**, **yyyy/MM/dd**, **yyyy/MM/dd/HH**, **yyyy/MM/dd/HH/mm**.

* `obs_path` - (Optional, String, ForceNew) Specifies the obs path of the smart connect task.
  Obs path is separated by a slash.

* `topics` - (Optional, String, ForceNew) Specifies the topic names separated by a comma of the smart connect task.
  One of topics and topics_regex is required.

* `topics_regex` - (Optional, String, ForceNew) Specifies the regular expression of topic name for the smart connect task.
  One of topics and topics_regex is required.

* `destination_file_type` - (Optional, String, ForceNew) Specifies the file type of the smart connect task. Defaults to **TEXT**.

* `record_delimiter` - (Optional, String, ForceNew) Specifies the file type of the smart connect task.
  The value can be : **,**, **;**, **|**, **\n**, **""**. Defaults to **\n**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `topics` - Indicates the topic config of the smart connect task. There are 2 kinds of values.
  One is topic name string separated by a comma, the other is a regular expression of topic name.

* `status` - Indicates the status of the smart connect task. The value can be : **RUNNING**, **PAUSED**.

* `created_at` - Indicates the creation time of the smart connect task.

## Import

The kafka smart connect task can be imported using the kafka instance `connector_id` and `id` separated by a slash, e.g.

```
$ terraform import huaweicloud_dms_kafka_smart_connect_task.test <connector_id>/<id>
```
