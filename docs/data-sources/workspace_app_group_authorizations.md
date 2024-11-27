---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_group_authorizations"
description: |-
  Use this data source to get the list of the application group authorizations within HuaweiCloud.
---

# huaweicloud_workspace_app_group_authorizations

Use this data source to get the list of the application group authorizations within HuaweiCloud.

## Example Usage

### Query all application group authorizations

```hcl
data "huaweicloud_workspace_app_group_authorizations" "test" {}
```

### Query the application group authorizations that contains the same name segment and the account type is USER

```hcl
variable "account_name_prefix" {}

data "huaweicloud_workspace_app_group_authorizations" "test" {
  account      = var.account_name_prefix
  account_type = "USER"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the authorizations.
  If omitted, the provider-level region will be used.

* `app_group_id` - (Optional, String) Specifies the authorized application group ID.

* `account` - (Optional, String) Specifies the name of the authorized account. Fuzzy search is supported.

* `account_type` - (Optional, String) Specifies the type of the authorized account.  
  The valid values are as follows:
  + **USER**
  + **USER_GROUP**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `authorizations` - All authorizations that match the filter parameters.

  The [authorizations](#app_group_authorizations) structure is documented below.

<a name="app_group_authorizations"></a>
The `authorizations` block supports:

* `id` - The authorized ID.

* `account_id` - The ID of the authorized account.

* `account` - The name of the authorized account.

* `account_type` - The type of the authorized account.

* `app_group_id` - The application group ID corresponding to the authorized account.

* `app_group_name` - The application group name corresponding to the authorized account.

* `created_at` - The time when the account is authorized to the specified application group, in RFC3339 format.
