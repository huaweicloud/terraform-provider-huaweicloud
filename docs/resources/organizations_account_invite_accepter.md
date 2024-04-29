---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_account_invite_accepter"
description: ""
---

# huaweicloud_organizations_account_invite_accepter

Manages an Organizations account invite accepter resource within HuaweiCloud.

## Example Usage

```hcl
variable invitation_id {}

resource "huaweicloud_organizations_account_invite_accepter" "test"{
  invitation_id = var.invitation_id
}
```

## Argument Reference

The following arguments are supported:

* `invitation_id` - (Required, String, ForceNew) Specifies the unique ID of an invitation (handshake).

  Changing this parameter will create a new resource.

* `leave_organization_on_destroy` - (Optional, Bool) Specifies whether to leave the organization when delete the
  invitation (handshake). Defaults to false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - Indicates the uniform resource name of the invitation

* `account_id` - Indicates the ID of the target account.

* `master_account_id` - Indicates the unique ID of the organization's management account.

* `master_account_name` - Indicates the name of the organization's management account.

* `organization_id` - Indicates the ID of the organization.

* `created_at` - Indicates the date and time when an invitation (handshake) request was made.

* `updated_at` - Indicates the date and time when an invitation (handshake) request was accepted, cancelled,
  declined, or expired.

* `status` - Indicates the current state of the invitation (handshake).

## Import

The Organizations account invite accepter can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_organizations_account_invite_accepter.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `leave_organization_on_destroy`. It is
generally recommended running `terraform plan` after importing an account invite accepter. You can then decide if
changes should be applied to the account invite accepter, or the resource definition should be updated to align with
the account invite accepter. Also you can ignore changes as below.

```hcl
resource "huaweicloud_organizations_account_invite_accepter" "test" {
  ...

  lifecycle {
    ignore_changes = [
      leave_organization_on_destroy,
    ]
  }
}
```
