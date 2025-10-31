---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_organizational_unit_register"
description: |-
  Manages an RGC organizational unit register resource within HuaweiCloud.
---

# huaweicloud_rgc_organizational_unit_register

Manages an RGC organizational unit register resource within HuaweiCloud.

## Example Usage

```hcl
variable "organizational_unit_id" {}

resource "huaweicloud_rgc_organizational_unit_register" "test"{
  organizational_unit_id = var.organizational_unit_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `organizational_unit_id` - (Required, String, NonUpdatable) Specifies the ID of the organizational unit.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `organizational_unit_name` - Indicates the name of registered organizational unit.

* `parent_organizational_unit_id` - Indicates the ID of registered parent organizational unit.

* `parent_organizational_unit_name` - Indicates the name of registered parent organizational unit.

* `manage_account_id` - Indicates the ID of organization manage account.

* `organizational_unit_status` - Indicates the status of registered parent organizational unit.

* `organizational_unit_type` - Indicates the type of registered parent organizational unit.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 6 hours.

* `delete` - Default is 6 hours.

## Import

The RGC organizational unit register can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rgc_organizational_unit_register.test <id>
```
