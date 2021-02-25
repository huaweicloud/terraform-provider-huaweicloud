---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud\_identity\_user

Use this data source to get an available HuaweiCloud IAM user.

```hcl
variable "identity_user_name" {}

data "huaweicloud_identity_user" "user" {
  name = var.identity_user_name
}
```

## Argument Reference

* `name` - (Required, String) Specifies the name of the IAM user.

* `enabled` - (Optional, Bool) Specifies whether the IAM user is enabled or disabled.
    Valid values are `true` and `false`, defaults to `true`.

* `domain_id` - (Optional, String) Specifies the domain which IAM user belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A data source ID in UUID format.

* `description` - A description of the IAM user.