---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_snapshot_policies"
description: |-
  Use this data source to query the list of snapshot policies under specified DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_snapshot_policies

Use this data source to query the list of snapshot policies under specified DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_dws_snapshot_policies" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the DWS cluster ID to which the snapshot policies belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `keep_day` - The number of days to retain the generated automated snapshot.

* `device_name` - The device on which the snapshots are stored.

* `server_ips` - The shared IP addresses of the NFS corresponding to the snapshots.

* `policies` - All automated snapshot policies that match the filter parameters.

  The [policies](#policies_struct) structure is documented below.

<a name="policies_struct"></a>
The `policies` block supports:

* `id` - The ID of the snapshot policy.

* `name` - The name of the snapshot policy.

* `type` - The type of the snapshot policy.
  + **full**
  + **increment**

* `strategy` - The execution strategy of the snapshot.

* `backup_level` - The backup level of the snapshot.

* `next_fire_time` - The start time for doing next snapshot, in RFC3339 format.

* `updated_at` - The latest update time of the snapshot policy, in RFC3339 format.
