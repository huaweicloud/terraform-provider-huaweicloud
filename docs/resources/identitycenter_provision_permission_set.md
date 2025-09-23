---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_provision_permission_set"
description: |-
  Manages an Identity Center provision permission set resource within HuaweiCloud.
---

# huaweicloud_identitycenter_provision_permission_set

Manages an Identity Center provision permission set resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "permission_set_id" {}
variable "account_id" {}

resource "huaweicloud_identitycenter_provision_permission_set" "test" {
  instance_id       = var.instance_id
  permission_set_id = var.permission_set_id
  account_id        = var.account_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of an IAM Identity Center instance.

* `permission_set_id` - (Required, String, NonUpdatable) Specifies the ID of a permission set.

* `account_id` - (Required, String, NonUpdatable) Specifies the account ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The authorization status of a permission set.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

## Import

The Identity Center provision permission set can be imported using the `instance_id` and `id`(request ID)
separated by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_provision_permission_set.test <instance_id>/<id>
```
