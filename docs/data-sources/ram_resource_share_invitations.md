---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_share_invitations"
description: |-
  Use this data source to get the list of RAM resource share invitations.
---

# huaweicloud_ram_resource_share_invitations

Use this data source to get the list of resource share invitations.

## Example Usage

```hcl
var "status" {}

data "huaweicloud_ram_resource_share_invitations" "test" {
  status = var.status
}
```

## Argument Reference

The following arguments are supported:

* `resource_share_ids` - (Optional, List) Specifies the list of the resource share IDs.

* `resource_share_invitation_ids` - (Optional, List) Specifies the list of the resource share invitation IDs.

* `status` - (Optional, String) Specifies the status of the resource share invitation.
  The valid values are **pending**, **accepted** and **rejected**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resource_share_invitations` - The list of shared resources.
  The [resource_share_invitations](#ram_resource_share_invitations) structure is documented below.

<a name="ram_resource_share_invitations"></a>
The `resource_share_invitations` block supports:

* `id` - The ID of the resource share invitation.

* `resource_share_id` - The ID of the resource share.

* `resource_share_name` - The name of the resource share.

* `receiver_account_id` - The ID of the account that receives the resource share invitation.

* `sender_account_id` - The ID of the account that sends the resource share invitation.

* `status` - The status of the resource share invitation.

* `created_at` - The creation time of the resource share invitation.

* `updated_at` - The latest update time of the resource share invitation.
