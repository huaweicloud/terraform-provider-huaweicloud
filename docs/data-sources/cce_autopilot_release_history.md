---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_release_history"
description: |-
  Use this data source to get the list of historical records of a CCE Autopilot chart release within HuaweiCloud.
---

# huaweicloud_cce_autopilot_release_history

Use this data source to get the list of historical records of a CCE Autopilot chart release within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "name" {}
variable "namespace" {}

data "huaweicloud_cce_autopilot_release_history" "test" {
  cluster_id = var.cluster_id
  name       = var.name
  namespace  = var.namespace
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the CCE Autopilot cluster.

* `name` - (Required, String) Specifies the name of the CCE Autopilot chart release.

* `namespace` - (Required, String) Specifies the namespace to which the CCE Autopilot chart release belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID.

* `releases` - The releases data in the cce cluster.

  The [releases](#releases_struct) structure is documented below.

<a name="releases_struct"></a>
The `releases` block supports:

* `chart_name` - The name of the chart.

* `chart_public` - Whether the chart is public.

* `chart_version` - The version of the chart.

* `cluster_id` - The ID of the CCE Autopilot cluster.

* `cluster_name` - The name of the CCE Autopilot cluster.

* `create_at` - The creation time of the chart release.

* `description` - The description of the chart release.

* `name` - The name of the chart release.

* `namespace` - The namespace to which the chart release belongs.

* `parameters` - The parameters of the chart release in JSON format.

* `resources` - The resources required by the chart release in JSON format.

* `status` - The status of the chart release.
  The valid values are as follows:
  + **DEPLOYED**: The release is normal.
  + **DELETED**: The release has been deleted.
  + **FAILED**: The release fails to be deployed.
  + **DELETING**: The release is being deleted.
  + **PENDING_INSTALL**: The release is waiting to be installed.
  + **PENDING_UPGRADE**: The release is waiting to be upgraded.
  + **PENDING_ROLLBACK**: The release is waiting for rollback.
  + **UNKNOWN**: The release status is unknown, indicating that the release is abnormal.

* `status_description` - The status description of the chart release.

* `update_at` - The update time of the chart release.

* `values` - The values of the chart release in JSON format.

* `version` - The version number of the chart release.
