---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_accounts"
description: |-
  Use this data source to get the GeminiDB Redis database account list.
---

# huaweicloud_geminidb_accounts

Use this data source to get the GeminiDB Redis database account list.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_accounts" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the scaling policies.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - All database accounts.
  The [users](#users_struct) structure is documented below.

<a name="users_struct"></a>
The `users` block supports:

* `name` - Account name.

* `type` - Account type.

* `privilege` - Account permissions.

* `databases` - All databases that the account has access permissions for.
