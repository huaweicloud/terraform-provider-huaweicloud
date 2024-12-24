---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestage_repo_token_authorization"
description: ""
---

# huaweicloud_servicestage_repo_token_authorization

This resource is used for the ServiceStage service to establish the authorization relationship through personal access
token with various types of the Open-Source repository.

## Example Usage

```hcl
variable "authorization_name" {}
variable "personal_access_token" {}

resource "huaweicloud_servicestage_repo_token_authorization" "test" {
  type  = "github"
  name  = var.authorization_name
  token = var.personal_access_token
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
  + **github**
  + **gitlab**
  + **gitee**

  Changing this parameter will create a new authorization.

* `token` - (Required, String, ForceNew) Specified the personal access token of the repository.
  Changing this parameter will create a new authorization.

<!-- markdownlint-disable MD034 -->
* `host` - (Optional, String, ForceNew) Specified the host name of the repository, e.g. **https://api.github.com**.
  Changing this parameter will create a new authorization.
<!-- markdownlint-enable MD034 -->

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Authorizations can be imported using their `id` or `name`, e.g.:

```bash
$ terraform import huaweicloud_servicestage_repo_token_authorization.test terraform-test
```
