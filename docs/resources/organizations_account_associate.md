---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_account_associate"
description: ""
---

# huaweicloud_organizations_account_associate

Manages an Organizations account associate resource within HuaweiCloud.

## Example Usage

```hcl
variable account_id {}
variable parent_id {}

resource "huaweicloud_organizations_account_associate" "test"{
  account_id = var.account_id
  parent_id  = var.parent_id
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required, String, ForceNew) Specifies the ID of the account.

  Changing this parameter will create a new resource.

* `parent_id` - (Required, String) Specifies the ID of root or organizational unit in which you want to move the account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `name` - Indicates the name of the account.

* `urn` - Indicates the uniform resource name of the account.

* `joined_at` - Indicates the time when the account was created.

* `joined_method` - Indicates how an account joined an organization.

## Import

The Organizations account associate can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_organizations_account_associate.test <id>
```
