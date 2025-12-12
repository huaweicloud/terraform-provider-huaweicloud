---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_distinct_shared_principals"
description: |
  Use this data source to get the list of RAM distinct shared principals.
---

# huaweicloud_ram_distinct_shared_principals

Use this data source to get the list of RAM distinct shared principals.

## Example Usage

```hcl
variable resource_owner {}
variable resource_urn {}
variable principals {}

data "huaweicloud_ram_shared_principals" "test" {
  resource_owner = var.resource_owner
  resource_urn   = var.resource_urn
  principals     = var.principals
}
```

## Argument Reference

The following arguments are supported:

* `resource_owner` - (Required, String) Specifies the owner associated with the RAM share.
  Value options are as follows:
  + **self**: Shared to other users by myself.
  + **other-accounts**: Shared to me by other users.

* `principals` - (Optional, List) Specifies the principal associated with the RAM share.
  The principal could be account ID or organization ID.
  + If set to account ID, please make sure the account ID is not your owner account ID.
  + If set to organization ID, you first need to use the RAM console to enable sharing with Organization.

* `resource_urn` - (Optional, String) Specifies the resources urn associated with the
  RAM share. The format of URN is: `<service-name>:<region>:<account-id>:<type-name>:<resource-path>`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `distinct_shared_principals` - List of distinct roles.
  + `id` - Account ID of the principal or resource owner, or URN of the resource in the resource share.
  + `updated_at` - The latest update time of the RAM distinct shared principals.
