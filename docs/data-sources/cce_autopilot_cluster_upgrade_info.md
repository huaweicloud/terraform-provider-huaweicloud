---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_cluster_upgrade_info"
description: |-
  Use this data source to get the list of CCE Autopilot clusters upgrade info within huaweicloud.
---

# huaweicloud_cce_autopilot_cluster_upgrade_info

Use this data source to get the list of CCE Autopilot clusters upgrade info within huaweicloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cce_autopilot_cluster_upgrade_info" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the Cluster ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `kind` - The API type.

* `api_version` - The API version.

* `metadata` - Upgrade status information.

  The [metadata](#metadata_struct) structure is documented below.

* `spec` - Upgrade configuration information.

  The [spec](#spec_struct) structure is documented below.

* `status` - Upgrade status information.

  The [status](#status_struct) structure is documented below.
  
<a name="metadata_struct"></a>
The `metadata` block supports:

* `uid` - Unique ID identifier.

* `name` - Resource name.

* `labels` - Resource labels, in key/value pair format. These are reserved fields in the interface
  and will not take effect.

* `annotations` - Resource annotations, which consist of key/value pairs.

* `creation_timestamp` - Creation time.

* `update_timestamp` - Update time.

<a name="spec_struct"></a>
The `spec` block supports:

* `last_upgrade_info` - Last cluster upgrade information.

  The [last_upgrade_info](#spec_last_upgrade_info_struct) structure is documented below.

* `version_info` - Version information.

  The [version_info](#spec_version_info_struct) structure is documented below.

* `upgrade_feature_gates` - Cluster upgrade feature switches.

  The [upgrade_feature_gates](#spec_upgrade_feature_gates_struct) structure is documented below.

<a name="spec_last_upgrade_info_struct"></a>
The `last_upgrade_info` block supports:

* `completion_time` - Upgrade task end time.

* `phase` - Upgrade task status. Possible values: Init, Running, Pause, Success, Failed.

* `progress` - Upgrade task progress.

<a name="spec_upgrade_feature_gates_struct"></a>
The `upgrade_feature_gates` block supports:

* `support_upgrade_page_v4` - Whether the cluster upgrade Console interface supports version V4.
  This field is generally used by the CCE Console.

<a name="spec_version_info_struct"></a>
The `version_info` block supports:

* `patch` - Patch version number, e.g., r0.

* `release` - Official version number, e.g., v1.19.10.

* `suggest_patch` - Recommended target patch version number, e.g., r0.

* `target_versions` - Upgrade target version collection.

<a name="status_struct"></a>
The `status` block supports:

* `completion_time` - Upgrade task end time.

* `phase` - Upgrade task status. Possible values: Init, Running, Pause, Success, Failed.

* `progress` - Upgrade task progress.
