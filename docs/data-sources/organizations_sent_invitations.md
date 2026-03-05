---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_sent_invitations"
description: |-
  Use this data source to get the list of sent invitations within HuaweiCloud.
---

# huaweicloud_organizations_sent_invitations

Use this data source to get the list of sent invitations within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_organizations_sent_invitations" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the sent invitations are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `handshakes` - Indicates the list of invitations (handshakes).

  The [handshakes](#handshakes_struct) structure is documented below.

<a name="handshakes_struct"></a>
The `handshakes` block supports:

* `id` - Indicates the unique ID of an invitation (handshake).

* `urn` - Indicates the uniform resource name of the invitation (handshake).

* `status` - Indicates the current state of the invitation (handshake).
  + **pending**
  + **accepted**
  + **cancelled**
  + **declined**
  + **expired**

* `organization_id` - Indicates the unique ID of an organization.

* `management_account_id` - Indicates the unique ID of the organization's management account.

* `management_account_name` - Indicates the name of the organization's management account.

* `target` - Indicates the target information to be invited.

  The [target](#handshakes_target_struct) structure is documented below.

* `created_at` - Indicates the date and time when an invitation (handshake) request was made.

* `updated_at` - Indicates the date and time when an invitation (handshake) request was updated.

* `notes` - Indicates the additional information included in the email to the recipient account owner.

<a name="handshakes_target_struct"></a>
The `target` block supports:

* `type` - Indicates the type of the account.

* `entity` - Indicates the ID of the account.
