---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_flavors"
description: ""
---

# huaweicloud_cbh_flavors

Use this data source to get the list of CBH specifications.

## Example Usage

```hcl
data "huaweicloud_cbh_flavors" "test" {
  type = "basic"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `action` - (Optional, String) Specifies the action of querying instances specification information.
  The valid values are as follows:
  + **create**: Query instance specification information that can be created.
  + **update**: Query instance specification information that can be updated.

  If omitted, the CBH specifications that can be created will be queried.

* `spec_code` - (Optional, String) Specifies the ID of the CBH specification, the query result shows all specifications
  that can be changed by this specification. This parameter is required when `action` is set to **update**.

* `flavor_id` - (Optional, String) Specifies the ID of the specification of CBH.
  At present, CBH provides two functional versions: standard version and professional version.
  The standard version is equipped with asset specifications of 10(for example the `flavor_id` is: **cbh.basic.10**),
  20, 50, 100, 200, 500, 1000, 2000, 5000, and 10000.
  The professional version is equipped with 10(for example the `flavor_id` is: **cbh.enhance.10**),
  20, 50, 100, 200, 500, 1000, 2000, 5000, 10000 asset specifications.
  The specification 'enhance' is more advanced than the specification 'basic'.

* `type` - (Optional, String) Specifies the type of CBH specification. The value can be:
  + **basic**: Standard version.
  + **enhance**: Professional version.

* `asset` - (Optional, Int) Specifies the number of CBH assets.

* `memory` - (Optional, Int) Specifies the memory size of the CBH, in GB.

* `vcpus` - (Optional, Int) Specifies the number of CPU cores of the CBH.

* `max_connection` - (Optional, Int) Specifies the maximum number of connections to the CBH.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Indicates the list of CBH specification.
  The [flavor](#CbhFlavors_flavor) structure is documented below.

<a name="CbhFlavors_flavor"></a>
The `flavor` block supports:

* `id` - Indicates the ID of the specification.

* `ecs_system_data_size` - The disk size of the CBH system disk, in GB.

* `vcpus` - The number of CPU cores of the CBH.

* `memory` - The memory size of the CBH, in GB.

* `asset` - The number of CBH assets.

* `max_connection` - The maximum number of connections to the CBH.

* `type` - The type of CBH specification. The value can be:
  + **basic**: Standard version.
  + **enhance**: Professional version.

* `data_disk_size` - The size of the CBH data disk, in TB.
