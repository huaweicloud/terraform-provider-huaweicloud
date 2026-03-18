---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_ai_ops_setting"
description: |-
  Manages an ai-ops setting resource within HuaweiCloud.
---

# huaweicloud_css_ai_ops_setting

Manages an ai-ops setting resource within HuaweiCloud.

-> This resource only supports elasticsearch or opensearch engine.

## Example Usage

```hcl
variable "cluster_id" {}
variable "check_type" {}
variable "period" {}

resource "huaweicloud_css_ai_ops_setting" "test" {
  cluster_id = var.cluster_id
  check_type = var.check_type
  period     = var.period
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the cluster resource. If omitted, the
  provider-level region will be used. Changing this creates a new cluster resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the cluster ID.

* `check_type` - (Required, String) Specifies the check type.
  The valid values are as follows:
  + **full_detection**: Indicates all check items.
  + **unavailability_detection**: Indicates cluster unavailability check items.
  + **partial_detection**: Indicates partial check items.

* `period` - (Required, String) Specifies the intelligent O&M automatic check time,
  The format is **HH:mm GTM+08:00**, e.g. **00:00 GMT+08:00**.

* `check_items` - (Optional, List) Specifies the ID list of partial check items.
  This parameter is required when the `check_type` is set to **partial_detection**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also is the cluster ID.

## Import

The ai-ops setting can be imported by `id`, e.g.

```bash
$ terraform import huaweicloud_css_ai_ops_setting.test <id>
```
