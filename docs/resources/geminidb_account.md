---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_account"
description: |-
  Manages a GeminiDB account resource within HuaweiCloud.
---

# huaweicloud_geminidb_account

Manages a GeminiDB account resource within HuaweiCloud.

## Example Usage

```hcl
var "instance_id" {}

data "huaweicloud_geminidb_account" "test" {
  instance_id = var.instance_id
  name        = "test_account"
  password    = "Test@1234567"
  privilege   = "ReadWrite"
  databases   = [ "1", "2" ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Geminidb instance.

* `name` - (Required, String, NonUpdatable) Specifies the username of the Geminidb account.

* `password` - (Required, String) Specifies the password of the Geminidb account.

* `privilege` - (Required, String) Specifies the privilege of the Geminidb account.

* `databases` - (Required, List) Specifies the list of databases.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `type` - The type of the Geminidb account.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The GeminiDB account can be imported using the `instance_id` and `name`, separated by a slash, e.g.

```
$ terraform import huaweicloud_geminidb_account.test <instance_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `password` .
It is generally recommended running `terraform plan` after importing an account.
You can then decide if changes should be applied to the account, or the resource definition should be updated to
align with the account. Also you can ignore changes as below.

```hcl
resource "huaweicloud_geminidb_account" "test" {
  ...

  lifecycle {
    ignore_changes = [
      password, 
    ]
  }
}
```
