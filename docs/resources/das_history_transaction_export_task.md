---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_history_transaction_export_task"
description: |-
  Manages a DAS history transaction export task resource within HuaweiCloud.
---

# huaweicloud_das_history_transaction_export_task

Manages a DAS history transaction export task resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "bucket_name" {}

resource "huaweicloud_das_history_transaction_export_task" "test" {
  instance_id = var.instance_id
  bucket_name = var.bucket_name
  start_time  = "2000-06-01T00:00:00+08:00"
  end_time    = "2099-06-02T00:00:00+08:00"
  time_zone   = "GMT+8"
}
```

### With Optional Parameters

```hcl
variable "instance_id" {}
variable "bucket_name" {}

resource "huaweicloud_das_history_transaction_export_task" "test" {
  instance_id  = var.instance_id
  bucket_name  = var.bucket_name
  start_time   = "2000-06-01T00:00:00+08:00"
  end_time     = "2099-06-02T00:00:00+08:00"
  time_zone    = "GMT+8"
  file_path    = "export/history"
  order_field  = "collectTime"
  order_by     = "asc"
  last_sec_min = 0
  last_sec_max = 3600
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the history transaction export task is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID.

* `bucket_name` - (Required, String, NonUpdatable) Specifies the OBS bucket name.  
  The maximum length is `63` characters.

* `start_time` - (Required, String, NonUpdatable) Specifies the start time, in RFC3339 format.

* `end_time` - (Required, String, NonUpdatable) Specifies the end time, in RFC3339 format.

* `file_path` - (Optional, String, NonUpdatable) Specifies the OBS file directory.  
  The maximum length is `1024` characters.

* `time_zone` - (Optional, String, NonUpdatable) Specifies the time zone.  
  The maximum length is `64` characters.  
  The valid values are any legal timezone definition, such as **UTC**, **GMT**.

* `order_field` - (Optional, String, NonUpdatable) Specifies the sort field.  
  The valid values are as follows:
  + **collectTime**: Collection time.
  + **occurrenceTime**: Occurrence time.
  + **lastSec**: Transaction duration.
  + **waitLockStructCount**: Number of held locks.
  + **holdLockStructCount**: Number of waiting locks.

* `order_by` - (Optional, String, NonUpdatable) Specifies the sort order.  
  The valid values are as follows:
  + **asc**: Ascending.
  + **desc**: Descending.

* `last_sec_min` - (Optional, Int, NonUpdatable) Specifies the minimum duration.

* `last_sec_max` - (Optional, Int, NonUpdatable) Specifies the maximum duration.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The export task ID.

* `status` - The task status.
  + **-1**: Quota full.
  + **0**: Waiting.
  + **1**: Running.
  + **2**: Failed.
  + **3**: Successful.
  + **4**: Time out.
  + **5**: OBS file deleted.

* `created_time` - The task creation time, in RFC3339 format.

* `export_line_num` - The number of exported lines.

* `download_url` - The download URL of the exported file.

## Import

The DAS history transaction export task can be imported using `<instance_id>/<id>`, e.g.

```bash
$ terraform import huaweicloud_das_history_transaction_export_task.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `start_time`, `end_time`, `bucket_name`, `file_path`, `time_zone`, `order_field`,
`order_by`, `last_sec_min`, `last_sec_max`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_das_history_transaction_export_task" "test" {
  ...

  lifecycle {
    ignore_changes = [
      start_time, end_time, bucket_name, file_path, time_zone, order_field, order_by, last_sec_min, last_sec_max,
    ]
  }
}
```
