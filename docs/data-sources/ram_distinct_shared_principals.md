---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_distinct_shared_principals"
description: |-
  Use this data source to get the list of RAM distinct shared principals.
---

# huaweicloud_ram_distinct_shared_principals

Use this data source to get the list of RAM distinct shared principals.

## Example Usage

```hcl
variable resource_owner {}

data "huaweicloud_ram_distinct_shared_principals" "test" {
  resource_owner = var.resource_owner
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
  + If set to organization ID, you first need to use the RAM console to enable sharing with organization.

* `resource_urn` - (Optional, String) Specifies the resources urn associated with the
  RAM share. The format of URN is: `<service-name>:<region>:<account-id>:<type-name>:<resource-path>`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `distinct_shared_principals` - The list of information for different roles.

  The [distinct_shared_principals](#distinct_shared_principals_struct) structure is documented below.

<a name="distinct_shared_principals_struct"></a>
The `distinct_shared_principals` block supports:

* `id` - The account ID or URN of the creator or user of the resource sharing instance.

* `updated_at` - The last time the resource sharing instance was updated.
