---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_user_token"
description: ""
---

# huaweicloud_identity_user_token

Manages an IAM user token resource within HuaweiCloud.

->**Note** The token can not be destroyed. It will be invalid after expiration time. If password or AK/SK is changed,
the token valid time will last less than 30 minutes.

## Example Usage

```hcl
variable "account_name" {}
variable "user_name" {}
variable "password" {}

resource "huaweicloud_identity_user_token" "test" {
  account_name = var.account_name
  user_name    = var.user_name
  password     = var.password
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, String, ForceNew) Specifies the account name to which the IAM user belongs.
  Changing this will create a new token.

* `user_name` - (Required, String, ForceNew) Specifies the IAM user name. Changing this will create a new token.

* `password` - (Required, String, ForceNew) Specifies the IAM user password. Changing this will create a new token.

* `project_name` - (Optional, String, ForceNew) Specifies the project name. If it is blank, the token applies to global
  services, otherwise the token applies to project-level services. Changing this will create a new token.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID in format `<account_name>/<user_name>`.

* `token` - The token. Validity period is 24 hours.



