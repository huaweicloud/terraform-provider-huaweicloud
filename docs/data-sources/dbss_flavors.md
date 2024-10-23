---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_flavors"
description: ""
---

# huaweicloud_dbss_flavors

Use this data source to get the list of DBSS flavors.

## Example Usage

```hcl
data "huaweicloud_dbss_flavors" "test" {
  level = "high"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `flavor_id` - (Optional, String) Specifies the ID of the flavor.

* `availability_zone` - (Optional, String) Specifies the availability zone which the flavor belongs to.

* `level` - (Optional, String) Specifies the level of the flavor. Value options:
  + **entry**: Starter edition.
  + **low**: Basic edition.
  + **medium**: Professional edition.
  + **high**: Premium edition.

* `memory` - (Optional, Float) Specifies the memory size(GB) in the flavor.

* `vcpus` - (Optional, Int) Specifies the number of CPUs.

* `proxy` - (Optional, Int) Specifies the maximum supported database instances.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Indicates the list of DBSS flavors.
  The [flavors](#DbssFlavors_Flavor) structure is documented below.

<a name="DbssFlavors_Flavor"></a>
The `flavors` block supports:

* `id` - Indicates the ID of the flavor.

* `level` - Indicates the level of the flavor.

* `proxy` - Indicates the maximum supported database instances.

* `vcpus` - Indicates the number of CPUs.

* `memory` - Indicates the memory size(GB) in the flavor.

* `availability_zones` - Indicates the availability zones which the flavor belongs to
