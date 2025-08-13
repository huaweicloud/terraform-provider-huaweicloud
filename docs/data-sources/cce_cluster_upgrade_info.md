---
subcategory: "cce"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster_upgrade_info"
description: |-
  Use this data source to get the CCE cluster upgrade info.
---

# huaweicloud_cce_cluster_upgrade_info

Use this data source to get the CCE cluster upgrade info.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cce_cluster_upgrade_info" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `metadata` - Basic information, as an element type of collection class,
  contains a set of properties defined by different names.

  The [metadata](#metadata_struct) structure is documented below.

* `spec` - Upgrade configuration related information.

  The [spec](#spec_struct) structure is documented below.

* `status` - Upgrade status information.

  The [status](#status_struct) structure is documented below.

<a name="metadata_struct"></a>
The `metadata` block supports:

* `annotations` - Resource annotations, composed of key/value pairs.

* `creation_timestamp` - Creation time.

* `labels` - Resource tags, in key/value pair format, reserved fields for the interface.

* `name` - Resource name.

* `uid` - Unique ID identifier.

* `update_timestamp` - Update time.

<a name="spec_struct"></a>
The `spec` block supports:

* `last_upgrade_info` - Last cluster upgrade information.

  The [last_upgrade_info](#spec_last_upgrade_info_struct) structure is documented below.

* `upgrade_feature_gates` - Cluster upgrade feature flags.

  The [upgrade_feature_gates](#spec_upgrade_feature_gates_struct) structure is documented below.

* `version_info` - Version information.

  The [version_info](#spec_version_info_struct) structure is documented below.

<a name="spec_last_upgrade_info_struct"></a>
The `last_upgrade_info` block supports:

* `completion_time` - Upgrade task end time.

* `phase` - Upgrade task status. The value can be: **Init**, **Running**, **Pause**, **Success** and **Failed**.

* `progress` - Upgrade task progress.

<a name="spec_upgrade_feature_gates_struct"></a>
The `upgrade_feature_gates` block supports:

* `support_upgrade_page_v4` - Whether the cluster upgrade console supports V4 version, generally used by CCE Console.

<a name="spec_version_info_struct"></a>
The `version_info` block supports:

* `patch` - Patch version number, e.g. **r0**.

* `release` - Formal version number, e.g. **v1.19.10**.

* `suggest_patch` - Recommended target patch version for upgrade, e.g. **r0**.

* `target_versions` - Target versions for upgrade.

<a name="status_struct"></a>
The `status` block supports:

* `completion_time` - Upgrade task end time.

* `phase` - Upgrade task status. The value can be: **Init**, **Running**, **Pause**, **Success** and **Failed**.

* `progress` - Upgrade task progress.
