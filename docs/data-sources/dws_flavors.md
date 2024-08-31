---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_flavors"
description: |-
  Use this data source to get available flavors of DWS cluster node.
---

# huaweicloud_dws_flavors

Use this data source to get available flavors of DWS cluster node.

## Example Usage

```hcl
data "huaweicloud_dws_flavors" "flavor" {
  vcpus = 8
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `availability_zone` - (Optional, String) The availability zone name.

* `vcpus` - (Optional, Int) The vcpus of the dws node flavor.

* `memory` - (Optional, Int) The ram of the dws node flavor in GB.

* `datastore_type` - (Optional, String) The type of datastore.  
  The options are as follows:
    - **dws**: OLAP, elastic scaling, unlimited scaling of compute and storage capacity.
    - **hybrid**: a single data warehouse used for transaction and analytics workloads,
       in single-node or cluster mode.
    - **stream**: built-in time series operators; up to 40:1 compression ratio; applicable to IoT services.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - The list of flavor detail.
  The [Flavors](#dwsFlavors_Flavors) structure is documented below.

<a name="dwsFlavors_Flavors"></a>
The `Flavors` block supports:

* `flavor_id` - The name of the dws node flavor.  
 It is referenced by `node_type` in `huaweicloud_dws_flavors`.

* `datastore_type` - The type of datastore.  
  The options are as follows:
    - **dws**: OLAP, elastic scaling, unlimited scaling of compute and storage capacity.
    - **hybrid**: a single data warehouse used for transaction and analytics workloads,
       in single-node or cluster mode.
    - **stream**: built-in time series operators; up to 40:1 compression ratio; applicable to IoT services.

* `datastore_version` - The version of datastore.

* `vcpus` - The vcpus of the dws node flavor.

* `memory` - The ram of the dws node flavor in GB.

* `volumetype` - Disk type.  
  The options are as follows:
    - **LOCAL_DISK**:common I/O disk.
    - **SSD**: ultra-high I/O disk.

* `size` - The default disk size in GB.

* `availability_zones` - The list of availability zones.

* `elastic_volume_specs` - The [ElasticVolumeSpec](#dwsFlavors_FlavorsElasticVolumeSpec) structure is documented below.

<a name="dwsFlavors_FlavorsElasticVolumeSpec"></a>
The `FlavorsElasticVolumeSpec` block supports:

* `step` - Disk size increment step.

* `min_size` - Minimum disk size.

* `max_size` - Maximum disk size.
