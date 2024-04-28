---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_delegated_administrator"
description: ""
---

# huaweicloud_organizations_delegated_administrator

Manages an Organizations delegated administrator resource within HuaweiCloud.

## Example Usage

```hcl
variable "account_id" {}
variable "service_principal" {}

resource "huaweicloud_organizations_delegated_administrator" "test"{
  account_id        = var.account_id
  service_principal = var.service_principal
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required, String, ForceNew) Specifies the unique ID of an account.

  Changing this parameter will create a new resource.

* `service_principal` - (Required, String, ForceNew) Specifies the name of the service principal.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The Organizations delegated administrator can be imported using the `account_id` and `service_principal` separated by
a slash, e.g.

```bash
$ terraform import huaweicloud_organizations_delegated_administrator.test <account_id>/<service_principal>
```
