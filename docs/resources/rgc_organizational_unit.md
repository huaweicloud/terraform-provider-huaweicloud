---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_organizational_unit"
description: |-
  Manages an RGC organizational unit resource within HuaweiCloud.
---

# huaweicloud_rgc_account_enroll

Manages an RGC organizational unit resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "organizational_unit_name" {}
variable "parent_organizational_unit_id" {}

resource "huaweicloud_rgc_account_enroll" "test"{
  organizational_unit_name      = var.organizational_unit_name
  parent_organizational_unit_id = var.parent_organizational_unit_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `organizational_unit_name` - (Required, String, NonUpdatable) Specifies the ID of the account.

* `parent_organizational_unit_id` - (Required, String, NonUpdatable) Specifies parent organizational
  unit ID of enrolled account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `parent_organizational_unit_name` - Indicates the name of parent organizational unit.

* `manage_account_id` - Indicates the ID of organization manage account.

* `organizational_unit_id` - Indicates the ID of created organizational unit.

* `organizational_unit_status` - Indicates the state of created organizational unit.

* `organizational_unit_type` - Indicates the type of created organizational unit.

## Import

The RGC organizational unit can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rgc_organizational_unit.test <id>
```
