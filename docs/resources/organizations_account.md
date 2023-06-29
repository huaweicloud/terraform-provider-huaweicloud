---
subcategory: "Organizations"
---

# huaweicloud_organizations_account

Manages an Organizations account resource within HuaweiCloud.

-> **NOTE:** Deleting Organizations account is not support. If you destroy a resource of Organizations account,
the Organizations account is only removed from the state, but it remains in the cloud.

## Example Usage

```hcl
resource "huaweicloud_organizations_account" "test"{
  name = "account_test_name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the account.

  Changing this parameter will create a new resource.

* `parent_id` - (Optional, String) Specifies the ID of the root or organization unit in which you want to create a new
  account. The default is root ID.

* `tags` - (Optional, Map) Specifies the key/value to attach to the account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - Indicates the uniform resource name of the account.

* `joined_at` - Indicates the time when the account was created.

* `joined_method` - Indicates how an account joined an organization.

## Import

The organizations account can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_organizations_account.test <id>
```
