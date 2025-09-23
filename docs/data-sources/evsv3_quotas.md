---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evsv3_quotas"
description: |-
  Use this data source to get EVS v3 quota details within HuaweiCloud.
---

# huaweicloud_evsv3_quotas

Use this data source to get EVS v3 quota details within HuaweiCloud.

## Example Usage

```hcl
variable "usage" {}

data "huaweicloud_evsv3_quotas" "test" {
  usage = var.usage
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider - level region will be used.

* `usage` - (Required, Boolean) Specifies whether to query quota details. Only value **true** is supported currently.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID of EVS quota details.

* `quota_set` - The returned quota information.
  The [quota_set](#quota_set_struct) structure is documented below.

<a name="quota_set_struct"></a>
The `quota_set` block supports:

* `backup_gigabytes` - The backup size, in GiB.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `backups` - The number of backups.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `gigabytes` - The total capacity, in GiB.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `snapshots` - The number of snapshots.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `volumes` - The number of disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `gigabytes_sata` - The capacity (GiB) for common I/O disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `snapshots_sata` - The number of snapshots for common I/O disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `volumes_sata` - The number of common I/O disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `gigabytes_sas` - The capacity (GiB) for high I/O disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `snapshots_sas` - The number of snapshots for high I/O disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `volumes_sas` - The number of high I/O disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `gigabytes_ssd` - The capacity (GiB) for ultra - high I/O disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `snapshots_ssd` - The number of snapshots for ultra - high I/O disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `volumes_ssd` - The number of ultra - high I/O disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `gigabytes_gpssd` - The capacity (GiB) for general purpose ssd disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `snapshots_gpssd` - The number of snapshots for general purpose ssd disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `volumes_gpssd` - The number of general purpose ssd disks.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

* `per_volume_gigabytes` - The capacity quota of a disk.
  The [quota_set_sub](#quota_set_sub_struct) structure is documented below.

<a name="quota_set_sub_struct"></a>
The `quota_set_sub` block supports:

* `limit` - Maximum quota.

* `in_use` - Used quota.
