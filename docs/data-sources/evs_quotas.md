---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_quotas"
description: |-
  Use this data source to query the user's quotas of EVS within HuaweiCloud.
---

# huaweicloud_evs_quotas

Use this data source to query the user's quotas of EVS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_evs_quotas" "test" {
  usage = "True"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `usage` - (Required, String) Specifies whether to query quota details. Only value **True** is supported currently.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quota_set` - The returned quota information.

  The [quota_set](#quota_set_struct) structure is documented below.

<a name="quota_set_struct"></a>
The `quota_set` block supports:

* `backups` - The number of backups.

  The [backups](#evs_quota_sub_struct) structure is documented below.

* `snapshots_ssd` - The number of snapshots for ultra-high I/O disks.

  The [snapshots_ssd](#evs_quota_sub_struct) structure is documented below.

* `volumes_gpssd` - The number of general purpose SSD disks.

  The [volumes_gpssd](#evs_quota_sub_struct) structure is documented below.

* `volumes` - The number of disks.

  The [volumes](#evs_quota_sub_struct) structure is documented below.

* `snapshots_sata` - The number of snapshots for common I/O disks.

  The [snapshots_sata](#evs_quota_sub_struct) structure is documented below.

* `gigabytes_sata` - The capacity (GiB) for common I/O disks.

  The [gigabytes_sata](#evs_quota_sub_struct) structure is documented below.

* `volumes_sata` - The number of common I/O disks.

  The [volumes_sata](#evs_quota_sub_struct) structure is documented below.

* `snapshots_sas` - The number of snapshots for high I/O disks.

  The [snapshots_sas](#evs_quota_sub_struct) structure is documented below.

* `volumes_ssd` - The number of ultra-high I/O disks.

  The [volumes_ssd](#evs_quota_sub_struct) structure is documented below.

* `backup_gigabytes` - The backup size, in GiB.

  The [backup_gigabytes](#evs_quota_sub_struct) structure is documented below.

* `gigabytes` - The total capacity, in GiB.

  The [gigabytes](#evs_quota_sub_struct) structure is documented below.

* `gigabytes_sas` - The capacity (GiB) for high I/O disks.

  The [gigabytes_sas](#evs_quota_sub_struct) structure is documented below.

* `volumes_sas` - The number of high I/O disks.

  The [volumes_sas](#evs_quota_sub_struct) structure is documented below.

* `gigabytes_ssd` - The capacity (GiB) for ultra-high I/O disks.

  The [gigabytes_ssd](#evs_quota_sub_struct) structure is documented below.

* `gigabytes_gpssd` - The capacity (GiB) for general purpose SSD disks.

  The [gigabytes_gpssd](#evs_quota_sub_struct) structure is documented below.

* `snapshots_gpssd` - The number of snapshots for general purpose SSD disks.

  The [snapshots_gpssd](#evs_quota_sub_struct) structure is documented below.

* `per_volume_gigabytes` - The capacity quota of a disk.

  The [per_volume_gigabytes](#evs_quota_sub_struct) structure is documented below.

* `id` - The project ID.

* `snapshots` - The number of snapshots.

  The [snapshots](#evs_quota_sub_struct) structure is documented below.

* `volumes_essd` - The number of extreme SSD disks.

  The [volumes_essd](#evs_quota_sub_struct) structure is documented below.

* `gigabytes_essd` - The capacity (GiB) for extreme SSD disks.

  The [gigabytes_essd](#evs_quota_sub_struct) structure is documented below.

* `snapshots_essd` - The number of snapshots for extreme SSD disks.

  The [snapshots_essd](#evs_quota_sub_struct) structure is documented below.

<a name="evs_quota_sub_struct"></a>
The `backups`, `snapshots_ssd`, `volumes_gpssd`, `volumes`, `snapshots_sata`, `gigabytes_sata`, `volumes_sata`,
`snapshots_sas`, `volumes_ssd`, `backup_gigabytes`, `gigabytes`, `gigabytes_sas`, `volumes_sas`, `gigabytes_ssd`,
`gigabytes_gpssd`, `snapshots_gpssd`, `per_volume_gigabytes`, `snapshots`, `volumes_essd`, `gigabytes_essd`, and
`snapshots_essd` blocks supports:

* `in_use` - The used quota.

* `limit` - The maximum quota.

* `reserved` - The reserved field.
