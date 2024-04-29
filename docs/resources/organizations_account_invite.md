---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_account_invite"
description: ""
---

# huaweicloud_organizations_account_invite

Manages an Organizations account invite resource within HuaweiCloud.

## Example Usage

```hcl
variable account_id {}

resource "huaweicloud_organizations_account_invite" "test"{
  account_id = var.account_id
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required, String, ForceNew) Specifies the ID of the target account.

  Changing this parameter will create a new resource.

* `remove_account_on_destroy` - (Optional, Bool) Specifies whether to remove the invited account when delete the
  invitation (handshake). Defaults to false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - Indicates the uniform resource name of the invitation

* `master_account_id` - Indicates the unique ID of the organization's management account.

* `master_account_name` - Indicates the name of the organization's management account.

* `organization_id` - Indicates the ID of the organization.

* `created_at` - Indicates the date and time when an invitation (handshake) request was made.

* `updated_at` - Indicates the date and time when an invitation (handshake) request was accepted, canceled,
  declined, or expired.

* `status` - Indicates the current state of the invitation (handshake).

## Import

The Organizations account invite can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_organizations_account_invite.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `remove_account_on_destroy`. It is
generally recommended running `terraform plan` after importing an account invite. You can then decide if changes should
be applied to the account invite, or the resource definition should be updated to align with the account invite.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_organizations_account_invite" "test" {
  ...

  lifecycle {
    ignore_changes = [
      remove_account_on_destroy,
    ]
  }
}
```
