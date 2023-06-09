---
subcategory: "Resource Access Manager (RAM)"
---

# huaweicloud_ram_resource_share

Manages a RAM resource share resource within HuaweiCloud.

## Example Usage

```hcl
variable "account_id" {}
variable "resource_urn" {}

data "huaweicloud_ram_resource_permissions" "test" {
  resource_type = "vpc:subnets"
}

resource "huaweicloud_ram_resource_share" "test" {
  name        = "demo-share"
  description = "test description information"

  resource_urns  = [var.resource_urn]
  principals     = [var.account_id]
  permission_ids = [huaweicloud_ram_resource_permissions.test.permissions[0].id]

  tags = {
    foo = "bar"
    key = "value"
  }
} 
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the resource share.

* `principals` - (Required, List, ForceNew) Specifies the list of one or more principals associated with the resource
  share. The principals could be account IDs and organization IDs. You can put account IDs and organization IDs to this
  field together.
  + If set to account IDs, please make sure the account ID is not your owner account ID.
  + If set to organization IDs, you first need to use the RAM console to enable sharing with Organizations. Please refer
  to the [document](https://support.huaweicloud.com/intl/en-us/qs-ram/ram_02_0004.html).

  Changing this parameter will create a new resource.

* `resource_urns` - (Required, List, ForceNew) Specifies the list of URNs of one or more resources associated with the
  resource share. The format of URN is: `<service-name>:<region>:<account-id>:<type-name>:<resource-path>`.
  Sharable cloud services and resource types refer to
  [document](https://support.huaweicloud.com/intl/en-us/productdesc-ram/ram_01_0007.html).

  Changing this parameter will create a new resource.

* `permission_ids` - (Optional, List, ForceNew) Specifies the list of RAM permissions associated with the resource
  share. A resource type can be associated with only one RAM permission. If you do not specify a permission ID,
  RAM automatically associates the default permission for each resource type.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the resource share.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the resource share.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `owning_account_id` - The owning account ID of the RAM share.

* `status` - The status of the RAM share.

* `created_at` - The creation time of the RAM share.

* `updated_at` - The latest update time of the RAM share.

* `associated_permissions` - The associated permissions of the RAM share.
  The [associated_permissions](#RAMShare_associated_permissions) structure is documented below.

<a name="RAMShare_associated_permissions"></a>
The `associated_permissions` block supports:

* `permission_id` - The permission ID.

* `permission_name` - The permission name.

* `resource_type` - The resource type of the permission.

* `status` - The status of the permission.

## Import

The ram share can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ram_resource_share.test <id>
```
