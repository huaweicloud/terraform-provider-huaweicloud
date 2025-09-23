---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_permissions"
description: ""
---

# huaweicloud_dms_kafka_permissions

Use the resource to grant user permissions of a kafka topic within HuaweiCloud.

## Example Usage

```hcl
variable "kafka_instance_id" {}
variable "kafka_topic_name" {}
variable "user_1" {}
variable "user_2" {}

resource "huaweicloud_dms_kafka_permissions" "test" {
  instance_id = var.kafka_instance_id
  topic_name  = var.kafka_topic_name
  policies {
    user_name     = var.user_1
    access_policy = "all"
  }

  policies {
    user_name     = var.user_2
    access_policy = "pub"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DMS kafka permissions resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DMS kafka instance to which the permissions belongs.
  Changing this creates a new resource.

* `topic_name` - (Required, String, ForceNew) Specifies the name of the topic to which the permissions belongs.
  Changing this creates a new resource.

* `policies` - (Required, List) Specifies the permissions policies. The [object](#dms_kafka_policies) structure is
  documented below.

<a name="dms_kafka_policies"></a>
The `policies` block supports:

* `user_name` - (Required, String) Specifies the username.

* `access_policy` - (Required, String) Specifies the permissions type. The value can be:
  + **all**: publish and subscribe permissions.
  + **pub**: publish permissions.
  + **sub**: subscribe permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<topic_name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

DMS kafka permissions can be imported using the kafka instance ID and topic name separated by a slash, e.g.:

```bash
terraform import huaweicloud_dms_kafka_permissions.permissions c8057fe5-23a8-46ef-ad83-c0055b4e0c5c/topic_1
```
