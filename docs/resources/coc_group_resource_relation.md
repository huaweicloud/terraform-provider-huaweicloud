---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_group_resource_relation"
description: |-
  Manages a COC group resource relation resource within HuaweiCloud.
---

# huaweicloud_coc_group_resource_relation

Manages a COC group resource relation resource within HuaweiCloud.

## Example Usage

```hcl
variable "group_id" {}
variable "cmdb_resource_id" {}

resource "huaweicloud_coc_group_resource_relation" "test" {
  group_id         = var.group_id
  cmdb_resource_id = var.cmdb_resource_id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, NonUpdatable) Specifies the group ID.

* `cmdb_resource_id` - (Required, String, NonUpdatable) Specifies the CMDB resource ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
