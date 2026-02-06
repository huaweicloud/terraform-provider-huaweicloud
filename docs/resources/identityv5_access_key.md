---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_access_key"
description: |-
  Manages a permanent access key resource within HuaweiCloud.
---

# huaweicloud_identityv5_access_key

Manages a permanent access key resource within HuaweiCloud.

## Example Usage

```hcl
variable "user_id" {}

resource "huaweicloud_identityv5_access_key" "test" {
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String, NonUpdatable) Specifies the ID of the user.

* `status` - (Optional, String) Specifies the status of the access key.  
  Defaults to **active**.  
  The valid values are as follows:
  + **active**
  + **inactive**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The access key ID.

* `access_key_id` - The ID of the generated access key.

* `secret_access_key` - The generated secret access key.  

* `created_at` - The creation time of the access key.

* `last_used_at` - The time when the access key was last used.  
  If it is an empty string, it means that the AK and SK have never been used.

## Import

The IAM v5 access key can be imported using the `user_id` and `id` ( also the `access_key_id`),
separated by a slash, e.g.

```bash
$ terraform import huaweicloud_identityv5_access_key.test <user_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `secret_access_key`.
It is generally recommended running `terraform plan` after importing the resource. You can then decide
if changes should be applied to the resource, or the resource definition should be updated to align with the resource.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_identityv5_access_key" "test" {
  ...

  lifecycle {
    ignore_changes = [
      secret_access_key,
    ]
  }
}
```
