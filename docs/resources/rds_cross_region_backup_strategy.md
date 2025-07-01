---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_cross_region_backup_strategy"
description: ""
---

# huaweicloud_rds_cross_region_backup_strategy

Manages RDS cross-region backup strategy resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "destination_region" {}
variable "destination_project_id" {}

resource "huaweicloud_rds_cross_region_backup_strategy" "test" {
  instance_id            = var.instance_id
  backup_type            = "all"
  keep_days              = 5
  destination_region     = var.destination_region
  destination_project_id = var.destination_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance.

* `backup_type` - (Required, String) Specifies the backup type. Value options:
    + **auto**: open automated full backup.
    + **all**: open both automated full backup and automated incremental backup.

  Only **all** is supported for SQL server.

* `keep_days` - (Required, Int) Specifies the number of days to retain the generated backup files.
  Value ranges from `1` to `1,825`.

* `destination_region` - (Required, String, NonUpdatable) Specifies the target region ID for the cross-region backup policy.

* `destination_project_id` - (Required, String, NonUpdatable) Specifies the target project ID for the cross-region backup
  policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The RDS cross-region backup strategy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rds_cross_region_backup_strategy.test <id>
```
