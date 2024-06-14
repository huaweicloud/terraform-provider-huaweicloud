---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_es_core_upgrade"
description: |-
  Manages CSS ElasticSearch core upgrade resource within HuaweiCloud.
---

# huaweicloud_css_es_core_upgrade

Manages CSS ElasticSearch core upgrade resource within HuaweiCloud.

-> **NOTE:** After the upgrade is successful, the `engine_version` field of the managed CSS cluster
(huaweicloud_css_cluster) has changed. You need to manually synchronize the field value in the script,
otherwise **forceNew** will be triggered.

## Example Usage

```hcl
variable "cluster_id" {}
variable "target_image_id" {}
variable "agency" {}

resource "huaweicloud_css_es_core_upgrade" "test" {
  cluster_id           = var.cluster_id
  target_image_id      = var.target_image_id
  upgrade_type         = "cross"
  agency               = var.agency
  indices_backup_check = true
  cluster_load_check   = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the CSS cluster.

* `target_image_id` - (Required, String, NonUpdatable) Specifies the upgradeable target image ID.

* `upgrade_type` - (Required, String, NonUpdatable) Specifies the upgrade type.
  The value can be **same**, **cross** and **crossEngine**.

* `agency` - (Required, String, NonUpdatable) Specifies the IAM agency used to access CSS.

* `indices_backup_check` - (Required, Bool, NonUpdatable) Specifies whether to perform backup verification.

* `cluster_load_check` - (Optional, Bool, NonUpdatable) Whether to verify the load. Default is **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `upgrade_detail` - The upgrade task detail.
  The [upgrade_detail](#css_upgrade_detail) structure is documented below.

<a name="css_upgrade_detail"></a>
The `upgrade_detail` block supports:

* `id` The job ID of the upgrade task.

* `start_time` - The start time.

* `end_time` - The end time.

* `status` - The status.

* `agency` - The IAM agency used to access CSS.

* `total_nodes` - The all nodes.

* `retry_times` - The retry times.

* `datastore` - The data store.
  The [datastore](#css_upgrade_detail_datastore) structure is documented below.

<a name="css_upgrade_detail_datastore"></a>
The `datastore` block supports:

* `type` - The type of the data store.

* `version` - The version of the data store.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
