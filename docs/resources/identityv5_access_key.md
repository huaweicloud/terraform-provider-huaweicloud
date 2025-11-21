---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_access_key"
description: |-
  Manages a permanent Access Key resource within HuaweiCloud IAM V5 service.
---

# huaweicloud_identityv5_access_key

Manages a permanent Access Key resource within HuaweiCloud IAM V5 service.

## Example Usage

```hcl
variable "name" {}

resource "huaweicloud_identityv5_user" "user_1" {
  name        = var.name
  description = "tested by terraform"
}

resource "huaweicloud_identityv5_access_key" "key_1" {
  user_id = huaweicloud_identityv5_user.user_1.id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String， NonUpdatable) Specifies the ID of the user.

* `status` - (Optional, String) Specifies the status of the access key can be *active* or *inactive*.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The access key ID.

* `created_at` - Indicates access key creation time.

* `last_used_at` - Indicates the last usage time of the access key. If it does not exist,
  it indicates that it has never been used.

* `access_key_id` - Indicates the access key ID.

* `secret_access_key` - Indicates the access secret key.  

## Import

The IAM v5 access key can be imported using the user_id and access_key_id separated by a slash, e.g:

```bash
$ terraform import huaweicloud_identityv5_access_key.access_key <user_id>/<id>
```
