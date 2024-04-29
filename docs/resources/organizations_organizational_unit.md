---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_organizational_unit"
description: ""
---

# huaweicloud_organizations_organizational_unit

Manages an Organizations organizational unit resource within HuaweiCloud.

## Example Usage

```hcl
  variable "parent_id" {}
  
  resource "huaweicloud_organizations_organizational_unit" "test"{
    name      = "organizational_unit_test_name"
    parent_id = var.parent_id
  }
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the organizational unit.

* `parent_id` - (Required, String, ForceNew) Specifies the ID of the root or organizational unit in which
  you want to create a new organizational unit.

  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value to attach to the organizational unit.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - Indicates the uniform resource name of the organizational unit.

* `created_at` - Indicates the time when the organizational unit was created.

## Import

The Organizations organizational unit can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_organizations_organizational_unit.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `parent_id`. It is generally recommended
running `terraform plan` after importing the organizational unit. You can then decide if changes should be applied to
the organizational unit, or the resource definition should be updated to align with the organizational unit. Also you
can ignore changes as below.

```hcl
resource "huaweicloud_organizations_organizational_unit" "instance_1" {
    ...

  lifecycle {
    ignore_changes = [
      parent_id,
    ]
  }
}
```
