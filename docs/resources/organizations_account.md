---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_account"
description: |-
  Manages an Organizations account resource within HuaweiCloud.
---

# huaweicloud_organizations_account

Manages an Organizations account resource within HuaweiCloud.

-> **NOTE:** Deleting Organizations account is not support. If you destroy a resource of Organizations account,
the Organizations account is only removed from the state, but it remains in the cloud.

## Example Usage

### Account in International Website

```hcl
resource "huaweicloud_organizations_account" "test"{
  name  = "account_test_name"
  email = "account_test@demo.com"
}
```

### Account in Chinese Mainland Website

```hcl
resource "huaweicloud_organizations_account" "test"{
  name  = "account_test_name"
  phone = "138********"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the account.
  Changing this parameter will create a new resource.

* `email` - (Optional, String, ForceNew) Specifies the email address of the account.
  This parameter is mandatory in **International Website**.
  Changing this parameter will create a new resource.

* `phone` - (Optional, String, ForceNew) Specifies the mobile number of the account.
  This parameter is mandatory in **Chinese Mainland Website**.
  Changing this parameter will create a new resource.

-> At least one of `email` and `phone` must be specified.

* `agency_name` - (Optional, String, ForceNew) Specifies the agency name of the account.
  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the account.

* `parent_id` - (Optional, String) Specifies the ID of the root or organization unit in which you want to create a new
  account. The default is root ID.

* `tags` - (Optional, Map) Specifies the key/value to attach to the account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `intl_number_prefix` - Indicates the prefix of a mobile number.

* `urn` - Indicates the uniform resource name of the account.

* `joined_at` - Indicates the time when the account was created.

* `joined_method` - Indicates how an account joined an organization.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.

## Import

The Organizations account can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_organizations_account.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include `agency_name`. It is generally recommended
running `terraform plan` after importing an account. You can then decide if changes should be applied to the account,
or the resource definition should be updated to align with the account. Also you can ignore changes as below.

```hcl
resource "huaweicloud_organizations_account" "test" {
  ...

  lifecycle {
    ignore_changes = [
      agency_name,
    ]
  }
}
```
