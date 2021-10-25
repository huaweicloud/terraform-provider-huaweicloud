---
subcategory: "Data Warehouse Service (DWS)"
---

# huaweicloud_dws_flavors

Use this data source to get available flavors of HuaweiCloud dws cluster node.

## Example Usage

```hcl
data "huaweicloud_dws_flavors" "flavor" {
  vcpus = 8
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the dws cluster client. If omitted, the
  provider-level region will be used.

* `availability_zone` - (Optional, String) Specifies the availability zone name.

* `vcpus` - (Optional, String) Specifies the vcpus of the dws node flavor.

* `memory` - (Optional, String) Specifies the ram of the dws node flavor in GB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID in UUID format.

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `flavor_id` - The name of the dws node flavor. It is referenced by `node_type` in `huaweicloud_dws_flavors`.
* `vcpus` - Indicates the vcpus of the dws node flavor.
* `memory` - Indicates the ram of the dws node flavor in GB.
* `volumetype` - Indicates Disk type.
  + **LOCAL_DISK**: common I/O disk
  + **SSD**: ultra-high I/O disk
* `size` - Indicates the Disk size in GB.
* `availability_zone` - Indicates the availability zone where the node resides.
