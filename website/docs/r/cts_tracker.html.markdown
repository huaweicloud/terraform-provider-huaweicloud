---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: resource_huaweicloud_cts_tracker"
sidebar_current: "docs-huaweicloud-resource-cts-tracker"
description: |-
   CTS tracker allows you to collect, store, and query cloud resource operation records and use these records for security analysis, compliance auditing, resource tracking, and fault locating.
---

# huaweicloud\_cts\_tracker

Allows you to collect, store, and query cloud resource operation records.
This is an alternative to `huaweicloud_cts_tracker_v1`

## Example Usage

 ```hcl
variable "bucket_name" {}
variable "topic_id" {}

resource "huaweicloud_cts_tracker" "tracker_v1" {
  bucket_name               = "${var.bucket_name}"
  file_prefix_name          = "yO8Q"
  is_support_smn            = true
  topic_id                  = "${var.topic_id}"
  is_send_all_key_operation = false
  operations                = ["login"]
  need_notify_user_list     = ["user1"]
}

 ```
## Argument Reference
The following arguments are supported:

* `bucket_name` - (Required) The OBS bucket name for a tracker.

* `file_prefix_name` - (Optional) The prefix of a log that needs to be stored in an OBS bucket. 

* `is_support_smn` - (Required) Specifies whether SMN is supported. When the value is false, topic_id and operations can be left empty.

* `topic_id` - (Required)The theme of the SMN service, Is obtained from SMN and in the format of **urn:smn:([a-z]|[A-Z]|[0-9]|\-){1,32}:([a-z]|[A-Z]|[0-9]){32}:([a-z]|[A-Z]|[0-9]|\-|\_){1,256}**.

* `operations` - (Required) Trigger conditions for sending a notification.

* `is_send_all_key_operation` - (Required) When the value is **false**, operations cannot be left empty.

* `need_notify_user_list` - (Optional) The users using the login function. When these users log in, notifications will be sent.



## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `status` - The status of a tracker. The value is **enabled**.

* `tracker_name` - The tracker name. Currently, only tracker **system** is available.


## Import

CTS tracker can be imported using  `tracker_name`, e.g.

```
$ terraform import huaweicloud_cts_tracker_v1.tracker system
```




