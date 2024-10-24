---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_shared_principals"
description: |
  Use this data source to get the list of RAM shared principals.
---

# huaweicloud_ram_shared_principals

Use this data source to get the list of RAM shared principals.

## Example Usage

```hcl
variable "resource_urn" {}

data "huaweicloud_ram_shared_principals" "test" {
  resource_owner = "self"
  resource_urn   = var.resource_urn
}
```

## Argument Reference

The following arguments are supported:

* `resource_owner` - (Required, String) Specifies the owner associated with the RAM share.
  Value options are as follows:
  + **self**: Shared to other users by myself.
  + **other-accounts**: Shared to me by other users.

* `principal` - (Optional, String) Specifies the principal associated with the RAM share.
  The principal could be account ID or organization ID.
  + If set to account ID, please make sure the account ID is not your owner account ID.
  + If set to organization ID, you first need to use the RAM console to enable sharing with Organization. Please refer
  to the [document](https://support.huaweicloud.com/intl/en-us/qs-ram/ram_02_0004.html).

* `resource_urn` - (Optional, String) Specifies the resources urn associated with the
  RAM share. The format of URN is: `<service-name>:<region>:<account-id>:<type-name>:<resource-path>`.

* `resource_share_id` - (Optional, String) Specifies the ID of resource share.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `shared_principals` - The list of shared principals.
  The [shared_principals](#attrblock-shared_principals) structure is documented below.

<a name="attrblock-shared_principals"></a>
The `shared_principals` block supports:

* `id` - The ID of shared principal.

* `resource_share_id` - The resource share ID.

* `created_at` - The creation time of the RAM share.

* `updated_at` - The latest update time of the RAM share.
