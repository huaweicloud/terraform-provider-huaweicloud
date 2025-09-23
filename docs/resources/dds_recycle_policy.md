---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_recycle_policy"
description: |-
  Manages a DDS recycle policy resource within HuaweiCloud.
---

# huaweicloud_dds_recycle_policy

Manages a DDS recycle policy resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_dds_recycle_policy" "test" {
  retention_period_in_days = 7
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `retention_period_in_days` - (Required, Int) Specifies the policy retention duration in days.
  Value ranges from **1** to **7**.

  -> `retention_period_in_days` defaults to **7**. Deleting recycle policy is unsupported. The resource is removed from
  the state, and the retention period is reset to **7**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
