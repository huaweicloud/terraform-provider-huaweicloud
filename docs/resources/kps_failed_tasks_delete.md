---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kps_failed_tasks_delete"
description: |-
  Manages a KPS failed tasks delete resource within HuaweiCloud.
---

# huaweicloud_kps_failed_tasks_delete

Manages a KPS failed tasks delete resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the deleted tasks,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_kps_failed_tasks_delete" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
