---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_snapshot_flavors"
description: |-
  Use this data source to query the DWS cluster snapshot flavors within HuaweiCloud.
---

# huaweicloud_dws_cluster_snapshot_flavors

Use this data source to query the DWS cluster snapshot flavors within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "snapshot_id" {}

data "huaweicloud_dws_cluster_snapshot_flavors" "test" {
  snapshot_id = var.snapshot_id
}
```

### Query with filters

```hcl
variable "snapshot_id" {}

data "huaweicloud_dws_cluster_snapshot_flavors" "test" {
  snapshot_id          = var.snapshot_id
  type                 = "restore"
  fine_grained_restore = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the snapshot flavors are located.  
  If omitted, the provider-level region will be used.

* `snapshot_id` - (Required, String) Specifies the ID of the snapshot.

* `type` - (Optional, String) Specifies the type of the snapshot flavor.  
  The valid values are as follows:
  + **snapshot** - Only query the flavor information of the snapshot.
  + **restore** - Query both the flavor information of the snapshot and the flavors available for restoration.

* `available_zones` - (Optional, String) Specifies the availability zone code for restoration.  
  When restoring a **3AZ** cluster, you need to pass three availability zone codes, separated by commas (no spaces).

* `fine_grained_restore` - (Optional, Bool) Specifies whether it is fine-grained backup restoration.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - The list of snapshot flavors that matched filter parameters.  
  The [flavors](#cluster_snapshot_flavors_attr) structure is documented below.

<a name="cluster_snapshot_flavors_attr"></a>
The `flavors` block supports:

* `id` - The ID of the flavor.

* `classify` - The classification of the flavor.

* `version` - The version of the flavor.  
  + **v1.0**
  + **v2.0**

* `default_node` - The default node number of the flavor.

* `max_node` - The maximum node number of the flavor.

* `flavor_id` - The underlying flavor ID. This is different from the `id` field and is generally not used.

* `code` - The code of the flavor.

* `status` - The status of the flavor.

* `attributes` - The attributes of the flavor.  
  The [attributes](#cluster_snapshot_flavors_attributes) structure is documented below.

* `min_node` - The minimum node number of the flavor.

* `flavor_code` - The underlying flavor code of the flavor.

* `product_versions` - The product version list of the flavor.  
  The [product_versions](#cluster_snapshot_flavors_product_versions) structure is documented below.

* `volume_num` - The volume number of the flavor.

* `default_capacity` - The default capacity of the flavor.

* `scenario` - The scenario of the flavor.

* `duplicate` - The duplicate of the flavor.

* `volume_used` - The volume used information of the flavor.  
  The [volume_used](#cluster_snapshot_flavors_volume_used) structure is documented below.

<a name="cluster_snapshot_flavors_attributes"></a>
The `attributes` block supports:

* `code` - The code of the attribute.

* `value` - The value of the attribute.

<a name="cluster_snapshot_flavors_product_versions"></a>
The `product_versions` block supports:

* `datastore_version` - The datastore version of the product version.

* `min_cn` - The minimum CN number of the product version.

* `max_cn` - The maximum CN number of the product version.

* `version_type` - The version type of the product version.
  + **1** - Stable version.
  + **0** - Latest version.

<a name="cluster_snapshot_flavors_volume_used"></a>
The `volume_used` block supports:

* `volume_type` - The volume type of the volume used.
  + **HIGH** - SAS disk.
  + **ULTRAHIGH** - SSD cloud disk.
  + **COMMON** - SATA disk.
  + **LOCAL_DISK** - Local disk.

* `volume_num` - The volume number of the volume used.

* `capacity` - The capacity of the volume used.

* `volume_size` - The volume size of the volume used.
