---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_recycle_policy"
description: |-
  Manages an ECS recycle policy resource within HuaweiCloud.
---

# huaweicloud_compute_recycle_policy

Manages an ECS recycle policy resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_compute_recycle_policy" "test" {
  retention_hour        = 10
  recycle_threshold_day = 50
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the volume resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `retention_hour` - (Required, Int) Specifies how long an instance can be retained in the recycle bin before being
  permanently deleted.

* `recycle_threshold_day` - (Required, Int) Specifies how long an instance can be moved to the recycle bin after it is created.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the project ID.

## Import

The ECS recycle policy resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_compute_recycle_policy.test <id>
```
