---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_export_acl_rule"
description: |-
  Manages a resource to export ACL rules within HuaweiCloud.
---

# huaweicloud_cfw_export_acl_rule

Manages a resource to export ACL rules within HuaweiCloud.

-> 1. This resource is a one-time action resource used to export ACL rule. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "object_id" {}

resource "huaweicloud_cfw_export_acl_rule" "test" {
  object_id = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `object_id` - (Required, String, NonUpdatable) Specifies the protection object ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `object_id`.
