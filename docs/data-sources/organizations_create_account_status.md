---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_create_account_status"
description: |-
  Use this data source to get the list of the account creation requests for an organization.
---

# huaweicloud_organizations_create_account_status

Use this data source to get the list of the account creation requests for an organization.

## Example Usage

```hcl
data "huaweicloud_organizations_create_account_status" "test" {}
```

## Argument Reference

The following arguments are supported:

* `states` - (Optional, List) Specifies the list of states.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `create_account_statuses` - Indicates the list of create account statuses.

  The [create_account_statuses](#create_account_statuses_struct) structure is documented below.

<a name="create_account_statuses_struct"></a>
The `create_account_statuses` block supports:

* `id` - Indicates the ID of a request.

* `state` - Indicates the status of the asynchronous request for creating an account.

* `account_id` - Indicates the ID of the newly created account if any.

* `account_name` - Indicates the account name.

* `failure_reason` - Indicates the reason for a request failure.

* `created_at` - Indicates the date and time when the create account request was made.

* `completed_at` - Indicates the date and time when the account was created and the request was completed.
