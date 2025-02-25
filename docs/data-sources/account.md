# huaweicloud_account

layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_account"
description: ""
Use this data source to get information about the current account.

## Example Usage

```hcl
data "huaweicloud_account" "current" {}

output "current_account_id" {
  value = data.huaweicloud_account.current.id
}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The domain ID.

* `name` - The domain name.

* `current_project_id` - The Project ID currently used.

* `username` - The username.

* `user_id` - The user ID.

-> **NOTE:** The `username` and `user_id` might be empty due to insufficient account permissions.
