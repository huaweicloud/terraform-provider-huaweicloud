---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_obs_targets"
description: |-
  Use this data source to get the list of the OBS targets.
---

# huaweicloud_sfs_turbo_obs_targets

Use this data source to get the list of the OBS targets.

## Example Usage

```hcl
variable "share_id" {}
variable "target_id" {}

data "huaweicloud_sfs_turbo_obs_targets" "test" {
  share_id  = var.share_id
  target_id = var.target_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `share_id` - (Required, String) Specifies the ID of the SFS Turbo file system to which the OBS target belongs.

* `target_id` - (Optional, String) Specifies the ID of the OBS target.

* `status` - (Optional, String) Specifies the status of the OBS target.
  The valid values are **AVAILABLE**, **MISCONFIGURED** and **FAILED**.

* `bucket` - (Optional, String) Specifies the name of the OBS bucket associated with the OBS target.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `targets` - The list of OBS targets.

  The [targets](#targets_struct) structure is documented below.

<a name="targets_struct"></a>
The `targets` block supports:

* `id` - The ID of the OBS target.

* `file_system_path` - The linkage directory name of the OBS target.

* `status` - The status of the OBS target.

* `obs` - The detail of the OBS bucket.

  The [obs](#targets_obs_struct) structure is documented below.

* `created_at` - The creation time of the OBS target.

<a name="targets_obs_struct"></a>
The `obs` block supports:

* `bucket` - The name of the OBS bucket associated with the OBS target.

* `endpoint` - The domain name of the region where the OBS bucket belongs.
