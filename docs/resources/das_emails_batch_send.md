---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_emails_batch_send"
description: |-
  Use this resource to batch send emails within HuaweiCloud.
---

# huaweicloud_das_emails_batch_send

Use this resource to batch send emails within HuaweiCloud.

-> This resource is a one-time action resource for batch sending emails. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Batch send emails with email address

```hcl
variable "task_ids" {
  type = list(string)
}

resource "huaweicloud_das_emails_batch_send" "test" {
  task_ids = var.task_ids
  email    = "test@example.com"
}
```

### Batch send emails with SMN topic

```hcl
variable "task_ids" {
  type = list(string)
}
var "topic_id" {}
var "topic_urn" {}

resource "huaweicloud_das_emails_batch_send" "test" {
  task_ids  = var.task_ids
  topic     = var.topic_id
  topic_urn = var.topic_urn
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the email templates are located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `task_ids` - (Required, List, NonUpdatable) Specifies the list of report IDs.

* `email` - (Optional, String, NonUpdatable) Specifies the email address.  

* `topic` - (Optional, String, NonUpdatable) Specifies the topic ID.  

* `topic_urn` - (Optional, String, NonUpdatable) Specifies the topic URN.  

-> 1.Use either `email` or `topic` to receive notifications.  
   2.The `topic` and `topic_urn` must be used at the same time.

* `obs_bucket_name` - (Optional, String, NonUpdatable) Specifies the OBS bucket name.  
  The diagnostic report will be uploaded to the target OBS bucket, and the email will provide a download link for the
  diagnostic report.

* `service_uri` - (Optional, String, NonUpdatable) Specifies the service URI of the OBS bucket.

* `access_key` - (Optional, String, NonUpdatable) Specifies the access key used to access the OBS bucket.

* `secret_key` - (Optional, String, NonUpdatable) Specifies the secret key used to access the OBS bucket.

-> The `access_key` and `secret_key` is **Required** when `obs_bucket_name` is not empty.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
