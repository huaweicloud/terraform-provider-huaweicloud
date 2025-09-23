---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestage_repo_password_authorization"
description: ""
---

# huaweicloud_servicestage_repo_password_authorization

This resource is used for the ServiceStage service to establish the authorization relationship through username
and password with various types of the Open-Source repository.

## Example Usage

```hcl
variable "authorization_name"
variable "domain_name"
variable "tenant_name"
variable "password"

resource "huaweicloud_servicestage_repo_password_authorization" "test" {
  type      = "devcloud"
  name      = var.authorization_name
  user_name = format("%s/%s", var.domain_name, var.tenant_name)
  password  = var.password
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specified the region in which to create the repository authorization.
  If omitted, the provider-level region will be used. Changing this parameter will create a new authorization.

* `name` - (Required, String, ForceNew) Specified the authorization name.  
  The name can contain of `4` to `63` characters, only letters, digits, underscores (_), hyphens (-) and dots (.) are
  allowed. Changing this parameter will create a new authorization.

* `type` - (Required, String, ForceNew) Specified the repository type. The valid values are as follows:
  + **devcloud**
  + **bitbucket**

  Changing this parameter will create a new authorization.

* `user_name` - (Required, String, ForceNew) Specified the user name of the repository.
  The format for each type is as follows:
  + **devcloud**: `{domain name}/{tenant name}`
  + **bitbucket**: `{account name}`

  Changing this parameter will create a new authorization.

* `password` - (Required, String, ForceNew) Specified the repository password.
  Changing this parameter will create a new authorization.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Authorizations can be imported using their `id` or `name`, e.g.:

```bash
$ terraform import huaweicloud_servicestage_repo_password_authorization.test terraform-test
```
