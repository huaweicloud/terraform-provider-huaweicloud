---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_reviewer"
description: ""
---

# huaweicloud_dataarts_architecture_reviewer

Manages DataArts Studio architecture reviewer resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "user_name" {}
variable "user_id" {}
variable "email" {}
variable "phone_number" {}

resource "huaweicloud_dataarts_architecture_reviewer" "test" {
  workspace_id = var.workspace_id
  user_name    = var.user_name
  user_id      = var.user_id
  email        = var.email
  phone_number = var.phone_number
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the architecture reviewer resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID to which the architecture reviewer belongs.
  Changing this parameter will create a new resource.

* `user_name` - (Required, String, ForceNew) Specifies the user name of the architecture reviewer.
  Changing this parameter will create a new resource.

* `user_id` - (Required, String, ForceNew) Specifies the user ID of the architecture reviewer.
  Changing this parameter will create a new resource.

* `email` - (Optional, String, ForceNew) Specifies the email of the architecture reviewer.
  Changing this parameter will create a new resource.

* `phone_number` - (Optional, String, ForceNew) Specifies the phone number of the architecture reviewer.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, the value is the `user_name`.

* `reviewer_id` - The ID of the reviewer.

## Import

The DataArts architecture reviewer can be imported using the `workspace_id` and `user_name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_architecture_reviewer.test <workspace_id>/<user_name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `email` and `phone_number`.
It is generally recommended running `terraform plan` after importing a reviewer.
You can then decide if changes should be applied to the reviewer, or the resource definition should be updated to
align with the reviewer. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dataarts_architecture_reviewer" "test"{
    ...

  lifecycle {
    ignore_changes = [
      email, phone_number,
    ]
  }
}
```
