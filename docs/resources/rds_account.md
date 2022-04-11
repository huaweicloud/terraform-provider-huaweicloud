---
subcategory: "Relational Database Service (RDS)"
---

# huaweicloud_rds_account

Manages RDS Mysql account resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_account" "test" {
  instance_id = var.instance_id
  name        = "test"
  password    = "Test@12345678"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds account resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the rds instance id. Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the username of the db account. Only lowercase letters, digits,
  hyphens (-), and userscores (_) are allowed. Changing this will create a new resource.
  + If the database version is MySQL 5.6, the username consists of 1 to 16 characters.
  + If the database version is MySQL 5.7 or 8.0, the username consists of 1 to 32 characters.

* `password` - (Required, String) Specifies the password of the db account. The parameter must be 8 to 32 characters
  long and contain only letters(case-sensitive), digits, and special characters(~!@#$%^*-_=+?,()&). The value must be
  different from name or name spelled backwards.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of account which is formatted `<instance_id>/<account_name>`.

## Import

RDS account can be imported using the `instance id` and `account name`, e.g.:

```
$ terraform import huaweicloud_rds_account.user_1 instance_id/account_name
```
