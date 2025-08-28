---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_group_sync"
description: |-
  Manages a COC group sync resource within HuaweiCloud.
---

# huaweicloud_coc_group_sync

Manages a COC group sync resource within HuaweiCloud.

~> Deleting group sync resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "group_id" {}

resource "huaweicloud_coc_group_sync" "test" {
  group_id = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, NonUpdatable) Specifies the group ID.

* `cloud_service_name` - (Optional, String, NonUpdatable) Specifies the resource provider.
  The value can be **ecs**, **cce**, **rds** and so on.

* `type` - (Optional, String, NonUpdatable) Specifies the cloud resource type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The group ID, which equals to `group_id`.
