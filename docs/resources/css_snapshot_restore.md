---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_snapshot_restore"
description: |-
  Manages CSS cluster snapshot restore resource within HuaweiCloud.
---

# huaweicloud_css_snapshot_restore

Manages CSS cluster snapshot restore resource within HuaweiCloud.

## Example Usage

```hcl
variable "target_cluster_id" {}
variable "source_cluster_id" {}
variable "snapshot_id" {}

resource "huaweicloud_css_snapshot_restore" "test" {
  source_cluster_id = var.target_cluster_id
  target_cluster_id = var.source_cluster_id
  snapshot_id       = var.snapshot_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `target_cluster_id` - (Required, String, NonUpdatable) Specifies the target cluster ID.

* `source_cluster_id` - (Required, String, NonUpdatable) Specifies the source cluster ID.

* `snapshot_id` - (Required, String, NonUpdatable) Specifies the ID of the snapshot to be restored.

* `indices` - (Optional, String, NonUpdatable) Name of an index to be restored. Multiple indexes are separated by
  commas (,). By default, all indexes are restored.You can use `*` to match multiple indexes. For example, if you
  specify `2018-06*`, then the data of the indexes with the prefix 2018-06 will be restored. The value can contain
  **0** to **1,024** characters. Uppercase letters, spaces, and the following special characters are not allowed:
  **"\<|>/?**.

* `rename_pattern` - (Optional, String, NonUpdatable) Rule for defining the indexes to be restored. The value can
  contain `0` to `1,024` characters. Uppercase letters, spaces, and the following special characters are not allowed:
  **"\<|>/?**. Indexes that match this rule will be restored. The filtering condition must be a regular expression.

* `rename_replacement` - (Optional, String, NonUpdatable) Rule for renaming an index. The value can contain **0** to
  **1,024** characters. Uppercase letters, spaces, and the following special characters are not allowed: **"\<|>/?**.
  For example, **restored_index_$1** indicates adding the **restored_index_** prefix to the names of all the restored
  indexes. The `rename_replacement` parameter takes effect only if rename_pattern has been enabled.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
