---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_user"
description: |-
  Manages a TaurusDB HTAP StarRocks database user resource within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_user

Manages a TaurusDB HTAP StarRocks database user resource within HuaweiCloud.

## Example Usage

```hcl
variable "htap_instance_id" {}
variable "user_password" {}

resource "huaweicloud_taurusdb_htap_starrocks_user" "test" {
  instance_id = var.htap_instance_id
  user_name   = "user_test"
  password    = var.user_password
  dml         = 2
  ddl         = 0
  databases   = ["*"]

  lifecycle {
    ignore_changes = [
      password
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the StarRocks database user.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NoneUpdatable) Specifies the StarRocks instance ID.

* `user_name` - (Required, String, NoneUpdatable) Specifies the database account name.
  The name must be `2` to `32` characters long, start with a lowercase letter, end with a lowercase letter or digit,
  and can contain lowercase letters, digits, and underscores.

* `password` - (Required, String) Specifies the password of the database account.
  The password must be `8` to `32` characters long, cannot be the same as the username or the reversed username,
  and must contain at least three of the following character types: uppercase letters, lowercase letters,
  digits, and special characters (~!@#%^*-_=+?).

* `databases` - (Required, List) Specifies the list of authorized databases names.

* `dml` - (Optional, Int) Specifies the DML permission type.
  The valid values are as follows:
  + **0**: read and write permissions
  + **1**: read-only permission
  + **2**: read-only and setting permissions
  + **3**: read-write and setting permissions
  Defaults to **2**.

* `ddl` - (Optional, Int) Specifies the DDL permission type.
  The valid values are as follows:
  + **0**: no DDL permission
  + **1**: DDL permission
  Defaults to **0**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<instance_id>/<user_name>`.

## Import

The StarRocks database user can be imported using the `instance_id` and `user_name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_taurusdb_htap_starrocks_user.test <instance_id>/<user_name>
```

Note that the imported state may not be identical to your resource definition, the attribute `password` missing from
the API response due to security reason. It is generally recommended running `terraform plan` after importing
a resource. You can then decide if changes should be applied to the resource, or the resource definition should be
updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_taurusdb_htap_starrocks_user" "test" {
  ...

  lifecycle {
    ignore_changes = [
      password
    ]
  }
}
```
