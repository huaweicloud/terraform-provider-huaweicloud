---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_email_template"
description: |-
  Manages a DAS email template resource within HuaweiCloud.
---

# huaweicloud_das_email_template

Manages a DAS email template resource within HuaweiCloud.

## Example Usage

### Create an email template via email

```hcl
variable "name" {}
variable "groups" {
  type = list(string)
}
variable "email" {}

resource "huaweicloud_das_email_template" "test" {
  datastore_type  = "MySQL"
  name            = var.name
  groups          = var.groups
  health_rank     = ["HIGH", "MEDIUM"]
  inspection_time = "00:00-00:00"
  send_time       = "08:00-10:00"
  time_zone       = "Asia/Shanghai"
  email           = var.email
}
```

### Create an email template via SMN topic

```hcl
variable "name" {}
variable "groups" {
  type = list(string)
}
variable "topic_id" {}
variable "topic_urn" {}

resource "huaweicloud_das_email_template" "test" {
  datastore_type  = "MySQL"
  name            = var.name
  groups          = var.groups
  health_rank     = ["HIGH", "MEDIUM"]
  inspection_time = "00:00-00:00"
  send_time       = "08:00-10:00"
  time_zone       = "Asia/Shanghai"
  topic           = var.topic
  topic_urn       = var.topic_urn
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the email template is located.  
  If omitted, the provider-level region will be used.  
  Changing this creates a new resource.

* `datastore_type` - (Required, String, NonUpdatable) Specifies the database type.  
  The valid values are as follows:
  + **MySQL**
  + **TaurusDB**
  + **GaussDB**
  + **MariaDB**.

* `name` - (Required, String) Specifies the name of the email template.

* `groups` - (Required, List) Specifies the list of instance group IDs.

* `health_rank` - (Required, List) Specifies the list of health ranks.  
  The valid values are as follows:
  + **dangerous**
  + **sub_healthy**
  + **healthy**
  + **high_risk**

* `inspection_time` - (Required, String) Specifies the diagnosis time.  
  The valid values are as follows:
  + **00:00-00:00**
  + **12:00-12:00**.

* `send_time` - (Required, String) Specifies the send time.  
  The valid values are as follows:
  + **00:00-02:00**
  + **02:00-04:00**
  + **04:00-06:00**
  + **06:00-08:00**
  + **08:00-10:00**
  + **10:00-12:00**
  + **12:00-14:00**
  + **14:00-16:00**
  + **16:00-18:00**
  + **18:00-20:00**
  + **20:00-22:00**
  + **22:00-24:00**

* `time_zone` - (Required, String) Specifies the time zone.

* `email` - (Optional, String) Specifies the email address.  
  Use either `email` or `topic`/`topic_urn` to receive notifications.

* `topic` - (Optional, String) Specifies the topic ID.  
  Use either `email` or `topic`/`topic_urn` to receive notifications.

* `topic_urn` - (Optional, String) Specifies the topic URN.  
  Use either `email` or `topic`/`topic_urn` to receive notifications.

  -> The `topic` and `topic_urn` must be used at the same time.

* `obs_bucket_name` - (Optional, String) Specifies the OBS bucket name.  
  The diagnosis report will be uploaded to the target OBS bucket.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the template ID.

* `updated_at` - The update time, in RFC3339 format.

* `user_id` - The ID of the user who last modified the template.

* `status` - The status of the email template.

## Import

The email template can be imported using `template_id`, e.g.

```bash
$ terraform import huaweicloud_das_email_template.test <id>
```
