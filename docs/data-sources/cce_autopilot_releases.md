---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_releases"
description: |- 
  Use this data source to get the list of Autopilot CCE releases within HuaweiCloud.
---

# huaweicloud_cce_autopilot_releases

Use this data source to get the list of Autopilot CCE releases within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cce_autopilot_releases" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the Clsuter ID

* `chart_id` - (Optional, String) Specifies the Chart ID

* `namespace` - (Optional, String) Specifies the Namespace corresponding to the template

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID.

* `releases` - The releases data in the cce cluster.

  The [object](#releases) structure is documented below.

<a name="releases"></a>
The `releases` block supports:

* `chart_name` - The chart name.

* `chart_public` - The chart is public or not.

* `chart_version` - The chart version.

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `description` - The release description.

* `name` - The release name.

* `namespace` - The namespace of release.

* `parameters` - The release parameters.

* `resources` - The release resources.

* `status` - The release status.

* `status_description` - The release status description.

* `create_at` - The create time of release.

* `update_at` - The update time of release.

* `values` - The release values.

* `version` - The release version.
