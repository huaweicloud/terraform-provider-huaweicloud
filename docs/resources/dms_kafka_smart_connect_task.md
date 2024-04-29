---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_smart_connect_task"
description: ""
---

# huaweicloud_dms_kafka_smart_connect_task

Manage DMS kafka smart connect task resource within HuaweiCloud.

## Example Usage

```hcl
variable "connector_id" {}
variable "access_key" {}
variable "secret_key" {}

resource "huaweicloud_dms_kafka_smart_connect_task" "test" {
  connector_id          = var.connector_id
  source_type           = "BLOB"
  task_name             = "task_test"
  destination_type      = "OBS"
  topics                = "topic-test"
  consumer_strategy     = "latest"
  destination_file_type = "TEXT"
  access_key            = var.access_key
  secret_key            = var.secret_key
  obs_bucket_name       = "bucket-test"
  obs_path              = "path-test"
  partition_format      = "yyyy/MM/dd/HH/mm"
  record_delimiter      = "\n"
  deliver_time_interval = 300
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `connector_id` - (Required, String, ForceNew) Specifies the connector ID of the kafka instance.
  Changing this parameter will create a new resource.

* `source_type` - (Required, String, ForceNew) Specifies the source type of the smart connect task.
  Only **BLOB** is supported. Changing this parameter will create a new resource.

* `task_name` - (Required, String, ForceNew) Specifies the name of the smart connect task.
  Changing this parameter will create a new resource.

* `destination_type` - (Required, String, ForceNew) Specifies the destination type of the smart connect task.
  Only **OBS** is supported. Changing this parameter will create a new resource.

* `access_key` - (Required, String, ForceNew) Specifies the access key used to access the OBS bucket.
  Changing this parameter will create a new resource.

* `secret_key` - (Required, String, ForceNew) Specifies the secret access key used to access the OBS bucket.
  Changing this parameter will create a new resource.

* `consumer_strategy` - (Required, String, ForceNew) Specifies the consumer strategy of the smart connect task.
  Changing this parameter will create a new resource.
  Value options:
  + **latest**: Read the latest data.
  + **earliest**: Read the earliest data.

* `deliver_time_interval` - (Required, Int, ForceNew) Specifies the deliver time interval of the smart connect task.
  The value should be between 30 and 900. Changing this parameter will create a new resource.

* `obs_bucket_name` - (Required, String, ForceNew) Specifies the obs bucket name of the smart connect task.
  Changing this parameter will create a new resource.

* `partition_format` - (Required, String, ForceNew) Specifies the time directory format of the smart connect task.
  Value options: **yyyy**, **yyyy/MM**, **yyyy/MM/dd**, **yyyy/MM/dd/HH**, **yyyy/MM/dd/HH/mm**.
  Changing this parameter will create a new resource.

* `obs_path` - (Optional, String, ForceNew) Specifies the obs path of the smart connect task.
  Obs path is separated by a slash. Changing this parameter will create a new resource.

* `topics` - (Optional, String, ForceNew) Specifies the topic names separated by a comma of the smart connect task.
  Changing this parameter will create a new resource.

* `topics_regex` - (Optional, String, ForceNew) Specifies the regular expression of topic name for the smart connect task.
  Changing this parameter will create a new resource.

-> **NOTE:** Exactly one of `topics`, `topics_regex` should be specified.

* `destination_file_type` - (Optional, String, ForceNew) Specifies the destination file type of the smart connect task.
  Only **TEXT** is supported. Changing this parameter will create a new resource.

* `record_delimiter` - (Optional, String, ForceNew) Specifies the record delimiter of the smart connect task.
  Value options: **,**, **;**, **|**, **\n**, **""**. Defaults to **\n**.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the status of the smart connect task. The value can be **RUNNING**, **PAUSED**.

* `created_at` - Indicates the creation time of the smart connect task.

## Timeouts

This resource provides the following timeout configuration options:

* `create` - Default is 50 minutes.

## Import

The kafka smart connect task can be imported using the kafka instance `connector_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dms_kafka_smart_connect_task.test <connector_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from
the API response. The missing attributes include: `source_type`, `access_key` and `secret_key`.
It is generally recommended running `terraform plan` after importing a kafka smart connect task.
You can then decide if changes should be applied to the kafka smart connect task, or the resource definition
should be updated to align with the kafka smart connect task. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dms_kafka_smart_connect_task" "test" {
  ...

  lifecycle {
    ignore_changes = [
      source_type, access_key, secret_key,
    ]
  }
}
```
