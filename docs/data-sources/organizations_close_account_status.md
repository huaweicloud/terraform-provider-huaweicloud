---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_close_account_status"
description: |-
  Use this data source to get the list of the account close requests for an organization.
---

# huaweicloud_organizations_close_account_status

Use this data source to get the list of the account close requests for an organization.

## Example Usage

```hcl
data "huaweicloud_organizations_close_account_status" "test" {}
```

## Argument Reference

The following arguments are supported:

* `states` - (Optional, List) Specifies the list of one or more states that you want to include in the response.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `close_account_statuses` - Indicates the list of close account statuses.

  The [close_account_statuses](#close_account_statuses_struct) structure is documented below.

<a name="close_account_statuses_struct"></a>
The `close_account_statuses` block supports:

* `account_id` - Indicates the ID of an account.

* `state` - Indicates the Status of the close account request.

* `organization_id` - Indicates the ID of an organization.

* `failure_reason` - Indicates the reason for a request failure.

* `created_at` - Indicates the date and time when the close account request was made.

* `updated_at` - Indicates the date and time when the status of close account request was updated.
