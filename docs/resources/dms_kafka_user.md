---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_user"
description: ""
---

# huaweicloud_dms_kafka_user

Manages a DMS kafka user resource within HuaweiCloud.

## Example Usage

```hcl
variable "kafka_instance_id" {}
variable "user_password" {}

resource "huaweicloud_dms_kafka_user" "user" {
  instance_id = var.kafka_instance_id
  name        = "user_1"
  password    = var.user_password
  description = "test_description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DMS kafka user resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DMS kafka instance to which the user belongs.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the user. Changing this creates a new resource.

* `password` - (Required, String) Specifies the password of the user. The parameter must be 8 to 32 characters
  long and contain only letters(case-sensitive), digits, and special characters(`~!@#$%^&*()-_=+|[{}]:'",<.>/?).
  The value must be different from name.

* `description` - (Optional, String) Specifies the description of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<user_name>`.

* `role` - Indicates the user role.

* `default_app` - Indicates whether the application is the default application.

* `created_at` - Indicates the create time.

## Import

DMS kafka users can be imported using the kafka instance ID and user name separated by a slash, e.g.

```bash
terraform import huaweicloud_dms_kafka_user.user c8057fe5-23a8-46ef-ad83-c0055b4e0c5c/user_1
```
