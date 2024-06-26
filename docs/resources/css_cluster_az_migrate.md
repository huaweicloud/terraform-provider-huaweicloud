---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_cluster_az_migrate"
description: |-
  Manages CSS cluster az migrate resource within HuaweiCloud.
---

# huaweicloud_css_cluster_az_migrate

Manages CSS cluster az migrate resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "agency" {}

resource "huaweicloud_css_cluster_az_migrate" "test" {
  cluster_id           = var.cluster_id
  instance_type        = "all"
  source_az            = "cn-south-1e"
  target_az            = "cn-south-1e,cn-south-1c"
  migrate_type         = "multi_az_change"
  agency               = var.agency
  indices_backup_check = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String) Specifies the ID of the CSS cluster.

* `instance_type` - (Required, String, NonUpdatable) Specifies the node type of the AZ to be switched.
  The value can be **ess**, **ess-cold**, **ess-master**, **ess-client** and **all**.

* `source_az` - (Required, String, NonUpdatable) Specifies the AZ where the node is currently located.

* `target_az` - (Required, String, NonUpdatable) Specifies the AZ where the node is finally distributed.

* `migrate_type` - (Required, String, NonUpdatable) Specifies the migration type of AZ.
  The value can be **multi_az_change** and **az_migrate**.

* `agency` - (Required, String, NonUpdatable) Specifies the IAM agency used to access CSS.

* `indices_backup_check` - (Optional, Bool, NonUpdatable) Specifies whether to perform backup verification.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
